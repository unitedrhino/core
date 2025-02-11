package relationDB

import (
	"context"
	"gitee.com/unitedrhino/share/stores"
	"gorm.io/gorm"
)

type OperLogRepo struct {
	db *gorm.DB
}

func NewOperLogRepo(in any) *OperLogRepo {
	return &OperLogRepo{db: stores.GetCommonConn(in)}
}

type OperLogFilter struct {
	TenantCode   string
	OperName     string
	OperUserName string
	BusinessType int64
	AppCode      string
	Code         int64
	OperUserID   int64
	CreateTime   *stores.Cmp
}

func (p OperLogRepo) fmtFilter(ctx context.Context, f OperLogFilter) *gorm.DB {
	db := p.db.WithContext(ctx)
	db = f.CreateTime.Where(db, "created_time")
	if f.TenantCode != "" {
		db = db.Where("tenant_code = ?", f.TenantCode)
	}
	if f.OperName != "" {
		db = db.Where("oper_name = ?", f.OperName)
	}
	if f.OperUserName != "" {
		db = db.Where("oper_user_name = ?", f.OperUserName)
	}
	if f.BusinessType > 0 {
		db = db.Where("business_type= ?", f.BusinessType)
	}
	if f.AppCode != "" {
		db = db.Where("app_code = ?")
	}
	if f.Code != 0 {
		db = db.Where("code = ?")
	}
	if f.OperUserID != 0 {
		db = db.Where("oper_user_id = ?")
	}
	return db
}

func (p OperLogRepo) Insert(ctx context.Context, data *SysOperLog) error {
	result := p.db.WithContext(ctx).Create(data)
	return stores.ErrFmt(result.Error)
}

func (p OperLogRepo) FindOneByFilter(ctx context.Context, f OperLogFilter) (*SysOperLog, error) {
	var result SysOperLog
	db := p.fmtFilter(ctx, f)
	err := db.First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}
func (p OperLogRepo) FindByFilter(ctx context.Context, f OperLogFilter, page *stores.PageInfo) ([]*SysOperLog, error) {
	var results []*SysOperLog
	db := p.fmtFilter(ctx, f).Model(&SysOperLog{})
	db = page.ToGorm(db)
	err := db.Find(&results).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return results, nil
}

func (p OperLogRepo) CountByFilter(ctx context.Context, f OperLogFilter) (size int64, err error) {
	db := p.fmtFilter(ctx, f).Model(&SysOperLog{})
	err = db.Count(&size).Error
	return size, stores.ErrFmt(err)
}

func (p OperLogRepo) Update(ctx context.Context, data *SysOperLog) error {
	err := p.db.WithContext(ctx).Where("id = ?", data.ID).Save(data).Error
	return stores.ErrFmt(err)
}

func (p OperLogRepo) DeleteByFilter(ctx context.Context, f OperLogFilter) error {
	db := p.fmtFilter(ctx, f)
	err := db.Delete(&SysOperLog{}).Error
	return stores.ErrFmt(err)
}

func (p OperLogRepo) Delete(ctx context.Context, id int64) error {
	err := p.db.WithContext(ctx).Where("id = ?", id).Delete(&SysOperLog{}).Error
	return stores.ErrFmt(err)
}
func (p OperLogRepo) FindOne(ctx context.Context, id int64) (*SysOperLog, error) {
	var result SysOperLog
	err := p.db.WithContext(ctx).Where("id = ?", id).First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}
