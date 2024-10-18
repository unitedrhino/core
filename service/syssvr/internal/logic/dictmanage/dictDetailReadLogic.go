package dictmanagelogic

import (
	"context"
	"gitee.com/unitedrhino/core/service/syssvr/internal/repo/relationDB"
	"gitee.com/unitedrhino/share/def"
	"gitee.com/unitedrhino/share/utils"

	"gitee.com/unitedrhino/core/service/syssvr/internal/svc"
	"gitee.com/unitedrhino/core/service/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type DictDetailReadLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDictDetailReadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DictDetailReadLogic {
	return &DictDetailReadLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DictDetailReadLogic) DictDetailRead(in *sys.DictDetailReadReq) (*sys.DictDetail, error) {
	var po *relationDB.SysDictDetail
	var err error
	if in.Id <= def.RootNode {
		po = &relationDB.SysDictDetail{
			ID:       def.RootNode,
			DictCode: in.DictCode,
			Label:    "根节点",
			Status:   def.True,
			Sort:     1,
		}
		if in.WithChildren {
			pos, err := relationDB.NewDictDetailRepo(l.ctx).FindByFilter(l.ctx,
				relationDB.DictDetailFilter{DictCode: in.DictCode, ParentID: def.RootNode}, nil)
			if err != nil {
				return nil, err
			}
			po.Children = pos
		}
	} else {
		po, err = relationDB.NewDictDetailRepo(l.ctx).FindOneByFilter(l.ctx, relationDB.DictDetailFilter{
			ID:           in.Id,
			DictCode:     in.DictCode,
			WithChildren: in.WithChildren,
			Value:        in.Value,
		})
		if err != nil {
			return nil, err
		}
	}

	ret := utils.Copy[sys.DictDetail](po)
	if po.ID != def.RootNode && in.WithFather {
		fatherIDs := utils.GetIDPath(po.IDPath)
		if len(fatherIDs) > 1 {
			fs, err := relationDB.NewDictDetailRepo(l.ctx).FindByFilter(l.ctx, relationDB.DictDetailFilter{IDs: fatherIDs}, nil)
			if err != nil {
				return nil, err
			}
			var fsMap = map[int64]*sys.DictDetail{}
			for _, v := range fs {
				fsMap[v.ID] = utils.Copy[sys.DictDetail](v)
			}
			FillFather(ret, fsMap)
		}

	}
	return ret, nil
}
func FillFather(in *sys.DictDetail, fsMap map[int64]*sys.DictDetail) {
	f := fsMap[in.ParentID]
	if f != nil {
		in.Parent = f
		FillFather(f, fsMap)
	}
}
