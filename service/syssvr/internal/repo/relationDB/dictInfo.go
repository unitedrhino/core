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
1. 将DictInfo全局替换为模型的表名
2. 完善todo
*/

type DictInfoRepo struct {
	db *gorm.DB
}

func NewDictInfoRepo(in any) *DictInfoRepo {
	return &DictInfoRepo{db: stores.GetCommonConn(in)}
}

type DictInfoFilter struct {
	ID    int64
	Name  string
	Group string
	Code  string
}

func (p DictInfoRepo) fmtFilter(ctx context.Context, f DictInfoFilter) *gorm.DB {
	db := p.db.WithContext(ctx)
	if f.ID != 0 {
		db = db.Where("id = ?", f.ID)
	}
	if f.Name != "" {
		db = db.Where("name like ?", "%"+f.Name+"%")
	}
	if f.Code != "" {
		db = db.Where("code = ?", f.Code)
	}
	if f.Group != "" {
		db = db.Where("group = ?", f.Group)
	}
	return db
}

func (p DictInfoRepo) Insert(ctx context.Context, data *SysDictInfo) error {
	result := p.db.WithContext(ctx).Create(data)
	return stores.ErrFmt(result.Error)
}

func (p DictInfoRepo) FindOneByFilter(ctx context.Context, f DictInfoFilter) (*SysDictInfo, error) {
	var result SysDictInfo
	db := p.fmtFilter(ctx, f)
	err := db.First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}
func (p DictInfoRepo) FindByFilter(ctx context.Context, f DictInfoFilter, page *stores.PageInfo) ([]*SysDictInfo, error) {
	var results []*SysDictInfo
	db := p.fmtFilter(ctx, f).Model(&SysDictInfo{})
	db = page.ToGorm(db)
	err := db.Find(&results).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return results, nil
}

func (p DictInfoRepo) CountByFilter(ctx context.Context, f DictInfoFilter) (size int64, err error) {
	db := p.fmtFilter(ctx, f).Model(&SysDictInfo{})
	err = db.Count(&size).Error
	return size, stores.ErrFmt(err)
}

func (p DictInfoRepo) Update(ctx context.Context, data *SysDictInfo) error {
	err := p.db.WithContext(ctx).Where("id = ?", data.ID).Save(data).Error
	return stores.ErrFmt(err)
}

func (p DictInfoRepo) DeleteByFilter(ctx context.Context, f DictInfoFilter) error {
	db := p.fmtFilter(ctx, f)
	err := db.Delete(&SysDictInfo{}).Error
	return stores.ErrFmt(err)
}

func (p DictInfoRepo) Delete(ctx context.Context, id int64) error {
	err := p.db.WithContext(ctx).Where("id = ?", id).Delete(&SysDictInfo{}).Error
	return stores.ErrFmt(err)
}
func (p DictInfoRepo) FindOne(ctx context.Context, id int64) (*SysDictInfo, error) {
	var result SysDictInfo
	err := p.db.WithContext(ctx).Where("id = ?", id).First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}

// 批量插入 LightStrategyDevice 记录
func (p DictInfoRepo) MultiInsert(ctx context.Context, data []*SysDictInfo) error {
	err := p.db.WithContext(ctx).Clauses(clause.OnConflict{UpdateAll: true}).Model(&SysDictInfo{}).Create(data).Error
	return stores.ErrFmt(err)
}
