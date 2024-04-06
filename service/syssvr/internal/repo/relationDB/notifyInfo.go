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
1. 将NotifyInfo全局替换为模型的表名
2. 完善todo
*/

type NotifyInfoRepo struct {
	db *gorm.DB
}

func NewNotifyInfoRepo(in any) *NotifyInfoRepo {
	return &NotifyInfoRepo{db: stores.GetCommonConn(in)}
}

type NotifyInfoFilter struct {
	ID    int64
	Code  string
	Group string
	Name  string
}

func (p NotifyInfoRepo) fmtFilter(ctx context.Context, f NotifyInfoFilter) *gorm.DB {
	db := p.db.WithContext(ctx)
	if f.Code != "" {
		db = db.Where("code=?", f.Code)
	}
	if f.ID != 0 {
		db = db.Where("id=?", f.ID)
	}
	if f.Group != "" {
		db = db.Where("group=?", f.Group)
	}
	if f.Name != "" {
		db = db.Where("name like ?", "%"+f.Name+"%")
	}
	return db
}

func (p NotifyInfoRepo) Insert(ctx context.Context, data *SysNotifyInfo) error {
	result := p.db.WithContext(ctx).Create(data)
	return stores.ErrFmt(result.Error)
}

func (p NotifyInfoRepo) FindOneByFilter(ctx context.Context, f NotifyInfoFilter) (*SysNotifyInfo, error) {
	var result SysNotifyInfo
	db := p.fmtFilter(ctx, f)
	err := db.First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}
func (p NotifyInfoRepo) FindByFilter(ctx context.Context, f NotifyInfoFilter, page *def.PageInfo) ([]*SysNotifyInfo, error) {
	var results []*SysNotifyInfo
	db := p.fmtFilter(ctx, f).Model(&SysNotifyInfo{})
	db = page.ToGorm(db)
	err := db.Find(&results).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return results, nil
}

func (p NotifyInfoRepo) CountByFilter(ctx context.Context, f NotifyInfoFilter) (size int64, err error) {
	db := p.fmtFilter(ctx, f).Model(&SysNotifyInfo{})
	err = db.Count(&size).Error
	return size, stores.ErrFmt(err)
}

func (p NotifyInfoRepo) Update(ctx context.Context, data *SysNotifyInfo) error {
	err := p.db.WithContext(ctx).Where("id = ?", data.ID).Save(data).Error
	return stores.ErrFmt(err)
}

func (p NotifyInfoRepo) DeleteByFilter(ctx context.Context, f NotifyInfoFilter) error {
	db := p.fmtFilter(ctx, f)
	err := db.Delete(&SysNotifyInfo{}).Error
	return stores.ErrFmt(err)
}

func (p NotifyInfoRepo) Delete(ctx context.Context, id int64) error {
	err := p.db.WithContext(ctx).Where("id = ?", id).Delete(&SysNotifyInfo{}).Error
	return stores.ErrFmt(err)
}
func (p NotifyInfoRepo) FindOne(ctx context.Context, id int64) (*SysNotifyInfo, error) {
	var result SysNotifyInfo
	err := p.db.WithContext(ctx).Where("id = ?", id).First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}

// 批量插入 LightStrategyDevice 记录
func (p NotifyInfoRepo) MultiInsert(ctx context.Context, data []*SysNotifyInfo) error {
	err := p.db.WithContext(ctx).Clauses(clause.OnConflict{UpdateAll: true}).Model(&SysNotifyInfo{}).Create(data).Error
	return stores.ErrFmt(err)
}
