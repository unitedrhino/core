package cache

import (
	"context"
	"github.com/maypok86/otter"
	"github.com/zeromicro/go-zero/core/logx"
	"time"
)

var (
	deptAuthCache otter.Cache[int64, []DeptAuth]
)

type DeptAuth struct {
	DeptID   int64
	DeptPath string
}

func init() {
	cache, err := otter.MustBuilder[int64, []DeptAuth](10_000).
		CollectStats().
		Cost(func(key int64, value []DeptAuth) uint32 {
			return 1
		}).
		WithTTL(time.Minute + time.Second*5).
		Build()
	logx.Must(err)
	deptAuthCache = cache
}

func ClearDeptAuth(userID int64) {
	deptAuthCache.Delete(userID)
}

func GetDeptAuth(ctx context.Context, userID int64, roleIDs []int64) ([]DeptAuth, error) {
	ret, ok := deptAuthCache.Get(userID)
	if ok {
		return ret, nil
	}
	//
	//filter := relationDB.DeptUserFilter{
	//	UserID: userID,
	//}

	//areas, err := relationDB.NewDeptUserRepo(ctx).FindByFilter(ctxs.WithAllProject(ctx), filter, nil)
	//if err != nil {
	//	logx.WithContext(ctx).Error(err)
	//	return nil, err
	//}
	var deptAuth = []DeptAuth{}
	//var deptPathAuth = make(map[string]int64, len(areas))
	//for _, po := range areas {
	//	if po.IsAuthChildren != def.True {
	//		old, ok := deptAuth[po.DeptID]
	//		if !ok  { //取权限大的
	//			deptAuth[po.DeptID] =
	//		}
	//		continue
	//	}
	//	old, ok := auth.AreaPath[po.AreaIDPath]
	//	if !ok || po.AuthType < old { //取权限大的
	//		auth.AreaPath[po.AreaIDPath] = po.AuthType
	//	}
	//}
	//deptAuthCache.Set(userID, deptAuth)
	return deptAuth, nil
}
