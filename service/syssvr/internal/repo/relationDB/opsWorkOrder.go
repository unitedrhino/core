package relationDB

import (
	"context"
	"gitee.com/i-Things/share/domain/ops"
	"gitee.com/i-Things/share/stores"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
)

/*
这个是参考样例
使用教程:
1. 将example全局替换为模型的表名
2. 完善todo
*/

type OpsWorkOrderRepo struct {
	db *gorm.DB
}

func NewOpsWorkOrderRepo(in any) *OpsWorkOrderRepo {
	return &OpsWorkOrderRepo{db: stores.GetCommonConn(in)}
}

type OpsWorkOrderFilter struct {
	Status    ops.WorkOrderStatus
	Type      string
	StartTime time.Time
	EndTime   time.Time
	AreaID    int64
	Number    string
}

func (p OpsWorkOrderRepo) fmtFilter(ctx context.Context, f OpsWorkOrderFilter) *gorm.DB {
	db := p.db.WithContext(ctx)
	if f.Type != "" {
		db = db.Where("type = ?", f.Type)
	}
	if f.Number != "" {
		db = db.Where("number = ?", f.Number)
	}
	if f.AreaID != 0 {
		db = db.Where("area_id=?", f.AreaID)
	}
	if f.Status != 0 {
		db = db.Where("status = ?", f.Status)
	}
	if !f.StartTime.IsZero() {
		db = db.Where("created_time >= ?", f.StartTime)
	}
	if !f.EndTime.IsZero() {
		db = db.Where("created_time <= ?", f.EndTime)
	}
	return db
}

func (p OpsWorkOrderRepo) Insert(ctx context.Context, data *SysOpsWorkOrder) error {
	result := p.db.WithContext(ctx).Create(data)
	return stores.ErrFmt(result.Error)
}

func (p OpsWorkOrderRepo) FindOneByFilter(ctx context.Context, f OpsWorkOrderFilter) (*SysOpsWorkOrder, error) {
	var result SysOpsWorkOrder
	db := p.fmtFilter(ctx, f)
	err := db.First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}
func (p OpsWorkOrderRepo) FindByFilter(ctx context.Context, f OpsWorkOrderFilter, page *stores.PageInfo) ([]*SysOpsWorkOrder, error) {
	var results []*SysOpsWorkOrder
	db := p.fmtFilter(ctx, f).Model(&SysOpsWorkOrder{})
	db = page.ToGorm(db)
	err := db.Find(&results).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return results, nil
}

func (p OpsWorkOrderRepo) CountByFilter(ctx context.Context, f OpsWorkOrderFilter) (size int64, err error) {
	db := p.fmtFilter(ctx, f).Model(&SysOpsWorkOrder{})
	err = db.Count(&size).Error
	return size, stores.ErrFmt(err)
}

func (p OpsWorkOrderRepo) Update(ctx context.Context, data *SysOpsWorkOrder) error {
	err := p.db.WithContext(ctx).Where("id = ?", data.ID).Save(data).Error
	return stores.ErrFmt(err)
}

func (p OpsWorkOrderRepo) DeleteByFilter(ctx context.Context, f OpsWorkOrderFilter) error {
	db := p.fmtFilter(ctx, f)
	err := db.Delete(&SysOpsWorkOrder{}).Error
	return stores.ErrFmt(err)
}

func (p OpsWorkOrderRepo) Delete(ctx context.Context, id int64) error {
	err := p.db.WithContext(ctx).Where("id = ?", id).Delete(&SysOpsWorkOrder{}).Error
	return stores.ErrFmt(err)
}
func (p OpsWorkOrderRepo) FindOne(ctx context.Context, id int64) (*SysOpsWorkOrder, error) {
	var result SysOpsWorkOrder
	err := p.db.WithContext(ctx).Where("id = ?", id).First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}

// 批量插入 LightStrategyDevice 记录
func (p OpsWorkOrderRepo) MultiInsert(ctx context.Context, data []*SysOpsWorkOrder) error {
	err := p.db.WithContext(ctx).Clauses(clause.OnConflict{UpdateAll: true}).Model(&SysOpsWorkOrder{}).Create(data).Error
	return stores.ErrFmt(err)
}
