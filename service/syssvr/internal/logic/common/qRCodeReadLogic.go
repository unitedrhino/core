package commonlogic

import (
	"context"
	"gitee.com/unitedrhino/share/errors"
	"github.com/alibabacloud-go/tea/tea"

	"gitee.com/unitedrhino/core/service/syssvr/internal/svc"
	"gitee.com/unitedrhino/core/service/syssvr/pb/sys"
	"github.com/silenceper/wechat/v2/miniprogram/qrcode"

	"github.com/zeromicro/go-zero/core/logx"
)

type QRCodeReadLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewQRCodeReadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *QRCodeReadLogic {
	return &QRCodeReadLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *QRCodeReadLogic) QRCodeRead(in *sys.QRCodeReadReq) (*sys.QRCodeReadResp, error) {
	cli, er := l.svcCtx.Cm.GetClients(l.ctx, "")
	if er != nil || cli.MiniProgram == nil {
		return nil, errors.System.AddDetail(er)
	}
	ret, err := cli.MiniProgram.GetQRCode().GetWXACodeUnlimit(qrcode.QRCoder{Page: in.Page, Scene: in.Scene, EnvVersion: in.EnvVersion, CheckPath: tea.Bool(false)})
	if err != nil {
		return nil, errors.System.AddDetail(err)
	}
	return &sys.QRCodeReadResp{Buffer: ret}, nil
}
