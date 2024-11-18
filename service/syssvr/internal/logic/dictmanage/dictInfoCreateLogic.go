package dictmanagelogic

import (
	"context"
	"gitee.com/unitedrhino/core/service/syssvr/internal/repo/relationDB"
	"gitee.com/unitedrhino/core/service/syssvr/internal/svc"
	"gitee.com/unitedrhino/core/service/syssvr/pb/sys"
	"gitee.com/unitedrhino/share/ctxs"
	"gitee.com/unitedrhino/share/errors"
	"github.com/zeromicro/go-zero/core/logx"
)

type DictInfoCreateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDictInfoCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DictInfoCreateLogic {
	return &DictInfoCreateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DictInfoCreateLogic) DictInfoCreate(in *sys.DictInfo) (*sys.WithID, error) {
	if err := ctxs.IsRoot(l.ctx); err != nil {
		return nil, err
	}
	if in.Code == "" || in.Group == "" {
		return &sys.WithID{}, errors.Parameter.AddMsg("code or group is empty")
	}
	po := &relationDB.SysDictInfo{
		Group:      in.Group,
		Code:       in.Code,
		Name:       in.Name,
		StructType: in.StructType,
		Desc:       in.Desc.GetValue(),
		Body:       in.Body.GetValue(),
	}
	err := relationDB.NewDictInfoRepo(l.ctx).Insert(l.ctx, po)
	if err != nil {
		return nil, err
	}
	return &sys.WithID{Id: po.ID}, err
}
