package service

import (
	"card-game-server-prototype/pkg/core"
	"card-game-server-prototype/pkg/game/txpoker/api"
)

type ClubModeUserService struct {
	*BaseUserService

	memApi api.ClubMemberAPI
}

func ProvideClubModeUserService(
	baseUserService *BaseUserService,

	memApi api.ClubMemberAPI,
) *ClubModeUserService {
	return &ClubModeUserService{
		BaseUserService: baseUserService,

		memApi: memApi,
	}
}

func (s *ClubModeUserService) FetchFromRepo(uids ...core.Uid) error {
	var newUids []core.Uid
	for _, v := range uids {
		if v.String() == "" {
			continue
		}
		newUids = append(newUids, v)
	}

	if len(newUids) == 0 {
		return nil
	}

	if err := s.BaseUserService.FetchFromRepo(newUids...); err != nil {
		return err
	}

	resp, err := s.memApi.FetchDetail(s.roomInfo.ClubId, newUids...)
	if err != nil {
		return err
	}

	for _, v := range resp.Data {
		if _, ok := s.userGroup.Data[core.Uid(v.Uid)]; ok {
			s.userGroup.Data[core.Uid(v.Uid)].Cash = v.Gold
		}

		if _, ok := s.userCacheGroup.Data[core.Uid(v.Uid)]; ok {
			s.userCacheGroup.Data[core.Uid(v.Uid)].Cash = v.Gold
		}
	}

	return nil
}
