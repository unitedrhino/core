package appmanagelogic

import (
	"context"
	"gitee.com/unitedrhino/core/service/syssvr/internal/repo/relationDB"
	"gitee.com/unitedrhino/share/ctxs"
	"gitee.com/unitedrhino/share/def"
	"gitee.com/unitedrhino/share/errors"

	"gitee.com/unitedrhino/core/service/syssvr/internal/svc"
	"gitee.com/unitedrhino/core/service/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type AppInfoDeleteLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAppInfoDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AppInfoDeleteLogic {
	return &AppInfoDeleteLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *AppInfoDeleteLogic) AppInfoDelete(in *sys.WithIDCode) (*sys.Empty, error) {
	if err := ctxs.IsRoot(l.ctx); err != nil {
		return nil, err
	}
	f := relationDB.AppInfoFilter{ID: in.Id}
	if in.Code != "" {
		f.Codes = []string{in.Code}
	}
	info, err := relationDB.NewAppInfoRepo(l.ctx).FindOneByFilter(l.ctx, f)
	if err != nil {
		return nil, err
	}
	if info.Code == def.AppCore {
		return nil, errors.Parameter.AddMsg("core应用不允许删除")
	}
	err = relationDB.NewAppInfoRepo(l.ctx).DeleteByFilter(l.ctx, f)
	return &sys.Empty{}, err
}
