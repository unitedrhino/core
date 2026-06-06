package relationDB

import (
	"context"

	"gitee.com/unitedrhino/core/share/dataType"
	"gitee.com/unitedrhino/share/def"
	"gitee.com/unitedrhino/share/stores"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type UserPushClientRepo struct {
	db *gorm.DB
}

func NewUserPushClientRepo(in any) *UserPushClientRepo {
	return &UserPushClientRepo{db: stores.GetCommonConn(in)}
}

type UserPushClientFilter struct {
	UserIDs      []int64
	UserID       int64
	PushClientID string
	IsActive     int64
}

func (p UserPushClientRepo) fmtFilter(ctx context.Context, f UserPushClientFilter) *gorm.DB {
	db := p.db.WithContext(ctx)
	if f.UserID != 0 {
		db = db.Where("user_id = ?", f.UserID)
	}
	if len(f.UserIDs) != 0 {
		db = db.Where("user_id IN ?", f.UserIDs)
	}
	if f.PushClientID != "" {
		db = db.Where("push_client_id = ?", f.PushClientID)
	}
	if f.IsActive != 0 {
		db = db.Where("is_active = ?", f.IsActive)
	}
	return db
}

func (p UserPushClientRepo) Upsert(ctx context.Context, data *SysUserPushClient) error {
	data.IsActive = def.True
	err := p.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   stores.SetColumnsWithPg(p.db, &SysUserPushClient{}, "idx_sys_user_push_client_tc_un"),
		DoUpdates: clause.AssignmentColumns([]string{"platform", "app_id", "app_version", "is_active", "updated_time"}),
	}).Create(data).Error
	return stores.ErrFmt(err)
}

// DeactivateForUser 将指定用户下的推送 cid 标记为失效；pushClientID 为空则失效该用户全部 cid。
func (p UserPushClientRepo) DeactivateForUser(ctx context.Context, tenantCode dataType.TenantCode, userID int64, pushClientID string) error {
	db := p.db.WithContext(ctx).Model(&SysUserPushClient{}).
		Where("tenant_code = ? AND user_id = ? AND is_active = ?", tenantCode, userID, def.True)
	if pushClientID != "" {
		db = db.Where("push_client_id = ?", pushClientID)
	}
	err := db.Update("is_active", def.False).Error
	return stores.ErrFmt(err)
}

// DeactivateOtherCidsForUser 同一用户仅保留当前 cid 有效，避免历史设备 cid 抢收或推送到旧机。
func (p UserPushClientRepo) DeactivateOtherCidsForUser(ctx context.Context, tenantCode dataType.TenantCode, userID int64, keepPushClientID string) error {
	if keepPushClientID == "" {
		return nil
	}
	err := p.db.WithContext(ctx).Model(&SysUserPushClient{}).
		Where("tenant_code = ? AND user_id = ? AND push_client_id <> ? AND is_active = ?",
			tenantCode, userID, keepPushClientID, def.True).
		Update("is_active", def.False).Error
	return stores.ErrFmt(err)
}

// DeactivateOtherUsersByCid 同租户下其他用户绑定同一 cid 的记为失效（本机切换账号后仅当前用户可收推送）。
func (p UserPushClientRepo) DeactivateOtherUsersByCid(ctx context.Context, tenantCode dataType.TenantCode, keepUserID int64, pushClientID string) error {
	if pushClientID == "" {
		return nil
	}
	err := p.db.WithContext(ctx).Model(&SysUserPushClient{}).
		Where("tenant_code = ? AND push_client_id = ? AND user_id <> ? AND is_active = ?",
			tenantCode, pushClientID, keepUserID, def.True).
		Update("is_active", def.False).Error
	return stores.ErrFmt(err)
}

func (p UserPushClientRepo) FindByFilter(ctx context.Context, f UserPushClientFilter) ([]*SysUserPushClient, error) {
	var results []*SysUserPushClient
	err := p.fmtFilter(ctx, f).Model(&SysUserPushClient{}).Find(&results).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	return results, nil
}

func (p UserPushClientRepo) ListActiveClientIDs(ctx context.Context, userIDs []int64) ([]string, error) {
	if len(userIDs) == 0 {
		return nil, nil
	}
	var list []*SysUserPushClient
	err := p.db.WithContext(ctx).Model(&SysUserPushClient{}).
		Where("user_id IN ? AND is_active = ?", userIDs, def.True).
		Order("updated_time DESC").
		Find(&list).Error
	if err != nil {
		return nil, stores.ErrFmt(err)
	}
	seen := map[string]struct{}{}
	var out []string
	// 每用户最多取 2 个最新 cid，避免历史残留设备抢收
	perUser := map[int64]int{}
	for _, row := range list {
		if row == nil || row.PushClientID == "" {
			continue
		}
		if perUser[row.UserID] >= 2 {
			continue
		}
		if _, ok := seen[row.PushClientID]; ok {
			continue
		}
		seen[row.PushClientID] = struct{}{}
		perUser[row.UserID]++
		out = append(out, row.PushClientID)
	}
	return out, nil
}
