package projectmanagelogic

import (
	"context"
	"gitee.com/unitedrhino/core/service/syssvr/internal/logic"
	areamanagelogic "gitee.com/unitedrhino/core/service/syssvr/internal/logic/areamanage"
	"gitee.com/unitedrhino/core/service/syssvr/internal/repo/relationDB"
	"gitee.com/unitedrhino/core/service/syssvr/internal/svc"
	"gitee.com/unitedrhino/core/service/syssvr/pb/sys"
	"gitee.com/unitedrhino/share/oss/common"
	"gitee.com/unitedrhino/share/utils"
	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func ProjectInfoToPb(ctx context.Context, svcCtx *svc.ServiceContext, po *relationDB.SysProjectInfo) *sys.ProjectInfo {
	if po == nil {
		return nil
	}
	if po.ProjectImg != "" {
		var err error
		po.ProjectImg, err = svcCtx.OssClient.PrivateBucket().SignedGetUrl(ctx, po.ProjectImg, 24*60*60, common.OptionKv{})
		if err != nil {
			logx.WithContext(ctx).Errorf("%s.SignedGetUrl err:%v", utils.FuncName(), err)
		}
	}
	pb := &sys.ProjectInfo{
		TenantCode:   string(po.TenantCode),
		CreatedTime:  po.CreatedTime.Unix(),
		ProjectID:    int64(po.ProjectID),
		ProjectName:  po.ProjectName,
		ProjectImg:   po.ProjectImg,
		AdminUserID:  po.AdminUserID,
		IsSysCreated: po.IsSysCreated,
		Ppsm:         po.Ppsm,
		Tags:         po.Tags,
		Area:         &wrapperspb.FloatValue{Value: po.Area},
		Desc:         utils.ToRpcNullString(po.Desc),
		Position:     logic.ToSysPoint(po.Position),
		AreaCount:    po.AreaCount,
		UserCount:    po.UserCount,
		Address:      utils.ToRpcNullString(po.Address),
		DeviceCount:  utils.ToRpcNullInt64(po.DeviceCount),
		Areas:        areamanagelogic.AreaInfosToPb(ctx, svcCtx, po.Areas),
	}
	return pb
}
