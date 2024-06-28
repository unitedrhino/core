package relationDB

import (
	"context"
	"gitee.com/i-Things/share/stores"
	"gorm.io/gorm"
)

type ApiInfoRepo struct {
	db *gorm.DB
}

func NewApiInfoRepo(in any) *ApiInfoRepo {
	return &ApiInfoRepo{db: stores.GetCommonConn(in)}
}

type ApiInfoFilter struct {
	ApiIDs       []int64
	Route        string
	Method       string
	Name         string
	AccessCode   string
	IsAuthTenant int64
}

func (p ApiInfoRepo) fmtFilter(ctx context.Context, f ApiInfoFilter) *gorm.DB {
	db := p.db.WithContext(ctx)
	if f.ApiIDs != nil {
		db = db.Where("id in ?", f.ApiIDs)
	}
	if f.IsAuthTenant != 0 {
		db = db.Where("is_auth_tenant =?", f.IsAuthTenant)
	}
	if f.Route != "" {
		db = db.Where("route like ?", "%"+f.Route+"%")
	}
	if f.AccessCode != "" {
		db = db.Where("api_scope_code =?", f.AccessCode)
	}
	if f.Method != "" {
		db = db.Where("method = ?", f.Method)
	}
	if f.Name != "" {
		db = db.Where("name like ?", "%"+f.Name+"%")
	}
	return db
}

func (p ApiInfoRepo) Insert(ctx context.Context, data *SysApiInfo) error {
	result := p.db.WithContext(ctx).Create(data)
	return stores.ErrFmt(result.Error)
}

func (p ApiInfoRepo) FindOneByFilter(ctx context.Context, f ApiInfoFilter) (*SysApiInfo, error) {
	var result SysApiInfo
	db := p.fmtFilter(ctx, f)
	err := db.First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}

func (p ApiInfoRepo) FindByFilter(ctx context.Context, f ApiInfoFilter, page *stores.PageInfo) ([]*SysApiInfo, error) {
	var results []*SysApiInfo
	db := p.fmtFilter(ctx, f).Model(&SysApiInfo{})
	db = page.ToGorm(db)
	err := db.Find(&results).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return results, nil
}

func (p ApiInfoRepo) CountByFilter(ctx context.Context, f ApiInfoFilter) (size int64, err error) {
	db := p.fmtFilter(ctx, f).Model(&SysApiInfo{})
	err = db.Count(&size).Error
	return size, stores.ErrFmt(err)
}

func (p ApiInfoRepo) Update(ctx context.Context, data *SysApiInfo) error {
	err := p.db.WithContext(ctx).Where("id = ?", data.ID).Save(data).Error
	return stores.ErrFmt(err)
}

func (p ApiInfoRepo) DeleteByFilter(ctx context.Context, f ApiInfoFilter) error {
	db := p.fmtFilter(ctx, f)
	err := db.Delete(&SysApiInfo{}).Error
	return stores.ErrFmt(err)
}
func (p ApiInfoRepo) Delete(ctx context.Context, id int64) error {
	err := p.db.WithContext(ctx).Where("id = ?", id).Delete(&SysApiInfo{}).Error
	return stores.ErrFmt(err)
}

func (p ApiInfoRepo) FindOne(ctx context.Context, id int64) (*SysApiInfo, error) {
	var result SysApiInfo
	err := p.db.WithContext(ctx).Where("id = ?", id).First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}
