package relationDB

import (
	"context"
	"gitee.com/unitedrhino/share/stores"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type MenuInfoRepo struct {
	db *gorm.DB
}

func NewMenuInfoRepo(in any) *MenuInfoRepo {
	return &MenuInfoRepo{db: stores.GetCommonConn(in)}
}

type MenuInfoFilter struct {
	ModuleCode string
	Name       string
	Path       string
	Paths      []string
	MenuIDs    []int64
	ParentID   int64
	ParentIDs  []int64
	IsCommon   int64
}

func (p MenuInfoRepo) fmtFilter(ctx context.Context, f MenuInfoFilter) *gorm.DB {
	db := p.db.WithContext(ctx)
	if f.IsCommon != 0 {
		db = db.Where("is_common=?", f.IsCommon)
	}
	if f.ParentID != 0 {
		db = db.Where("parent_id=?", f.ParentID)
	}
	if len(f.ParentIDs) != 0 {
		db = db.Where("parent_id IN ?", f.ParentIDs)
	}
	if f.ModuleCode != "" {
		db = db.Where("module_code =?", f.ModuleCode)
	}
	if f.Name != "" {
		db = db.Where("name like ?", "%"+f.Name+"%")
	}
	if f.Path != "" {
		db = db.Where("path like ?", "%"+f.Path+"%")
	}
	if len(f.Path) != 0 {
		db = db.Where("path in ?", f.Paths)
	}
	if len(f.MenuIDs) != 0 {
		db = db.Where("id in ?", f.MenuIDs)
	}
	return db
}

func (p MenuInfoRepo) Insert(ctx context.Context, data *SysModuleMenu) error {
	result := p.db.WithContext(ctx).Create(data)
	return stores.ErrFmt(result.Error)
}

func (p MenuInfoRepo) FindOneByFilter(ctx context.Context, f MenuInfoFilter) (*SysModuleMenu, error) {
	var result SysModuleMenu
	db := p.fmtFilter(ctx, f)
	err := db.First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}
func (p MenuInfoRepo) FindByFilter(ctx context.Context, f MenuInfoFilter, page *stores.PageInfo) ([]*SysModuleMenu, error) {
	var results []*SysModuleMenu
	db := p.fmtFilter(ctx, f).Model(&SysModuleMenu{})
	db = page.ToGorm(db).Order(stores.Col("order"))
	err := db.Find(&results).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return results, nil
}

func (p MenuInfoRepo) CountByFilter(ctx context.Context, f MenuInfoFilter) (size int64, err error) {
	db := p.fmtFilter(ctx, f).Model(&SysModuleMenu{})
	err = db.Count(&size).Error
	return size, stores.ErrFmt(err)
}

func (p MenuInfoRepo) Update(ctx context.Context, data *SysModuleMenu) error {
	err := p.db.WithContext(ctx).Where("id = ?", data.ID).Save(data).Error
	return stores.ErrFmt(err)
}

func (p MenuInfoRepo) DeleteByFilter(ctx context.Context, f MenuInfoFilter) error {
	db := p.fmtFilter(ctx, f)
	err := db.Delete(&SysModuleMenu{}).Error
	return stores.ErrFmt(err)
}

func (p MenuInfoRepo) Delete(ctx context.Context, id int64) error {
	err := p.db.WithContext(ctx).Where("id = ?", id).Delete(&SysModuleMenu{}).Error
	return stores.ErrFmt(err)
}
func (p MenuInfoRepo) FindOne(ctx context.Context, id int64) (*SysModuleMenu, error) {
	var result SysModuleMenu
	err := p.db.WithContext(ctx).Where("id = ?", id).First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}

// 批量插入 LightStrategyDevice 记录
func (m MenuInfoRepo) MultiInsert(ctx context.Context, data []*SysModuleMenu) error {
	err := m.db.WithContext(ctx).Clauses(clause.OnConflict{UpdateAll: true}).Model(&SysModuleMenu{}).Create(data).Error
	return stores.ErrFmt(err)
}

func (m MenuInfoRepo) MultiInsertOnly(ctx context.Context, data []*SysModuleMenu) error {
	err := m.db.WithContext(ctx).Clauses(clause.OnConflict{DoNothing: true}).Model(&SysModuleMenu{}).Create(data).Error
	return stores.ErrFmt(err)
}
