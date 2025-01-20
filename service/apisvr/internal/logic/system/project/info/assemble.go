package info

import (
	"gitee.com/unitedrhino/core/service/apisvr/internal/logic"
	"gitee.com/unitedrhino/core/service/apisvr/internal/types"
	"gitee.com/unitedrhino/core/service/syssvr/pb/sys"
	"gitee.com/unitedrhino/share/utils"
)

func ToProjectPb(in *types.ProjectInfo) *sys.ProjectInfo {
	if in == nil {
		return nil
	}
	return &sys.ProjectInfo{
		Tags:               in.Tags,
		ProjectID:          in.ProjectID,
		ProjectName:        in.ProjectName,
		ProjectImg:         in.ProjectImg,
		IsUpdateProjectImg: in.IsUpdateProjectImg,
		AdminUserID:        in.AdminUserID,
		Position:           logic.ToSysPointRpc(in.Position),
		Desc:               utils.ToRpcNullString(in.Desc),
		AreaCount:          in.AreaCount,
		Area:               utils.ToRpcNullFloat32(in.Area),
		Ppsm:               in.Ppsm,
		IsSysCreated:       in.IsSysCreated,
		Address:            utils.ToRpcNullString(in.Address),
	}
}
