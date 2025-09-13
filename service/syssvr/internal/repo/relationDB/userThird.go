package relationDB

import (
	"context"

	"gitee.com/unitedrhino/share/def"
	"gitee.com/unitedrhino/share/stores"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

/*
这个是参考样例
使用教程:
1. 将example全局替换为模型的表名
2. 完善todo
*/

type UserThirdRepo struct {
	db *gorm.DB
}

func NewUserThirdRepo(in any) *UserThirdRepo {
	return &UserThirdRepo{db: stores.GetCommonConn(in)}
}

type UserThirdFilter struct {
	TenantCode string
	AppType    def.ThirdType
	AppID      string
	UserID     int64
	UnionID    string // 微信union id
	OpenID     string // 钉钉里是UserID
	OpenIDs    []string
	WithUser   bool
}

func (p UserThirdRepo) fmtFilter(ctx context.Context, f UserThirdFilter) *gorm.DB {
	db := p.db.WithContext(ctx)
	// 根据过滤条件构建查询
	if f.TenantCode != "" {
		db = db.Where("tenant_code = ?", f.TenantCode)
	}
	if f.UserID != 0 {
		db = db.Where("user_id = ?", f.UserID)
	}
	if f.AppType != "" {
		db = db.Where("app_type = ?", f.AppType)
	}
	if f.AppID != "" {
		db = db.Where("app_id = ?", f.AppID)
	}
	var isThirdID bool
	thirdOr := db
	if f.OpenID != "" {
		isThirdID = true
		thirdOr = thirdOr.Or("open_id = ?", f.OpenID)
	}
	if len(f.OpenIDs) > 0 {
		isThirdID = true
		thirdOr = thirdOr.Or("open_id in ?", f.OpenIDs)
	}
	if f.UnionID != "" {
		isThirdID = true
		thirdOr = thirdOr.Or("union_id = ?", f.UnionID)
	}
	if isThirdID {
		db = db.Where(thirdOr)
	}
	if f.WithUser {
		db = db.Preload("User")
	}
	return db
}

func (p UserThirdRepo) Insert(ctx context.Context, data *SysUserThird) error {
	u := data.User
	data.User = nil
	defer func() {
		data.User = u
	}()
	result := p.db.WithContext(ctx).Create(data)
	return stores.ErrFmt(result.Error)
}

func (p UserThirdRepo) FindOneByFilter(ctx context.Context, f UserThirdFilter) (*SysUserThird, error) {
	var result SysUserThird
	db := p.fmtFilter(ctx, f)
	err := db.First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}
func (p UserThirdRepo) FindByFilter(ctx context.Context, f UserThirdFilter, page *stores.PageInfo) ([]*SysUserThird, error) {
	var results []*SysUserThird
	db := p.fmtFilter(ctx, f).Model(&SysUserThird{})
	db = page.ToGorm(db)
	err := db.Find(&results).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return results, nil
}

func (p UserThirdRepo) CountByFilter(ctx context.Context, f UserThirdFilter) (size int64, err error) {
	db := p.fmtFilter(ctx, f).Model(&SysUserThird{})
	err = db.Count(&size).Error
	return size, stores.ErrFmt(err)
}

//func (p UserThirdRepo) Update(ctx context.Context, data *SysUserThird) error {
//	err := p.db.WithContext(ctx).Where("id = ?", data.ID).Save(data).Error
//	return stores.ErrFmt(err)
//}

func (p UserThirdRepo) DeleteByFilter(ctx context.Context, f UserThirdFilter) error {
	db := p.fmtFilter(ctx, f)
	err := db.Delete(&SysUserThird{}).Error
	return stores.ErrFmt(err)
}

func (p UserThirdRepo) Delete(ctx context.Context, id int64) error {
	err := p.db.WithContext(ctx).Where("id = ?", id).Delete(&SysUserThird{}).Error
	return stores.ErrFmt(err)
}
func (p UserThirdRepo) FindOne(ctx context.Context, id int64) (*SysUserThird, error) {
	var result SysUserThird
	err := p.db.WithContext(ctx).Where("id = ?", id).First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}

// 批量插入 LightStrategyDevice 记录
func (p UserThirdRepo) MultiInsert(ctx context.Context, data []*SysUserThird) error {
	err := p.db.WithContext(ctx).Clauses(clause.OnConflict{UpdateAll: true}).Model(&SysUserThird{}).Create(data).Error
	return stores.ErrFmt(err)
}

func (p UserThirdRepo) UpdateWithField(ctx context.Context, f UserThirdFilter, updates map[string]any) error {
	db := p.fmtFilter(ctx, f)
	err := db.Model(&SysUserThird{}).Updates(updates).Error
	return stores.ErrFmt(err)
}
