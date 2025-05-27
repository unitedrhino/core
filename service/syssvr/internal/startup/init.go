package startup

import (
	"context"
	"fmt"
	"gitee.com/unitedrhino/core/service/syssvr/internal/domain/dept"
	"gitee.com/unitedrhino/core/service/syssvr/internal/domain/module"
	"gitee.com/unitedrhino/core/service/syssvr/internal/event/day"
	"gitee.com/unitedrhino/core/service/syssvr/internal/event/deptSync"
	"gitee.com/unitedrhino/core/service/syssvr/internal/logic"
	accessmanagelogic "gitee.com/unitedrhino/core/service/syssvr/internal/logic/accessmanage"
	areamanagelogic "gitee.com/unitedrhino/core/service/syssvr/internal/logic/areamanage"
	departmentmanagelogic "gitee.com/unitedrhino/core/service/syssvr/internal/logic/departmentmanage"
	dictmanagelogic "gitee.com/unitedrhino/core/service/syssvr/internal/logic/dictmanage"
	modulemanagelogic "gitee.com/unitedrhino/core/service/syssvr/internal/logic/modulemanage"
	projectmanagelogic "gitee.com/unitedrhino/core/service/syssvr/internal/logic/projectmanage"
	tenantmanagelogic "gitee.com/unitedrhino/core/service/syssvr/internal/logic/tenantmanage"
	usermanagelogic "gitee.com/unitedrhino/core/service/syssvr/internal/logic/usermanage"
	"gitee.com/unitedrhino/core/service/syssvr/internal/repo/cache"
	"gitee.com/unitedrhino/core/service/syssvr/internal/repo/relationDB"
	"gitee.com/unitedrhino/core/service/syssvr/internal/svc"
	"gitee.com/unitedrhino/core/service/syssvr/pb/sys"
	"gitee.com/unitedrhino/core/service/timed/timedjobsvr/client/timedmanage"
	coreCache "gitee.com/unitedrhino/core/share/caches"
	"gitee.com/unitedrhino/core/share/domain/tenant"
	"gitee.com/unitedrhino/core/share/topics"
	"gitee.com/unitedrhino/share/caches"
	"gitee.com/unitedrhino/share/ctxs"
	"gitee.com/unitedrhino/share/def"
	"gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/share/utils"
	"github.com/zeromicro/go-zero/core/logx"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func Init(svcCtx *svc.ServiceContext) {
	ctx := ctxs.WithRoot(context.Background())
	utils.Go(ctx, func() {
		list, err := relationDB.NewTenantInfoRepo(ctx).FindByFilter(ctx, relationDB.TenantInfoFilter{}, nil)
		logx.Must(err)
		err = coreCache.InitTenant(ctx, logic.ToTenantInfoCaches(list)...)
		logx.Must(err)
	})
	VersionUpdate(svcCtx)
	InitCache(svcCtx)
	TableInit(svcCtx)
	InitEventBus(svcCtx)
	TimerInit(svcCtx)
	usermanagelogic.Init()
	InitSync(svcCtx)
}

func VersionUpdate(svcCtx *svc.ServiceContext) {
	ctx := ctxs.WithRoot(context.Background())
	{ //1.2.0->1.3.0
		po := relationDB.SysSlotInfo{
			Code: "userSubscribe", SubCode: def.UserSubscribeDevicePropertyReport2, SlotCode: "ithings", Uri: "/api/v1/things/slot/user/subscribe", Hosts: []string{"http://localhost:7788"}, AuthType: def.AppCore}
		relationDB.NewSlotInfoRepo(ctx).Insert(ctx, &po)
	}
}

func TableInit(svcCtx *svc.ServiceContext) {
	if !relationDB.NeedInitColumn {
		return
	}
	{
		root := "./etc/init/dict/"
		err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() {
				return nil
			}
			if !strings.HasSuffix(info.Name(), ".json") {
				return nil
			}
			body, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			ret, err := dictmanagelogic.NewDictMultiImportLogic(ctxs.WithRoot(context.TODO()), svcCtx).DictMultiImport(&sys.DictMultiImportReq{
				Dicts: string(body),
			})
			logx.Info("DictMultiImport", info.Name(), ret, err)

			return nil
		})
		if err != nil {
			logx.Error(err)
		}
	}
	{
		root := "./etc/init/module/"
		err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() {
				return nil
			}
			if !strings.HasSuffix(info.Name(), ".json") {
				return nil
			}
			body, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			moduleCode, _ := strings.CutSuffix(info.Name(), ".json")
			ret, err := modulemanagelogic.NewModuleMenuMultiImportLogic(ctxs.WithRoot(context.TODO()), svcCtx).ModuleMenuMultiImport(&sys.MenuMultiImportReq{
				ModuleCode: moduleCode,
				Mode:       module.MenuImportModeAll,
				Menu:       string(body),
			})
			logx.Info("ModuleMenuMultiImport", info.Name(), ret, err)

			return nil
		})
		if err != nil {
			logx.Error(err)
		}
	}
	{
		root := "./etc/init/access/"
		err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() {
				return nil
			}
			if !strings.HasSuffix(info.Name(), ".json") {
				return nil
			}
			body, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			moduleCode, _ := strings.CutSuffix(info.Name(), ".json")
			ret, err := accessmanagelogic.NewAccessInfoMultiImportLogic(ctxs.WithRoot(context.TODO()), svcCtx).AccessInfoMultiImport(&sys.AccessInfoMultiImportReq{
				Module: moduleCode,
				Access: string(body),
			})
			logx.Info("CommonSchemaMultiImport", info.Name(), ret, err)
			return nil
		})
		if err != nil {
			logx.Error(err)
		}
	}

	return
}

