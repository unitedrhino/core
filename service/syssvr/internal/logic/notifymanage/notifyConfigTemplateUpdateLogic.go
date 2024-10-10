package notifymanagelogic

import (
	"context"
	"gitee.com/unitedrhino/core/service/syssvr/internal/repo/relationDB"
	"gitee.com/unitedrhino/share/def"
	"gitee.com/unitedrhino/share/stores"
	"gitee.com/unitedrhino/share/utils"
	"gorm.io/gorm"

	"gitee.com/unitedrhino/core/service/syssvr/internal/svc"
	"gitee.com/unitedrhino/core/service/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type NotifyConfigTemplateUpdateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewNotifyConfigTemplateUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *NotifyConfigTemplateUpdateLogic {
	return &NotifyConfigTemplateUpdateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 租户通知配置
func (l *NotifyConfigTemplateUpdateLogic) NotifyConfigTemplateUpdate(in *sys.NotifyConfigTemplate) (*sys.Empty, error) {
	err := stores.GetTenantConn(l.ctx).Transaction(func(tx *gorm.DB) error {
		db := relationDB.NewNotifyConfigTemplateRepo(tx)
		err := db.DeleteByFilter(l.ctx, relationDB.NotifyConfigTemplateFilter{
			NotifyCode: in.NotifyCode,
			Type:       in.Type,
		})
		if err != nil {
			return err
		}
		po := relationDB.SysNotifyConfigTemplate{
			NotifyCode: in.NotifyCode,
			Type:       in.Type,
			TemplateID: in.TemplateID,
		}
		err = db.Save(l.ctx, &po)
		if err != nil {
			return err
		}
		err = InitConfigEnableTypes(l.ctx, tx, in.NotifyCode)
		return err
	})

	return &sys.Empty{}, err
}
func InitConfigEnableTypes(ctx context.Context, tx *gorm.DB, notifyCode string) error {
	tps, err := relationDB.NewNotifyConfigTemplateRepo(tx).FindByFilter(ctx, relationDB.NotifyConfigTemplateFilter{NotifyCode: notifyCode}, nil)
	if err != nil {
		return err
	}
	var enableTypes = []def.NotifyType{}
	for _, v := range tps {
		if v.Template == nil {
			continue
		}
		enableTypes = append(enableTypes, v.Template.Type)
	}
	return relationDB.NewNotifyConfigRepo(tx).UpdateWithField(ctx, relationDB.NotifyConfigFilter{Code: notifyCode}, map[string]any{
		"enable_types": utils.MarshalNoErr(enableTypes),
	})
}
