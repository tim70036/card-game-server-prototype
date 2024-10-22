package service

import (
	"card-game-server-prototype/pkg/core"
	model2 "card-game-server-prototype/pkg/game/txpoker/model"
	action2 "card-game-server-prototype/pkg/game/txpoker/type/action"
	"card-game-server-prototype/pkg/game/txpoker/type/seatstatus"
	"card-game-server-prototype/pkg/game/txpoker/type/stage"
	"card-game-server-prototype/pkg/util"

	"github.com/samber/lo"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ActionHintService struct {
	actionHintGroup *model2.ActionHintGroup
	seatStatusGroup *model2.SeatStatusGroup
	playerGroup     *model2.PlayerGroup
	chipCacheGroup  *model2.ChipCacheGroup
	gameSetting     *model2.GameSetting
	table           *model2.Table
	replay          *model2.Replay

	logger *zap.Logger
}

func ProvideActionHintService(
	actionHintGroup *model2.ActionHintGroup,
	seatStatusGroup *model2.SeatStatusGroup,
	playerGroup *model2.PlayerGroup,
	chipCacheGroup *model2.ChipCacheGroup,
	gameSetting *model2.GameSetting,
	table *model2.Table,
	replay *model2.Replay,

	loggerFactory *util.LoggerFactory,
) *ActionHintService {
	return &ActionHintService{
		actionHintGroup: actionHintGroup,
		seatStatusGroup: seatStatusGroup,
		playerGroup:     playerGroup,
		chipCacheGroup:  chipCacheGroup,
		gameSetting:     gameSetting,
		table:           table,
		replay:          replay,
		logger:          loggerFactory.Create("ActionHintService"),
	}
}

func (s *ActionHintService) Pass(passUid core.Uid, passType action2.ActionType) error {
	if !lo.Contains([]action2.ActionType{action2.Fold, action2.Check}, passType) {
		return status.Errorf(codes.InvalidArgument, "invalid action type %s", passType.String())
	}

	if passType == action2.Check && s.actionHintGroup.RaiserHint != nil &&
		s.actionHintGroup.RaiserHint.Action != action2.BB { // BB can check itself.
		return status.Errorf(codes.FailedPrecondition, "cannot check when there is a raiser")
	}

	actionHint, ok := s.actionHintGroup.Hints[passUid]
	if !ok {
		return status.Errorf(codes.NotFound, "uid %s not found in action hint group", passUid)
	}

	actionHint.Action = passType

	var (
		actionRecord action2.ActionRecord = nil
		role                              = s.playerGroup.Data[actionHint.Uid].Role
		curBetStage                       = s.table.BetStageFSM.MustState().(stage.Stage)
	)

	switch passType {
	case action2.Fold:
		actionRecord = &action2.FoldRecord{Uid: actionHint.Uid, Role: role}
	case action2.Check:
		actionRecord = &action2.CheckRecord{Uid: actionHint.Uid, Role: role}
	}

	s.replay.ActionLog[curBetStage] = append(
		s.replay.ActionLog[curBetStage],
		actionRecord,
	)

	return nil
}

func (s *ActionHintService) OpenBet(openerUid core.Uid, betType action2.ActionType, betChip int) error {
	if !lo.Contains([]action2.ActionType{action2.Bet, action2.Raise, action2.BB, action2.AllIn}, betType) {
		return status.Errorf(codes.InvalidArgument, "invalid action type %s", betType.String())
	}

	actionHint, seatStatus, err := s.fetch(openerUid)
	if err != nil {
		return err
	}

	if betChip <= 0 || betChip > seatStatus.Chip {
		return status.Errorf(codes.InvalidArgument, "invalid bet chip %d with seat status chip %d", betChip, seatStatus.Chip)
	}

	if betChip < s.gameSetting.BigBlind && betType != action2.AllIn {
		return status.Errorf(codes.InvalidArgument, "invalid bet chip %d smaller then bb %d", betChip, s.gameSetting.BigBlind)
	}

	// If call/bet/raise and the remaining seat status chip == 0 then it's all in. (Edge case: BB with 0 chip cannot do any other action.)
	if betChip == seatStatus.Chip && betType != action2.AllIn && betType != action2.BB {
		return status.Errorf(codes.InvalidArgument, "invalid action type %s, should be all in with bet chip %d and seat status chip %d", betType.String(), betChip, seatStatus.Chip)
	}

	if betType == action2.Raise && betChip < actionHint.MinRaisingChip {
		return status.Errorf(codes.InvalidArgument, "invalid raise chip %d smaller then min raising chip %d with calling chip %d", betChip, actionHint.MinRaisingChip, actionHint.CallingChip)
	}

	// Edge case: pay BB but is all in.
	if betType == action2.BB && betChip == seatStatus.Chip {
		actionHint.IsBBAllIn = true
	}

	seatStatus.Chip -= betChip
	s.chipCacheGroup.SeatStatusChips[openerUid] -= betChip
	actionHint.BetChip += betChip
	actionHint.Action = betType

	// Record raise chip for the following min raise check.
	if betType == action2.BB || s.actionHintGroup.RaiserHint == nil {
		// For first bet (includes BB), raise chip is equal to bet chip.
		actionHint.RaiseChip = actionHint.BetChip
	} else {
		actionHint.RaiseChip = actionHint.BetChip - s.actionHintGroup.RaiserHint.BetChip
	}

	var (
		actionRecord action2.ActionRecord = nil
		role                              = s.playerGroup.Data[actionHint.Uid].Role
		curBetStage                       = s.table.BetStageFSM.MustState().(stage.Stage)
	)

	switch betType {
	case action2.Bet:
		actionRecord = &action2.BetRecord{Uid: actionHint.Uid, Role: role, Chip: betChip}
	case action2.Raise:
		actionRecord = &action2.RaiseRecord{Uid: actionHint.Uid, Role: role, Chip: betChip}
	case action2.BB:
		actionRecord = &action2.BBRecord{Uid: actionHint.Uid, Role: role, Chip: betChip}
	case action2.AllIn:
		actionRecord = &action2.AllInRecord{
			Uid:     actionHint.Uid,
			Role:    role,
			BetType: lo.Ternary(s.actionHintGroup.RaiserHint == nil, action2.Bet, action2.Raise),
			Chip:    betChip,
		}
	}

	s.replay.ActionLog[curBetStage] = append(
		s.replay.ActionLog[curBetStage],
		actionRecord,
	)

	s.changeRaiser(actionHint)

	return nil
}

// All other players still in the pot must either call the full
// amount of the bet or raise if they wish to remain in, the only
// exceptions being when a player does not have sufficient stake
// remaining to call the full amount of the bet (in which case
// they may either call with their remaining stake to go "all-in"
// or fold) or when the player is already all-in.
func (s *ActionHintService) FollowBet(followUid core.Uid, followType action2.ActionType) error {
	if !lo.Contains([]action2.ActionType{action2.Call, action2.AllIn, action2.SB, action2.BB}, followType) {
		return status.Errorf(codes.InvalidArgument, "invalid action type %s", followType.String())
	}

	if s.actionHintGroup.RaiserHint == nil {
		return status.Errorf(codes.FailedPrecondition, "no raiser yet, cannot follow bet with action type %s", followType.String())
	}

	actionHint, seatStatus, err := s.fetch(followUid)
	if err != nil {
		return err
	}

	followingChip := lo.If(followType == action2.Call, actionHint.CallingChip).
		ElseIf(followType == action2.AllIn, seatStatus.Chip).
		ElseIf(followType == action2.SB, s.gameSetting.SmallBlind).
		ElseIf(followType == action2.BB, s.gameSetting.BigBlind).
		Else(0)

	if followingChip <= 0 || followingChip > seatStatus.Chip {
		return status.Errorf(codes.InvalidArgument, "invalid following chip %d with seat status chip %d", followingChip, seatStatus.Chip)
	}

	// If call/bet/raise and the remaining seat status chip == 0 then it's all in. (Edge case: BB with 0 chip cannot do any other action.)
	if followingChip == seatStatus.Chip && followType != action2.AllIn && followType != action2.BB {
		return status.Errorf(codes.InvalidArgument, "invalid action type %s, should be all in with following chip %d and seat status chip %d", followType.String(), followingChip, seatStatus.Chip)
	}

	seatStatus.Chip -= followingChip
	s.chipCacheGroup.SeatStatusChips[followUid] -= followingChip
	actionHint.BetChip += followingChip
	actionHint.Action = followType

	var (
		actionRecord action2.ActionRecord = nil
		role                              = s.playerGroup.Data[actionHint.Uid].Role
		curBetStage                       = s.table.BetStageFSM.MustState().(stage.Stage)
	)

	switch followType {
	case action2.Call:
		actionRecord = &action2.CallRecord{Uid: actionHint.Uid, Role: role, Chip: followingChip}
	case action2.AllIn:
		actionRecord = &action2.AllInRecord{Uid: actionHint.Uid, Role: role, BetType: action2.Call, Chip: followingChip}
	case action2.SB:
		actionRecord = &action2.SBRecord{Uid: actionHint.Uid, Role: role, Chip: followingChip}
	case action2.BB:
		actionRecord = &action2.BBRecord{Uid: actionHint.Uid, Role: role, Chip: followingChip}
	}

	s.replay.ActionLog[curBetStage] = append(
		s.replay.ActionLog[curBetStage],
		actionRecord,
	)

	return nil
}

// All in is very special. It could be a follow or a open bet and even
// cause incomplete raise:
//
// If all in >= min raising chip or no current bet exists then it's
// opening bet round.
//
// If all in > cur bet chip but < min raising chip then it's following
// a bet but causes an incomplete raise.
//
// If all in <= cur bet chip then it's simply following a bet.
//
// If call/bet/raise and the remaining seat status chip == 0 then it
// should be an all in action.
func (s *ActionHintService) AllIn(allInUid core.Uid) error {
	actionHint, seatStatus, err := s.fetch(allInUid)
	if err != nil {
		return err
	}

	allInChip := seatStatus.Chip
	if s.actionHintGroup.RaiserHint == nil ||
		allInChip >= actionHint.MinRaisingChip {
		return s.OpenBet(allInUid, action2.AllIn, allInChip)
	}

	// Incomplete raise. If a player raises and that puts him all-in,
	// i.e. he bets his whole stack, but does not have enough money or
	// chips to make the minimum raise, then his raise is called
	// incomplete or irregular.
	//
	// An incomplete raise does not reopen the betting, i.e. players
	// who is the current raiser (RaiserHint) cannot raise further, he
	// can only call the incomplete raise or fold.
	if s.actionHintGroup.RaiserHint != nil &&
		allInChip > s.actionHintGroup.RaiserHint.BetChip &&
		allInChip < actionHint.MinRaisingChip {

		if err := s.FollowBet(allInUid, action2.AllIn); err != nil {
			return err
		}

		// All other player that does not all-in/fold should call the
		// this incomplete raise. (including current raiser). Note
		// that unless we complete the raise (see the following), the
		// current raiser can only call cannot raise further.
		liveHints := s.fetchLiveHints()
		for _, liveHint := range liveHints {
			liveHint.Action = action2.Undefined
		}

		// Edge case: When there is an incomplete raise going on, some
		// player's bet chip will be greater than raiser's bet chip. In
		// this case, we need to constantly check if the total excess bet
		// chip is greater than raiser's chip. If so, we said the current
		// caller complete the raise and he becomes the new raiser.
		//
		// Note that, under normal case, the total excess bet chip should
		// be 0, since everyone should call the same amount or perform
		// all-in/SB that has less amount of bet chip than raiser.
		totalExcessBetChip := lo.Sum(lo.FilterMap(lo.Values(s.actionHintGroup.Hints), func(h *model2.ActionHint, _ int) (int, bool) {
			excessRaiseChip := h.BetChip - s.actionHintGroup.RaiserHint.BetChip
			return excessRaiseChip, excessRaiseChip > 0
		}))

		if totalExcessBetChip > s.actionHintGroup.RaiserHint.RaiseChip {
			s.changeRaiser(actionHint)
		}

		return nil
	}

	if s.actionHintGroup.RaiserHint != nil &&
		allInChip <= s.actionHintGroup.RaiserHint.BetChip {
		return s.FollowBet(allInUid, action2.AllIn)
	}

	return status.Errorf(codes.Internal, "unhandled case with all in chip %d", allInChip)
}

func (s *ActionHintService) Eval() {
	for _, hint := range s.actionHintGroup.Hints {
		hint.AvailableActions = []action2.ActionType{}
		hint.CallingChip = 0
		hint.MinRaisingChip = 0
	}

	liveHints := s.fetchLiveHints()
	actionableHints := lo.Filter(liveHints, func(h *model2.ActionHint, _ int) bool {
		return h.Action == action2.Undefined ||
			h.Action == action2.SB ||
			h.Action == action2.BB
	})

	for _, hint := range actionableHints {
		hint.AvailableActions = append(hint.AvailableActions, action2.Fold)
	}

	// Some players might not be in playing state. They can only fold.
	// All other playing players can fold and all in.
	playingHints := lo.Filter(actionableHints, func(h *model2.ActionHint, _ int) bool {
		status, ok := s.seatStatusGroup.Status[h.Uid]
		return ok && status.FSM.MustState().(seatstatus.SeatStatusState) == seatstatus.PlayingState
	})

	for _, hint := range playingHints {
		hint.AvailableActions = append(hint.AvailableActions, action2.AllIn)
	}

	// If current betting round (stage) is not opened yet then
	// everyone can bet or check. Betting round is opened until
	// someone has placed a bet(incldues BB/SB).
	if s.actionHintGroup.RaiserHint == nil {
		for _, hint := range playingHints {
			hint.AvailableActions = append(hint.AvailableActions, action2.Check)

			// If can check then player should not fold.
			hint.AvailableActions = lo.Without(hint.AvailableActions, action2.Fold)

			if s.seatStatusGroup.Status[hint.Uid].Chip >= s.gameSetting.BigBlind {
				hint.AvailableActions = append(hint.AvailableActions, action2.Bet)
			}
		}
		return
	}

	// If current betting round (stage) is opened then, everyone can
	// call based on their chip.
	for _, hint := range playingHints {
		seatStatus := s.seatStatusGroup.Status[hint.Uid]

		maxBetChip := lo.Max(lo.Map(lo.Values(s.actionHintGroup.Hints), func(h *model2.ActionHint, _ int) int { return h.BetChip }))
		hint.CallingChip = maxBetChip - hint.BetChip

		// BB can only check himself.
		// Edge case: In club mode, SB can be placed same chip with BB. So, SB can check himself.
		if (hint.Action == action2.BB || hint.Action == action2.SB) && hint.CallingChip <= 0 {
			hint.AvailableActions = append(hint.AvailableActions, action2.Check)
			continue
		}

		if seatStatus.Chip > hint.CallingChip {
			hint.AvailableActions = append(hint.AvailableActions, action2.Call)
		}
	}

	// If only one is live, he can only call. Cannot raise.
	if s.OnlyOneIsLive() {
		return
	}

	for _, hint := range playingHints {
		seatStatus := s.seatStatusGroup.Status[hint.Uid]

		// Edge case: When there is incomplete raise going on, current
		// raiser's bet chip will not be equal to max bet chip on the
		// table. Also, under this case, current raiser cannot raise
		// further. (until the raise is completed, and the raiser is
		// changed to other player)
		if hint != s.actionHintGroup.RaiserHint || hint.Action == action2.BB { // Note that, BB player has the right to raise on the first round
			// The minimum raise is equal to the size of the previous bet
			// or raise. If someone wishes to re-raise, they must raise at
			// least the amount of the previous raise. For example, if the
			// big blind is $2 and there is a raise of $6 to a total of $8,
			// a re-raise must be at least $6 more for a total of $14.
			callingRaiserChip := s.actionHintGroup.RaiserHint.BetChip - hint.BetChip
			hint.MinRaisingChip = callingRaiserChip + s.actionHintGroup.RaiserHint.RaiseChip

			if seatStatus.Chip >= hint.MinRaisingChip {
				hint.AvailableActions = append(hint.AvailableActions, action2.Raise)
			}
		}

	}
}

// If all other opponents fold a player's bet or raise. Or if all
// others fold till only one player not fold.
func (s *ActionHintService) OnlyOneNotFold() bool {
	foldCount := lo.CountBy(
		lo.Values(s.actionHintGroup.Hints),
		func(h *model2.ActionHint) bool { return h.Action == action2.Fold },
	)

	return len(s.actionHintGroup.Hints)-foldCount <= 1
}

func (s *ActionHintService) OnlyOneIsLive() bool {
	return len(s.fetchLiveHints()) <= 1
}

func (s *ActionHintService) AllHaveActed() bool {
	return lo.EveryBy(
		lo.Values(s.actionHintGroup.Hints),
		func(h *model2.ActionHint) bool { return len(h.AvailableActions) == 0 },
	)
}

func (s *ActionHintService) ClearAction() {
	s.actionHintGroup.RaiserHint = nil
	for _, liveHint := range s.fetchLiveHints() {
		liveHint.Action = action2.Undefined
		liveHint.AvailableActions = []action2.ActionType{}
		liveHint.RaiseChip = 0
	}
}

func (s *ActionHintService) fetch(uid core.Uid) (*model2.ActionHint, *model2.SeatStatus, error) {
	seatStatus, ok := s.seatStatusGroup.Status[uid]
	if !ok {
		return nil, nil, status.Errorf(codes.NotFound, "uid %s not found in seat status group", uid)
	}

	if seatStatus.FSM.MustState().(seatstatus.SeatStatusState) != seatstatus.PlayingState {
		return nil, nil, status.Errorf(codes.FailedPrecondition, "uid %s is not in playing state but in %s state", uid, seatStatus.FSM.MustState().(seatstatus.SeatStatusState).String())
	}

	actionHint, ok := s.actionHintGroup.Hints[uid]
	if !ok {
		return nil, nil, status.Errorf(codes.NotFound, "uid %s not found in action hint group", uid)
	}

	return actionHint, seatStatus, nil
}

func (s *ActionHintService) changeRaiser(raiserHint *model2.ActionHint) {
	if s.actionHintGroup.RaiserHint != nil {
		s.logger.Debug("changing raiser",
			zap.Object("CurrentRaiserHint", s.actionHintGroup.RaiserHint),
			zap.Object("NextRaiserHint", raiserHint),
		)
	}

	s.actionHintGroup.RaiserHint = raiserHint

	for _, liveHint := range s.fetchLiveHints() {
		if liveHint != s.actionHintGroup.RaiserHint {
			liveHint.Action = action2.Undefined
			liveHint.AvailableActions = []action2.ActionType{}
			liveHint.RaiseChip = 0
		}
	}
}

func (s *ActionHintService) fetchLiveHints() []*model2.ActionHint {
	return lo.Filter(
		lo.Values(s.actionHintGroup.Hints),
		func(h *model2.ActionHint, _ int) bool {
			return h.Action != action2.Fold &&
				h.Action != action2.AllIn &&
				(h.Action != action2.BB || s.chipCacheGroup.SeatStatusChips[h.Uid] > 0) // Edge case: BB with 0 chip cannot do any other action.
		},
	)
}
