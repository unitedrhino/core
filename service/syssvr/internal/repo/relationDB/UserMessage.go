package relationDB

import (
	"context"
	"fmt"
	"gitee.com/i-Things/share/ctxs"
	"gitee.com/i-Things/share/def"
	"gitee.com/i-Things/share/stores"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

/*
这个是参考样例
使用教程:
1. 将UserMessage全局替换为模型的表名
2. 完善todo
*/

type UserMessageRepo struct {
	db *gorm.DB
}

func NewUserMessageRepo(in any) *UserMessageRepo {
	return &UserMessageRepo{db: stores.GetCommonConn(in)}
}

type UserMessageFilter struct {
	MessageID   int64
	WithMessage bool
	Group       string
	NotifyCode  string
	IsRead      int64
	Str1        string
	Str2        string
	Str3        string
}

func (p UserMessageRepo) fmtFilter(ctx context.Context, f UserMessageFilter) *gorm.DB {
	db := p.db.WithContext(ctx)
	if f.MessageID != 0 {
		db = db.Where("message_id=?", f.MessageID)
	}
	if f.Group != "" {
		db = db.Where(fmt.Sprintf("%s = ?", stores.Col("group")), f.Group)
	}
	db = db.Where("user_id=?", ctxs.GetUserCtxNoNil(ctx).UserID)
	if f.NotifyCode != "" {
		db = db.Where("notify_code = ?", f.NotifyCode)
	}
	if f.IsRead != 0 {
		db = db.Where("is_read = ?", f.IsRead)
	}
	if f.Str1 != "" {
		db = db.Where("str1 = ?", f.Str1)
	}
	if f.Str2 != "" {
		db = db.Where("str2 = ?", f.Str2)
	}
	if f.Str3 != "" {
		db = db.Where("str3 = ?", f.Str3)
	}

	if f.WithMessage {
		db = db.Preload("Message")
	}
	return db
}

func (p UserMessageRepo) Insert(ctx context.Context, data *SysUserMessage) error {
	result := p.db.WithContext(ctx).Create(data)
	return stores.ErrFmt(result.Error)
}

func (p UserMessageRepo) FindOneByFilter(ctx context.Context, f UserMessageFilter) (*SysUserMessage, error) {
	var result SysUserMessage
	db := p.fmtFilter(ctx, f)
	err := db.First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}

func (p UserMessageRepo) FindByFilter(ctx context.Context, f UserMessageFilter, page *stores.PageInfo) ([]*SysUserMessage, error) {
	var results []*SysUserMessage
	db := p.fmtFilter(ctx, f).Model(&SysUserMessage{})
	db = page.ToGorm(db)
	err := db.Find(&results).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return results, nil
}

func (p UserMessageRepo) CountByFilter(ctx context.Context, f UserMessageFilter) (size int64, err error) {
	db := p.fmtFilter(ctx, f).Model(&SysUserMessage{})
	err = db.Count(&size).Error
	return size, stores.ErrFmt(err)
}

func (p UserMessageRepo) Update(ctx context.Context, data *SysUserMessage) error {
	err := p.db.WithContext(ctx).Where("id = ?", data.ID).Save(data).Error
	return stores.ErrFmt(err)
}

func (p UserMessageRepo) MultiIsRead(ctx context.Context, userID int64, ids []int64) error {
	err := p.db.WithContext(ctx).Model(&SysUserMessage{}).Where("user_id = ? and id in ?", userID, ids).Update("is_read", def.True).Error
	return stores.ErrFmt(err)
}
func (p UserMessageRepo) CountNotRead(ctx context.Context, userID int64) (map[string]int64, error) {
	var ret []struct {
		Group string `json:"group"`
		C     int64  `json:"c"`
	}
	err := p.db.WithContext(ctx).Model(&SysUserMessage{}).Group(stores.Col("group")).
		Select(fmt.Sprintf("%s,count(1) as c", stores.Col("group"))).
		Where("user_id = ? and is_read = ?", userID, def.False).Find(&ret).Error
	var val = map[string]int64{}
	for _, v := range ret {
		val[v.Group] = v.C
	}
	return val, stores.ErrFmt(err)
}

func (p UserMessageRepo) DeleteByFilter(ctx context.Context, f UserMessageFilter) error {
	db := p.fmtFilter(ctx, f)
	err := db.Delete(&SysUserMessage{}).Error
	return stores.ErrFmt(err)
}

func (p UserMessageRepo) Delete(ctx context.Context, id int64) error {
	err := p.db.WithContext(ctx).Where("id = ?", id).Delete(&SysUserMessage{}).Error
	return stores.ErrFmt(err)
}
func (p UserMessageRepo) FindOne(ctx context.Context, id int64) (*SysUserMessage, error) {
	var result SysUserMessage
	err := p.db.WithContext(ctx).Where("id = ?", id).First(&result).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return &result, nil
}

// 批量插入 LightStrategyDevice 记录
func (p UserMessageRepo) MultiInsert(ctx context.Context, data []*SysUserMessage) error {
	err := p.db.WithContext(ctx).Clauses(clause.OnConflict{UpdateAll: true}).Model(&SysUserMessage{}).Create(data).Error
	return stores.ErrFmt(err)
}
