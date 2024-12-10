package departmentmanagelogic

import (
	"context"
	"fmt"
	"gitee.com/unitedrhino/core/service/syssvr/internal/repo/relationDB"
	"gitee.com/unitedrhino/share/ctxs"
	"gitee.com/unitedrhino/share/def"
	"gitee.com/unitedrhino/share/errors"

	"gitee.com/unitedrhino/core/service/syssvr/internal/svc"
	"gitee.com/unitedrhino/core/service/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeptInfoCreateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeptInfoCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeptInfoCreateLogic {
	return &DeptInfoCreateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeptInfoCreateLogic) DeptInfoCreate(in *sys.DeptInfo) (*sys.WithID, error) {
	if err := ctxs.IsAdmin(l.ctx); err != nil {
		return nil, err
	}
	po := relationDB.SysDeptInfo{
		ParentID: in.ParentID,
		Name:     in.Name,
		Status:   in.Status,
		Sort:     in.Sort,
		Desc:     in.Desc.GetValue(),
	}
	var parent = &relationDB.SysDeptInfo{
		ID: def.RootNode,
	}
	var err error
	if in.ParentID > def.RootNode {
		parent, err = relationDB.NewDeptInfoRepo(l.ctx).FindOne(l.ctx, in.ParentID)
		if err != nil {
			if errors.Cmp(err, errors.NotFind) {
				return nil, errors.Parameter.AddMsg("未找到父节点")
			}
			return nil, err
		}
	} else {
		po.ParentID = def.RootNode
	}
	err = relationDB.NewDeptInfoRepo(l.ctx).Insert(l.ctx, &po)
	if err != nil {
		if errors.Cmp(err, errors.Duplicate) {
			return nil, errors.Parameter.AddMsg("名称重复")
		}
		return nil, err
	}
	po.IDPath = fmt.Sprintf("%s%v-", parent.IDPath, po.ID)
	err = relationDB.NewDeptInfoRepo(l.ctx).Update(l.ctx, &po)
	if err != nil {
		return nil, err
	}

	return &sys.WithID{Id: po.ID}, err
}
