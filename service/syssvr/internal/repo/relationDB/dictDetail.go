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
	ID           int64
	IDs          []int64
	DictCode     string
	WithChildren bool
	IDPath       string
	ParentID     int64
	Status       int64
	Label        string
	Value        string
	Values       []string
}

func (p DictDetailRepo) fmtFilter(ctx context.Context, f DictDetailFilter) *gorm.DB {
	db := p.db.WithContext(ctx)
	if f.DictCode != "" {
		db = db.Where("dict_code = ?", f.DictCode)
	}
	if f.Label != "" {
		db = db.Where("label like ?", "%"+f.Label+"%")
	}
	if f.Status != 0 {
		db = db.Where("status = ?", f.Status)
	}
	if f.Value != "" {
		db = db.Where("value = ?", f.Value)
	}
	if len(f.Values) != 0 {
		db = db.Where("values IN ?", f.Values)
	}
	if f.IDPath != "" {
		db = db.Where("id_path like ?", f.IDPath+"%")
	}
	if f.ParentID != 0 {
		db = db.Where("parent_id = ?", f.ParentID)
	}
	if f.ID != 0 {
		db = db.Where("id = ?", f.ID)
	}
	if len(f.IDs) > 0 {
		db = db.Where("id in ?", f.IDs)
	}

	if f.WithChildren {
		db = db.Preload("Children")
	}

	return db
}

func (p DictDetailRepo) Insert(ctx context.Context, data *SysDictDetail) error {
	result := p.db.WithContext(ctx).Create(data)
	return stores.ErrFmt(result.Error)
}

func (d DictDetailRepo) UpdateWithField(ctx context.Context, f DictDetailFilter, updates map[string]any) error {
	db := d.fmtFilter(ctx, f)
	err := db.Model(&SysDictDetail{}).Updates(updates).Error
	return stores.ErrFmt(err)
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
func (p DictDetailRepo) FindByFilter(ctx context.Context, f DictDetailFilter, page *stores.PageInfo) ([]*SysDictDetail, error) {
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
	err := p.db.WithContext(ctx).Clauses(clause.OnConflict{UpdateAll: true,
		Columns: stores.SetColumnsWithPg(p.db, &SysDictDetail{}, "idx_sys_dict_detail_value")}).Model(&SysDictDetail{}).CreateInBatches(data, 100).Error
	return stores.ErrFmt(err)
}
