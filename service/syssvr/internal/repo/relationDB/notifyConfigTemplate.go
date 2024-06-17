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

type NotifyConfigTemplateRepo struct {
	db *gorm.DB
}

func NewNotifyConfigTemplateRepo(in any) *NotifyConfigTemplateRepo {
	return &NotifyConfigTemplateRepo{db: stores.GetCommonConn(in)}
}

type NotifyConfigTemplateFilter struct {
	NotifyCode string
	Type       string
}

func (p NotifyConfigTemplateRepo) fmtFilter(ctx context.Context, f NotifyConfigTemplateFilter) *gorm.DB {
	db := p.db.WithContext(ctx)
	if f.NotifyCode != "" {
		db = db.Where("notify_code=?", f.NotifyCode)
	}
	if f.Type != "" {
		db = db.Where(fmt.Sprintf("%v=?", stores.Col("type")), f.Type)
	}
	return db
}

func (p NotifyConfigTemplateRepo) Insert(ctx context.Context, data *SysNotifyConfigTemplate) error {
	result := p.db.WithContext(ctx).Create(data)
	return stores.ErrFmt(result.Error)
}

func (p NotifyConfigTemplateRepo) FindOneByFilter(ctx context.Context, f NotifyConfigTemplateFilter) (*SysNotifyConfigTemplate, error) {
	var result SysNotifyConfigTemplate
	db := p.fmtFilter(ctx, f).Preload("Template").Preload("Template.Channel").Preload("Config")
	err := db.First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}
func (p NotifyConfigTemplateRepo) FindByFilter(ctx context.Context, f NotifyConfigTemplateFilter, page *def.PageInfo) ([]*SysNotifyConfigTemplate, error) {
	var results []*SysNotifyConfigTemplate
	db := p.fmtFilter(ctx, f).Model(&SysNotifyConfigTemplate{})
	db = page.ToGorm(db).Preload("Template").Preload("Config")
	err := db.Find(&results).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return results, nil
}

func (p NotifyConfigTemplateRepo) CountByFilter(ctx context.Context, f NotifyConfigTemplateFilter) (size int64, err error) {
	db := p.fmtFilter(ctx, f).Model(&SysNotifyConfigTemplate{})
	err = db.Count(&size).Error
	return size, stores.ErrFmt(err)
}

func (p NotifyConfigTemplateRepo) Update(ctx context.Context, data *SysNotifyConfigTemplate) error {
	err := p.db.WithContext(ctx).Where("id = ?", data.ID).Save(data).Error
	return stores.ErrFmt(err)
}

func (p NotifyConfigTemplateRepo) DeleteByFilter(ctx context.Context, f NotifyConfigTemplateFilter) error {
	db := p.fmtFilter(ctx, f)
	err := db.Delete(&SysNotifyConfigTemplate{}).Error
	return stores.ErrFmt(err)
}

func (p NotifyConfigTemplateRepo) Delete(ctx context.Context, id int64) error {
	err := p.db.WithContext(ctx).Where("id = ?", id).Delete(&SysNotifyConfigTemplate{}).Error
	return stores.ErrFmt(err)
}
func (p NotifyConfigTemplateRepo) FindOne(ctx context.Context, id int64) (*SysNotifyConfigTemplate, error) {
	var result SysNotifyConfigTemplate
	err := p.db.WithContext(ctx).Where("id = ?", id).First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}

// 批量插入 LightStrategyDevice 记录
func (p NotifyConfigTemplateRepo) MultiInsert(ctx context.Context, data []*SysNotifyConfigTemplate) error {
	err := p.db.WithContext(ctx).Clauses(clause.OnConflict{UpdateAll: true}).Model(&SysNotifyConfigTemplate{}).Create(data).Error
	return stores.ErrFmt(err)
}

func (p NotifyConfigTemplateRepo) MultiUpdate(ctx context.Context, pos []*SysNotifyConfigTemplate) error {
	for i := range pos {
		pos[i].ID = 0
	}
	err := p.db.Transaction(func(tx *gorm.DB) error {
		rm := NewNotifyConfigTemplateRepo(tx)
		err := rm.DeleteByFilter(ctx, NotifyConfigTemplateFilter{})
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

// 批量插入 LightStrategyDevice 记录
func (p NotifyConfigTemplateRepo) Save(ctx context.Context, data *SysNotifyConfigTemplate) error {
	err := p.db.WithContext(ctx).Clauses(clause.OnConflict{UpdateAll: true}).Model(&SysNotifyConfigTemplate{}).Create(data).Error
	return stores.ErrFmt(err)
}
