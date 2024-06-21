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
1. 将AreaProfile全局替换为模型的表名
2. 完善todo
*/

type AreaProfileRepo struct {
	db *gorm.DB
}

func NewAreaProfileRepo(in any) *AreaProfileRepo {
	return &AreaProfileRepo{db: stores.GetCommonConn(in)}
}

type AreaProfileFilter struct {
	Codes  []string
	Code   string
	AreaID int64
}

func (p AreaProfileRepo) fmtFilter(ctx context.Context, f AreaProfileFilter) *gorm.DB {
	db := p.db.WithContext(ctx)
	if len(f.Codes) != 0 {
		db = db.Where("code in ?", f.Codes)
	}
	if f.Code != "" {
		db = db.Where("code = ?", f.Code)
	}
	db = db.Where("area_id =?", f.AreaID)
	return db
}

func (p AreaProfileRepo) Insert(ctx context.Context, data *SysAreaProfile) error {
	result := p.db.WithContext(ctx).Create(data)
	return stores.ErrFmt(result.Error)
}

func (p AreaProfileRepo) FindOneByFilter(ctx context.Context, f AreaProfileFilter) (*SysAreaProfile, error) {
	var result SysAreaProfile
	db := p.fmtFilter(ctx, f)
	err := db.First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}
func (p AreaProfileRepo) FindByFilter(ctx context.Context, f AreaProfileFilter, page *def.PageInfo) ([]*SysAreaProfile, error) {
	var results []*SysAreaProfile
	db := p.fmtFilter(ctx, f).Model(&SysAreaProfile{})
	db = page.ToGorm(db)
	err := db.Find(&results).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return results, nil
}

func (p AreaProfileRepo) CountByFilter(ctx context.Context, f AreaProfileFilter) (size int64, err error) {
	db := p.fmtFilter(ctx, f).Model(&SysAreaProfile{})
	err = db.Count(&size).Error
	return size, stores.ErrFmt(err)
}

func (p AreaProfileRepo) Update(ctx context.Context, data *SysAreaProfile) error {
	err := p.db.WithContext(ctx).Where("area_id = ? and code = ?", data.AreaID, data.Code).Save(data).Error
	return stores.ErrFmt(err)
}

func (p AreaProfileRepo) DeleteByFilter(ctx context.Context, f AreaProfileFilter) error {
	db := p.fmtFilter(ctx, f)
	err := db.Delete(&SysAreaProfile{}).Error
	return stores.ErrFmt(err)
}

func (p AreaProfileRepo) Delete(ctx context.Context, id int64) error {
	err := p.db.WithContext(ctx).Where("id = ?", id).Delete(&SysAreaProfile{}).Error
	return stores.ErrFmt(err)
}
func (p AreaProfileRepo) FindOne(ctx context.Context, id int64) (*SysAreaProfile, error) {
	var result SysAreaProfile
	err := p.db.WithContext(ctx).Where("id = ?", id).First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}

// 批量插入 LightStrategyDevice 记录
func (p AreaProfileRepo) MultiInsert(ctx context.Context, data []*SysAreaProfile) error {
	err := p.db.WithContext(ctx).Clauses(clause.OnConflict{UpdateAll: true}).Model(&SysAreaProfile{}).Create(data).Error
	return stores.ErrFmt(err)
}
