package info

import (
	"context"
	"gitee.com/i-Things/core/service/apisvr/internal/logic"
	"gitee.com/i-Things/core/service/syssvr/pb/sys"
	"gitee.com/i-Things/core/shared/def"
	"gitee.com/i-Things/core/shared/errors"
	"gitee.com/i-Things/core/shared/utils"

	"gitee.com/i-Things/core/service/apisvr/internal/svc"
	"gitee.com/i-Things/core/service/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateLogic {
	return &CreateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateLogic) Create(req *types.AreaInfo) (*types.AreaWithID, error) {
	if req.AreaName == "" || req.ParentAreaID == 0 || ////root节点不为0
		req.ParentAreaID == def.NotClassified { //未分类不能有下属的区域
		return nil, errors.Parameter
	}
	//if req.ParentAreaID != def.RootNode {
	//	dmRep, err := l.svcCtx.DeviceM.DeviceInfoIndex(l.ctx, &dm.DeviceInfoIndexReq{
	//		Page:    &dm.PageInfo{Page: 1, Size: 2}, //只需要知道是否有设备即可
	//		AreaIDs: []int64{req.ParentAreaID}})
	//	if err != nil {
	//		return nil, err
	//	}
	//	if len(dmRep.List) != 0 {
	//		return nil, errors.Parameter.AddMsg("父级区域已绑定了设备，不允许再添加子区域")
	//	}
	//}

	dmReq := &sys.AreaInfo{
		ParentAreaID: req.ParentAreaID,
		ProjectID:    req.ProjectID,
		AreaName:     req.AreaName,
		Position:     logic.ToSysPointRpc(req.Position),
		Desc:         utils.ToRpcNullString(req.Desc),
	}
	resp, err := l.svcCtx.AreaM.AreaInfoCreate(l.ctx, dmReq)
	if er := errors.Fmt(err); er != nil {
		l.Errorf("%s.rpc.AreaManage req=%v err=%v", utils.FuncName(), req, er)
		return nil, er
	}
	return &types.AreaWithID{AreaID: resp.AreaID}, nil
}
