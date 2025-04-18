// Code generated by goctl. DO NOT EDIT.
// goctl 1.7.1
// Source: sys.proto

package server

import (
	"context"

	"gitee.com/unitedrhino/core/service/syssvr/internal/logic/common"
	"gitee.com/unitedrhino/core/service/syssvr/internal/svc"
	"gitee.com/unitedrhino/core/service/syssvr/pb/sys"
)

type CommonServer struct {
	svcCtx *svc.ServiceContext
	sys.UnimplementedCommonServer
}

func NewCommonServer(svcCtx *svc.ServiceContext) *CommonServer {
	return &CommonServer{
		svcCtx: svcCtx,
	}
}

func (s *CommonServer) Config(ctx context.Context, in *sys.Empty) (*sys.ConfigResp, error) {
	l := commonlogic.NewConfigLogic(ctx, s.svcCtx)
	return l.Config(in)
}

func (s *CommonServer) QrCodeRead(ctx context.Context, in *sys.QRCodeReadReq) (*sys.QRCodeReadResp, error) {
	l := commonlogic.NewQrCodeReadLogic(ctx, s.svcCtx)
	return l.QrCodeRead(in)
}

func (s *CommonServer) WeatherRead(ctx context.Context, in *sys.WeatherReadReq) (*sys.WeatherReadResp, error) {
	l := commonlogic.NewWeatherReadLogic(ctx, s.svcCtx)
	return l.WeatherRead(in)
}

func (s *CommonServer) SlotInfoIndex(ctx context.Context, in *sys.SlotInfoIndexReq) (*sys.SlotInfoIndexResp, error) {
	l := commonlogic.NewSlotInfoIndexLogic(ctx, s.svcCtx)
	return l.SlotInfoIndex(in)
}

func (s *CommonServer) SlotInfoCreate(ctx context.Context, in *sys.SlotInfo) (*sys.WithID, error) {
	l := commonlogic.NewSlotInfoCreateLogic(ctx, s.svcCtx)
	return l.SlotInfoCreate(in)
}

func (s *CommonServer) SlotInfoUpdate(ctx context.Context, in *sys.SlotInfo) (*sys.Empty, error) {
	l := commonlogic.NewSlotInfoUpdateLogic(ctx, s.svcCtx)
	return l.SlotInfoUpdate(in)
}

func (s *CommonServer) SlotInfoDelete(ctx context.Context, in *sys.WithID) (*sys.Empty, error) {
	l := commonlogic.NewSlotInfoDeleteLogic(ctx, s.svcCtx)
	return l.SlotInfoDelete(in)
}

func (s *CommonServer) SlotInfoRead(ctx context.Context, in *sys.WithID) (*sys.SlotInfo, error) {
	l := commonlogic.NewSlotInfoReadLogic(ctx, s.svcCtx)
	return l.SlotInfoRead(in)
}

func (s *CommonServer) ServiceInfoRead(ctx context.Context, in *sys.WithCode) (*sys.ServiceInfo, error) {
	l := commonlogic.NewServiceInfoReadLogic(ctx, s.svcCtx)
	return l.ServiceInfoRead(in)
}

func (s *CommonServer) ServiceInfoUpdate(ctx context.Context, in *sys.ServiceInfo) (*sys.Empty, error) {
	l := commonlogic.NewServiceInfoUpdateLogic(ctx, s.svcCtx)
	return l.ServiceInfoUpdate(in)
}

func (s *CommonServer) ThirdDeptRead(ctx context.Context, in *sys.ThirdDeptInfoReadReq) (*sys.DeptInfo, error) {
	l := commonlogic.NewThirdDeptReadLogic(ctx, s.svcCtx)
	return l.ThirdDeptRead(in)
}

func (s *CommonServer) ThirdDeptIndex(ctx context.Context, in *sys.ThirdDeptInfoIndexReq) (*sys.DeptInfoIndexResp, error) {
	l := commonlogic.NewThirdDeptIndexLogic(ctx, s.svcCtx)
	return l.ThirdDeptIndex(in)
}
