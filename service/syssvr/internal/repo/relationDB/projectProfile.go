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
1. 将ProjectProfile全局替换为模型的表名
2. 完善todo
*/

type ProjectProfileRepo struct {
	db *gorm.DB
}

func NewProjectProfileRepo(in any) *ProjectProfileRepo {
	return &ProjectProfileRepo{db: stores.GetCommonConn(in)}
}

type ProjectProfileFilter struct {
	Codes     []string
	Code      string
	ProjectID int64
}

func (p ProjectProfileRepo) fmtFilter(ctx context.Context, f ProjectProfileFilter) *gorm.DB {
	db := p.db.WithContext(ctx)
	if len(f.Codes) != 0 {
		db = db.Where("code in ?", f.Codes)
	}
	if f.Code != "" {
		db = db.Where("code = ?", f.Code)
	}
	db = db.Where("project_id =?", f.ProjectID)
	return db
}

func (p ProjectProfileRepo) Insert(ctx context.Context, data *SysProjectProfile) error {
	result := p.db.WithContext(ctx).Create(data)
	return stores.ErrFmt(result.Error)
}

func (p ProjectProfileRepo) FindOneByFilter(ctx context.Context, f ProjectProfileFilter) (*SysProjectProfile, error) {
	var result SysProjectProfile
	db := p.fmtFilter(ctx, f)
	err := db.First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}
func (p ProjectProfileRepo) FindByFilter(ctx context.Context, f ProjectProfileFilter, page *stores.PageInfo) ([]*SysProjectProfile, error) {
	var results []*SysProjectProfile
	db := p.fmtFilter(ctx, f).Model(&SysProjectProfile{})
	db = page.ToGorm(db)
	err := db.Find(&results).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return results, nil
}

func (p ProjectProfileRepo) CountByFilter(ctx context.Context, f ProjectProfileFilter) (size int64, err error) {
	db := p.fmtFilter(ctx, f).Model(&SysProjectProfile{})
	err = db.Count(&size).Error
	return size, stores.ErrFmt(err)
}

func (p ProjectProfileRepo) Update(ctx context.Context, data *SysProjectProfile) error {
	err := p.db.WithContext(ctx).Where("project_id = ? and code = ?", data.ProjectID, data.Code).Save(data).Error
	return stores.ErrFmt(err)
}

func (p ProjectProfileRepo) DeleteByFilter(ctx context.Context, f ProjectProfileFilter) error {
	db := p.fmtFilter(ctx, f)
	err := db.Delete(&SysProjectProfile{}).Error
	return stores.ErrFmt(err)
}

func (p ProjectProfileRepo) Delete(ctx context.Context, id int64) error {
	err := p.db.WithContext(ctx).Where("id = ?", id).Delete(&SysProjectProfile{}).Error
	return stores.ErrFmt(err)
}
func (p ProjectProfileRepo) FindOne(ctx context.Context, id int64) (*SysProjectProfile, error) {
	var result SysProjectProfile
	err := p.db.WithContext(ctx).Where("id = ?", id).First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}

// 批量插入 LightStrategyDevice 记录
func (p ProjectProfileRepo) MultiInsert(ctx context.Context, data []*SysProjectProfile) error {
	err := p.db.WithContext(ctx).Clauses(clause.OnConflict{UpdateAll: true}).Model(&SysProjectProfile{}).Create(data).Error
	return stores.ErrFmt(err)
}
