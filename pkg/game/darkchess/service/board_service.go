package service

import (
	"errors"
	"card-game-server-prototype/pkg/core"
	"card-game-server-prototype/pkg/game/darkchess/constant"
	model2 "card-game-server-prototype/pkg/game/darkchess/model"
	"card-game-server-prototype/pkg/game/darkchess/type/piece"
	"card-game-server-prototype/pkg/util"
	"github.com/samber/lo"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"math/rand"
	"sort"
)

type BoardService struct {
	logger *zap.Logger

	board           *model2.Board
	capturedPieces  *model2.CapturedPieces
	playerGroup     *model2.PlayerGroup
	actionHintGroup *model2.ActionHintGroup
}

func ProvideBoardService(
	loggerFactory *util.LoggerFactory,

	board *model2.Board,
	capturedPieces *model2.CapturedPieces,
	playerGroup *model2.PlayerGroup,
	actionHintGroup *model2.ActionHintGroup,
) *BoardService {
	return &BoardService{
		logger: loggerFactory.Create("BoardService"),

		board:           board,
		capturedPieces:  capturedPieces,
		playerGroup:     playerGroup,
		actionHintGroup: actionHintGroup,
	}
}

func (s *BoardService) GetCellPos(p piece.Piece) (int, int, bool) {
	for x, row := range s.board.Cells {
		for y, v := range row {
			if v.Piece.IsSame(p) {
				return x, y, true
			}
		}
	}

	return 0, 0, false
}

func (s *BoardService) RevealPiece(x, y int) {
	s.board.Cells[x][y].IsPieceRevealed = true
}

func (s *BoardService) MovePiece(x, y, toX, toY int) {
	s.board.Cells[toX][toY] = s.board.Cells[x][y]
	s.board.Cells[x][y] = model2.NewEmptyCell()
}

func (s *BoardService) CapturePiece(x, y, targetX, targetY int) {
	// collect dead piece
	s.capturedPieces.Pieces = append(s.capturedPieces.Pieces,
		s.board.Cells[targetX][targetY].Piece)

	// Move hunter to target position
	s.board.Cells[targetX][targetY] = s.board.Cells[x][y]
	s.board.Cells[x][y] = model2.NewEmptyCell()
}

func (s *BoardService) ValidRevealRules(x, y int) error {
	if err := s.validPos(x, y); err != nil {
		return status.Errorf(codes.OutOfRange, "invalid reveal position: %v", err)
	}

	if s.board.Cells[x][y].IsEmptyCell {
		return status.Errorf(codes.NotFound, "invalid reveal cell: is empty")
	}

	if s.board.Cells[x][y].IsPieceRevealed {
		return status.Errorf(codes.AlreadyExists, "invalid reveal cell: piece already revealed")
	}

	return nil
}

func (s *BoardService) ValidMoveRules(x, y, toX, toY int) error {

	if err := s.validPos(x, y); err != nil {
		return status.Errorf(codes.OutOfRange, "invalid move position: %v", err)
	}

	if err := s.validPos(toX, toY); err != nil {
		return status.Errorf(codes.OutOfRange, "invalid move to target position: %v", err)
	}

	if s.board.Cells[x][y].Piece.IsInvalid() {
		return status.Errorf(codes.FailedPrecondition, "invalid move, invalid piece")
	}

	if s.board.Cells[x][y].IsEmptyCell {
		return status.Errorf(codes.NotFound, "invalid move, move cell is empty")
	}

	if !s.board.Cells[toX][toY].IsEmptyCell {
		return status.Errorf(codes.NotFound, "invalid move, target cell has a piece")
	}

	// Basic move rules

	if !s.IsNextInAxis(x, y, toX, toY) {
		return status.Errorf(codes.InvalidArgument, "invalid move, noly allow to move 1 cell")
	}

	return nil
}

