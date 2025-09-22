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
1. 将ProjectCurl全局替换为模型的表名
2. 完善todo
*/

type ProjectCurlRepo struct {
	db *gorm.DB
}

func NewProjectCurlRepo(in any) *ProjectCurlRepo {
	return &ProjectCurlRepo{db: stores.GetCommonConn(in)}
}

type ProjectCurlFilter struct {
	Purpose string `gorm:"column:purpose;type:VARCHAR(50);index:idx_sys_project_profile_tc_un;NOT NULL"` //用途必填
	Params  map[string]*stores.Compare
}

func (p ProjectCurlRepo) fmtFilter(ctx context.Context, f ProjectCurlFilter) *gorm.DB {
	db := p.db.WithContext(ctx)
	db = db.Where("purpose = ?", f.Purpose)
	for k, v := range f.Params {
		db = stores.GetCmp(v.CmpType, v.Value).Where(db, stores.Cast(stores.JsonCol("params", k), v.CastTo))
	}
	return db
}

func (p ProjectCurlRepo) Insert(ctx context.Context, data *SysProjectCrud) error {
	result := p.db.WithContext(ctx).Create(data)
	return stores.ErrFmt(result.Error)
}

func (p ProjectCurlRepo) FindOneByFilter(ctx context.Context, f ProjectCurlFilter) (*SysProjectCrud, error) {
	var result SysProjectCrud
	db := p.fmtFilter(ctx, f)
	err := db.First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}
func (p ProjectCurlRepo) FindByFilter(ctx context.Context, f ProjectCurlFilter, page *stores.PageInfo) ([]*SysProjectCrud, error) {
	var results []*SysProjectCrud
	db := p.fmtFilter(ctx, f).Model(&SysProjectCrud{})
	db = page.ToGorm(db)
	err := db.Find(&results).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return results, nil
}

func (p ProjectCurlRepo) CountByFilter(ctx context.Context, f ProjectCurlFilter) (size int64, err error) {
	db := p.fmtFilter(ctx, f).Model(&SysProjectCrud{})
	err = db.Count(&size).Error
	return size, stores.ErrFmt(err)
}

func (p ProjectCurlRepo) Update(ctx context.Context, data *SysProjectCrud) error {
	err := p.db.WithContext(ctx).Where("id = ?", data.ID).Save(data).Error
	return stores.ErrFmt(err)
}

func (p ProjectCurlRepo) DeleteByFilter(ctx context.Context, f ProjectCurlFilter) error {
	db := p.fmtFilter(ctx, f)
	err := db.Delete(&SysProjectCrud{}).Error
	return stores.ErrFmt(err)
}

func (p ProjectCurlRepo) Delete(ctx context.Context, id int64) error {
	err := p.db.WithContext(ctx).Where("id = ?", id).Delete(&SysProjectCrud{}).Error
	return stores.ErrFmt(err)
}
func (p ProjectCurlRepo) FindOne(ctx context.Context, id int64) (*SysProjectCrud, error) {
	var result SysProjectCrud
	err := p.db.WithContext(ctx).Where("id = ?", id).First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}

// 批量插入 LightStrategyDevice 记录
func (p ProjectCurlRepo) MultiInsert(ctx context.Context, data []*SysProjectCrud) error {
	err := p.db.WithContext(ctx).Clauses(clause.OnConflict{UpdateAll: true}).Model(&SysProjectCrud{}).Create(data).Error
	return stores.ErrFmt(err)
}
