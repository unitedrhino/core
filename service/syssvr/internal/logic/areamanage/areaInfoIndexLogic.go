package areamanagelogic

import (
	"context"
	"gitee.com/i-Things/core/service/syssvr/internal/logic"
	"gitee.com/i-Things/core/service/syssvr/internal/repo/relationDB"
	"gitee.com/i-Things/core/service/syssvr/internal/svc"
	"gitee.com/i-Things/core/service/syssvr/pb/sys"
	"gitee.com/i-Things/share/stores"
	"github.com/zeromicro/go-zero/core/logx"
)

type AreaInfoIndexLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	AiDB *relationDB.AreaInfoRepo
}

func NewAreaInfoIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AreaInfoIndexLogic {
	return &AreaInfoIndexLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
		AiDB:   relationDB.NewAreaInfoRepo(ctx),
	}
}

// 获取区域信息列表
func (l *AreaInfoIndexLogic) AreaInfoIndex(in *sys.AreaInfoIndexReq) (*sys.AreaInfoIndexResp, error) {
	var (
		poArr []*relationDB.SysAreaInfo
		f     = relationDB.AreaInfoFilter{
			ProjectID: in.ProjectID, AreaIDs: in.AreaIDs, ParentAreaID: in.ParentAreaID, IsLeaf: in.IsLeaf}
	)

	if in.DeviceCount != nil {
		f.DeviceCount = stores.GetCmp(in.DeviceCount.CmpType, in.DeviceCount.Value)
	}
	if in.GroupCount != nil {
		f.GroupCount = stores.GetCmp(in.GroupCount.CmpType, in.GroupCount.Value)
	}
	poArr, err := l.AiDB.FindByFilter(l.ctx,
		f, logic.ToPageInfo(in.Page))
	if err != nil {
		l.Errorf("AreaInfoIndex find menu_info err,menuIds:%d,err:%v", in.AreaIDs, err)
		return nil, err
	}
	total, err := l.AiDB.CountByFilter(l.ctx, f)
	if err != nil {
		l.Errorf("AreaInfoIndex find menu_info err,menuIds:%d,err:%v", in.AreaIDs, err)
		return nil, err
	}

	return &sys.AreaInfoIndexResp{List: AreaInfosToPb(l.ctx, l.svcCtx, poArr), Total: total}, nil
}
