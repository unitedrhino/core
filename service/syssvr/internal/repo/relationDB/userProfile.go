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
1. 将UserProfile全局替换为模型的表名
2. 完善todo
*/

type UserProfileRepo struct {
	db *gorm.DB
}

func NewUserProfileRepo(in any) *UserProfileRepo {
	return &UserProfileRepo{db: stores.GetCommonConn(in)}
}

type UserProfileFilter struct {
	Codes  []string
	Code   string
	UserID int64
}

func (p UserProfileRepo) fmtFilter(ctx context.Context, f UserProfileFilter) *gorm.DB {
	db := p.db.WithContext(ctx)
	if len(f.Codes) != 0 {
		db = db.Where("code in ?", f.Codes)
	}
	if f.Code != "" {
		db = db.Where("code = ?", f.Code)
	}
	db = db.Where("user_id =?", f.UserID)
	return db
}

func (p UserProfileRepo) Insert(ctx context.Context, data *SysUserProfile) error {
	result := p.db.WithContext(ctx).Create(data)
	return stores.ErrFmt(result.Error)
}

func (p UserProfileRepo) FindOneByFilter(ctx context.Context, f UserProfileFilter) (*SysUserProfile, error) {
	var result SysUserProfile
	db := p.fmtFilter(ctx, f)
	err := db.First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}
func (p UserProfileRepo) FindByFilter(ctx context.Context, f UserProfileFilter, page *stores.PageInfo) ([]*SysUserProfile, error) {
	var results []*SysUserProfile
	db := p.fmtFilter(ctx, f).Model(&SysUserProfile{})
	db = page.ToGorm(db)
	err := db.Find(&results).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return results, nil
}

func (p UserProfileRepo) CountByFilter(ctx context.Context, f UserProfileFilter) (size int64, err error) {
	db := p.fmtFilter(ctx, f).Model(&SysUserProfile{})
	err = db.Count(&size).Error
	return size, stores.ErrFmt(err)
}

func (p UserProfileRepo) Update(ctx context.Context, data *SysUserProfile) error {
	err := p.db.WithContext(ctx).Where("user_id = ? and code = ?", data.UserID, data.Code).Save(data).Error
	return stores.ErrFmt(err)
}

func (p UserProfileRepo) DeleteByFilter(ctx context.Context, f UserProfileFilter) error {
	db := p.fmtFilter(ctx, f)
	err := db.Delete(&SysUserProfile{}).Error
	return stores.ErrFmt(err)
}

func (p UserProfileRepo) Delete(ctx context.Context, id int64) error {
	err := p.db.WithContext(ctx).Where("id = ?", id).Delete(&SysUserProfile{}).Error
	return stores.ErrFmt(err)
}
func (p UserProfileRepo) FindOne(ctx context.Context, id int64) (*SysUserProfile, error) {
	var result SysUserProfile
	err := p.db.WithContext(ctx).Where("id = ?", id).First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}

// 批量插入 LightStrategyDevice 记录
func (p UserProfileRepo) MultiInsert(ctx context.Context, data []*SysUserProfile) error {
	err := p.db.WithContext(ctx).Clauses(clause.OnConflict{UpdateAll: true}).Model(&SysUserProfile{}).Create(data).Error
	return stores.ErrFmt(err)
}
