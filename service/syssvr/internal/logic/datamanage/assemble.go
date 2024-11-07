package datamanagelogic

import (
	"context"
	"gitee.com/unitedrhino/core/service/syssvr/internal/repo/relationDB"
	"gitee.com/unitedrhino/core/service/syssvr/internal/svc"
	"gitee.com/unitedrhino/core/service/syssvr/pb/sys"
	"gitee.com/unitedrhino/share/domain/userDataAuth"
)

func transAreaPoToPb(po *relationDB.SysDataArea) *sys.DataArea {
	return &sys.DataArea{
		AreaID:         int64(po.AreaID),
		AuthType:       po.AuthType,
		IsAuthChildren: po.IsAuthChildren,
		TargetID:       po.TargetID,
		TargetType:     po.TargetType,
	}
}

func transProjectPoToPb(po *relationDB.SysDataProject) *sys.DataProject {
	return &sys.DataProject{
		ProjectID: int64(po.ProjectID),
		AuthType:  po.AuthType,
		TargetID:  po.TargetID,
	}
}

func ToAuthAreaDo(ctx context.Context, svcCtx *svc.ServiceContext, area *sys.DataArea) *userDataAuth.Area {
	if area == nil {
		return nil
	}
	a, err := svcCtx.AreaCache.GetData(ctx, area.AreaID)
	if err != nil {
		return nil
	}
	return &userDataAuth.Area{AreaID: area.AreaID, AreaIDPath: a.AreaIDPath, AuthType: area.AuthType, IsAuthChildren: area.IsAuthChildren}
}
func ToAuthAreaDos(ctx context.Context, svcCtx *svc.ServiceContext, areas []*sys.DataArea) (ret []*userDataAuth.Area) {
	if len(areas) == 0 {
		return
	}
	for _, v := range areas {
		a := ToAuthAreaDo(ctx, svcCtx, v)
		if a == nil {
			continue
		}
		ret = append(ret, a)
	}
	return
}

func DBToAuthAreaDo(area *relationDB.SysDataArea) *userDataAuth.Area {
	if area == nil {
		return nil
	}
	return &userDataAuth.Area{AreaID: int64(area.AreaID), AuthType: area.AuthType, IsAuthChildren: area.IsAuthChildren}
}
func DBToAuthAreaDos(areas []*relationDB.SysDataArea) (ret []*userDataAuth.Area) {
	if len(areas) == 0 {
		return
	}
	for _, v := range areas {
		ret = append(ret, DBToAuthAreaDo(v))
	}
	return
}

func ToAuthProjectDo(area *sys.DataProject) *userDataAuth.Project {
	if area == nil {
		return nil
	}
	return &userDataAuth.Project{ProjectID: area.ProjectID, AuthType: area.AuthType}
}
func ToAuthProjectDos(areas []*sys.DataProject) (ret []*userDataAuth.Project) {
	if len(areas) == 0 {
		return
	}
	for _, v := range areas {
		ret = append(ret, ToAuthProjectDo(v))
	}
	return
}

func DBToAuthProjectDo(area *relationDB.SysDataProject) *userDataAuth.Project {
	if area == nil {
		return nil
	}
	return &userDataAuth.Project{ProjectID: int64(area.ProjectID), AuthType: area.AuthType}
}
func DBToAuthProjectDos(areas []*relationDB.SysDataProject) (ret []*userDataAuth.Project) {
	if len(areas) == 0 {
		return
	}
	for _, v := range areas {
		ret = append(ret, DBToAuthProjectDo(v))
	}
	return
}

func ToUserAreaApplyInfos(in []*relationDB.SysUserAreaApply) (ret []*sys.UserAreaApplyInfo) {
	for _, v := range in {
		ret = append(ret, &sys.UserAreaApplyInfo{
			Id:          v.ID,
			UserID:      v.UserID,
			AreaID:      int64(v.AreaID),
			AuthType:    v.AuthType,
			CreatedTime: v.CreatedTime.Unix(),
		})
	}
	return
}
