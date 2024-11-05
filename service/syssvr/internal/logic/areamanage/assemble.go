package areamanagelogic

import (
	"context"

	"gitee.com/unitedrhino/core/service/syssvr/internal/logic"
	"gitee.com/unitedrhino/core/service/syssvr/internal/repo/relationDB"
	"gitee.com/unitedrhino/core/service/syssvr/internal/svc"
	"gitee.com/unitedrhino/core/service/syssvr/pb/sys"
	"gitee.com/unitedrhino/share/def"
	"gitee.com/unitedrhino/share/oss/common"
	"gitee.com/unitedrhino/share/utils"
	"github.com/zeromicro/go-zero/core/logx"
)

func transPoArrToPbTree(ctx context.Context, svcCtx *svc.ServiceContext, root *relationDB.SysAreaInfo, poArr []*relationDB.SysAreaInfo) *sys.AreaInfo {
	pbList := make([]*sys.AreaInfo, 0, len(poArr))
	for _, po := range poArr {
		pbList = append(pbList, TransPoToPb(ctx, po, svcCtx))
	}
	return buildPbTree(TransPoToPb(ctx, root, svcCtx), pbList)
}

func buildPbTree(rootArea *sys.AreaInfo, pbList []*sys.AreaInfo) *sys.AreaInfo {
	// 将所有节点按照 parentID 分组
	nodeMap := make(map[int64][]*sys.AreaInfo)
	for _, pbOne := range pbList {
		nodeMap[pbOne.ParentAreaID] = append(nodeMap[pbOne.ParentAreaID], pbOne)
	}

	// 递归生成子树
	buildPbSubtree(rootArea, nodeMap)

	return rootArea
}

func TransPoToPb(ctx context.Context, po *relationDB.SysAreaInfo, svcCtx *svc.ServiceContext) *sys.AreaInfo {
	parentAreaID := po.ParentAreaID
	if parentAreaID == 0 {
		parentAreaID = def.RootNode
	}
	if po.AreaImg != "" {
		var err error
		po.AreaImg, err = svcCtx.OssClient.PrivateBucket().SignedGetUrl(ctx, po.AreaImg, 24*60*60, common.OptionKv{})
		if err != nil {
			logx.WithContext(ctx).Errorf("%s.SignedGetUrl err:%v", utils.FuncName(), err)
		}
	}
	if po.ConfigFile != "" {
		var err error
		po.ConfigFile, err = svcCtx.OssClient.PrivateBucket().SignedGetUrl(ctx, po.ConfigFile, 24*60*60, common.OptionKv{})
		if err != nil {
			logx.WithContext(ctx).Errorf("%s.SignedGetUrl err:%v", utils.FuncName(), err)
		}
	}
	return &sys.AreaInfo{
		TenantCode:      string(po.TenantCode),
		CreatedTime:     po.CreatedTime.Unix(),
		AreaID:          int64(po.AreaID),
		ParentAreaID:    parentAreaID,
		ProjectID:       int64(po.ProjectID),
		AreaName:        po.AreaName,
		AreaNamePath:    po.AreaNamePath,
		AreaIDPath:      po.AreaIDPath,
		Position:        logic.ToSysPoint(po.Position),
		Desc:            utils.ToRpcNullString(po.Desc),
		IsLeaf:          po.IsLeaf,
		IsSysCreated:    po.IsSysCreated,
		LowerLevelCount: po.LowerLevelCount,
		ChildrenAreaIDs: po.ChildrenAreaIDs,
		DeviceCount:     utils.ToRpcNullInt64(po.DeviceCount),
		GroupCount:      utils.ToRpcNullInt64(po.GroupCount),
		UseBy:           po.UseBy,
		AreaImg:         po.AreaImg,
	}
}
func AreaInfosToPb(ctx context.Context, svcCtx *svc.ServiceContext, pos []*relationDB.SysAreaInfo) (ret []*sys.AreaInfo) {
	if pos == nil {
		return nil
	}
	for _, po := range pos {
		ret = append(ret, TransPoToPb(ctx, po, svcCtx))
	}
	return
}

func buildPbSubtree(node *sys.AreaInfo, nodeMap map[int64][]*sys.AreaInfo) {
	// 找到当前节点的子节点数组
	children := nodeMap[node.AreaID]

	// 递归生成子树
	for _, child := range children {
		buildPbSubtree(child, nodeMap)
	}

	// 将生成的子树数组作为当前节点的子节点数组
	node.Children = children
}
