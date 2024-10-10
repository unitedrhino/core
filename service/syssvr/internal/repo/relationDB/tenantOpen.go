package relationDB

import (
	"context"
	"gitee.com/unitedrhino/share/stores"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

/*
这个是参考样例
使用教程:
1. 将TenantOpen全局替换为模型的表名
2. 完善todo
*/

type DataOpenAccessRepo struct {
	db *gorm.DB
}

func NewDataOpenAccessRepo(in any) *DataOpenAccessRepo {
	return &DataOpenAccessRepo{db: stores.GetCommonConn(in)}
}

type DataOpenAccessFilter struct {
	TenantCode string
	UserID     int64
	Code       string
}

func (p DataOpenAccessRepo) fmtFilter(ctx context.Context, f DataOpenAccessFilter) *gorm.DB {
	db := p.db.WithContext(ctx)
	if f.TenantCode != "" {
		db = db.Where("tenant_code = ?", f.TenantCode)
	}
	if f.Code != "" {
		db = db.Where("code = ?", f.Code)
	}
	if f.UserID != 0 {
		db = db.Where("user_id = ?", f.UserID)
	}
	return db
}

func (p DataOpenAccessRepo) Insert(ctx context.Context, data *SysDataOpenAccess) error {
	result := p.db.WithContext(ctx).Create(data)
	return stores.ErrFmt(result.Error)
}

func (p DataOpenAccessRepo) FindOneByFilter(ctx context.Context, f DataOpenAccessFilter) (*SysDataOpenAccess, error) {
	var result SysDataOpenAccess
	db := p.fmtFilter(ctx, f)
	err := db.First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}
func (p DataOpenAccessRepo) FindByFilter(ctx context.Context, f DataOpenAccessFilter, page *stores.PageInfo) ([]*SysDataOpenAccess, error) {
	var results []*SysDataOpenAccess
	db := p.fmtFilter(ctx, f).Model(&SysDataOpenAccess{})
	db = page.ToGorm(db)
	err := db.Find(&results).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return results, nil
}

func (p DataOpenAccessRepo) CountByFilter(ctx context.Context, f DataOpenAccessFilter) (size int64, err error) {
	db := p.fmtFilter(ctx, f).Model(&SysDataOpenAccess{})
	err = db.Count(&size).Error
	return size, stores.ErrFmt(err)
}

func (p DataOpenAccessRepo) Update(ctx context.Context, data *SysDataOpenAccess) error {
	err := p.db.WithContext(ctx).Where("id = ?", data.ID).Save(data).Error
	return stores.ErrFmt(err)
}

func (p DataOpenAccessRepo) DeleteByFilter(ctx context.Context, f DataOpenAccessFilter) error {
	db := p.fmtFilter(ctx, f)
	err := db.Delete(&SysDataOpenAccess{}).Error
	return stores.ErrFmt(err)
}

func (p DataOpenAccessRepo) Delete(ctx context.Context, id int64) error {
	err := p.db.WithContext(ctx).Where("id = ?", id).Delete(&SysDataOpenAccess{}).Error
	return stores.ErrFmt(err)
}
func (p DataOpenAccessRepo) FindOne(ctx context.Context, id int64) (*SysDataOpenAccess, error) {
	var result SysDataOpenAccess
	err := p.db.WithContext(ctx).Where("id = ?", id).First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}

// 批量插入 LightStrategyDevice 记录
func (p DataOpenAccessRepo) MultiInsert(ctx context.Context, data []*SysDataOpenAccess) error {
	err := p.db.WithContext(ctx).Clauses(clause.OnConflict{UpdateAll: true}).Model(&SysDataOpenAccess{}).Create(data).Error
	return stores.ErrFmt(err)
}
