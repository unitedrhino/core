package dictmanagelogic

import (
	"context"
	"fmt"
	"gitee.com/unitedrhino/core/service/syssvr/internal/repo/relationDB"
	"gitee.com/unitedrhino/share/ctxs"
	"gitee.com/unitedrhino/share/def"
	"gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/share/stores"
	"gorm.io/gorm"

	"gitee.com/unitedrhino/core/service/syssvr/internal/svc"
	"gitee.com/unitedrhino/core/service/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type DictDetailMultiCreateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDictDetailMultiCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DictDetailMultiCreateLogic {
	return &DictDetailMultiCreateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DictDetailMultiCreateLogic) DictDetailMultiCreate(in *sys.DictDetailMultiCreateReq) (*sys.Empty, error) {
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
	var pos []*relationDB.SysDictDetail
	for _, v := range in.List {
		po, err := l.DictDetailPbToPo(in.DictCode, v, nil)
		if err != nil {
			return nil, err
		}
		pos = append(pos, po)
	}
	err = stores.GetCommonConn(l.ctx).Transaction(func(tx *gorm.DB) error {
		err = relationDB.NewDictDetailRepo(tx).MultiInsert(l.ctx, pos)
		if err != nil {
			return err
		}
		var children []*relationDB.SysDictDetail

		for i, p := range in.List {
			parent := pos[i]
			for _, c := range p.Children {
				c.ParentID = parent.ID
				po, err := l.DictDetailPbToPo(in.DictCode, c, pos[i])
				if err != nil {
					return err
				}
				children = append(children, po)
			}
		}
		if len(children) > 0 {
			err = relationDB.NewDictDetailRepo(tx).MultiInsert(l.ctx, children)
			if err != nil {
				return err
			}
		}
		return nil
	})
	return &sys.Empty{}, err
}
func (l *DictDetailMultiCreateLogic) DictDetailPbToPo(dictCode string, v *sys.DictDetail, parent *relationDB.SysDictDetail) (*relationDB.SysDictDetail, error) {
	var err error
	po := relationDB.SysDictDetail{
		DictCode: dictCode,
		Label:    v.Label,
		Value:    v.Value,
		ParentID: v.ParentID,
		Status:   v.Status,
		Sort:     v.Sort,
		Desc:     v.Desc.GetValue(),
		Body:     v.Body.GetValue(),
	}
	if parent == nil {
		parent = &relationDB.SysDictDetail{
			ID: def.RootNode,
		}
		if po.ParentID > def.RootNode {
			parent, err = relationDB.NewDictDetailRepo(l.ctx).FindOne(l.ctx, po.ParentID)
			if err != nil {
				return nil, err
			}
		}
	}
	if po.ParentID < def.RootNode {
		po.ParentID = def.RootNode
	}

	err = relationDB.NewDictDetailRepo(l.ctx).Insert(l.ctx, &po)
	if err == nil && parent != nil {
		po.IDPath = fmt.Sprintf("%s%v-", parent.IDPath, po.ID)
	}
	//if len(v.Children) > 0 {
	//	for _, v := range v.Children {
	//		c, err := l.DictDetailPbToPo(dictCode, v)
	//		if err != nil {
	//			return nil, err
	//		}
	//		po.Children = append(po.Children, c)
	//	}
	//}
	return &po, nil
}
