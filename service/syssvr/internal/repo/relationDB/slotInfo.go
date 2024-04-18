package relationDB

import (
	"context"
	"gitee.com/i-Things/share/def"
	"gitee.com/i-Things/share/stores"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

/*
这个是参考样例
使用教程:
1. 将SlotInfo全局替换为模型的表名
2. 完善todo
*/

type SlotInfoRepo struct {
	db *gorm.DB
}

func NewSlotInfoRepo(in any) *SlotInfoRepo {
	return &SlotInfoRepo{db: stores.GetCommonConn(in)}
}

type SlotInfoFilter struct {
	SlotCode string
	Code     string
	SubCode  string
}

func (p SlotInfoRepo) fmtFilter(ctx context.Context, f SlotInfoFilter) *gorm.DB {
	db := p.db.WithContext(ctx)
	if f.Code != "" {
		db = db.Where("code = ?", f.Code)
	}
	if f.SubCode != "" {
		db = db.Where("sub_code = ?", f.SubCode)
	}
	if f.SlotCode != "" {
		db = db.Where("slot_code = ?", f.SlotCode)
	}
	return db
}

func (p SlotInfoRepo) Insert(ctx context.Context, data *SysSlotInfo) error {
	result := p.db.WithContext(ctx).Create(data)
	return stores.ErrFmt(result.Error)
}

func (p SlotInfoRepo) FindOneByFilter(ctx context.Context, f SlotInfoFilter) (*SysSlotInfo, error) {
	var result SysSlotInfo
	db := p.fmtFilter(ctx, f)
	err := db.First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}
func (p SlotInfoRepo) FindByFilter(ctx context.Context, f SlotInfoFilter, page *def.PageInfo) ([]*SysSlotInfo, error) {
	var results []*SysSlotInfo
	db := p.fmtFilter(ctx, f).Model(&SysSlotInfo{})
	db = page.ToGorm(db)
	err := db.Find(&results).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return results, nil
}

func (p SlotInfoRepo) CountByFilter(ctx context.Context, f SlotInfoFilter) (size int64, err error) {
	db := p.fmtFilter(ctx, f).Model(&SysSlotInfo{})
	err = db.Count(&size).Error
	return size, stores.ErrFmt(err)
}

func (p SlotInfoRepo) Update(ctx context.Context, data *SysSlotInfo) error {
	err := p.db.WithContext(ctx).Where("id = ?", data.ID).Save(data).Error
	return stores.ErrFmt(err)
}

func (p SlotInfoRepo) DeleteByFilter(ctx context.Context, f SlotInfoFilter) error {
	db := p.fmtFilter(ctx, f)
	err := db.Delete(&SysSlotInfo{}).Error
	return stores.ErrFmt(err)
}

func (p SlotInfoRepo) Delete(ctx context.Context, id int64) error {
	err := p.db.WithContext(ctx).Where("id = ?", id).Delete(&SysSlotInfo{}).Error
	return stores.ErrFmt(err)
}
func (p SlotInfoRepo) FindOne(ctx context.Context, id int64) (*SysSlotInfo, error) {
	var result SysSlotInfo
	err := p.db.WithContext(ctx).Where("id = ?", id).First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}

// 批量插入 LightStrategyDevice 记录
func (p SlotInfoRepo) MultiInsert(ctx context.Context, data []*SysSlotInfo) error {
	err := p.db.WithContext(ctx).Clauses(clause.OnConflict{UpdateAll: true}).Model(&SysSlotInfo{}).Create(data).Error
	return stores.ErrFmt(err)
}
