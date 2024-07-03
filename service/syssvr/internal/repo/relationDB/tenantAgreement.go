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
1. 将TenantAgreement全局替换为模型的表名
2. 完善todo
*/

type TenantAgreementRepo struct {
	db *gorm.DB
}

func NewTenantAgreementRepo(in any) *TenantAgreementRepo {
	return &TenantAgreementRepo{db: stores.GetCommonConn(in)}
}

type TenantAgreementFilter struct {
	ID   int64
	Code string
}

func (p TenantAgreementRepo) fmtFilter(ctx context.Context, f TenantAgreementFilter) *gorm.DB {
	db := p.db.WithContext(ctx)
	if f.ID != 0 {
		db = db.Where("id=?", f.ID)
	}
	if f.Code != "" {
		db = db.Where("code=?", f.Code)
	}
	return db
}

func (p TenantAgreementRepo) Insert(ctx context.Context, data *SysTenantAgreement) error {
	result := p.db.WithContext(ctx).Create(data)
	return stores.ErrFmt(result.Error)
}

func (p TenantAgreementRepo) FindOneByFilter(ctx context.Context, f TenantAgreementFilter) (*SysTenantAgreement, error) {
	var result SysTenantAgreement
	db := p.fmtFilter(ctx, f)
	err := db.First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}
func (p TenantAgreementRepo) FindByFilter(ctx context.Context, f TenantAgreementFilter, page *stores.PageInfo) ([]*SysTenantAgreement, error) {
	var results []*SysTenantAgreement
	db := p.fmtFilter(ctx, f).Model(&SysTenantAgreement{})
	db = page.ToGorm(db)
	err := db.Omit("Content").Find(&results).Error //列表不返回正文,太长了
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return results, nil
}

func (p TenantAgreementRepo) CountByFilter(ctx context.Context, f TenantAgreementFilter) (size int64, err error) {
	db := p.fmtFilter(ctx, f).Model(&SysTenantAgreement{})
	err = db.Count(&size).Error
	return size, stores.ErrFmt(err)
}

func (p TenantAgreementRepo) Update(ctx context.Context, data *SysTenantAgreement) error {
	err := p.db.WithContext(ctx).Where("id = ?", data.ID).Save(data).Error
	return stores.ErrFmt(err)
}

func (p TenantAgreementRepo) DeleteByFilter(ctx context.Context, f TenantAgreementFilter) error {
	db := p.fmtFilter(ctx, f)
	err := db.Delete(&SysTenantAgreement{}).Error
	return stores.ErrFmt(err)
}

func (p TenantAgreementRepo) Delete(ctx context.Context, id int64) error {
	err := p.db.WithContext(ctx).Where("id = ?", id).Delete(&SysTenantAgreement{}).Error
	return stores.ErrFmt(err)
}
func (p TenantAgreementRepo) FindOne(ctx context.Context, id int64) (*SysTenantAgreement, error) {
	var result SysTenantAgreement
	err := p.db.WithContext(ctx).Where("id = ?", id).First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}

// 批量插入 LightStrategyDevice 记录
func (p TenantAgreementRepo) MultiInsert(ctx context.Context, data []*SysTenantAgreement) error {
	err := p.db.WithContext(ctx).Clauses(clause.OnConflict{UpdateAll: true}).Model(&SysTenantAgreement{}).Create(data).Error
	return stores.ErrFmt(err)
}

func (d TenantAgreementRepo) UpdateWithField(ctx context.Context, f TenantAgreementFilter, updates map[string]any) error {
	db := d.fmtFilter(ctx, f)
	err := db.Model(&SysTenantAgreement{}).Updates(updates).Error
	return stores.ErrFmt(err)
}
