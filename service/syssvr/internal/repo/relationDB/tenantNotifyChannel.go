package relationDB

import (
	"context"
	"gitee.com/i-Things/share/def"
	"gitee.com/i-Things/share/stores"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

/*
这个是参考样例
使用教程:
1. 将TenantNotifyChannel全局替换为模型的表名
2. 完善todo
*/

type TenantNotifyChannelRepo struct {
	db *gorm.DB
}

func NewTenantNotifyChannelRepo(in any) *TenantNotifyChannelRepo {
	return &TenantNotifyChannelRepo{db: stores.GetCommonConn(in)}
}

type TenantNotifyChannelFilter struct {
	Name string
	Type string
}

func (p TenantNotifyChannelRepo) fmtFilter(ctx context.Context, f TenantNotifyChannelFilter) *gorm.DB {
	db := p.db.WithContext(ctx)
	if f.Name != "" {
		db = db.Where("name =?", f.Name)
	}
	if f.Type != "" {
		db = db.Where("type =?", f.Type)
	}
	return db
}

func (p TenantNotifyChannelRepo) Insert(ctx context.Context, data *SysTenantNotifyChannel) error {
	result := p.db.WithContext(ctx).Create(data)
	return stores.ErrFmt(result.Error)
}

func (p TenantNotifyChannelRepo) FindOneByFilter(ctx context.Context, f TenantNotifyChannelFilter) (*SysTenantNotifyChannel, error) {
	var result SysTenantNotifyChannel
	db := p.fmtFilter(ctx, f)
	err := db.First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}
func (p TenantNotifyChannelRepo) FindByFilter(ctx context.Context, f TenantNotifyChannelFilter, page *def.PageInfo) ([]*SysTenantNotifyChannel, error) {
	var results []*SysTenantNotifyChannel
	db := p.fmtFilter(ctx, f).Model(&SysTenantNotifyChannel{})
	db = page.ToGorm(db)
	err := db.Find(&results).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return results, nil
}

func (p TenantNotifyChannelRepo) CountByFilter(ctx context.Context, f TenantNotifyChannelFilter) (size int64, err error) {
	db := p.fmtFilter(ctx, f).Model(&SysTenantNotifyChannel{})
	err = db.Count(&size).Error
	return size, stores.ErrFmt(err)
}

func (p TenantNotifyChannelRepo) Update(ctx context.Context, data *SysTenantNotifyChannel) error {
	err := p.db.WithContext(ctx).Where("id = ?", data.ID).Save(data).Error
	return stores.ErrFmt(err)
}

func (p TenantNotifyChannelRepo) DeleteByFilter(ctx context.Context, f TenantNotifyChannelFilter) error {
	db := p.fmtFilter(ctx, f)
	err := db.Delete(&SysTenantNotifyChannel{}).Error
	return stores.ErrFmt(err)
}

func (p TenantNotifyChannelRepo) Delete(ctx context.Context, id int64) error {
	err := p.db.WithContext(ctx).Where("id = ?", id).Delete(&SysTenantNotifyChannel{}).Error
	return stores.ErrFmt(err)
}
func (p TenantNotifyChannelRepo) FindOne(ctx context.Context, id int64) (*SysTenantNotifyChannel, error) {
	var result SysTenantNotifyChannel
	err := p.db.WithContext(ctx).Where("id = ?", id).First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}

// 批量插入 LightStrategyDevice 记录
func (p TenantNotifyChannelRepo) MultiInsert(ctx context.Context, data []*SysTenantNotifyChannel) error {
	err := p.db.WithContext(ctx).Clauses(clause.OnConflict{UpdateAll: true}).Model(&SysTenantNotifyChannel{}).Create(data).Error
	return stores.ErrFmt(err)
}
