package projectmanagelogic

import (
	"gitee.com/i-Things/core/service/syssvr/internal/logic"
	"gitee.com/i-Things/core/service/syssvr/internal/repo/relationDB"
	"gitee.com/i-Things/core/service/syssvr/pb/sys"
	"gitee.com/i-Things/share/utils"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func ProjectInfoToPb(po *relationDB.SysProjectInfo) *sys.ProjectInfo {
	pb := &sys.ProjectInfo{
		CreatedTime: po.CreatedTime.Unix(),
		ProjectID:   int64(po.ProjectID),
		ProjectName: po.ProjectName,
		AdminUserID: po.AdminUserID,
		Ppsm:        po.Ppsm,
		Area:        &wrapperspb.FloatValue{Value: po.Area},
		Desc:        utils.ToRpcNullString(po.Desc),
		Position:    logic.ToSysPoint(po.Position),
		AreaCount:   po.AreaCount,
	}
	return pb
}
