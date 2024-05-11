package info

import (
	"context"
	"fmt"
	"gitee.com/i-Things/core/service/syssvr/pb/sys"
	"gitee.com/i-Things/share/utils"

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

func (l *TreeLogic) Tree(req *types.NotifyInfoIndexReq) (resp *types.NotifyInfoTreeResp, err error) {
	ret, err := l.svcCtx.NotifyM.NotifyInfoIndex(l.ctx, utils.Copy[sys.NotifyInfoIndexReq](req))
	list := utils.CopySlice[types.NotifyInfo](ret.List)
	var retMap = map[string][]*types.NotifyInfo{}
	for _, v := range list {
		retMap[v.Group] = append(retMap[v.Group], v)
	}
	var retList []*types.NotifyGroupInfo
	var groupID int64
	for k, v := range retMap {
		groupID++
		code := fmt.Sprintf("group%d", groupID)
		retList = append(retList, &types.NotifyGroupInfo{
			ID:       code,
			Code:     code,
			Name:     k,
			Children: v,
		})
	}
	return &types.NotifyInfoTreeResp{List: retList}, err
}
