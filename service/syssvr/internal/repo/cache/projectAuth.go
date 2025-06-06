package cache

import (
	"context"
	"gitee.com/unitedrhino/core/service/syssvr/internal/repo/relationDB"
	"gitee.com/unitedrhino/core/service/syssvr/pb/sys"
	"gitee.com/unitedrhino/share/ctxs"
	"gitee.com/unitedrhino/share/def"
	"github.com/maypok86/otter"
	"github.com/zeromicro/go-zero/core/logx"
	"time"
)

var (
	projectAuthCache otter.Cache[int64, map[int64]*sys.ProjectAuth]
)

func init() {
	cache, err := otter.MustBuilder[int64, map[int64]*sys.ProjectAuth](10_000).
		CollectStats().
		Cost(func(key int64, value map[int64]*sys.ProjectAuth) uint32 {
			return 1
		}).
		WithTTL(time.Minute + time.Second*5).
		Build()
	logx.Must(err)
	projectAuthCache = cache
}

func ClearProjectAuth(userID int64) {
	projectAuthCache.Delete(userID)
}

func GetProjectAuth(ctx context.Context, userID int64, roleIDs []int64) (map[int64]*sys.ProjectAuth, error) {
	ret, ok := projectAuthCache.Get(userID)
	if ok {
		return ret, nil
	}
	filter := relationDB.DataProjectFilter{
		Targets: []*relationDB.Target{{Type: def.TargetUser, ID: userID}},
	}
	for _, role := range roleIDs {
		filter.Targets = append(filter.Targets, &relationDB.Target{
			Type: def.TargetRole,
			ID:   role,
		})
	}
	poArr, err := relationDB.NewDataProjectRepo(ctx).FindByFilter(ctx, filter, nil)
	if err != nil {
		logx.WithContext(ctx).Error(err)
		return nil, err
	}
	var projectAuth = map[int64]*sys.ProjectAuth{}
	for _, po := range poArr {
		old := projectAuth[po.ProjectID]
		if old == nil || po.AuthType < old.AuthType { //取权限大的
			projectAuth[po.ProjectID] = &sys.ProjectAuth{
				Area:     nil,
				AuthType: po.AuthType,
			}
		}
	}
	for projectID, auth := range projectAuth {
		if auth.AuthType == def.AuthAdmin { //项目有管理权限不限制区域
			continue
		}
		filter := relationDB.DataAreaFilter{
			Targets:   []*relationDB.Target{{Type: def.TargetUser, ID: userID}},
			ProjectID: projectID,
		}
		for _, role := range roleIDs {
			filter.Targets = append(filter.Targets, &relationDB.Target{
				Type: def.TargetRole,
				ID:   role,
			})
		}
		areas, err := relationDB.NewDataAreaRepo(ctx).FindByFilter(ctxs.WithAllProject(ctx), filter, nil)
		if err != nil {
			logx.WithContext(ctx).Error(err)
			return nil, err
		}
		auth.Area = make(map[int64]int64, len(areas))
		auth.AreaPath = make(map[string]int64, len(areas))
		for _, po := range areas {
			if po.IsAuthChildren != def.True {
				old, ok := auth.Area[po.AreaID]
				if !ok || po.AuthType < old { //取权限大的
					auth.Area[po.AreaID] = po.AuthType
				}
				continue
			}
			old, ok := auth.AreaPath[po.AreaIDPath]
			if !ok || po.AuthType < old { //取权限大的
				auth.AreaPath[po.AreaIDPath] = po.AuthType
			}
		}
	}
	{ //项目的所有者拥有项目的最高权限
		adminPis, err := relationDB.NewProjectInfoRepo(ctx).FindByFilter(ctxs.WithAllProject(ctx), relationDB.ProjectInfoFilter{AdminUserID: userID}, nil)
		if err != nil {
			logx.WithContext(ctx).Error(err)
		}
		for _, po := range adminPis {
			projectAuth[int64(po.ProjectID)] = &sys.ProjectAuth{
				AuthType: def.AuthAdmin,
			}
		}
	}

	projectAuthCache.Set(userID, projectAuth)
	return projectAuth, nil
}
