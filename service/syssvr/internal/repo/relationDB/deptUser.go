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
1. 将DeptUser全局替换为模型的表名
2. 完善todo
*/

type DeptUserRepo struct {
	db *gorm.DB
}

func NewDeptUserRepo(in any) *DeptUserRepo {
	return &DeptUserRepo{db: stores.GetCommonConn(in)}
}

type DeptUserFilter struct {
	UserID     int64
	DeptID     int64
	DeptIDs    []int64
	DeptIDPath string
	WithUser   bool
}

func (p DeptUserRepo) fmtFilter(ctx context.Context, f DeptUserFilter) *gorm.DB {
	db := p.db.WithContext(ctx)
	if f.UserID != 0 {
		db = db.Where("user_id =?", f.UserID)
	}
	if f.DeptID != 0 {
		db = db.Where("dept_id =?", f.DeptID)
	}
	if f.DeptIDs != nil {
		db = db.Where("dept_id IN ?", f.DeptIDs)
	}
	if f.DeptIDPath != "" {
		db = db.Where("dept_id_path  like ?", f.DeptIDPath+"%")
	}
	if f.WithUser {
		db = db.Preload("User")
	}
	return db
}

func (p DeptUserRepo) Insert(ctx context.Context, data *SysDeptUser) error {
	result := p.db.WithContext(ctx).Create(data)
	return stores.ErrFmt(result.Error)
}

func (p DeptUserRepo) FindOneByFilter(ctx context.Context, f DeptUserFilter) (*SysDeptUser, error) {
	var result SysDeptUser
	db := p.fmtFilter(ctx, f)
	err := db.First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}

func (p DeptUserRepo) FindByFilter(ctx context.Context, f DeptUserFilter, page *stores.PageInfo) ([]*SysDeptUser, error) {
	var results []*SysDeptUser
	db := p.fmtFilter(ctx, f).Model(&SysDeptUser{})
	db = page.ToGorm(db)
	err := db.Find(&results).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return results, nil
}

func (p DeptUserRepo) CountByFilter(ctx context.Context, f DeptUserFilter) (size int64, err error) {
	db := p.fmtFilter(ctx, f).Model(&SysDeptUser{})
	err = db.Count(&size).Error
	return size, stores.ErrFmt(err)
}

func (p DeptUserRepo) Update(ctx context.Context, data *SysDeptUser) error {
	err := p.db.WithContext(ctx).Where("id = ?", data.ID).Save(data).Error
	return stores.ErrFmt(err)
}

func (p DeptUserRepo) DeleteByFilter(ctx context.Context, f DeptUserFilter) error {
	db := p.fmtFilter(ctx, f)
	err := db.Delete(&SysDeptUser{}).Error
	return stores.ErrFmt(err)
}

func (p DeptUserRepo) Delete(ctx context.Context, id int64) error {
	err := p.db.WithContext(ctx).Where("id = ?", id).Delete(&SysDeptUser{}).Error
	return stores.ErrFmt(err)
}
func (p DeptUserRepo) FindOne(ctx context.Context, id int64) (*SysDeptUser, error) {
	var result SysDeptUser
	err := p.db.WithContext(ctx).Where("id = ?", id).First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}

// 批量插入 LightStrategyDevice 记录
func (p DeptUserRepo) MultiInsert(ctx context.Context, data []*SysDeptUser) error {
	err := p.db.WithContext(ctx).Clauses(clause.OnConflict{UpdateAll: true,
		Columns: stores.SetColumnsWithPg(p.db, &SysDeptUser{}, "idx_sys_dept_user_ri_mi")}).Model(&SysDeptUser{}).Create(data).Error
	return stores.ErrFmt(err)
}

func (p DeptUserRepo) MultiUpdate(ctx context.Context, userID int64, roleIDs []int64) error {
	var datas []*SysDeptUser
	for _, v := range roleIDs {
		datas = append(datas, &SysDeptUser{
			DeptID: v,
			UserID: userID,
		})
	}
	err := p.db.Transaction(func(tx *gorm.DB) error {
		rm := NewDeptUserRepo(tx)
		err := rm.DeleteByFilter(ctx, DeptUserFilter{UserID: userID})
		if err != nil {
			return err
		}
		if len(datas) != 0 {
			err = rm.MultiInsert(ctx, datas)
			if err != nil {
				return err
			}
		}
		return nil
	})
	return stores.ErrFmt(err)
}
