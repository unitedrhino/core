package project

import (
	"gitee.com/unitedrhino/core/service/viewsvr/internal/repo/relationDB"
	"gitee.com/unitedrhino/core/service/viewsvr/internal/types"
)

func ToProjectInfoTypes(p *relationDB.ViewProjectInfo) *types.ProjectInfo {
	return &types.ProjectInfo{
		IndexImage:    p.IndexImage,
		Name:          p.Name,
		Desc:          p.Desc,
		CreatedUserID: p.CreatedUserID,
		Status:        p.Status,
	}
}
