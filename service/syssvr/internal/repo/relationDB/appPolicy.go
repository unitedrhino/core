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
1. 将AppPolicy全局替换为模型的表名
2. 完善todo
*/

type AppPolicyRepo struct {
	db *gorm.DB
}

func NewAppPolicyRepo(in any) *AppPolicyRepo {
	return &AppPolicyRepo{db: stores.GetCommonConn(in)}
}

type AppPolicyFilter struct {
	ID      int64
	Code    string
	AppCode string
}

func (p AppPolicyRepo) fmtFilter(ctx context.Context, f AppPolicyFilter) *gorm.DB {
	db := p.db.WithContext(ctx)
	if f.ID != 0 {
		db = db.Where("id=?", f.ID)
	}
	if f.AppCode != "" {
		db = db.Where("app_code=?", f.AppCode)
	}
	if f.Code != "" {
		db = db.Where("code=?", f.Code)
	}
	return db
}

func (p AppPolicyRepo) Insert(ctx context.Context, data *SysAppPolicy) error {
	result := p.db.WithContext(ctx).Create(data)
	return stores.ErrFmt(result.Error)
}

func (p AppPolicyRepo) FindOneByFilter(ctx context.Context, f AppPolicyFilter) (*SysAppPolicy, error) {
	var result SysAppPolicy
	db := p.fmtFilter(ctx, f)
	err := db.First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}
func (p AppPolicyRepo) FindByFilter(ctx context.Context, f AppPolicyFilter, page *def.PageInfo) ([]*SysAppPolicy, error) {
	var results []*SysAppPolicy
	db := p.fmtFilter(ctx, f).Model(&SysAppPolicy{})
	db = page.ToGorm(db)
	err := db.Find(&results).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return results, nil
}

func (p AppPolicyRepo) CountByFilter(ctx context.Context, f AppPolicyFilter) (size int64, err error) {
	db := p.fmtFilter(ctx, f).Model(&SysAppPolicy{})
	err = db.Count(&size).Error
	return size, stores.ErrFmt(err)
}

func (p AppPolicyRepo) Update(ctx context.Context, data *SysAppPolicy) error {
	err := p.db.WithContext(ctx).Where("id = ?", data.ID).Save(data).Error
	return stores.ErrFmt(err)
}

func (p AppPolicyRepo) DeleteByFilter(ctx context.Context, f AppPolicyFilter) error {
	db := p.fmtFilter(ctx, f)
	err := db.Delete(&SysAppPolicy{}).Error
	return stores.ErrFmt(err)
}

func (p AppPolicyRepo) Delete(ctx context.Context, id int64) error {
	err := p.db.WithContext(ctx).Where("id = ?", id).Delete(&SysAppPolicy{}).Error
	return stores.ErrFmt(err)
}
func (p AppPolicyRepo) FindOne(ctx context.Context, id int64) (*SysAppPolicy, error) {
	var result SysAppPolicy
	err := p.db.WithContext(ctx).Where("id = ?", id).First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}

// 批量插入 LightStrategyDevice 记录
func (p AppPolicyRepo) MultiInsert(ctx context.Context, data []*SysAppPolicy) error {
	err := p.db.WithContext(ctx).Clauses(clause.OnConflict{UpdateAll: true}).Model(&SysAppPolicy{}).Create(data).Error
	return stores.ErrFmt(err)
}
