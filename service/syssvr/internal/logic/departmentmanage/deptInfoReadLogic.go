package departmentmanagelogic

import (
	"context"
	"gitee.com/unitedrhino/core/service/syssvr/internal/repo/relationDB"
	"gitee.com/unitedrhino/share/ctxs"
	"gitee.com/unitedrhino/share/def"
	"gitee.com/unitedrhino/share/utils"

	"gitee.com/unitedrhino/core/service/syssvr/internal/svc"
	"gitee.com/unitedrhino/core/service/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeptInfoReadLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeptInfoReadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeptInfoReadLogic {
	return &DeptInfoReadLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeptInfoReadLogic) DeptInfoRead(in *sys.DeptInfoReadReq) (*sys.DeptInfo, error) {
	var po *relationDB.SysDeptInfo
	var err error
	if in.Id <= def.RootNode {
		po = &relationDB.SysDeptInfo{
			ID:     def.RootNode,
			Name:   "根节点",
			Status: def.True,
			Sort:   1,
		}
		uc := ctxs.GetUserCtx(l.ctx)
		ti, _ := l.svcCtx.TenantCache.GetData(l.ctx, uc.TenantCode)
		if ti != nil {
			po.Name = ti.Name
		}
		t, err := relationDB.NewUserInfoRepo(l.ctx).CountByFilter(l.ctx, relationDB.UserInfoFilter{})
		if err != nil {
			return nil, err
		}
		po.UserCount = t
		if in.WithChildren {
			pos, err := relationDB.NewDeptInfoRepo(l.ctx).FindByFilter(l.ctx,
				relationDB.DeptInfoFilter{ParentID: def.RootNode}, nil)
			if err != nil {
				return nil, err
			}
			po.Children = pos
		}
	} else {
		po, err = relationDB.NewDeptInfoRepo(l.ctx).FindOneByFilter(l.ctx, relationDB.DeptInfoFilter{
			ID:           in.Id,
			WithChildren: in.WithChildren,
		})
		if err != nil {
			return nil, err
		}
	}

	ret := utils.Copy[sys.DeptInfo](po)
	if po.ID != def.RootNode && in.WithFather {
		fatherIDs := utils.GetIDPath(po.IDPath)
		if len(fatherIDs) > 1 {
			fs, err := relationDB.NewDeptInfoRepo(l.ctx).FindByFilter(l.ctx, relationDB.DeptInfoFilter{IDs: fatherIDs}, nil)
			if err != nil {
				return nil, err
			}
			var fsMap = map[int64]*sys.DeptInfo{}
			for _, v := range fs {
				fsMap[v.ID] = utils.Copy[sys.DeptInfo](v)
			}
			FillFather(ret, fsMap)
		}
	}
	return ret, nil
}

func FillFather(in *sys.DeptInfo, fsMap map[int64]*sys.DeptInfo) {
	f := fsMap[in.ParentID]
	if f != nil {
		in.Parent = f
		FillFather(f, fsMap)
	}
}
