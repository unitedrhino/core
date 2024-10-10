package relationDB

import (
	"context"
	"gitee.com/unitedrhino/share/domain/ops"
	"gitee.com/unitedrhino/share/stores"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

/*
这个是参考样例
使用教程:
1. 将OpsFeedback全局替换为模型的表名
2. 完善todo
*/

type OpsFeedbackRepo struct {
	db *gorm.DB
}

func NewOpsFeedbackRepo(in any) *OpsFeedbackRepo {
	return &OpsFeedbackRepo{db: stores.GetCommonConn(in)}
}

type OpsFeedbackFilter struct {
	TenantCode string
	ProjectID  int64
	Type       string
	Status     ops.WorkOrderStatus
}

func (p OpsFeedbackRepo) fmtFilter(ctx context.Context, f OpsFeedbackFilter) *gorm.DB {
	db := p.db.WithContext(ctx)
	if f.TenantCode != "" {
		db = db.Where("tenant_code = ?", f.TenantCode)
	}
	if f.Status != 0 {
		db = db.Where("status = ?", f.Status)
	}
	if f.ProjectID != 0 {
		db = db.Where("project_id = ?", f.ProjectID)
	}
	if f.Type != "" {
		return db.Where("type = ?", f.Type)
	}
	return db
}

func (p OpsFeedbackRepo) Insert(ctx context.Context, data *SysOpsFeedback) error {
	result := p.db.WithContext(ctx).Create(data)
	return stores.ErrFmt(result.Error)
}

func (p OpsFeedbackRepo) FindOneByFilter(ctx context.Context, f OpsFeedbackFilter) (*SysOpsFeedback, error) {
	var result SysOpsFeedback
	db := p.fmtFilter(ctx, f)
	err := db.First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}
func (p OpsFeedbackRepo) FindByFilter(ctx context.Context, f OpsFeedbackFilter, page *stores.PageInfo) ([]*SysOpsFeedback, error) {
	var results []*SysOpsFeedback
	db := p.fmtFilter(ctx, f).Model(&SysOpsFeedback{})
	db = page.ToGorm(db)
	err := db.Find(&results).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return results, nil
}

func (p OpsFeedbackRepo) CountByFilter(ctx context.Context, f OpsFeedbackFilter) (size int64, err error) {
	db := p.fmtFilter(ctx, f).Model(&SysOpsFeedback{})
	err = db.Count(&size).Error
	return size, stores.ErrFmt(err)
}

func (p OpsFeedbackRepo) Update(ctx context.Context, data *SysOpsFeedback) error {
	err := p.db.WithContext(ctx).Where("id = ?", data.ID).Save(data).Error
	return stores.ErrFmt(err)
}

func (p OpsFeedbackRepo) DeleteByFilter(ctx context.Context, f OpsFeedbackFilter) error {
	db := p.fmtFilter(ctx, f)
	err := db.Delete(&SysOpsFeedback{}).Error
	return stores.ErrFmt(err)
}

func (p OpsFeedbackRepo) Delete(ctx context.Context, id int64) error {
	err := p.db.WithContext(ctx).Where("id = ?", id).Delete(&SysOpsFeedback{}).Error
	return stores.ErrFmt(err)
}
func (p OpsFeedbackRepo) FindOne(ctx context.Context, id int64) (*SysOpsFeedback, error) {
	var result SysOpsFeedback
	err := p.db.WithContext(ctx).Where("id = ?", id).First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}

// 批量插入 LightStrategyDevice 记录
func (p OpsFeedbackRepo) MultiInsert(ctx context.Context, data []*SysOpsFeedback) error {
	err := p.db.WithContext(ctx).Clauses(clause.OnConflict{UpdateAll: true}).Model(&SysOpsFeedback{}).Create(data).Error
	return stores.ErrFmt(err)
}