func (s *BoardService) ValidCaptureRules(x, y, toX, toY int) error {

	if err := s.validPos(x, y); err != nil {
		return status.Errorf(codes.OutOfRange, "invalid capture, hunter position: %v", err)
	}

	if err := s.validPos(toX, toY); err != nil {
		return status.Errorf(codes.OutOfRange, "invalid capture, target position: %v", err)
	}

	hunter := s.board.Cells[x][y]
	target := s.board.Cells[toX][toY]

	if hunter.Piece.IsInvalid() {
		return status.Errorf(codes.FailedPrecondition, "invalid capture, invalid hunter piece")
	}

	if target.Piece.IsInvalid() {
		return status.Errorf(codes.FailedPrecondition, "invalid capture, invalid target piece")
	}

	if hunter.IsEmptyCell {
		return status.Errorf(codes.NotFound, "invalid capture, hunter cell is empty")
	}

	if target.IsEmptyCell {
		return status.Errorf(codes.NotFound, "invalid capture, target cell is empty")
	}

	if hunter.Piece.IsCannon() {
		if err := s.ValidCanonAndTargetPos(x, y, toX, toY); err != nil {
			return status.Errorf(codes.FailedPrecondition, err.Error())
		}

		countPieces := s.CountPiecesBetween(x, y, toX, toY)

		if countPieces == 0 || countPieces > 1 {
			return status.Errorf(codes.FailedPrecondition,
				"invalid capture, there are 0 or more than 1 piece between canon and target: %v", countPieces)
		}

	} else {

		if !s.IsNextInAxis(x, y, toX, toY) {
			return status.Errorf(codes.FailedPrecondition, "invalid capture, only allow to capture next cell in axis")
		}

		if !s.IsAllowToCapture(hunter.Piece, target.Piece) {
			return status.Errorf(codes.FailedPrecondition, "invalid capture, hunter is not greater than target")
		}
	}

	return nil
}

func (s *BoardService) ValidCanonAndTargetPos(canonX, canonY, targetX, targetY int) error {
	if !s.isInAxis(canonX, canonY, targetX, targetY) {
		return errors.New("invalid capture, only allow to capture cell in axis")
	}

	if s.IsNextInAxis(canonX, canonY, targetX, targetY) {
		return errors.New("invalid capture, canon does not allow to capture next cell")
	}

	return nil
}

func (s *BoardService) CountPiecesBetween(x, y, x2, y2 int) int {
	return len(lo.Filter(s.posBetween(x, y, x2, y2), func(pos model2.Pos, i int) bool {
		return !s.board.Cells[pos.X][pos.Y].IsEmptyCell
	}))
}

func (s *BoardService) DistanceBetween(x, y, x2, y2 int) int {
	return len(s.posBetween(x, y, x2, y2))
}

func (s *BoardService) IsRedAllDead() bool {
	return lo.CountBy(s.capturedPieces.Pieces, func(p piece.Piece) bool {
		return p.IsRed()
	}) == constant.OneColorTotalPieceCnt
}

func (s *BoardService) IsBlackAllDead() bool {
	return lo.CountBy(s.capturedPieces.Pieces, func(p piece.Piece) bool {
		return p.IsBlack()
	}) == constant.OneColorTotalPieceCnt
}

func (s *BoardService) RandomPickUnrevealedPiece() (int, int, bool) {
	var actions []model2.Pos

	for x, row := range s.board.Cells {
		for y := range row {
			if err := s.ValidRevealRules(x, y); err != nil {
				continue
			}

			actions = append(actions, model2.Pos{X: x, Y: y})
		}
	}

	if len(actions) == 0 {
		return 0, 0, false
	}

	pickRandom := actions[rand.Intn(len(actions))]
	return pickRandom.X, pickRandom.Y, true
}

