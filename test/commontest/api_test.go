package commontest

import (
	commonapi "card-game-server-prototype/pkg/common/api"
	"card-game-server-prototype/pkg/core"
	"card-game-server-prototype/pkg/util"
	"testing"

	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestFetchUserDetail(t *testing.T) {
	assert := assert.New(t)
	logger := util.NewTestLogger()

	registry, err := BuildRegistry()
	assert.NoError(err)

	uids := core.UidList{
		core.Uid("001a1c72-1abc-4ae8-a574-41049a43926d"),
		core.Uid("000006a1-5c32-439e-b465-ee8d941645c0"),
		core.Uid("000428f3-4b95-4412-b266-2a7ce8ac105e"),
	}

	results := make(chan lo.Tuple2[*commonapi.UserDetailResponse, error], len(uids))
	for _, uid := range uids {
		go func(uid core.Uid) {
			resp, err := registry.api.baseUserAPI.FetchUserDetail(uid)
			results <- lo.T2(resp, err)
		}(uid)
	}

	for i := 0; i < len(uids); i++ {
		result := <-results
		resp, err := result.A, result.B
		assert.NoError(err)
		assert.Contains(uids, core.Uid(resp.Data.Uid))
		logger.Info("TestFetchUserDetail", zap.Any("resp", resp))
	}

	close(results)
}

func TestEnterLeaveRoom(t *testing.T) {
	assert := assert.New(t)
	logger := util.NewTestLogger()

	registry, err := BuildRegistry()
	assert.NoError(err)

	uid := core.Uid("001a1c72-1abc-4ae8-a574-41049a43926d")
	err = registry.api.baseRoomAPI.EnterRoom("001a1c72-1abc-4ae8-a574-41049a43926d", uid)
	logger.Info("TestEnterLeaveRoom", zap.Error(err))
}
