package areamanagelogic

import (
	"context"
	"gitee.com/i-Things/core/service/syssvr/internal/repo/relationDB"
	"gitee.com/i-Things/share/def"
	"github.com/samber/lo"

	"gitee.com/i-Things/core/service/syssvr/internal/svc"
	"gitee.com/i-Things/core/service/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type AreaInfoReadLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	AiDB *relationDB.AreaInfoRepo
}

func NewAreaInfoReadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AreaInfoReadLogic {
	return &AreaInfoReadLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
		AiDB:   relationDB.NewAreaInfoRepo(ctx),
	}
}

var (
	rootNode = relationDB.SysAreaInfo{
		AreaID:       def.RootNode,
		AreaIDPath:   "",
		AreaNamePath: "全部区域",
		AreaName:     "全部区域",
	}
	notClassifiedNode = relationDB.SysAreaInfo{
		AreaID:       def.NotClassified,
		AreaIDPath:   "2-",
		AreaNamePath: "未分类的区域",
		AreaName:     "未分类的区域",
	}
)

// 获取区域信息详情
func (l *AreaInfoReadLogic) AreaInfoRead(in *sys.AreaInfoReadReq) (*sys.AreaInfo, error) {
	var (
		po  *relationDB.SysAreaInfo
		err error
	)

	switch in.AreaID {
	case def.RootNode, 0:
		po = &rootNode
	case def.NotClassified:
		po = &notClassifiedNode
		return transPoToPb(l.ctx, po, l.svcCtx), nil
	default:
		po, err = l.AiDB.FindOne(l.ctx, in.AreaID, nil)
		if err != nil {
			return nil, err
		}
	}
	if !in.WithChildren {
		return transPoToPb(l.ctx, po, l.svcCtx), nil
	}
	poArr, err := l.AiDB.FindByFilter(l.ctx, relationDB.AreaInfoFilter{ProjectID: in.ProjectID, AreaIDPath: po.AreaIDPath}, nil)
	if err != nil {
		return nil, err
	}

	return transPoArrToPbTree(l.ctx, l.svcCtx, po, poArr), err
}

func (l *AreaInfoReadLogic) checkMissingParentIdMenuIndex(areaInfos []*relationDB.SysAreaInfo) []*relationDB.SysAreaInfo {
	missingParentIds := findMissingParentIds(areaInfos)
	if len(missingParentIds) > 0 {
		areaIDs := lo.Keys(missingParentIds)
		areaInfo, err := l.AiDB.FindByFilter(l.ctx, relationDB.AreaInfoFilter{AreaIDs: areaIDs}, nil)
		if err != nil {
			l.Errorf("MenuIndex find menu_info err,menuIds:%d,err:%v", areaIDs, err)
			return areaInfos
		}
		areaInfos = append(areaInfos, areaInfo...)
		return l.checkMissingParentIdMenuIndex(areaInfos) //多级嵌套需要遍历
	}
	return areaInfos
}

func findMissingParentIds(areaInfos []*relationDB.SysAreaInfo) map[int64]bool {
	missingParentIds := make(map[int64]bool)
	ids := make(map[int64]bool)
	for _, menu := range areaInfos {
		ids[int64(menu.AreaID)] = true
	}
	for _, menu := range areaInfos {
		if !ids[menu.ParentAreaID] && menu.ParentAreaID != 1 {
			missingParentIds[menu.ParentAreaID] = true
		}
	}
	return missingParentIds
}
