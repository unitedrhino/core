// Code generated by goctl. DO NOT EDIT.
// goctl 1.7.1
// Source: sys.proto

package server

import (
	"context"

	"gitee.com/unitedrhino/core/service/syssvr/internal/logic/projectmanage"
	"gitee.com/unitedrhino/core/service/syssvr/internal/svc"
	"gitee.com/unitedrhino/core/service/syssvr/pb/sys"
)

type ProjectManageServer struct {
	svcCtx *svc.ServiceContext
	sys.UnimplementedProjectManageServer
}

func NewProjectManageServer(svcCtx *svc.ServiceContext) *ProjectManageServer {
	return &ProjectManageServer{
		svcCtx: svcCtx,
	}
}

// 新增项目
func (s *ProjectManageServer) ProjectInfoCreate(ctx context.Context, in *sys.ProjectInfo) (*sys.ProjectWithID, error) {
	l := projectmanagelogic.NewProjectInfoCreateLogic(ctx, s.svcCtx)
	return l.ProjectInfoCreate(in)
}

// 更新项目
func (s *ProjectManageServer) ProjectInfoUpdate(ctx context.Context, in *sys.ProjectInfo) (*sys.Empty, error) {
	l := projectmanagelogic.NewProjectInfoUpdateLogic(ctx, s.svcCtx)
	return l.ProjectInfoUpdate(in)
}

// 删除项目
func (s *ProjectManageServer) ProjectInfoDelete(ctx context.Context, in *sys.ProjectWithID) (*sys.Empty, error) {
	l := projectmanagelogic.NewProjectInfoDeleteLogic(ctx, s.svcCtx)
	return l.ProjectInfoDelete(in)
}

// 获取项目信息详情
func (s *ProjectManageServer) ProjectInfoRead(ctx context.Context, in *sys.ProjectWithID) (*sys.ProjectInfo, error) {
	l := projectmanagelogic.NewProjectInfoReadLogic(ctx, s.svcCtx)
	return l.ProjectInfoRead(in)
}

// 获取项目信息列表
func (s *ProjectManageServer) ProjectInfoIndex(ctx context.Context, in *sys.ProjectInfoIndexReq) (*sys.ProjectInfoIndexResp, error) {
	l := projectmanagelogic.NewProjectInfoIndexLogic(ctx, s.svcCtx)
	return l.ProjectInfoIndex(in)
}

func (s *ProjectManageServer) ProjectProfileRead(ctx context.Context, in *sys.ProjectProfileReadReq) (*sys.ProjectProfile, error) {
	l := projectmanagelogic.NewProjectProfileReadLogic(ctx, s.svcCtx)
	return l.ProjectProfileRead(in)
}

func (s *ProjectManageServer) ProjectProfileUpdate(ctx context.Context, in *sys.ProjectProfile) (*sys.Empty, error) {
	l := projectmanagelogic.NewProjectProfileUpdateLogic(ctx, s.svcCtx)
	return l.ProjectProfileUpdate(in)
}

func (s *ProjectManageServer) ProjectProfileIndex(ctx context.Context, in *sys.ProjectProfileIndexReq) (*sys.ProjectProfileIndexResp, error) {
	l := projectmanagelogic.NewProjectProfileIndexLogic(ctx, s.svcCtx)
	return l.ProjectProfileIndex(in)
}
