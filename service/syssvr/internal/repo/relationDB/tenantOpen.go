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
1. 将TenantOpen全局替换为模型的表名
2. 完善todo
*/

type TenantOpenRepo struct {
	db *gorm.DB
}

func NewTenantOpenRepo(in any) *TenantOpenRepo {
	return &TenantOpenRepo{db: stores.GetCommonConn(in)}
}

type TenantOpenFilter struct {
	TenantCode string
	UserID     int64
	Code       string
}

func (p TenantOpenRepo) fmtFilter(ctx context.Context, f TenantOpenFilter) *gorm.DB {
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

func (p TenantOpenRepo) Insert(ctx context.Context, data *SysTenantOpenAccess) error {
	result := p.db.WithContext(ctx).Create(data)
	return stores.ErrFmt(result.Error)
}

func (p TenantOpenRepo) FindOneByFilter(ctx context.Context, f TenantOpenFilter) (*SysTenantOpenAccess, error) {
	var result SysTenantOpenAccess
	db := p.fmtFilter(ctx, f)
	err := db.First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}
func (p TenantOpenRepo) FindByFilter(ctx context.Context, f TenantOpenFilter, page *def.PageInfo) ([]*SysTenantOpenAccess, error) {
	var results []*SysTenantOpenAccess
	db := p.fmtFilter(ctx, f).Model(&SysTenantOpenAccess{})
	db = page.ToGorm(db)
	err := db.Find(&results).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return results, nil
}

func (p TenantOpenRepo) CountByFilter(ctx context.Context, f TenantOpenFilter) (size int64, err error) {
	db := p.fmtFilter(ctx, f).Model(&SysTenantOpenAccess{})
	err = db.Count(&size).Error
	return size, stores.ErrFmt(err)
}

func (p TenantOpenRepo) Update(ctx context.Context, data *SysTenantOpenAccess) error {
	err := p.db.WithContext(ctx).Where("id = ?", data.ID).Save(data).Error
	return stores.ErrFmt(err)
}

func (p TenantOpenRepo) DeleteByFilter(ctx context.Context, f TenantOpenFilter) error {
	db := p.fmtFilter(ctx, f)
	err := db.Delete(&SysTenantOpenAccess{}).Error
	return stores.ErrFmt(err)
}

func (p TenantOpenRepo) Delete(ctx context.Context, id int64) error {
	err := p.db.WithContext(ctx).Where("id = ?", id).Delete(&SysTenantOpenAccess{}).Error
	return stores.ErrFmt(err)
}
func (p TenantOpenRepo) FindOne(ctx context.Context, id int64) (*SysTenantOpenAccess, error) {
	var result SysTenantOpenAccess
	err := p.db.WithContext(ctx).Where("id = ?", id).First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}

// 批量插入 LightStrategyDevice 记录
func (p TenantOpenRepo) MultiInsert(ctx context.Context, data []*SysTenantOpenAccess) error {
	err := p.db.WithContext(ctx).Clauses(clause.OnConflict{UpdateAll: true}).Model(&SysTenantOpenAccess{}).Create(data).Error
	return stores.ErrFmt(err)
}
