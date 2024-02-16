package dictmanagelogic

import (
	"context"
	"gitee.com/i-Things/core/service/syssvr/internal/repo/relationDB"
	"gitee.com/i-Things/core/service/syssvr/internal/svc"
	"gitee.com/i-Things/core/service/syssvr/pb/sys"
	"gitee.com/i-Things/share/def"
	"github.com/spf13/cast"
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
	po := &relationDB.SysDictInfo{
		ParentID: in.ParentID,
		Name:     in.Name,
		Type:     in.Type,
		Desc:     in.Desc.GetValue(),
		Body:     in.Body.GetValue(),
	}

	err := relationDB.NewDictInfoRepo(l.ctx).Insert(l.ctx, po)
	if err != nil {
		return nil, err
	}
	po, err = relationDB.NewDictInfoRepo(l.ctx).FindOne(l.ctx, po.ID)
	if err != nil {
		return nil, err
	}
	po.DictIDPath = cast.ToString(po.ID) + "-"
	if po.ParentID != 0 && po.ParentID != def.RootNode {
		parent, err := relationDB.NewDictInfoRepo(l.ctx).FindOne(l.ctx, in.ParentID)
		if err != nil {
			return nil, err
		}
		po.DictIDPath = parent.DictIDPath + po.DictIDPath
	}
	err = relationDB.NewDictInfoRepo(l.ctx).Update(l.ctx, po)
	if err != nil {
		return nil, err
	}
	return &sys.WithID{Id: po.ID}, err
}
