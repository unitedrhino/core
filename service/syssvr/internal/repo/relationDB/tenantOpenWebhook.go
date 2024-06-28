package relationDB

import (
	"context"
	"gitee.com/i-Things/share/stores"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

/*
这个是参考样例
使用教程:
1. 将TenantOpenWebhook全局替换为模型的表名
2. 完善todo
*/

type TenantOpenWebhookRepo struct {
	db *gorm.DB
}

func NewTenantOpenWebhookRepo(in any) *TenantOpenWebhookRepo {
	return &TenantOpenWebhookRepo{db: stores.GetCommonConn(in)}
}

type TenantOpenWebhookFilter struct {
	Code string
}

func (p TenantOpenWebhookRepo) fmtFilter(ctx context.Context, f TenantOpenWebhookFilter) *gorm.DB {
	db := p.db.WithContext(ctx)
	if f.Code != "" {
		db = db.Where("code = ?", f.Code)
	}
	return db
}

func (p TenantOpenWebhookRepo) Insert(ctx context.Context, data *SysTenantOpenWebhook) error {
	result := p.db.WithContext(ctx).Create(data)
	return stores.ErrFmt(result.Error)
}

func (p TenantOpenWebhookRepo) FindOneByFilter(ctx context.Context, f TenantOpenWebhookFilter) (*SysTenantOpenWebhook, error) {
	var result SysTenantOpenWebhook
	db := p.fmtFilter(ctx, f)
	err := db.First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}
func (p TenantOpenWebhookRepo) FindByFilter(ctx context.Context, f TenantOpenWebhookFilter, page *stores.PageInfo) ([]*SysTenantOpenWebhook, error) {
	var results []*SysTenantOpenWebhook
	db := p.fmtFilter(ctx, f).Model(&SysTenantOpenWebhook{})
	db = page.ToGorm(db)
	err := db.Find(&results).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return results, nil
}

func (p TenantOpenWebhookRepo) CountByFilter(ctx context.Context, f TenantOpenWebhookFilter) (size int64, err error) {
	db := p.fmtFilter(ctx, f).Model(&SysTenantOpenWebhook{})
	err = db.Count(&size).Error
	return size, stores.ErrFmt(err)
}

func (p TenantOpenWebhookRepo) Update(ctx context.Context, data *SysTenantOpenWebhook) error {
	err := p.db.WithContext(ctx).Where("id = ?", data.ID).Save(data).Error
	return stores.ErrFmt(err)
}

func (p TenantOpenWebhookRepo) DeleteByFilter(ctx context.Context, f TenantOpenWebhookFilter) error {
	db := p.fmtFilter(ctx, f)
	err := db.Delete(&SysTenantOpenWebhook{}).Error
	return stores.ErrFmt(err)
}

func (p TenantOpenWebhookRepo) Delete(ctx context.Context, id int64) error {
	err := p.db.WithContext(ctx).Where("id = ?", id).Delete(&SysTenantOpenWebhook{}).Error
	return stores.ErrFmt(err)
}
func (p TenantOpenWebhookRepo) FindOne(ctx context.Context, id int64) (*SysTenantOpenWebhook, error) {
	var result SysTenantOpenWebhook
	err := p.db.WithContext(ctx).Where("id = ?", id).First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}

// 批量插入 LightStrategyDevice 记录
func (p TenantOpenWebhookRepo) MultiInsert(ctx context.Context, data []*SysTenantOpenWebhook) error {
	err := p.db.WithContext(ctx).Clauses(clause.OnConflict{UpdateAll: true}).Model(&SysTenantOpenWebhook{}).Create(data).Error
	return stores.ErrFmt(err)
}
