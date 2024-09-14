package config

import (
	"context"
	"fmt"
	"gitee.com/i-Things/core/service/syssvr/pb/sys"
	"gitee.com/i-Things/share/utils"
	"sort"

	"gitee.com/i-Things/core/service/apisvr/internal/svc"
	"gitee.com/i-Things/core/service/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type TreeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewTreeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TreeLogic {
	return &TreeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *TreeLogic) Tree(req *types.NotifyConfigIndexReq) (resp *types.NotifyConfigTreeResp, err error) {
	ret, err := l.svcCtx.NotifyM.NotifyConfigIndex(l.ctx, utils.Copy[sys.NotifyConfigIndexReq](req))
	list := utils.CopySlice[types.NotifyConfig](ret.List)
	var retMap = map[string][]*types.NotifyConfig{}
	for _, v := range list {
		retMap[v.Group] = append(retMap[v.Group], v)
	}
	var retList []*types.NotifyGroupInfo
	var groupID int64
	var keys = utils.SetToSlice(retMap)
	sort.Strings(keys)
	for _, k := range keys {
		v := retMap[k]
		groupID++
		code := fmt.Sprintf("group%d", groupID)
		retList = append(retList, &types.NotifyGroupInfo{
			ID:       code,
			Code:     code,
			Name:     k,
			Children: v,
		})
	}
	return &types.NotifyConfigTreeResp{List: retList}, err
}