func InitSync(svcCtx *svc.ServiceContext) {
	ctx := ctxs.WithRoot(context.Background())
	pos, err := relationDB.NewDeptSyncJobRepo(ctx).FindByFilter(ctx, relationDB.DeptSyncJobFilter{
		Direction: dept.SyncDirectionFrom, SyncMode: dept.SyncModeRealTime}, nil)
	logx.Must(err)
	for _, v := range pos {
		err := departmentmanagelogic.DeptSyncAddDing(ctx, svcCtx, v)
		if err != nil {
			logx.Error(err)
			continue
		}
	}
}

func InitCache(svcCtx *svc.ServiceContext) {
	{
		tenantCache, err := caches.NewCache(caches.CacheConfig[tenant.Info, string]{
			KeyType:   topics.ServerCacheKeySysTenantInfo,
			FastEvent: svcCtx.FastEvent,
			GetData: func(ctx context.Context, key string) (*tenant.Info, error) {
				db := relationDB.NewTenantInfoRepo(ctx)
				if key == "" {
					key = ctxs.GetUserCtxNoNil(ctx).TenantCode
				}
				pi, err := db.FindOneByFilter(ctx, relationDB.TenantInfoFilter{
					Codes: []string{key}})
				pb := logic.ToTenantInfoCache(pi)
				return pb, err
			},
			ExpireTime: 20 * time.Minute,
		})
		logx.Must(err)
		svcCtx.TenantCache = tenantCache
	}
	{
		tenantCache, err := caches.NewCache(caches.CacheConfig[sys.TenantConfig, string]{
			KeyType:   topics.ServerCacheKeySysTenantConfig,
			FastEvent: svcCtx.FastEvent,
			GetData: func(ctx context.Context, key string) (*sys.TenantConfig, error) {
				db := relationDB.NewTenantConfigRepo(ctx)
				if key == "" {
					key = ctxs.GetUserCtxNoNil(ctx).TenantCode
				}
				pi, err := db.FindOneByFilter(ctx, relationDB.TenantConfigFilter{
					TenantCode: key})
				pb := tenantmanagelogic.ToTenantConfigPb(ctx, svcCtx, pi)
				return pb, err
			},
			ExpireTime: 20 * time.Minute,
		})
		logx.Must(err)
		svcCtx.TenantConfigCache = tenantCache
	}
	{
		userCache, err := caches.NewCache(caches.CacheConfig[sys.UserInfo, int64]{
			KeyType:   topics.ServerCacheKeySysUserInfo,
			FastEvent: svcCtx.FastEvent,
			GetData: func(ctx context.Context, key int64) (*sys.UserInfo, error) {
				db := relationDB.NewUserInfoRepo(ctx)
				if key == 0 {
					key = ctxs.GetUserCtxNoNil(ctx).UserID
				}
				pi, err := db.FindOne(ctx, key)
				pb := usermanagelogic.UserInfoToPb(ctx, pi, svcCtx)
				return pb, err
			},
			ExpireTime: 20 * time.Minute,
		})
		logx.Must(err)
		svcCtx.UserCache = userCache
	}
	{
		AreaCache, err := caches.NewCache(caches.CacheConfig[sys.AreaInfo, int64]{
			KeyType:   topics.ServerCacheKeySysAreaInfo,
			FastEvent: svcCtx.FastEvent,
			GetData: func(ctx context.Context, key int64) (*sys.AreaInfo, error) {
				db := relationDB.NewAreaInfoRepo(ctx)
				if key == 0 {
					key = ctxs.GetUserCtxNoNil(ctx).UserID
				}
				pi, err := db.FindOne(ctx, key, nil)
				pb := areamanagelogic.TransPoToPb(ctx, pi, svcCtx)
				return pb, err
			},
			ExpireTime: 20 * time.Minute,
		})
		logx.Must(err)
		svcCtx.AreaCache = AreaCache
	}
	{
		projectCache, err := caches.NewCache(caches.CacheConfig[sys.ProjectInfo, int64]{
			KeyType:   topics.ServerCacheKeySysProjectInfo,
			FastEvent: svcCtx.FastEvent,
			GetData: func(ctx context.Context, key int64) (*sys.ProjectInfo, error) {
				db := relationDB.NewProjectInfoRepo(ctx)
				if key == 0 {
					key = ctxs.GetUserCtxNoNil(ctx).ProjectID
				}
				pi, err := db.FindOne(ctx, key)
				pb := projectmanagelogic.ProjectInfoToPb(ctx, svcCtx, pi)
				return pb, err
			},
			ExpireTime: 20 * time.Minute,
		})
		logx.Must(err)
		svcCtx.ProjectCache = projectCache
	}

	{
		c, err := caches.NewCache(caches.CacheConfig[relationDB.SysApiInfo, string]{
			KeyType:   topics.ServerCacheKeySysAccessApi,
			FastEvent: svcCtx.FastEvent,
			GetData: func(ctx context.Context, key string) (*relationDB.SysApiInfo, error) {
				method, path, _ := strings.Cut(key, ":")
				db := relationDB.NewApiInfoRepo(ctx)
				pi, err := db.FindOneByFilter(ctx, relationDB.ApiInfoFilter{
					Route:      path,
					Method:     method,
					WithAccess: true,
				})
				return pi, err
			},
			ExpireTime: 20 * time.Minute,
		})
		logx.Must(err)
		svcCtx.ApiCache = c
	}

	{
		c, err := caches.NewCache(caches.CacheConfig[map[int64]struct{}, string]{
			KeyType:   topics.ServerCacheKeySysRoleAccess,
			FastEvent: svcCtx.FastEvent,
			GetData: func(ctx context.Context, key string) (*map[int64]struct{}, error) {
				db := relationDB.NewRoleAccessRepo(ctx)
				pi, err := db.FindByFilter(ctx, relationDB.RoleAccessFilter{
					AccessCodes: []string{key},
				}, nil)
				if err != nil {
					return nil, err
				}
				var ret = make(map[int64]struct{})
				for _, v := range pi {
					ret[v.RoleID] = struct{}{}
				}
				return &ret, err
			},
			ExpireTime: 20 * time.Minute,
		})
		logx.Must(err)
		svcCtx.RoleAccessCache = c
	}
	{
		userTokenInfo, err := cache.NewUserCache(svcCtx.FastEvent, svcCtx.TenantCache, svcCtx.UserCache)
		logx.Must(err)
		svcCtx.UsersCache = userTokenInfo
	}
}

