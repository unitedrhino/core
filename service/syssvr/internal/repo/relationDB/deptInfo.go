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
1. 将DeptInfo全局替换为模型的表名
2. 完善todo
*/

type DeptInfoRepo struct {
	db *gorm.DB
}

func NewDeptInfoRepo(in any) *DeptInfoRepo {
	return &DeptInfoRepo{db: stores.GetCommonConn(in)}
}

type DeptInfoFilter struct {
	ID           int64
	IDs          []int64
	WithChildren bool
	IDPath       string
	ParentID     int64
	Status       int64
	Name         string
	Names        []string
}

func (p DeptInfoRepo) fmtFilter(ctx context.Context, f DeptInfoFilter) *gorm.DB {
	db := p.db.WithContext(ctx)

	if f.Name != "" {
		db = db.Where("name like ?", "%"+f.Name+"%")
	}
	if len(f.Names) > 0 {
		db = db.Where("name in ?", f.Names)
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

func (p DeptInfoRepo) Insert(ctx context.Context, data *SysDeptInfo) error {
	result := p.db.WithContext(ctx).Create(data)
	return stores.ErrFmt(result.Error)
}

func (p DeptInfoRepo) FindOneByFilter(ctx context.Context, f DeptInfoFilter) (*SysDeptInfo, error) {
	var result SysDeptInfo
	db := p.fmtFilter(ctx, f)
	err := db.First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}
func (p DeptInfoRepo) FindByFilter(ctx context.Context, f DeptInfoFilter, page *stores.PageInfo) ([]*SysDeptInfo, error) {
	var results []*SysDeptInfo
	db := p.fmtFilter(ctx, f).Model(&SysDeptInfo{})
	db = page.ToGorm(db)
	err := db.Find(&results).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return results, nil
}

func (p DeptInfoRepo) CountByFilter(ctx context.Context, f DeptInfoFilter) (size int64, err error) {
	db := p.fmtFilter(ctx, f).Model(&SysDeptInfo{})
	err = db.Count(&size).Error
	return size, stores.ErrFmt(err)
}

func (p DeptInfoRepo) Update(ctx context.Context, data *SysDeptInfo) error {
	err := p.db.WithContext(ctx).Where("id = ?", data.ID).Save(data).Error
	return stores.ErrFmt(err)
}

func (d DeptInfoRepo) UpdateWithField(ctx context.Context, f DeptInfoFilter, updates map[string]any) error {
	db := d.fmtFilter(ctx, f)
	err := db.Model(&SysDeptInfo{}).Updates(updates).Error
	return stores.ErrFmt(err)
}

func (p DeptInfoRepo) DeleteByFilter(ctx context.Context, f DeptInfoFilter) error {
	db := p.fmtFilter(ctx, f)
	err := db.Delete(&SysDeptInfo{}).Error
	return stores.ErrFmt(err)
}

func (p DeptInfoRepo) Delete(ctx context.Context, id int64) error {
	err := p.db.WithContext(ctx).Where("id = ?", id).Delete(&SysDeptInfo{}).Error
	return stores.ErrFmt(err)
}
func (p DeptInfoRepo) FindOne(ctx context.Context, id int64) (*SysDeptInfo, error) {
	var result SysDeptInfo
	err := p.db.WithContext(ctx).Where("id = ?", id).First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}

// 批量插入 LightStrategyDevice 记录
func (p DeptInfoRepo) MultiInsert(ctx context.Context, data []*SysDeptInfo) error {
	err := p.db.WithContext(ctx).Clauses(clause.OnConflict{UpdateAll: true}).Model(&SysDeptInfo{}).Create(data).Error
	return stores.ErrFmt(err)
}
