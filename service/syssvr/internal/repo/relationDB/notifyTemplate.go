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
1. 将NotifyTemplate全局替换为模型的表名
2. 完善todo
*/

type NotifyTemplateRepo struct {
	db *gorm.DB
}

func NewNotifyTemplateRepo(in any) *NotifyTemplateRepo {
	return &NotifyTemplateRepo{db: stores.GetCommonConn(in)}
}

type NotifyTemplateFilter struct {
	Name       string
	NotifyCode string
	Type       string
}

func (p NotifyTemplateRepo) fmtFilter(ctx context.Context, f NotifyTemplateFilter) *gorm.DB {
	db := p.db.WithContext(ctx)
	if f.Name != "" {
		db = db.Where("name =?", f.Name)
	}
	if f.NotifyCode != "" {
		db = db.Where("notify_code=?", f.NotifyCode)
	}
	if f.Type != "" {
		db = db.Where("type =?", f.Type)
	}
	return db
}

func (p NotifyTemplateRepo) Insert(ctx context.Context, data *SysNotifyTemplate) error {
	result := p.db.WithContext(ctx).Create(data)
	return stores.ErrFmt(result.Error)
}

func (p NotifyTemplateRepo) FindOneByFilter(ctx context.Context, f NotifyTemplateFilter) (*SysNotifyTemplate, error) {
	var result SysNotifyTemplate
	db := p.fmtFilter(ctx, f)
	err := db.First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}
func (p NotifyTemplateRepo) FindByFilter(ctx context.Context, f NotifyTemplateFilter, page *stores.PageInfo) ([]*SysNotifyTemplate, error) {
	var results []*SysNotifyTemplate
	db := p.fmtFilter(ctx, f).Model(&SysNotifyTemplate{})
	db = page.ToGorm(db)
	err := db.Preload("Channel").Find(&results).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return results, nil
}

func (p NotifyTemplateRepo) CountByFilter(ctx context.Context, f NotifyTemplateFilter) (size int64, err error) {
	db := p.fmtFilter(ctx, f).Model(&SysNotifyTemplate{})
	err = db.Count(&size).Error
	return size, stores.ErrFmt(err)
}

func (p NotifyTemplateRepo) Update(ctx context.Context, data *SysNotifyTemplate) error {
	err := p.db.WithContext(ctx).Where("id = ?", data.ID).Save(data).Error
	return stores.ErrFmt(err)
}

func (p NotifyTemplateRepo) DeleteByFilter(ctx context.Context, f NotifyTemplateFilter) error {
	db := p.fmtFilter(ctx, f)
	err := db.Delete(&SysNotifyTemplate{}).Error
	return stores.ErrFmt(err)
}

func (p NotifyTemplateRepo) Delete(ctx context.Context, id int64) error {
	err := p.db.WithContext(ctx).Where("id = ?", id).Delete(&SysNotifyTemplate{}).Error
	return stores.ErrFmt(err)
}

func (p NotifyTemplateRepo) FindOne(ctx context.Context, id int64) (*SysNotifyTemplate, error) {
	var result SysNotifyTemplate
	err := p.db.WithContext(ctx).Preload("Channel").Where("id = ?", id).First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}

// 批量插入 LightStrategyDevice 记录
func (p NotifyTemplateRepo) MultiInsert(ctx context.Context, data []*SysNotifyTemplate) error {
	err := p.db.WithContext(ctx).Clauses(clause.OnConflict{UpdateAll: true}).Model(&SysNotifyTemplate{}).Create(data).Error
	return stores.ErrFmt(err)
}

// 批量插入 LightStrategyDevice 记录
func (p NotifyTemplateRepo) Save(ctx context.Context, data *SysNotifyTemplate) error {
	err := p.db.WithContext(ctx).Clauses(clause.OnConflict{UpdateAll: true}).Model(&SysNotifyTemplate{}).Create(data).Error
	return stores.ErrFmt(err)
}