func InitEventBus(svcCtx *svc.ServiceContext) {
	err := svcCtx.FastEvent.QueueSubscribe(topics.CoreSyncHalfHour, func(ctx context.Context, t time.Time, body []byte) error {
		return deptSync.NewDeptSync(ctxs.WithRoot(ctx), svcCtx).Timing()
	})
	logx.Must(err)
	err = svcCtx.FastEvent.QueueSubscribe(topics.CoreSyncDay, func(ctx context.Context, t time.Time, body []byte) error {
		return day.NewDaySync(ctxs.WithRoot(ctx), svcCtx).Handle()
	})
	logx.Must(err)
}

func TimerInit(svcCtx *svc.ServiceContext) {
	ctx := context.Background()
	_, err := svcCtx.TimedM.TaskInfoCreate(ctx, &timedmanage.TaskInfo{
		GroupCode: def.TimedUnitedRhinoQueueGroupCode,                                  //组编码
		Type:      1,                                                                   //任务类型 1 定时任务 2 延时任务
		Name:      "联犀中台半小时同步",                                                         // 任务名称
		Code:      "coreSyncHalfHour",                                                  //任务编码
		Params:    fmt.Sprintf(`{"topic":"%s","payload":""}`, topics.CoreSyncHalfHour), // 任务参数,延时任务如果没有传任务参数会拿数据库的参数来执行
		CronExpr:  "@every 30m",                                                        // cron执行表达式
		Status:    def.StatusWaitRun,                                                   // 状态
		Priority:  3,                                                                   //优先级: 10:critical 最高优先级  3: default 普通优先级 1:low 低优先级
	})
	if err != nil && !errors.Cmp(errors.Fmt(err), errors.Duplicate) {
		logx.Must(err)
	}
	{
		_, err := svcCtx.TimedM.TaskInfoCreate(ctx, &timedmanage.TaskInfo{
			GroupCode: def.TimedUnitedRhinoQueueGroupCode,                             //组编码
			Type:      1,                                                              //任务类型 1 定时任务 2 延时任务
			Name:      "联犀中台一天同步",                                                     // 任务名称
			Code:      "coreSyncDay",                                                  //任务编码
			Params:    fmt.Sprintf(`{"topic":"%s","payload":""}`, topics.CoreSyncDay), // 任务参数,延时任务如果没有传任务参数会拿数据库的参数来执行
			CronExpr:  "1 0 * * *",                                                    // cron执行表达式 0点01 分进行统计
			Status:    def.StatusWaitRun,                                              // 状态
			Priority:  3,                                                              //优先级: 10:critical 最高优先级  3: default 普通优先级 1:low 低优先级
		})
		if err != nil && !errors.Cmp(errors.Fmt(err), errors.Duplicate) {
			logx.Must(err)
		}
	}

}
