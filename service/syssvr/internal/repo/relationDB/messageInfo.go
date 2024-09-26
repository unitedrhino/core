package relationDB

import (
	"context"
	"fmt"
	"gitee.com/i-Things/share/stores"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

/*
这个是参考样例
使用教程:
1. 将MessageInfo全局替换为模型的表名
2. 完善todo
*/

type MessageInfoRepo struct {
	db *gorm.DB
}

func NewMessageInfoRepo(in any) *MessageInfoRepo {
	return &MessageInfoRepo{db: stores.GetCommonConn(in)}
}

type MessageInfoFilter struct {
	NotifyCode       string
	Group            string
	IsGlobal         int64
	IsDirectNotify   int64 //是否是发送通知消息创建
	NotifyTime       *stores.Cmp
	WithNotifyConfig bool
}

func (p MessageInfoRepo) fmtFilter(ctx context.Context, f MessageInfoFilter) *gorm.DB {
	db := p.db.WithContext(ctx)
	db = f.NotifyTime.Where(db, "notify_time")
	if f.WithNotifyConfig {
		db = db.Preload("NotifyConfig")
	}
	if f.Group != "" {
		db = db.Where(fmt.Sprintf("%s=?", stores.Col("group")), f.Group)
	}
	if f.NotifyCode != "" {
		db = db.Where("notify_code=?", f.NotifyCode)
	}
	if f.IsGlobal != 0 {
		db = db.Where("is_global=?", f.IsGlobal)
	}
	if f.IsDirectNotify != 0 {
		db = db.Where("is_direct_notify=?", f.IsDirectNotify)
	}
	return db
}

func (p MessageInfoRepo) Insert(ctx context.Context, data *SysMessageInfo) error {
	result := p.db.WithContext(ctx).Create(data)
	return stores.ErrFmt(result.Error)
}

func (p MessageInfoRepo) FindOneByFilter(ctx context.Context, f MessageInfoFilter) (*SysMessageInfo, error) {
	var result SysMessageInfo
	db := p.fmtFilter(ctx, f)
	err := db.First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}
func (p MessageInfoRepo) FindByFilter(ctx context.Context, f MessageInfoFilter, page *stores.PageInfo) ([]*SysMessageInfo, error) {
	var results []*SysMessageInfo
	db := p.fmtFilter(ctx, f).Model(&SysMessageInfo{})
	db = page.ToGorm(db)
	err := db.Find(&results).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return results, nil
}

func (p MessageInfoRepo) CountByFilter(ctx context.Context, f MessageInfoFilter) (size int64, err error) {
	db := p.fmtFilter(ctx, f).Model(&SysMessageInfo{})
	err = db.Count(&size).Error
	return size, stores.ErrFmt(err)
}

func (p MessageInfoRepo) Update(ctx context.Context, data *SysMessageInfo) error {
	err := p.db.WithContext(ctx).Where("id = ?", data.ID).Save(data).Error
	return stores.ErrFmt(err)
}

func (p MessageInfoRepo) DeleteByFilter(ctx context.Context, f MessageInfoFilter) error {
	db := p.fmtFilter(ctx, f)
	err := db.Delete(&SysMessageInfo{}).Error
	return stores.ErrFmt(err)
}

func (p MessageInfoRepo) Delete(ctx context.Context, id int64) error {
	err := p.db.WithContext(ctx).Where("id = ?", id).Delete(&SysMessageInfo{}).Error
	return stores.ErrFmt(err)
}
func (p MessageInfoRepo) FindOne(ctx context.Context, id int64) (*SysMessageInfo, error) {
	var result SysMessageInfo
	err := p.db.WithContext(ctx).Where("id = ?", id).First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}

// 批量插入 LightStrategyDevice 记录
func (p MessageInfoRepo) MultiInsert(ctx context.Context, data []*SysMessageInfo) error {
	err := p.db.WithContext(ctx).Clauses(clause.OnConflict{UpdateAll: true}).Model(&SysMessageInfo{}).Create(data).Error
	return stores.ErrFmt(err)
}
