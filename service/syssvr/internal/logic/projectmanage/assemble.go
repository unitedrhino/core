package projectmanagelogic

import (
	"context"
	"encoding/json"
	"strings"

	"gitee.com/unitedrhino/core/service/syssvr/internal/logic"
	areamanagelogic "gitee.com/unitedrhino/core/service/syssvr/internal/logic/areamanage"
	"gitee.com/unitedrhino/core/service/syssvr/internal/repo/relationDB"
	"gitee.com/unitedrhino/core/service/syssvr/internal/svc"
	"gitee.com/unitedrhino/core/service/syssvr/pb/sys"
	"gitee.com/unitedrhino/share/oss"
	"gitee.com/unitedrhino/share/oss/common"
	"gitee.com/unitedrhino/share/utils"
	"github.com/spf13/cast"
	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func ProjectCrudsToPb(ctx context.Context, svcCtx *svc.ServiceContext, po []*relationDB.SysProjectCrud) []*sys.ProjectCrud {
	var pbs []*sys.ProjectCrud
	for _, p := range po {
		pbs = append(pbs, ProjectCrudToPb(ctx, svcCtx, p))
	}
	return pbs
}

func ProjectCrudToPb(ctx context.Context, svcCtx *svc.ServiceContext, po *relationDB.SysProjectCrud) *sys.ProjectCrud {
	if po == nil {
		return nil
	}
	pb := utils.Copy[sys.ProjectCrud](po)
	var params = map[string]interface{}{}
	if len(pb.Params) > 0 {
		err := json.Unmarshal([]byte(pb.Params), &params)
		if err != nil {
			logx.WithContext(ctx).Errorf("unmarshal params err:%v", err)
			return pb
		}
		for k, v := range params {
			if !(strings.HasSuffix(k, "Img") || strings.HasSuffix(k, "File")) {
				continue
			}
			url, err := svcCtx.OssClient.SignedGetUrl(ctx, cast.ToString(v), 24*60*60, common.OptionKv{})
			if err != nil {
				logx.WithContext(ctx).Error(po, k, v, err.Error())
				continue
			}
			params[k] = url
		}
		pb.Params = utils.MarshalNoErr(params)
	}

	return pb
}

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
	var at []*sys.Attachment

	for _, att := range po.Attachments {
		url, err := svcCtx.OssClient.PrivateBucket().SignedGetUrl(ctx, att.FilePath, 300, common.OptionKv{})
		if err != nil {
			logx.WithContext(ctx).Errorf("get url err:%v", err)
			continue
		}
		at = append(at, &sys.Attachment{
			Id:       att.ID,
			UseBy:    att.UseBy,
			FileUrl:  url,
			FileName: oss.GetFileNameWithPath(att.FilePath),
		})
	}
	pb := &sys.ProjectInfo{
		TenantCode:        string(po.TenantCode),
		CreatedTime:       po.CreatedTime.Unix(),
		ProjectID:         int64(po.ProjectID),
		ProjectName:       po.ProjectName,
		ProjectImg:        po.ProjectImg,
		AdminUserID:       po.AdminUserID,
		IsSysCreated:      po.IsSysCreated,
		Ppsm:              po.Ppsm,
		Tags:              po.Tags,
		Area:              &wrapperspb.FloatValue{Value: po.Area},
		Desc:              utils.ToRpcNullString(po.Desc),
		Position:          logic.ToSysPoint(po.Position),
		AreaCount:         po.AreaCount,
		UserCount:         po.UserCount,
		Address:           utils.ToRpcNullString(po.Address),
		DeviceCount:       utils.ToRpcNullInt64(po.DeviceCount),
		Areas:             areamanagelogic.AreaInfosToPb(ctx, svcCtx, po.Areas),
		Sort:              po.Sort,
		Status:            po.Status,
		AlarmStatus:       po.AlarmStatus,
		Type:              po.Type,
		Attachments:       at,
		DeviceOnlineCount: utils.ToRpcNullInt64(po.DeviceOnlineCount),
	}
	return pb
}
