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

type TenantNotifyRepo struct {
	db *gorm.DB
}

func NewTenantNotifyRepo(in any) *TenantNotifyRepo {
	return &TenantNotifyRepo{db: stores.GetCommonConn(in)}
}

type TenantNotifyConfigFilter struct {
	NotifyCode string
	Type       string
}

func (p TenantNotifyRepo) fmtFilter(ctx context.Context, f TenantNotifyConfigFilter) *gorm.DB {
	db := p.db.WithContext(ctx)
	if f.NotifyCode != "" {
		db = db.Where("notify_code=?", f.NotifyCode)
	}
	if f.Type != "" {
		db = db.Where(fmt.Sprintf("%v=?", stores.Col("type")), f.Type)
	}
	return db
}

func (p TenantNotifyRepo) Insert(ctx context.Context, data *SysTenantNotify) error {
	result := p.db.WithContext(ctx).Create(data)
	return stores.ErrFmt(result.Error)
}

func (p TenantNotifyRepo) FindOneByFilter(ctx context.Context, f TenantNotifyConfigFilter) (*SysTenantNotify, error) {
	var result SysTenantNotify
	db := p.fmtFilter(ctx, f).Preload("Template").Preload("Config")
	err := db.First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}
func (p TenantNotifyRepo) FindByFilter(ctx context.Context, f TenantNotifyConfigFilter, page *def.PageInfo) ([]*SysTenantNotify, error) {
	var results []*SysTenantNotify
	db := p.fmtFilter(ctx, f).Model(&SysTenantNotify{})
	db = page.ToGorm(db).Preload("Template").Preload("Config")
	err := db.Find(&results).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return results, nil
}

func (p TenantNotifyRepo) CountByFilter(ctx context.Context, f TenantNotifyConfigFilter) (size int64, err error) {
	db := p.fmtFilter(ctx, f).Model(&SysTenantNotify{})
	err = db.Count(&size).Error
	return size, stores.ErrFmt(err)
}

func (p TenantNotifyRepo) Update(ctx context.Context, data *SysTenantNotify) error {
	err := p.db.WithContext(ctx).Where("id = ?", data.ID).Save(data).Error
	return stores.ErrFmt(err)
}

func (p TenantNotifyRepo) DeleteByFilter(ctx context.Context, f TenantNotifyConfigFilter) error {
	db := p.fmtFilter(ctx, f)
	err := db.Delete(&SysTenantNotify{}).Error
	return stores.ErrFmt(err)
}

func (p TenantNotifyRepo) Delete(ctx context.Context, id int64) error {
	err := p.db.WithContext(ctx).Where("id = ?", id).Delete(&SysTenantNotify{}).Error
	return stores.ErrFmt(err)
}
func (p TenantNotifyRepo) FindOne(ctx context.Context, id int64) (*SysTenantNotify, error) {
	var result SysTenantNotify
	err := p.db.WithContext(ctx).Where("id = ?", id).First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}

// 批量插入 LightStrategyDevice 记录
func (p TenantNotifyRepo) MultiInsert(ctx context.Context, data []*SysTenantNotify) error {
	err := p.db.WithContext(ctx).Clauses(clause.OnConflict{UpdateAll: true}).Model(&SysTenantNotify{}).Create(data).Error
	return stores.ErrFmt(err)
}

func (p TenantNotifyRepo) MultiUpdate(ctx context.Context, appCode string, pos []*SysTenantNotify) error {
	for i := range pos {
		pos[i].ID = 0
	}
	err := p.db.Transaction(func(tx *gorm.DB) error {
		rm := NewTenantNotifyRepo(tx)
		err := rm.DeleteByFilter(ctx, TenantNotifyConfigFilter{})
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
