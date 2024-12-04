package relationDB

import (
	"context"
	"fmt"
	"gitee.com/unitedrhino/share/stores"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

/*
这个是参考样例
使用教程:
1. 将NotifyInfo全局替换为模型的表名
2. 完善todo
*/

type NotifyConfigRepo struct {
	db *gorm.DB
}

func NewNotifyConfigRepo(in any) *NotifyConfigRepo {
	return &NotifyConfigRepo{db: stores.GetCommonConn(in)}
}

type NotifyConfigFilter struct {
	ID            int64
	Code          string
	Group         string
	Name          string
	WithTemplates bool
}

func (p NotifyConfigRepo) fmtFilter(ctx context.Context, f NotifyConfigFilter) *gorm.DB {
	db := p.db.WithContext(ctx)
	if f.Code != "" {
		db = db.Where("code=?", f.Code)
	}
	if f.ID != 0 {
		db = db.Where("id=?", f.ID)
	}
	if f.Group != "" {
		db = db.Where(fmt.Sprintf("%s=?", stores.Col("group")), f.Group)
	}
	if f.Name != "" {
		db = db.Where("name like ?", "%"+f.Name+"%")
	}
	if f.WithTemplates {
		db = db.Preload("Templates")
	}
	return db
}

func (p NotifyConfigRepo) Insert(ctx context.Context, data *SysNotifyConfig) error {
	result := p.db.WithContext(ctx).Create(data)
	return stores.ErrFmt(result.Error)
}

func (p NotifyConfigRepo) FindOneByFilter(ctx context.Context, f NotifyConfigFilter) (*SysNotifyConfig, error) {
	var result SysNotifyConfig
	db := p.fmtFilter(ctx, f)
	err := db.First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}
func (p NotifyConfigRepo) FindByFilter(ctx context.Context, f NotifyConfigFilter, page *stores.PageInfo) ([]*SysNotifyConfig, error) {
	var results []*SysNotifyConfig
	db := p.fmtFilter(ctx, f).Model(&SysNotifyConfig{})
	db = page.ToGorm(db)
	err := db.Find(&results).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return results, nil
}

func (p NotifyConfigRepo) CountByFilter(ctx context.Context, f NotifyConfigFilter) (size int64, err error) {
	db := p.fmtFilter(ctx, f).Model(&SysNotifyConfig{})
	err = db.Count(&size).Error
	return size, stores.ErrFmt(err)
}

func (p NotifyConfigRepo) Update(ctx context.Context, data *SysNotifyConfig) error {
	err := p.db.WithContext(ctx).Where("id = ?", data.ID).Save(data).Error
	return stores.ErrFmt(err)
}

func (d NotifyConfigRepo) UpdateWithField(ctx context.Context, f NotifyConfigFilter, updates map[string]any) error {
	db := d.fmtFilter(ctx, f)
	err := db.Model(&SysNotifyConfig{}).Updates(updates).Error
	return stores.ErrFmt(err)
}

func (p NotifyConfigRepo) DeleteByFilter(ctx context.Context, f NotifyConfigFilter) error {
	db := p.fmtFilter(ctx, f)
	err := db.Delete(&SysNotifyConfig{}).Error
	return stores.ErrFmt(err)
}

func (p NotifyConfigRepo) Delete(ctx context.Context, id int64) error {
	err := p.db.WithContext(ctx).Where("id = ?", id).Delete(&SysNotifyConfig{}).Error
	return stores.ErrFmt(err)
}
func (p NotifyConfigRepo) FindOne(ctx context.Context, id int64) (*SysNotifyConfig, error) {
	var result SysNotifyConfig
	err := p.db.WithContext(ctx).Where("id = ?", id).First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}

// 批量插入 LightStrategyDevice 记录
func (p NotifyConfigRepo) MultiInsert(ctx context.Context, data []*SysNotifyConfig) error {
	err := p.db.WithContext(ctx).Clauses(clause.OnConflict{UpdateAll: true}).Model(&SysNotifyConfig{}).Create(data).Error
	return stores.ErrFmt(err)
}
