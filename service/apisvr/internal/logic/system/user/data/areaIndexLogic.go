package data

import (
	"context"

	"gitee.com/unitedrhino/core/service/apisvr/internal/logic"
	"gitee.com/unitedrhino/core/service/apisvr/internal/logic/system/data/area"
	"gitee.com/unitedrhino/core/service/apisvr/internal/svc"
	"gitee.com/unitedrhino/core/service/apisvr/internal/types"
	"gitee.com/unitedrhino/core/service/syssvr/pb/sys"
	"gitee.com/unitedrhino/share/utils"

	"github.com/zeromicro/go-zero/core/logx"
)

type AreaIndexLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取区域权限列表
func NewAreaIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AreaIndexLogic {
	return &AreaIndexLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AreaIndexLogic) AreaIndex(req *types.UserDataAreaIndexReq) (resp *types.DataAreaIndexResp, err error) {
	ret, err := l.svcCtx.UserRpc.UserDataAreaIndex(l.ctx, utils.Copy[sys.UserDataAreaIndexReq](req))
	if err != nil {
		l.Errorf("%s.rpc.DataAreaIndex req=%v err=%+v", utils.FuncName(), req, err)
		return nil, err
	}
	if len(ret.List) == 0 {
		return &types.DataAreaIndexResp{}, nil
	}
	var areaIDs []int64
	for _, v := range ret.List {
		areaIDs = append(areaIDs, v.AreaID)
	}
	areaInfos, err := l.svcCtx.AreaM.AreaInfoIndex(l.ctx, &sys.AreaInfoIndexReq{AreaIDs: areaIDs})
	if err != nil {
		return nil, err
	}
	var areaMap = map[int64]*sys.AreaInfo{}
	for _, v := range areaInfos.List {
		areaMap[v.AreaID] = v
	}
	list := area.ToDataAreaDetail(l.ctx, nil, ret.List, areaMap)
	return &types.DataAreaIndexResp{
		PageResp: logic.ToPageResp(req.Page, ret.Total),
		List:     list,
	}, nil
}