func (s *BoardService) RandomMovePiece(uid core.Uid) (int, int, int, int, bool) {
	type pos struct {
		before model2.Pos
		after  model2.Pos
	}

	freezePos, hasFreeze := s.getFreezePos()

	useColor := s.playerGroup.Data[uid].Color

	var actions []pos

	for x, row := range s.board.Cells {
		for y := range row {
			if hasFreeze && freezePos.X == x && freezePos.Y == y {
				continue
			}

			if s.board.Cells[x][y].Piece.GetColor() != useColor {
				continue
			}

			if newX, hasMove := s.moveLeftNCells(x, 1); hasMove {
				err := s.ValidMoveRules(x, y, newX, y)
				if err == nil {
					actions = append(actions, pos{
						before: model2.Pos{X: x, Y: y},
						after:  model2.Pos{X: newX, Y: y},
					})
				}
			}

			if newX, hasMove := s.moveRightNCells(x, 1); hasMove {
				err := s.ValidMoveRules(x, y, newX, y)
				if err == nil {
					actions = append(actions, pos{
						before: model2.Pos{X: x, Y: y},
						after:  model2.Pos{X: newX, Y: y},
					})
				}
			}

			if newY, hasMove := s.moveUpNCells(y, 1); hasMove {
				err := s.ValidMoveRules(x, y, x, newY)
				if err == nil {
					actions = append(actions, pos{
						before: model2.Pos{X: x, Y: y},
						after:  model2.Pos{X: x, Y: newY},
					})
				}
			}

			if newY, hasMove := s.moveDownNCells(y, 1); hasMove {
				err := s.ValidMoveRules(x, y, x, newY)
				if err == nil {
					actions = append(actions, pos{
						before: model2.Pos{X: x, Y: y},
						after:  model2.Pos{X: x, Y: newY},
					})
				}
			}
		}
	}

	if len(actions) == 0 {
		return 0, 0, 0, 0, false
	}

	s.logger.Debug("RandomMove", zap.Array("pos", zapcore.ArrayMarshalerFunc(func(enc zapcore.ArrayEncoder) error {
		for _, v := range actions {
			_ = enc.AppendObject(zapcore.ObjectMarshalerFunc(func(enc zapcore.ObjectEncoder) error {
				enc.AddInt("x", v.before.X)
				enc.AddInt("y", v.before.Y)
				enc.AddInt("toX", v.after.X)
				enc.AddInt("toY", v.after.Y)
				return nil
			}))
		}
		return nil
	})))

	pickRandom := actions[rand.Intn(len(actions))]
	return pickRandom.before.X, pickRandom.before.Y, pickRandom.after.X, pickRandom.after.Y, true
}

