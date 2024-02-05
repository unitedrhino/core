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
1. 将DictDetail全局替换为模型的表名
2. 完善todo
*/

type DictDetailRepo struct {
	db *gorm.DB
}

func NewDictDetailRepo(in any) *DictDetailRepo {
	return &DictDetailRepo{db: stores.GetCommonConn(in)}
}

type DictDetailFilter struct {
	DictID int64
}

func (p DictDetailRepo) fmtFilter(ctx context.Context, f DictDetailFilter) *gorm.DB {
	db := p.db.WithContext(ctx)
	if f.DictID != 0 {
		db = db.Where("dict_id=?", f.DictID)
	}
	return db
}

func (p DictDetailRepo) Insert(ctx context.Context, data *SysDictDetail) error {
	result := p.db.WithContext(ctx).Create(data)
	return stores.ErrFmt(result.Error)
}

func (p DictDetailRepo) FindOneByFilter(ctx context.Context, f DictDetailFilter) (*SysDictDetail, error) {
	var result SysDictDetail
	db := p.fmtFilter(ctx, f)
	err := db.First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}
func (p DictDetailRepo) FindByFilter(ctx context.Context, f DictDetailFilter, page *def.PageInfo) ([]*SysDictDetail, error) {
	var results []*SysDictDetail
	db := p.fmtFilter(ctx, f).Model(&SysDictDetail{})
	db = page.ToGorm(db)
	err := db.Find(&results).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return results, nil
}

func (p DictDetailRepo) CountByFilter(ctx context.Context, f DictDetailFilter) (size int64, err error) {
	db := p.fmtFilter(ctx, f).Model(&SysDictDetail{})
	err = db.Count(&size).Error
	return size, stores.ErrFmt(err)
}

func (p DictDetailRepo) Update(ctx context.Context, data *SysDictDetail) error {
	err := p.db.WithContext(ctx).Where("id = ?", data.ID).Save(data).Error
	return stores.ErrFmt(err)
}

func (p DictDetailRepo) DeleteByFilter(ctx context.Context, f DictDetailFilter) error {
	db := p.fmtFilter(ctx, f)
	err := db.Delete(&SysDictDetail{}).Error
	return stores.ErrFmt(err)
}

func (p DictDetailRepo) Delete(ctx context.Context, id int64) error {
	err := p.db.WithContext(ctx).Where("id = ?", id).Delete(&SysDictDetail{}).Error
	return stores.ErrFmt(err)
}
func (p DictDetailRepo) FindOne(ctx context.Context, id int64) (*SysDictDetail, error) {
	var result SysDictDetail
	err := p.db.WithContext(ctx).Where("id = ?", id).First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}

// 批量插入 LightStrategyDevice 记录
func (p DictDetailRepo) MultiInsert(ctx context.Context, data []*SysDictDetail) error {
	err := p.db.WithContext(ctx).Clauses(clause.OnConflict{UpdateAll: true}).Model(&SysDictDetail{}).Create(data).Error
	return stores.ErrFmt(err)
}
