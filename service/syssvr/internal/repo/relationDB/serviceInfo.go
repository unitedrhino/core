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
1. 将ServiceInfo全局替换为模型的表名
2. 完善todo
*/

type ServiceInfoRepo struct {
	db *gorm.DB
}

func NewServiceInfoRepo(in any) *ServiceInfoRepo {
	return &ServiceInfoRepo{db: stores.GetCommonConn(in)}
}

type ServiceInfoFilter struct {
	Code string
}

func (p ServiceInfoRepo) fmtFilter(ctx context.Context, f ServiceInfoFilter) *gorm.DB {
	db := p.db.WithContext(ctx)
	db = db.Where("code = ?", f.Code)

	return db
}

func (p ServiceInfoRepo) Insert(ctx context.Context, data *SysServiceInfo) error {
	result := p.db.WithContext(ctx).Create(data)
	return stores.ErrFmt(result.Error)
}

func (p ServiceInfoRepo) FindOneByFilter(ctx context.Context, f ServiceInfoFilter) (*SysServiceInfo, error) {
	var result SysServiceInfo
	db := p.fmtFilter(ctx, f)
	err := db.First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}
func (p ServiceInfoRepo) FindByFilter(ctx context.Context, f ServiceInfoFilter, page *stores.PageInfo) ([]*SysServiceInfo, error) {
	var results []*SysServiceInfo
	db := p.fmtFilter(ctx, f).Model(&SysServiceInfo{})
	db = page.ToGorm(db)
	err := db.Find(&results).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return results, nil
}

func (p ServiceInfoRepo) CountByFilter(ctx context.Context, f ServiceInfoFilter) (size int64, err error) {
	db := p.fmtFilter(ctx, f).Model(&SysServiceInfo{})
	err = db.Count(&size).Error
	return size, stores.ErrFmt(err)
}

func (p ServiceInfoRepo) Update(ctx context.Context, data *SysServiceInfo) error {
	err := p.db.WithContext(ctx).Where("code = ?", data.Code).Save(data).Error
	return stores.ErrFmt(err)
}

func (p ServiceInfoRepo) DeleteByFilter(ctx context.Context, f ServiceInfoFilter) error {
	db := p.fmtFilter(ctx, f)
	err := db.Delete(&SysServiceInfo{}).Error
	return stores.ErrFmt(err)
}

func (p ServiceInfoRepo) Delete(ctx context.Context, id int64) error {
	err := p.db.WithContext(ctx).Where("id = ?", id).Delete(&SysServiceInfo{}).Error
	return stores.ErrFmt(err)
}
func (p ServiceInfoRepo) FindOne(ctx context.Context, id int64) (*SysServiceInfo, error) {
	var result SysServiceInfo
	err := p.db.WithContext(ctx).Where("id = ?", id).First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}

// 批量插入 LightStrategyDevice 记录
func (p ServiceInfoRepo) MultiInsert(ctx context.Context, data []*SysServiceInfo) error {
	err := p.db.WithContext(ctx).Clauses(clause.OnConflict{UpdateAll: true}).Model(&SysServiceInfo{}).Create(data).Error
	return stores.ErrFmt(err)
}
