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

type NotifyChannelRepo struct {
	db *gorm.DB
}

func NewNotifyChannelRepo(in any) *NotifyChannelRepo {
	return &NotifyChannelRepo{db: stores.GetCommonConn(in)}
}

type NotifyChannelFilter struct {
	Name string
	Type string
}

func (p NotifyChannelRepo) fmtFilter(ctx context.Context, f NotifyChannelFilter) *gorm.DB {
	db := p.db.WithContext(ctx)
	if f.Name != "" {
		db = db.Where("name =?", f.Name)
	}
	if f.Type != "" {
		db = db.Where("type =?", f.Type)
	}
	return db
}

func (p NotifyChannelRepo) Insert(ctx context.Context, data *SysNotifyChannel) error {
	result := p.db.WithContext(ctx).Create(data)
	return stores.ErrFmt(result.Error)
}

func (p NotifyChannelRepo) FindOneByFilter(ctx context.Context, f NotifyChannelFilter) (*SysNotifyChannel, error) {
	var result SysNotifyChannel
	db := p.fmtFilter(ctx, f)
	err := db.First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}
func (p NotifyChannelRepo) FindByFilter(ctx context.Context, f NotifyChannelFilter, page *def.PageInfo) ([]*SysNotifyChannel, error) {
	var results []*SysNotifyChannel
	db := p.fmtFilter(ctx, f).Model(&SysNotifyChannel{})
	db = page.ToGorm(db)
	err := db.Find(&results).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return results, nil
}

func (p NotifyChannelRepo) CountByFilter(ctx context.Context, f NotifyChannelFilter) (size int64, err error) {
	db := p.fmtFilter(ctx, f).Model(&SysNotifyChannel{})
	err = db.Count(&size).Error
	return size, stores.ErrFmt(err)
}

func (p NotifyChannelRepo) Update(ctx context.Context, data *SysNotifyChannel) error {
	err := p.db.WithContext(ctx).Where("id = ?", data.ID).Save(data).Error
	return stores.ErrFmt(err)
}

func (p NotifyChannelRepo) DeleteByFilter(ctx context.Context, f NotifyChannelFilter) error {
	db := p.fmtFilter(ctx, f)
	err := db.Delete(&SysNotifyChannel{}).Error
	return stores.ErrFmt(err)
}

func (p NotifyChannelRepo) Delete(ctx context.Context, id int64) error {
	err := p.db.WithContext(ctx).Where("id = ?", id).Delete(&SysNotifyChannel{}).Error
	return stores.ErrFmt(err)
}
func (p NotifyChannelRepo) FindOne(ctx context.Context, id int64) (*SysNotifyChannel, error) {
	var result SysNotifyChannel
	err := p.db.WithContext(ctx).Where("id = ?", id).First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}

// 批量插入 LightStrategyDevice 记录
func (p NotifyChannelRepo) MultiInsert(ctx context.Context, data []*SysNotifyChannel) error {
	err := p.db.WithContext(ctx).Clauses(clause.OnConflict{UpdateAll: true}).Model(&SysNotifyChannel{}).Create(data).Error
	return stores.ErrFmt(err)
}
