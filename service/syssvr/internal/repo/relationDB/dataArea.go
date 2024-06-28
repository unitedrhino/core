package relationDB

import (
	"context"
	"gitee.com/i-Things/share/def"
	"gitee.com/i-Things/share/domain/userDataAuth"
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

type DataAreaRepo struct {
	db *gorm.DB
}

func NewDataAreaRepo(in any) *DataAreaRepo {
	return &DataAreaRepo{db: stores.GetCommonConn(in)}
}

type DataAreaFilter struct {
	ProjectID  int64
	AreaIDs    []int64
	Targets    []*Target
	TargetID   int64
	TargetType def.TargetType
	AuthType   def.AuthType
}

func (p DataAreaRepo) fmtFilter(ctx context.Context, f DataAreaFilter) *gorm.DB {
	db := p.db.WithContext(ctx)
	if len(f.AreaIDs) > 0 {
		db = db.Where("area_id in ?", f.AreaIDs)
	}
	if f.TargetType != "" {
		db = db.Where("target_type in ?", f.TargetType)
	}
	if len(f.Targets) != 0 {
		scope := func(db *gorm.DB) *gorm.DB {
			for i, d := range f.Targets {
				if i == 0 {
					db = db.Where("target_id = ? and target_type = ?", d.ID, d.Type)
					continue
				}
				db = db.Or("target_id = ? and target_type = ?", d.ID, d.Type)
			}
			return db
		}
		db = db.Where(scope(db))
	}
	if f.TargetID != 0 {
		db = db.Where("target_id = ? ", f.TargetID)
	}
	//if f.UserID != 0 {
	//	db = db.Where("user_id= ?", f.UserID)
	//}
	if f.ProjectID != 0 {
		db = db.Where("project_id= ?", f.ProjectID)
		//ctxs.SetMetaProjectID(ctx, f.ProjectID) //指定项目id的时候需要清除项目id
	}
	return db
}

func (g DataAreaRepo) Insert(ctx context.Context, data *SysDataArea) error {
	result := g.db.WithContext(ctx).Create(data)
	return stores.ErrFmt(result.Error)
}

func (g DataAreaRepo) FindOneByFilter(ctx context.Context, f DataAreaFilter) (*SysDataArea, error) {
	var result SysDataArea
	db := g.fmtFilter(ctx, f)
	err := db.First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}
func (p DataAreaRepo) FindByFilter(ctx context.Context, f DataAreaFilter, page *stores.PageInfo) ([]*SysDataArea, error) {
	var results []*SysDataArea
	db := p.fmtFilter(ctx, f).Model(&SysDataArea{})
	db = page.ToGorm(db)
	err := db.Find(&results).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return results, nil
}

func (p DataAreaRepo) CountByFilter(ctx context.Context, f DataAreaFilter) (size int64, err error) {
	db := p.fmtFilter(ctx, f).Model(&SysDataArea{})
	err = db.Count(&size).Error
	return size, stores.ErrFmt(err)
}

func (g DataAreaRepo) Update(ctx context.Context, data *SysDataArea) error {
	err := g.db.WithContext(ctx).Where("id = ?", data.ID).Save(data).Error
	return stores.ErrFmt(err)
}

func (g DataAreaRepo) DeleteByFilter(ctx context.Context, f DataAreaFilter) error {
	db := g.fmtFilter(ctx, f)
	err := db.Delete(&SysDataArea{}).Error
	return stores.ErrFmt(err)
}

func (g DataAreaRepo) Delete(ctx context.Context, id int64) error {
	err := g.db.WithContext(ctx).Where("id = ?", id).Delete(&SysDataArea{}).Error
	return stores.ErrFmt(err)
}
func (g DataAreaRepo) FindOne(ctx context.Context, id int64) (*SysDataArea, error) {
	var result SysDataArea
	err := g.db.WithContext(ctx).Where("id = ?", id).First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}

// 批量插入 LightStrategyDevice 记录
func (m DataAreaRepo) MultiInsert(ctx context.Context, data []*SysDataArea) error {
	err := m.db.WithContext(ctx).Clauses(clause.OnConflict{UpdateAll: true}).Model(&SysDataArea{}).Create(data).Error
	return stores.ErrFmt(err)
}
func (g DataAreaRepo) MultiUpdate(ctx context.Context, target *Target, projectID int64, areas []*userDataAuth.Area) error {
	var datas []*SysDataArea
	for _, v := range areas {
		datas = append(datas, &SysDataArea{
			TargetID:   target.ID,
			TargetType: target.Type,
			ProjectID:  stores.ProjectID(projectID),
			AreaID:     v.AreaID,
			AuthType:   v.AuthType,
		})
	}
	err := g.db.Transaction(func(tx *gorm.DB) error {
		rm := NewDataAreaRepo(tx)
		err := rm.DeleteByFilter(ctx, DataAreaFilter{Targets: []*Target{target}, ProjectID: projectID})
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
