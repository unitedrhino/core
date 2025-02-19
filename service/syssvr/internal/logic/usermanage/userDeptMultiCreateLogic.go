package usermanagelogic

import (
	"context"
	"gitee.com/unitedrhino/core/service/syssvr/internal/repo/cache"
	"gitee.com/unitedrhino/core/service/syssvr/internal/repo/relationDB"
	"gitee.com/unitedrhino/core/service/syssvr/internal/svc"
	"gitee.com/unitedrhino/core/service/syssvr/pb/sys"
	"gitee.com/unitedrhino/share/ctxs"
	"gitee.com/unitedrhino/share/def"
	"gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/share/utils"
	"github.com/spf13/cast"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserDeptMultiCreateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUserDeptMultiCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserDeptMultiCreateLogic {
	return &UserDeptMultiCreateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UserDeptMultiCreateLogic) UserDeptMultiCreate(in *sys.UserDeptMultiSaveReq) (*sys.Empty, error) {
	if err := ctxs.IsAdmin(l.ctx); err != nil {
		return nil, err
	}
	var idMap = make(map[int64]string)
	var idPaths []string
	if len(in.DeptIDs) != 0 {
		rs, err := relationDB.NewDeptInfoRepo(l.ctx).FindByFilter(l.ctx, relationDB.DeptInfoFilter{
			IDs: in.DeptIDs,
		}, nil)
		if err != nil {
			return nil, err
		}
		if len(rs) != len(in.DeptIDs) {
			return nil, errors.Parameter.WithMsg("有部门不存咋")
		}
		for _, v := range rs {
			idMap[v.ID] = v.IDPath
			idPaths = append(idPaths, v.IDPath)
		}
	}
	var datas []*relationDB.SysDeptUser
	for _, v := range in.DeptIDs {
		datas = append(datas, &relationDB.SysDeptUser{
			DeptID:     v,
			DeptIDPath: idMap[v],
			UserID:     in.UserID,
		})
	}
	err := relationDB.NewDeptUserRepo(l.ctx).MultiInsert(l.ctx, datas)
	if err == nil {
		l.svcCtx.UsersCache.SetData(l.ctx, in.UserID, nil)
		cache.ClearProjectAuth(in.UserID)
	}
	FillDeptUserCount(l.ctx, l.svcCtx, idPaths...)
	return &sys.Empty{}, err
}

func FillDeptUserCount(ctx context.Context, svcCtx *svc.ServiceContext, deptIDPaths ...string) error {
	logx.WithContext(ctx).Infof("FillDeptUserCount areaIDPaths:%v", deptIDPaths)
	defer utils.Recover(ctx)
	ctx = ctxs.WithRoot(ctx)
	log := logx.WithContext(ctx)
	var idMap = map[int64]struct{}{}
	for _, deptIDPath := range deptIDPaths {
		if deptIDPath == "" || deptIDPath == def.NotClassifiedPath {
			continue
		}
		ids := utils.GetIDPath(deptIDPath)
		var idPath string
		for _, id := range ids {
			idPath += cast.ToString(id) + "-"
			if _, ok := idMap[id]; ok {
				continue
			}
			idMap[id] = struct{}{}
			count, err := relationDB.NewDeptUserRepo(ctx).CountByFilter(ctx, relationDB.DeptUserFilter{DeptIDPath: idPath})
			if err != nil {
				log.Error(err)
				continue
			}
			err = relationDB.NewDeptInfoRepo(ctx).UpdateWithField(ctx, relationDB.DeptInfoFilter{ID: id}, map[string]any{
				"user_count": count,
			})
			if err != nil {
				log.Error(err)
				continue
			}
		}
	}

	return nil
}
