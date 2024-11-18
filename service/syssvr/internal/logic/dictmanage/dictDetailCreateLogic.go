package dictmanagelogic

import (
	"context"
	"fmt"
	"gitee.com/unitedrhino/core/service/syssvr/internal/repo/relationDB"
	"gitee.com/unitedrhino/core/service/syssvr/internal/svc"
	"gitee.com/unitedrhino/core/service/syssvr/pb/sys"
	"gitee.com/unitedrhino/share/ctxs"
	"gitee.com/unitedrhino/share/def"
	"gitee.com/unitedrhino/share/errors"

	"github.com/zeromicro/go-zero/core/logx"
)

type DictDetailCreateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDictDetailCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DictDetailCreateLogic {
	return &DictDetailCreateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DictDetailCreateLogic) DictDetailCreate(in *sys.DictDetail) (*sys.WithID, error) {
	if err := ctxs.IsRoot(l.ctx); err != nil {
		return nil, err
	}
	_, err := relationDB.NewDictInfoRepo(l.ctx).FindOneByFilter(l.ctx, relationDB.DictInfoFilter{Code: in.GetDictCode()})
	if err != nil {
		if errors.Cmp(err, errors.NotFind) {
			return nil, errors.Parameter.AddMsg("字典未定义")
		}
		return nil, err
	}

	po := relationDB.SysDictDetail{
		DictCode: in.DictCode,
		Label:    in.Label,
		Value:    in.Value,
		ParentID: in.ParentID,
		Status:   in.Status,
		Sort:     in.Sort,
		Desc:     in.Desc.GetValue(),
		Body:     in.Body.GetValue(),
	}
	var parent = &relationDB.SysDictDetail{
		ID: def.RootNode,
	}
	if in.ParentID > def.RootNode {
		parent, err = relationDB.NewDictDetailRepo(l.ctx).FindOne(l.ctx, in.ParentID)
		if err != nil {
			return nil, err
		}
	} else {
		po.ParentID = def.RootNode
	}
	err = relationDB.NewDictDetailRepo(l.ctx).Insert(l.ctx, &po)
	if err == nil && parent != nil {
		po.IDPath = fmt.Sprintf("%s%v-", parent.IDPath, po.ID)
	}
	err = relationDB.NewDictDetailRepo(l.ctx).Update(l.ctx, &po)
	return &sys.WithID{Id: po.ID}, err
}
