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
1. 将example全局替换为模型的表名
2. 完善todo
*/

type StatisticsInfoRepo struct {
	db *gorm.DB
}

func NewStatisticsInfoRepo(in any) *StatisticsInfoRepo {
	return &StatisticsInfoRepo{db: stores.GetCommonConn(in)}
}

type StatisticsInfoFilter struct {
	Code string
}

func (p StatisticsInfoRepo) fmtFilter(ctx context.Context, f StatisticsInfoFilter) *gorm.DB {
	db := p.db.WithContext(ctx)
	if f.Code != "" {
		db = db.Where("code = ?", f.Code)
	}
	return db
}

func (p StatisticsInfoRepo) Insert(ctx context.Context, data *DataStatisticsInfo) error {
	result := p.db.WithContext(ctx).Create(data)
	return stores.ErrFmt(result.Error)
}

func (p StatisticsInfoRepo) FindOneByFilter(ctx context.Context, f StatisticsInfoFilter) (*DataStatisticsInfo, error) {
	var result DataStatisticsInfo
	db := p.fmtFilter(ctx, f)
	err := db.First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}
func (p StatisticsInfoRepo) FindByFilter(ctx context.Context, f StatisticsInfoFilter, page *def.PageInfo) ([]*DataStatisticsInfo, error) {
	var results []*DataStatisticsInfo
	db := p.fmtFilter(ctx, f).Model(&DataStatisticsInfo{})
	db = page.ToGorm(db)
	err := db.Find(&results).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return results, nil
}

func (p StatisticsInfoRepo) CountByFilter(ctx context.Context, f StatisticsInfoFilter) (size int64, err error) {
	db := p.fmtFilter(ctx, f).Model(&DataStatisticsInfo{})
	err = db.Count(&size).Error
	return size, stores.ErrFmt(err)
}

func (p StatisticsInfoRepo) Update(ctx context.Context, data *DataStatisticsInfo) error {
	err := p.db.WithContext(ctx).Where("id = ?", data.ID).Save(data).Error
	return stores.ErrFmt(err)
}

func (p StatisticsInfoRepo) DeleteByFilter(ctx context.Context, f StatisticsInfoFilter) error {
	db := p.fmtFilter(ctx, f)
	err := db.Delete(&DataStatisticsInfo{}).Error
	return stores.ErrFmt(err)
}

func (p StatisticsInfoRepo) Delete(ctx context.Context, id int64) error {
	err := p.db.WithContext(ctx).Where("id = ?", id).Delete(&DataStatisticsInfo{}).Error
	return stores.ErrFmt(err)
}
func (p StatisticsInfoRepo) FindOne(ctx context.Context, id int64) (*DataStatisticsInfo, error) {
	var result DataStatisticsInfo
	err := p.db.WithContext(ctx).Where("id = ?", id).First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}

// 批量插入 LightStrategyDevice 记录
func (p StatisticsInfoRepo) MultiInsert(ctx context.Context, data []*DataStatisticsInfo) error {
	err := p.db.WithContext(ctx).Clauses(clause.OnConflict{UpdateAll: true}).Model(&DataStatisticsInfo{}).Create(data).Error
	return stores.ErrFmt(err)
}
