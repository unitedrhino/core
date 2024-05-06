package relationDB

import (
	"context"
	"fmt"
	"gitee.com/i-Things/share/def"
	"gitee.com/i-Things/share/stores"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

/*
这个是参考样例
使用教程:
1. 将TenantNotifyConfig全局替换为模型的表名
2. 完善todo
*/

type TenantNotifyTemplateRepo struct {
	db *gorm.DB
}

func NewTenantNotifyTemplateRepo(in any) *TenantNotifyTemplateRepo {
	return &TenantNotifyTemplateRepo{db: stores.GetCommonConn(in)}
}

type TenantNotifyTemplateFilter struct {
	NotifyCode string
	Type       string
}

func (p TenantNotifyTemplateRepo) fmtFilter(ctx context.Context, f TenantNotifyTemplateFilter) *gorm.DB {
	db := p.db.WithContext(ctx)
	if f.NotifyCode != "" {
		db = db.Where("notify_code=?", f.NotifyCode)
	}
	if f.Type != "" {
		db = db.Where(fmt.Sprintf("%v=?", stores.Col("type")), f.Type)
	}
	return db
}

func (p TenantNotifyTemplateRepo) Insert(ctx context.Context, data *SysTenantNotifyTemplate) error {
	result := p.db.WithContext(ctx).Create(data)
	return stores.ErrFmt(result.Error)
}

func (p TenantNotifyTemplateRepo) FindOneByFilter(ctx context.Context, f TenantNotifyTemplateFilter) (*SysTenantNotifyTemplate, error) {
	var result SysTenantNotifyTemplate
	db := p.fmtFilter(ctx, f).Preload("Template").Preload("Template.Channel").Preload("Config")
	err := db.First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}
func (p TenantNotifyTemplateRepo) FindByFilter(ctx context.Context, f TenantNotifyTemplateFilter, page *def.PageInfo) ([]*SysTenantNotifyTemplate, error) {
	var results []*SysTenantNotifyTemplate
	db := p.fmtFilter(ctx, f).Model(&SysTenantNotifyTemplate{})
	db = page.ToGorm(db).Preload("Template").Preload("Config")
	err := db.Find(&results).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return results, nil
}

func (p TenantNotifyTemplateRepo) CountByFilter(ctx context.Context, f TenantNotifyTemplateFilter) (size int64, err error) {
	db := p.fmtFilter(ctx, f).Model(&SysTenantNotifyTemplate{})
	err = db.Count(&size).Error
	return size, stores.ErrFmt(err)
}

func (p TenantNotifyTemplateRepo) Update(ctx context.Context, data *SysTenantNotifyTemplate) error {
	err := p.db.WithContext(ctx).Where("id = ?", data.ID).Save(data).Error
	return stores.ErrFmt(err)
}

func (p TenantNotifyTemplateRepo) DeleteByFilter(ctx context.Context, f TenantNotifyTemplateFilter) error {
	db := p.fmtFilter(ctx, f)
	err := db.Delete(&SysTenantNotifyTemplate{}).Error
	return stores.ErrFmt(err)
}

func (p TenantNotifyTemplateRepo) Delete(ctx context.Context, id int64) error {
	err := p.db.WithContext(ctx).Where("id = ?", id).Delete(&SysTenantNotifyTemplate{}).Error
	return stores.ErrFmt(err)
}
func (p TenantNotifyTemplateRepo) FindOne(ctx context.Context, id int64) (*SysTenantNotifyTemplate, error) {
	var result SysTenantNotifyTemplate
	err := p.db.WithContext(ctx).Where("id = ?", id).First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}

// 批量插入 LightStrategyDevice 记录
func (p TenantNotifyTemplateRepo) MultiInsert(ctx context.Context, data []*SysTenantNotifyTemplate) error {
	err := p.db.WithContext(ctx).Clauses(clause.OnConflict{UpdateAll: true}).Model(&SysTenantNotifyTemplate{}).Create(data).Error
	return stores.ErrFmt(err)
}

func (p TenantNotifyTemplateRepo) MultiUpdate(ctx context.Context, appCode string, pos []*SysTenantNotifyTemplate) error {
	for i := range pos {
		pos[i].ID = 0
	}
	err := p.db.Transaction(func(tx *gorm.DB) error {
		rm := NewTenantNotifyTemplateRepo(tx)
		err := rm.DeleteByFilter(ctx, TenantNotifyTemplateFilter{})
		if err != nil {
			return err
		}
		if len(pos) != 0 {
			err = rm.MultiInsert(ctx, pos)
			if err != nil {
				return err
			}
		}
		return nil
	})
	return stores.ErrFmt(err)
}