func (s *BoardService) RandomCapturePiece(uid core.Uid) (int, int, int, int, bool) {
	type pos struct {
		m model2.Pos
		c model2.Pos
	}

	useColor := s.playerGroup.Data[uid].Color
	s.logger.Debug("random capture, use color", zap.String("color", useColor.String()))

	var actions []pos

	freezePos, hasFreeze := s.getFreezePos()

	for x, row := range s.board.Cells {
		for y := range row {
			if hasFreeze && freezePos.X == x && freezePos.Y == y {
				continue
			}

			if s.board.Cells[x][y].Piece.GetColor() != useColor {
				s.logger.Debug("random capture, skip hunter color", zap.Int("x", x), zap.Int("y", y))
				continue
			}

			if s.board.Cells[x][y].Piece.IsCannon() {
				s.logger.Debug("random capture, skip todo canon", zap.Int("x", x), zap.Int("y", y))
				// todo: add cannon rules. 炮用以下規則會做出不合法的動作，需要另外處理
				continue
			}

			if newX, hasMove := s.moveLeftNCells(x, 1); hasMove {
				if s.board.Cells[newX][y].Piece.GetColor() == useColor {
					s.logger.Debug("random capture, skip target color - L", zap.Int("x", x), zap.Int("y", y))
				} else {
					if err := s.ValidCaptureRules(x, y, newX, y); err != nil {
						s.logger.Debug("random capture, valid fail - L",
							zap.Int("x", x), zap.Int("y", y),
							zap.Int("newX", newX), zap.Int("newY", y),
							zap.Error(err))
					} else {
						c := s.board.Cells[newX][y]
						if !c.IsEmptyCell &&
							s.IsAllowToCapture(s.board.Cells[x][y].Piece, c.Piece) {
							actions = append(actions, pos{
								m: model2.Pos{X: x, Y: y},
								c: model2.Pos{X: newX, Y: y},
							})
						}
					}
				}
			}

			if newX, hasMove := s.moveRightNCells(x, 1); hasMove {
				if s.board.Cells[newX][y].Piece.GetColor() == useColor {
					s.logger.Debug("random capture, skip target color - R", zap.Int("x", x), zap.Int("y", y))
				} else {
					if err := s.ValidCaptureRules(x, y, newX, y); err != nil {
						s.logger.Debug("random capture, valid fail - R",
							zap.Int("x", x), zap.Int("y", y),
							zap.Int("newX", newX), zap.Int("newY", y),
							zap.Error(err))
					} else {
						c := s.board.Cells[newX][y]
						if !c.IsEmptyCell &&
							s.IsAllowToCapture(s.board.Cells[x][y].Piece, c.Piece) {
							actions = append(actions, pos{
								m: model2.Pos{X: x, Y: y},
								c: model2.Pos{X: newX, Y: y},
							})
						}
					}
				}
			}

			if newY, hasMove := s.moveUpNCells(y, 1); hasMove {
				if s.board.Cells[x][newY].Piece.GetColor() == useColor {
					s.logger.Debug("random capture, skip target color - U", zap.Int("x", x), zap.Int("y", y))
				} else {

					if err := s.ValidCaptureRules(x, y, x, newY); err != nil {
						s.logger.Debug("random capture, valid fail - U",
							zap.Int("x", x), zap.Int("y", y),
							zap.Int("newX", x), zap.Int("newY", newY),
							zap.Error(err))
					} else {
						c := s.board.Cells[x][newY]
						if !c.IsEmptyCell &&
							s.IsAllowToCapture(s.board.Cells[x][y].Piece, c.Piece) {
							actions = append(actions, pos{
								m: model2.Pos{X: x, Y: y},
								c: model2.Pos{X: x, Y: newY},
							})
						}
					}
				}
			}

			if newY, hasMove := s.moveDownNCells(y, 1); hasMove {
				if s.board.Cells[x][newY].Piece.GetColor() == useColor {
					s.logger.Debug("random capture, skip target color - D", zap.Int("x", x), zap.Int("y", y))
				} else {
					if err := s.ValidCaptureRules(x, y, x, newY); err != nil {
						s.logger.Debug("random capture, valid fail - D",
							zap.Int("x", x), zap.Int("y", y),
							zap.Int("newX", x), zap.Int("newY", newY),
							zap.Error(err))
					} else {
						c := s.board.Cells[x][newY]
						if !c.IsEmptyCell &&
							s.IsAllowToCapture(s.board.Cells[x][y].Piece, c.Piece) {
							actions = append(actions, pos{
								m: model2.Pos{X: x, Y: y},
								c: model2.Pos{X: x, Y: newY},
							})
						}
					}
				}
			}
		}
	}

	if len(actions) == 0 {
		return 0, 0, 0, 0, false
	}

	s.logger.Debug("RandomCapture", zap.Array("pos", zapcore.ArrayMarshalerFunc(func(enc zapcore.ArrayEncoder) error {
		for _, v := range actions {
			_ = enc.AppendObject(zapcore.ObjectMarshalerFunc(func(enc zapcore.ObjectEncoder) error {
				enc.AddInt("x", v.m.X)
				enc.AddInt("y", v.m.Y)
				enc.AddInt("toX", v.c.X)
				enc.AddInt("toY", v.c.Y)
				return nil
			}))
		}
		return nil
	})))

	pickRandom := actions[rand.Intn(len(actions))]
	return pickRandom.m.X, pickRandom.m.Y, pickRandom.c.X, pickRandom.c.Y, true
}

func (s *BoardService) getFreezePos() (*model2.Pos, bool) {
	for _, hint := range s.actionHintGroup.Data {
		if hint.FreezeCell != nil {
			return &model2.Pos{
				X: hint.FreezeCell.Pos.X,
				Y: hint.FreezeCell.Pos.Y,
			}, true
		}
	}
	return nil, false
}

func (s *BoardService) validPos(x, y int) error {
	if x < 0 || x >= constant.BoardWidth || y < 0 || y >= constant.BoardHeight {
		return errors.New("out of board")
	}
	return nil
}

