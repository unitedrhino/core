// Code generated by goctl. DO NOT EDIT.
// Source: timedjob.proto

package server

import (
	"context"

	"gitee.com/unitedrhino/core/service/timed/timedjobsvr/internal/logic/timedmanage"
	"gitee.com/unitedrhino/core/service/timed/timedjobsvr/internal/svc"
	"gitee.com/unitedrhino/core/service/timed/timedjobsvr/pb/timedjob"
)

type TimedManageServer struct {
	svcCtx *svc.ServiceContext
	timedjob.UnimplementedTimedManageServer
}

func NewTimedManageServer(svcCtx *svc.ServiceContext) *TimedManageServer {
	return &TimedManageServer{
		svcCtx: svcCtx,
	}
}

func (s *TimedManageServer) TaskGroupCreate(ctx context.Context, in *timedjob.TaskGroup) (*timedjob.Response, error) {
	l := timedmanagelogic.NewTaskGroupCreateLogic(ctx, s.svcCtx)
	return l.TaskGroupCreate(in)
}

func (s *TimedManageServer) TaskGroupUpdate(ctx context.Context, in *timedjob.TaskGroup) (*timedjob.Response, error) {
	l := timedmanagelogic.NewTaskGroupUpdateLogic(ctx, s.svcCtx)
	return l.TaskGroupUpdate(in)
}

func (s *TimedManageServer) TaskGroupDelete(ctx context.Context, in *timedjob.WithCode) (*timedjob.Response, error) {
	l := timedmanagelogic.NewTaskGroupDeleteLogic(ctx, s.svcCtx)
	return l.TaskGroupDelete(in)
}

func (s *TimedManageServer) TaskGroupIndex(ctx context.Context, in *timedjob.TaskGroupIndexReq) (*timedjob.TaskGroupIndexResp, error) {
	l := timedmanagelogic.NewTaskGroupIndexLogic(ctx, s.svcCtx)
	return l.TaskGroupIndex(in)
}

func (s *TimedManageServer) TaskGroupRead(ctx context.Context, in *timedjob.WithCode) (*timedjob.TaskGroup, error) {
	l := timedmanagelogic.NewTaskGroupReadLogic(ctx, s.svcCtx)
	return l.TaskGroupRead(in)
}

func (s *TimedManageServer) TaskInfoCreate(ctx context.Context, in *timedjob.TaskInfo) (*timedjob.Response, error) {
	l := timedmanagelogic.NewTaskInfoCreateLogic(ctx, s.svcCtx)
	return l.TaskInfoCreate(in)
}

func (s *TimedManageServer) TaskInfoUpdate(ctx context.Context, in *timedjob.TaskInfo) (*timedjob.Response, error) {
	l := timedmanagelogic.NewTaskInfoUpdateLogic(ctx, s.svcCtx)
	return l.TaskInfoUpdate(in)
}

func (s *TimedManageServer) TaskInfoDelete(ctx context.Context, in *timedjob.WithGroupCode) (*timedjob.Response, error) {
	l := timedmanagelogic.NewTaskInfoDeleteLogic(ctx, s.svcCtx)
	return l.TaskInfoDelete(in)
}

func (s *TimedManageServer) TaskInfoIndex(ctx context.Context, in *timedjob.TaskInfoIndexReq) (*timedjob.TaskInfoIndexResp, error) {
	l := timedmanagelogic.NewTaskInfoIndexLogic(ctx, s.svcCtx)
	return l.TaskInfoIndex(in)
}

func (s *TimedManageServer) TaskInfoRead(ctx context.Context, in *timedjob.WithGroupCode) (*timedjob.TaskInfo, error) {
	l := timedmanagelogic.NewTaskInfoReadLogic(ctx, s.svcCtx)
	return l.TaskInfoRead(in)
}

func (s *TimedManageServer) TaskLogIndex(ctx context.Context, in *timedjob.TaskLogIndexReq) (*timedjob.TaskLogIndexResp, error) {
	l := timedmanagelogic.NewTaskLogIndexLogic(ctx, s.svcCtx)
	return l.TaskLogIndex(in)
}

// 发送延时请求,如果任务不存在,则会自动创建,但是自动创建的需要填写param
func (s *TimedManageServer) TaskSend(ctx context.Context, in *timedjob.TaskSendReq) (*timedjob.TaskWithTaskID, error) {
	l := timedmanagelogic.NewTaskSendLogic(ctx, s.svcCtx)
	return l.TaskSend(in)
}

func (s *TimedManageServer) TaskCancel(ctx context.Context, in *timedjob.TaskWithTaskID) (*timedjob.Response, error) {
	l := timedmanagelogic.NewTaskCancelLogic(ctx, s.svcCtx)
	return l.TaskCancel(in)
}
