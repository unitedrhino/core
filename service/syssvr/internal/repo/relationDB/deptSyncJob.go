package relationDB

import (
	"context"
	"gitee.com/unitedrhino/core/service/syssvr/domain/dept"
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

type DeptSyncJobRepo struct {
	db *gorm.DB
}

func NewDeptSyncJobRepo(in any) *DeptSyncJobRepo {
	return &DeptSyncJobRepo{db: stores.GetCommonConn(in)}
}

type DeptSyncJobFilter struct {
	Direction dept.SyncDirection
	SyncMode  dept.SyncMode
	SyncModes []dept.SyncMode
	ThirdType def.AppSubType
}

func (p DeptSyncJobRepo) fmtFilter(ctx context.Context, f DeptSyncJobFilter) *gorm.DB {
	db := p.db.WithContext(ctx)
	if f.Direction != 0 {
		db = db.Where("direction=?", f.Direction)
	}
	if f.SyncMode != 0 {
		db = db.Where("sync_mode=?", f.SyncMode)
	}
	if len(f.SyncModes) != 0 {
		db = db.Where("sync_modes in ?", f.SyncModes)
	}
	if f.ThirdType != "" {
		db = db.Where("third_type=?", f.ThirdType)
	}
	return db
}

func (p DeptSyncJobRepo) Insert(ctx context.Context, data *SysDeptSyncJob) error {
	result := p.db.WithContext(ctx).Create(data)
	return stores.ErrFmt(result.Error)
}

func (p DeptSyncJobRepo) FindOneByFilter(ctx context.Context, f DeptSyncJobFilter) (*SysDeptSyncJob, error) {
	var result SysDeptSyncJob
	db := p.fmtFilter(ctx, f)
	err := db.First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}
func (p DeptSyncJobRepo) FindByFilter(ctx context.Context, f DeptSyncJobFilter, page *stores.PageInfo) ([]*SysDeptSyncJob, error) {
	var results []*SysDeptSyncJob
	db := p.fmtFilter(ctx, f).Model(&SysDeptSyncJob{})
	db = page.ToGorm(db)
	err := db.Find(&results).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return results, nil
}

func (p DeptSyncJobRepo) CountByFilter(ctx context.Context, f DeptSyncJobFilter) (size int64, err error) {
	db := p.fmtFilter(ctx, f).Model(&SysDeptSyncJob{})
	err = db.Count(&size).Error
	return size, stores.ErrFmt(err)
}

func (p DeptSyncJobRepo) Update(ctx context.Context, data *SysDeptSyncJob) error {
	err := p.db.WithContext(ctx).Where("id = ?", data.ID).Save(data).Error
	return stores.ErrFmt(err)
}

func (p DeptSyncJobRepo) DeleteByFilter(ctx context.Context, f DeptSyncJobFilter) error {
	db := p.fmtFilter(ctx, f)
	err := db.Delete(&SysDeptSyncJob{}).Error
	return stores.ErrFmt(err)
}

func (p DeptSyncJobRepo) Delete(ctx context.Context, id int64) error {
	err := p.db.WithContext(ctx).Where("id = ?", id).Delete(&SysDeptSyncJob{}).Error
	return stores.ErrFmt(err)
}
func (p DeptSyncJobRepo) FindOne(ctx context.Context, id int64) (*SysDeptSyncJob, error) {
	var result SysDeptSyncJob
	err := p.db.WithContext(ctx).Where("id = ?", id).First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}

// 批量插入 LightStrategyDevice 记录
func (p DeptSyncJobRepo) MultiInsert(ctx context.Context, data []*SysDeptSyncJob) error {
	err := p.db.WithContext(ctx).Clauses(clause.OnConflict{UpdateAll: true}).Model(&SysDeptSyncJob{}).Create(data).Error
	return stores.ErrFmt(err)
}

func (d DeptSyncJobRepo) UpdateWithField(ctx context.Context, f DeptSyncJobFilter, updates map[string]any) error {
	db := d.fmtFilter(ctx, f)
	err := db.Model(&SysDeptSyncJob{}).Updates(updates).Error
	return stores.ErrFmt(err)
}