func (s *BoardService) moveRightNCells(x, n int) (int, bool) {
	newX := x + n
	if newX >= constant.BoardWidth {
		return constant.BoardWidth - 1, false
	}
	return newX, true
}

func (s *BoardService) moveLeftNCells(x, n int) (int, bool) {
	newX := x - n
	if newX < 0 {
		return 0, false
	}
	return newX, true
}

func (s *BoardService) moveUpNCells(y, n int) (int, bool) {
	newY := y + n
	if newY >= constant.BoardHeight {
		return constant.BoardHeight - 1, false
	}
	return newY, true
}

func (s *BoardService) moveDownNCells(y, n int) (int, bool) {
	newY := y - n
	if newY < 0 {
		return 0, false
	}
	return newY, true
}

func (s *BoardService) isSamePos(x, y, x2, y2 int) bool {
	diffX := x2 - x
	diffY := y2 - y
	return diffX == 0 && diffY == 0
}

// 十字 diffX*diffX+diffY*diffY == 1
// 斜角 diffX*diffX+diffY*diffY == 2
// 終點是否在起點的上下左右的 1 格上，不含起點。
func (s *BoardService) IsNextInAxis(x, y, x2, y2 int) bool {
	if s.isSamePos(x, y, x2, y2) {
		return false
	}

	diffX := x2 - x
	diffY := y2 - y
	return diffX*diffX+diffY*diffY == 1
}

// 終點是否在上下左右的軸線上，不含起點。
func (s *BoardService) isInAxis(x, y, x2, y2 int) bool {
	if s.isSamePos(x, y, x2, y2) {
		return false
	}

	diffX := x2 - x
	diffY := y2 - y

	return (diffX*diffX)*(diffY*diffY) == 0
}

// 計算 兩點之間的所有座標，不含起點與終點。
func (s *BoardService) posBetween(x, y, x2, y2 int) []model2.Pos {
	diffX := x2 - x
	diffY := y2 - y

	var result []model2.Pos

	// 檢查 y 位移
	if diffX == 0 {
		sorted := []int{y, y2}

		sort.Slice(sorted, func(i, j int) bool {
			return sorted[i] < sorted[j]
		})

		from := sorted[0]
		to := sorted[1]

		for i := from + 1; i < to; i++ {
			newY := lo.Ternary(i < 0, 0, i)
			newY = lo.Ternary(newY >= constant.BoardHeight, constant.BoardHeight-1, newY)
			result = append(result, model2.Pos{X: x, Y: newY})
		}
	}

	// 檢查 x 位移
	if diffY == 0 {
		sorted := []int{x, x2}

		sort.Slice(sorted, func(i, j int) bool {
			return sorted[i] < sorted[j]
		})

		from := sorted[0]
		to := sorted[1]

		for i := from + 1; i < to; i++ {
			newX := lo.Ternary(i < 0, 0, i)
			newX = lo.Ternary(newX >= constant.BoardWidth, constant.BoardWidth-1, newX)
			result = append(result, model2.Pos{X: newX, Y: y})
		}
	}

	return result
}

func (s *BoardService) IsAllowToCapture(hunter, target piece.Piece) bool {
	if hunter.IsInvalid() || target.IsInvalid() {
		return false
	}

	if hunter.IsCannon() {
		return true
	}

	if hunter.IsGeneral() && target.IsSoldier() {
		return false
	}

	if hunter.IsSoldier() && target.IsGeneral() {
		return true
	}

	if hunter.GetWeight() == target.GetWeight() {
		return true
	}

	return hunter.GetWeight() > target.GetWeight()
}

// IsChaseSamePiece 是否有長捉的情況
func (s *BoardService) IsChaseSamePiece(hunter, target piece.Piece) bool {
	if hunter.IsInvalid() || target.IsInvalid() {
		return false
	}

	if hunter.IsCannon() {
		return true
	}

	if hunter.IsGeneral() && target.IsSoldier() {
		return false
	}

	if hunter.IsSoldier() && target.IsGeneral() {
		return true
	}

	// 只認定大追小，同 type 不算。
	return hunter.GetWeight() > target.GetWeight()
}
