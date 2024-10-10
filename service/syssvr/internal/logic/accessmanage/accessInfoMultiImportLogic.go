package accessmanagelogic

import (
	"context"
	"encoding/json"
	"gitee.com/unitedrhino/core/service/syssvr/domain/access"
	"gitee.com/unitedrhino/core/service/syssvr/internal/repo/relationDB"
	"gitee.com/unitedrhino/core/service/syssvr/internal/svc"
	"gitee.com/unitedrhino/core/service/syssvr/pb/sys"
	"gitee.com/unitedrhino/share/def"
	"gitee.com/unitedrhino/share/errors"
	"github.com/zeromicro/go-zero/core/logx"
)

type AccessInfoMultiImportLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAccessInfoMultiImportLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AccessInfoMultiImportLogic {
	return &AccessInfoMultiImportLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *AccessInfoMultiImportLogic) AccessInfoMultiImport(in *sys.AccessInfoMultiImportReq) (*sys.AccessInfoMultiImportResp, error) {
	var ac access.Access
	err := json.Unmarshal([]byte(in.Access), &ac)
	if err != nil {
		return nil, errors.Parameter.AddMsg("json格式不对").AddDetail(err)
	}
	var (
		total          int64
		errCount       int64
		noCount        int64
		succCount      int64
		addAccessCodes []string
	)

	acDB := relationDB.NewAccessRepo(l.ctx)
	apiDB := relationDB.NewApiInfoRepo(l.ctx)
	for _, acc := range ac.Access {
		total += int64(len(acc.Apis))
		old, err := acDB.FindOneByFilter(l.ctx, relationDB.AccessFilter{Code: acc.Code})
		if err != nil && !errors.Cmp(err, errors.NotFind) {
			l.Errorf("find one by code(%d) err:%v", acc.Code, err)
			errCount += int64(len(acc.Apis))
			continue
		}
		if old == nil {
			old = &relationDB.SysAccessInfo{
				Name:       acc.Name,
				Module:     in.Module,
				Code:       acc.Code,
				Group:      acc.Group,
				IsNeedAuth: acc.IsNeedAuth,
				AuthType:   access.GetAuthType(acc.AuthType),
				Desc:       acc.Desc,
			}
			err = acDB.Insert(l.ctx, old)
			if err != nil {
				l.Errorf("insert access info failed, err:%v", err)
				errCount += int64(len(acc.Apis))
				continue
			}
		}
		addAccessCodes = append(addAccessCodes, acc.Code)
		for _, api := range acc.Apis {
			old, err := apiDB.FindOneByFilter(l.ctx, relationDB.ApiInfoFilter{
				Route:  api.Route,
				Method: api.Method,
			})
			if err != nil && !errors.Cmp(err, errors.NotFind) {
				l.Error("find one by code(%d) err:%v", api.Route, err)
				errCount++
				continue
			}
			if old == nil {
				err = apiDB.Insert(l.ctx, &relationDB.SysApiInfo{
					AccessCode:   acc.Code,
					Method:       api.Method,
					Route:        api.Route,
					Name:         api.Name,
					BusinessType: api.GetBusinessType(),
					Desc:         api.Desc,
				})
				if err != nil {
					errCount++
				} else {
					succCount++
				}
			} else {
				noCount++
			}
		}
	}
	if len(addAccessCodes) > 0 {
		var datas []*relationDB.SysTenantAccess
		for _, v := range addAccessCodes {
			datas = append(datas, &relationDB.SysTenantAccess{
				TenantCode: def.TenantCodeDefault,
				AccessCode: v,
			})
		}
		err = relationDB.NewTenantAccessRepo(l.ctx).MultiInsert(l.ctx, datas)
		if err != nil {
			l.Error(err)
		}
	}

	return &sys.AccessInfoMultiImportResp{
		Total:       total,
		ErrCount:    errCount,
		IgnoreCount: noCount,
		SuccCount:   succCount,
	}, nil
}
