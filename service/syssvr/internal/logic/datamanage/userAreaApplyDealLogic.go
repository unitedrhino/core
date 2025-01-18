package datamanagelogic

import (
	"context"
	"gitee.com/unitedrhino/core/service/syssvr/internal/repo/cache"
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

type UserAreaApplyDealLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUserAreaApplyDealLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserAreaApplyDealLogic {
	return &UserAreaApplyDealLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UserAreaApplyDealLogic) UserAreaApplyDeal(in *sys.UserAreaApplyDealReq) (*sys.Empty, error) {
	if !in.IsApprove {
		err := relationDB.NewUserAreaApplyRepo(l.ctx).DeleteByFilter(l.ctx, relationDB.UserAreaApplyFilter{IDs: in.Ids})
		return &sys.Empty{}, err
	}
	uc := ctxs.GetUserCtx(l.ctx)
	db := stores.GetTenantConn(l.ctx)
	var authSet = map[int64]struct{}{}
	err := db.Transaction(func(tx *gorm.DB) error {
		uaa := relationDB.NewUserAreaApplyRepo(tx)
		ua := relationDB.NewDataAreaRepo(tx)
		dp := relationDB.NewDataProjectRepo(tx)
		uaas, err := uaa.FindByFilter(l.ctx, relationDB.UserAreaApplyFilter{IDs: in.Ids}, nil)
		if err != nil {
			return err
		}
		if len(uaas) == 0 {
			return errors.Parameter.AddMsgf("未查询到授权的id")
		}
		var uas []*relationDB.SysDataArea
		var authUserIDs []int64
		for _, v := range uaas {
			uas = append(uas, &relationDB.SysDataArea{
				TargetType: def.TargetUser,
				TargetID:   v.UserID,
				ProjectID:  v.ProjectID,
				AreaID:     int64(v.AreaID),
				AuthType:   v.AuthType,
			})
			authUserIDs = append(authUserIDs, v.UserID)
		}
		authd, err := dp.FindByFilter(l.ctx, relationDB.DataProjectFilter{TargetType: def.TargetUser, TargetIDs: authUserIDs}, nil)
		if err != nil {
			return err
		}
		for _, v := range authd {
			authSet[v.TargetID] = struct{}{}
		}
		var needAuthUser []*relationDB.SysDataProject
		for _, v := range authUserIDs {
			if _, ok := authSet[v]; !ok {
				needAuthUser = append(needAuthUser, &relationDB.SysDataProject{
					ProjectID:  uc.ProjectID,
					TargetType: def.TargetUser,
					TargetID:   v,
					AuthType:   def.AuthRead,
				})
			}
		}
		if len(needAuthUser) != 0 {
			err = dp.MultiInsert(l.ctx, needAuthUser)
			if err != nil {
				return err
			}
		}
		err = ua.MultiInsert(l.ctx, uas)
		if err != nil {
			return err
		}
		err = uaa.DeleteByFilter(l.ctx, relationDB.UserAreaApplyFilter{IDs: in.Ids})
		if err != nil {
			return err
		}
		return nil
	})
	for userID := range authSet {
		cache.ClearProjectAuth(userID)
	}
	ProjectUserCount(l.ctx, uc.ProjectID)
	return &sys.Empty{}, err
}
