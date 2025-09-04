package relationDB

import (
	"context"

	"gitee.com/unitedrhino/share/def"
	"gitee.com/unitedrhino/share/stores"
	"gorm.io/gorm"
)

type UserInfoRepo struct {
	db *gorm.DB
}

func NewUserInfoRepo(in any) *UserInfoRepo {
	return &UserInfoRepo{db: stores.GetCommonConn(in)}
}

type UserInfoFilter struct {
	UserIDs      []int64
	UserNames    []string
	UserName     string
	NickName     string
	Phone        string
	Phones       []string
	Email        string
	Emails       []string
	Accounts     []string //账号查询 非模糊查询
	WithRoles    bool
	WithTenant   bool
	TenantStatus def.Bool
	UpdatedTime  *stores.Cmp
}

func (p UserInfoRepo) accountsFilter(db *gorm.DB, accounts []string) *gorm.DB {
	db = db.Where(db.Or("user_name in ?", accounts).
		Or("email in ?", accounts).
		Or("phone in ?", accounts))
	return db
}

func (p UserInfoRepo) fmtFilter(ctx context.Context, f UserInfoFilter) *gorm.DB {
	db := p.db.WithContext(ctx)
	db = f.UpdatedTime.Where(db, "updated_time")

	if f.WithRoles {
		db = db.Preload("Tenants.Roles").Preload("Tenants.Roles.Role")
	}
	if f.WithTenant {
		if f.TenantStatus != 0 {
			db = db.Preload("Tenants")
		} else {
			db = db.Preload("Tenants", "status = ?", f.TenantStatus)
		}
		db = db.Preload("Tenants.TenantInfo").Preload("Tenants.TenantConfig")
	}
	if len(f.UserIDs) != 0 {
		db = db.Where("user_id in?", f.UserIDs)
	}
	if len(f.UserNames) != 0 {
		db = db.Where("user_name in ?", f.UserNames)
	}
	if f.NickName != "" {
		db = db.Where("nick_name like ?", "%"+f.NickName+"%")
	}
	if len(f.Accounts) != 0 {
		db = p.accountsFilter(db, f.Accounts)
	}
	if f.UserName != "" {
		db = db.Where("user_name like ?", "%"+f.UserName+"%")
	}
	if f.Phone != "" {
		db = db.Where("phone like ?", "%"+f.Phone+"%")
	}
	if len(f.Phones) != 0 {
		db = db.Where("phone in ?", f.Phones)
	}
	if f.Email != "" {
		db = db.Where("email like ?", "%"+f.Email+"%")
	}
	if len(f.Emails) != 0 {
		db = db.Where("email in ?", f.Emails)
	}

	return db
}

func (p UserInfoRepo) Insert(ctx context.Context, data *SysUserInfo) error {
	result := p.db.WithContext(ctx).Create(data)
	return stores.ErrFmt(result.Error)
}

func (p UserInfoRepo) FindOneByFilter(ctx context.Context, f UserInfoFilter) (*SysUserInfo, error) {
	var result SysUserInfo
	db := p.fmtFilter(ctx, f)
	err := db.First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}
func (p UserInfoRepo) FindByFilter(ctx context.Context, f UserInfoFilter, page *stores.PageInfo) ([]*SysUserInfo, error) {
	var results []*SysUserInfo
	db := p.fmtFilter(ctx, f).Model(&SysUserInfo{})
	db = page.ToGorm(db)
	err := db.Find(&results).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return results, nil
}

func (p UserInfoRepo) CountByFilter(ctx context.Context, f UserInfoFilter) (size int64, err error) {
	db := p.fmtFilter(ctx, f).Model(&SysUserInfo{})
	err = db.Count(&size).Error
	return size, stores.ErrFmt(err)
}

func (p UserInfoRepo) Update(ctx context.Context, data *SysUserInfo) error {
	err := p.db.WithContext(ctx).Where("user_id = ?", data.UserID).Save(data).Error
	return stores.ErrFmt(err)
}
func (d UserInfoRepo) UpdateWithField(ctx context.Context, f UserInfoFilter, updates map[string]any) error {
	db := d.fmtFilter(ctx, f)
	err := db.Model(&SysUserInfo{}).Updates(updates).Error
	return stores.ErrFmt(err)
}

func (p UserInfoRepo) UpdateDeviceCount(ctx context.Context, userID int64) error {
	subQuery1 := p.db.Model(&SysProjectInfo{}).Select("sum(device_count)").Where("admin_user_id=?", userID)
	err := p.db.WithContext(ctx).Model(&SysUserInfo{}).Where("user_id = ?", userID).
		Update("device_count", subQuery1).Error
	return stores.ErrFmt(err)
}

func (p UserInfoRepo) DeleteByFilter(ctx context.Context, f UserInfoFilter) error {
	db := p.fmtFilter(ctx, f)
	err := db.Delete(&SysUserInfo{}).Error
	return stores.ErrFmt(err)
}

func (p UserInfoRepo) Delete(ctx context.Context, userID int64) error {
	err := p.db.WithContext(ctx).Where("user_id = ?", userID).Delete(&SysUserInfo{}).Error
	return stores.ErrFmt(err)
}
func (p UserInfoRepo) FindOne(ctx context.Context, userID int64) (*SysUserInfo, error) {
	var result SysUserInfo
	err := p.db.WithContext(ctx).Preload("Tenants").Where("user_id = ?", userID).First(&result).Error
	return &result, stores.ErrFmt(err)
}

func (p UserInfoRepo) FindUserCore(ctx context.Context, f UserInfoFilter) (ret []*SysUserInfo, err error) {
	var results []*SysUserInfo
	db := p.fmtFilter(ctx, f).Model(&SysUserInfo{})
	err = db.Select("user_id,user_name,email,phone,wechat_union_id,wechat_open_id,ding_talk_user_id").Find(&results).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return results, nil
}
