package relationDB

import (
	"context"
	"gitee.com/unitedrhino/share/stores"
	"gorm.io/gorm"
)

type LoginLogRepo struct {
	db *gorm.DB
}

func NewLoginLogRepo(in any) *LoginLogRepo {
	return &LoginLogRepo{db: stores.GetCommonConn(in)}
}

type DateRange struct {
	Start string
	End   string
}
type LoginLogFilter struct {
	TenantCode    string
	IpAddr        string
	LoginLocation string
	Data          *DateRange
	CreateTime    *stores.Cmp
	UserID        int64
	UserName      string
	Code          int64
}

func (p LoginLogRepo) fmtFilter(ctx context.Context, f LoginLogFilter) *gorm.DB {
	db := p.db.WithContext(ctx)
	db = f.CreateTime.Where(db, "created_time")
	if f.TenantCode != "" {
		db = db.Where("tenant_code = ?", f.TenantCode)
	}
	if f.UserID != 0 {
		db = db.Where("user_id = ?", f.UserID)
	}
	if f.UserName != "" {
		db = db.Where("user_name = ?", f.UserName)
	}
	if f.Code != 0 {
		db = db.Where("code = ?", f.Code)
	}
	if f.IpAddr != "" {
		db = db.Where("ip_addr= ?", f.IpAddr)
	}
	if f.LoginLocation != "" {
		db = db.Where("login_location like ?", "%"+f.LoginLocation+"%")
	}
	if f.Data != nil && f.Data.Start != "" && f.Data.End != "" {
		db = db.Where("created_time >= ? and created_time <= ?", f.Data.Start, f.Data.End)
	}
	return db
}

func (p LoginLogRepo) Insert(ctx context.Context, data *SysLoginLog) error {
	result := p.db.WithContext(ctx).Create(data)
	return stores.ErrFmt(result.Error)
}

func (p LoginLogRepo) FindOneByFilter(ctx context.Context, f LoginLogFilter) (*SysLoginLog, error) {
	var result SysLoginLog
	db := p.fmtFilter(ctx, f)
	err := db.First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}
func (p LoginLogRepo) FindByFilter(ctx context.Context, f LoginLogFilter, page *stores.PageInfo) ([]*SysLoginLog, error) {
	var results []*SysLoginLog
	db := p.fmtFilter(ctx, f).Model(&SysLoginLog{})
	db = page.ToGorm(db)
	err := db.Find(&results).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return results, nil
}

func (p LoginLogRepo) CountByFilter(ctx context.Context, f LoginLogFilter) (size int64, err error) {
	db := p.fmtFilter(ctx, f).Model(&SysLoginLog{})
	err = db.Count(&size).Error
	return size, stores.ErrFmt(err)
}

func (p LoginLogRepo) Update(ctx context.Context, data *SysLoginLog) error {
	err := p.db.WithContext(ctx).Where("id = ?", data.ID).Save(data).Error
	return stores.ErrFmt(err)
}

func (p LoginLogRepo) DeleteByFilter(ctx context.Context, f LoginLogFilter) error {
	db := p.fmtFilter(ctx, f)
	err := db.Delete(&SysLoginLog{}).Error
	return stores.ErrFmt(err)
}
func (p LoginLogRepo) Delete(ctx context.Context, id int64) error {
	err := p.db.WithContext(ctx).Where("id = ?", id).Delete(&SysLoginLog{}).Error
	return stores.ErrFmt(err)
}

func (p LoginLogRepo) FindOne(ctx context.Context, id int64) (*SysLoginLog, error) {
	var result SysLoginLog
	err := p.db.WithContext(ctx).Where("id = ?", id).First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}
