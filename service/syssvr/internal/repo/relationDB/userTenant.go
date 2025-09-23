package relationDB

import (
	"context"

	"gitee.com/unitedrhino/share/def"
	"gitee.com/unitedrhino/share/stores"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type UserTenantRepo struct {
	db *gorm.DB
}

func NewUserTenantRepo(in any) *UserTenantRepo {
	return &UserTenantRepo{db: stores.GetCommonConn(in)}
}

type UserTenantFilter struct {
	UserNames []string
	UserName  string
	NickName  string
	Phone     string
	Phones    []string
	Email     string
	Emails    []string
	Accounts  []string //账号查询 非模糊查询

	UserIDs         []int64
	HasAccessAreas  []int64
	TenantCode      string
	WechatOpenIDs   []string
	WechatUnionID   string
	WechatOpenID    string
	DingTalkUserID  string
	DingTalkUserIDs []string
	DingTalkUnionID string
	WithRoles       bool
	WithTenant      bool
	WithUser        bool
	RoleCode        string
	DeptID          int64
	UpdatedTime     *stores.Cmp
}

func (p UserTenantRepo) accountsFilter(db *gorm.DB, accounts []string) *gorm.DB {
	db = db.Where(db.Or("user_name in ?", accounts).
		Or("email in ?", accounts).
		Or("phone in ?", accounts))
	return db
}

func (p UserTenantRepo) fmtFilter(ctx context.Context, f UserTenantFilter) *gorm.DB {
	db := p.db.WithContext(ctx)
	db = f.UpdatedTime.Where(db, "updated_time")
	if f.HasAccessAreas != nil {
		if len(f.HasAccessAreas) == 0 {
			subQuery := p.db.Model(&SysDataArea{}).Select("target_id").Where("target_type=?", def.TargetUser)
			db = db.Where("user_id in (?)", subQuery)
		} else {
			subQuery := p.db.Model(&SysDataArea{}).Select("target_id").Where("target_type=? and area_id in ?", def.TargetUser, f.HasAccessAreas)
			db = db.Where("user_id in (?)", subQuery)
		}
	}
	if f.DeptID > 0 {
		subQuery := p.db.Model(&SysDeptUser{}).Select("user_id").Where("dept_id=?", f.DeptID)
		db = db.Where("user_id in (?)", subQuery)
	}
	if f.WithUser {
		db = db.Preload("User")
	}

	if f.WithRoles {
		db = db.Preload("Roles.Role")
	}
	if f.WithTenant {
		db = db.Preload("TenantInfo").Preload("TenantConfig")
	}
	if len(f.UserIDs) != 0 {
		db = db.Where("user_id in?", f.UserIDs)
	}

	dingOr := db
	var isDing bool
	if f.DingTalkUserID != "" {
		isDing = true
		dingOr = dingOr.Or("ding_talk_user_id = ?", f.DingTalkUserID)
	}
	if len(f.DingTalkUserIDs) != 0 {
		isDing = true
		dingOr = dingOr.Or("ding_talk_user_id in ?", f.DingTalkUserIDs)
	}
	if f.DingTalkUnionID != "" {
		isDing = true
		dingOr = dingOr.Or("ding_talk_union_id = ?", f.DingTalkUnionID)
	}
	if isDing {
		db = db.Where(dingOr)
	}
	wechatOr := db
	var isWechat bool
	if f.WechatUnionID != "" {
		isWechat = true
		wechatOr = wechatOr.Or("wechat_union_id = ?", f.WechatUnionID)
	}
	if f.WechatOpenID != "" {
		isWechat = true
		wechatOr = wechatOr.Or("wechat_open_id = ?", f.WechatOpenID)
	}
	if isWechat {
		db = db.Where(wechatOr)
	}
	if f.TenantCode != "" {
		db = db.Where("tenant_code =?", f.TenantCode)
	}
	if f.RoleCode != "" {
		subQuery1 := p.db.Model(&SysRoleInfo{}).Select("id").Where("code=?", f.RoleCode)
		subQuery2 := p.db.Model(&SysUserRole{}).Select("user_id").Where("role_id in (?)", subQuery1)
		db = db.Where("user_id in (?)", subQuery2)
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

func (p UserTenantRepo) Insert(ctx context.Context, data *SysUserTenant) error {
	result := p.db.WithContext(ctx).Create(data)
	return stores.ErrFmt(result.Error)
}

// 批量插入 LightStrategyDevice 记录
func (p UserTenantRepo) MultiInsert(ctx context.Context, data []*SysUserTenant) error {
	err := p.db.WithContext(ctx).Clauses(clause.OnConflict{DoNothing: true}).Model(&SysUserTenant{}).Create(data).Error
	return stores.ErrFmt(err)
}

func (p UserTenantRepo) FindOneByFilter(ctx context.Context, f UserTenantFilter) (*SysUserTenant, error) {
	var result SysUserTenant
	db := p.fmtFilter(ctx, f)
	err := db.First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}
func (p UserTenantRepo) FindByFilter(ctx context.Context, f UserTenantFilter, page *stores.PageInfo) ([]*SysUserTenant, error) {
	var results []*SysUserTenant
	db := p.fmtFilter(ctx, f).Model(&SysUserTenant{})
	db = page.ToGorm(db)
	err := db.Find(&results).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return results, nil
}

func (p UserTenantRepo) CountByFilter(ctx context.Context, f UserTenantFilter) (size int64, err error) {
	db := p.fmtFilter(ctx, f).Model(&SysUserTenant{})
	err = db.Count(&size).Error
	return size, stores.ErrFmt(err)
}

func (p UserTenantRepo) Update(ctx context.Context, data *SysUserTenant) error {
	err := p.db.WithContext(ctx).Where("user_id = ?", data.UserID).Save(data).Error
	return stores.ErrFmt(err)
}
func (d UserTenantRepo) UpdateWithField(ctx context.Context, f UserTenantFilter, updates map[string]any) error {
	db := d.fmtFilter(ctx, f)
	err := db.Model(&SysUserTenant{}).Updates(updates).Error
	return stores.ErrFmt(err)
}

func (p UserTenantRepo) UpdateDeviceCount(ctx context.Context, userID int64) error {
	subQuery1 := p.db.Model(&SysProjectInfo{}).Select("sum(device_count)").Where("admin_user_id=?", userID)
	err := p.db.WithContext(ctx).Model(&SysUserTenant{}).Where("user_id = ?", userID).
		Update("device_count", subQuery1).Error
	return stores.ErrFmt(err)
}

func (p UserTenantRepo) DeleteByFilter(ctx context.Context, f UserTenantFilter) error {
	db := p.fmtFilter(ctx, f)
	err := db.Delete(&SysUserTenant{}).Error
	return stores.ErrFmt(err)
}

func (p UserTenantRepo) Delete(ctx context.Context, userID int64) error {
	err := p.db.WithContext(ctx).Where("user_id = ?", userID).Delete(&SysUserTenant{}).Error
	return stores.ErrFmt(err)
}
func (p UserTenantRepo) FindOne(ctx context.Context, userID int64) (*SysUserTenant, error) {
	var result SysUserTenant
	err := p.db.WithContext(ctx).Preload("User").Where("user_id = ?", userID).First(&result).Error
	return &result, stores.ErrFmt(err)
}
