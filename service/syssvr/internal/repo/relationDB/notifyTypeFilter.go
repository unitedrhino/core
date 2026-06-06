package relationDB

import (
	"gitee.com/unitedrhino/core/service/syssvr/internal/defext"
	"gitee.com/unitedrhino/share/def"
	"gorm.io/gorm"
)

// applyMessageInfoNotifyTypeWhere 按通知渠道筛选；历史数据 notify_type 为空视为站内信 message。
func applyMessageInfoNotifyTypeWhere(db *gorm.DB, notifyType string, columns ...string) *gorm.DB {
	col := "notify_type"
	if len(columns) > 0 && columns[0] != "" {
		col = columns[0]
	}
	switch notifyType {
	case string(def.NotifyTypeMessage), "":
		return db.Where(col+" = ? OR "+col+" = '' OR "+col+" IS NULL", string(def.NotifyTypeMessage))
	case string(defext.NotifyTypeSystemNotice):
		return db.Where(col+" = ?", string(defext.NotifyTypeSystemNotice))
	default:
		return db.Where(col+" = ?", notifyType)
	}
}
