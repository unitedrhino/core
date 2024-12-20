// Code generated by goctl. DO NOT EDIT.
// goctl 1.7.1
// Source: sys.proto

package server

import (
	"context"

	"gitee.com/unitedrhino/core/service/syssvr/internal/logic/ops"
	"gitee.com/unitedrhino/core/service/syssvr/internal/svc"
	"gitee.com/unitedrhino/core/service/syssvr/pb/sys"
)

type OpsServer struct {
	svcCtx *svc.ServiceContext
	sys.UnimplementedOpsServer
}

func NewOpsServer(svcCtx *svc.ServiceContext) *OpsServer {
	return &OpsServer{
		svcCtx: svcCtx,
	}
}

// 维护工单  Work Order
func (s *OpsServer) OpsWorkOrderCreate(ctx context.Context, in *sys.OpsWorkOrder) (*sys.WithID, error) {
	l := opslogic.NewOpsWorkOrderCreateLogic(ctx, s.svcCtx)
	return l.OpsWorkOrderCreate(in)
}

func (s *OpsServer) OpsWorkOrderUpdate(ctx context.Context, in *sys.OpsWorkOrder) (*sys.Empty, error) {
	l := opslogic.NewOpsWorkOrderUpdateLogic(ctx, s.svcCtx)
	return l.OpsWorkOrderUpdate(in)
}

func (s *OpsServer) OpsWorkOrderIndex(ctx context.Context, in *sys.OpsWorkOrderIndexReq) (*sys.OpsWorkOrderIndexResp, error) {
	l := opslogic.NewOpsWorkOrderIndexLogic(ctx, s.svcCtx)
	return l.OpsWorkOrderIndex(in)
}

// 反馈
func (s *OpsServer) OpsFeedbackCreate(ctx context.Context, in *sys.OpsFeedback) (*sys.WithID, error) {
	l := opslogic.NewOpsFeedbackCreateLogic(ctx, s.svcCtx)
	return l.OpsFeedbackCreate(in)
}

func (s *OpsServer) OpsFeedbackUpdate(ctx context.Context, in *sys.OpsFeedback) (*sys.Empty, error) {
	l := opslogic.NewOpsFeedbackUpdateLogic(ctx, s.svcCtx)
	return l.OpsFeedbackUpdate(in)
}

func (s *OpsServer) OpsFeedbackIndex(ctx context.Context, in *sys.OpsFeedbackIndexReq) (*sys.OpsFeedbackIndexResp, error) {
	l := opslogic.NewOpsFeedbackIndexLogic(ctx, s.svcCtx)
	return l.OpsFeedbackIndex(in)
}
