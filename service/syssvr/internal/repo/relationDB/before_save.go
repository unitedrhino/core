package relationDB

import (
	"strings"

	"gitee.com/unitedrhino/share/def"
	"gorm.io/gorm"
)

func normalizeJSONMapStringString(in map[string]string) map[string]string {
	if in == nil {
		return map[string]string{}
	}
	return in
}

func normalizeJSONString(in string) string {
	if strings.TrimSpace(in) == "" {
		return "{}"
	}
	return in
}

func normalizeStringSlice(in []string) []string {
	if in == nil {
		return []string{}
	}
	return in
}

func (m *SysNotifyConfig) BeforeSave(tx *gorm.DB) error {
	if m == nil {
		return nil
	}
	if m.SupportTypes == nil {
		m.SupportTypes = []def.NotifyType{}
	}
	if m.EnableTypes == nil {
		m.EnableTypes = []def.NotifyType{}
	}
	m.Params = normalizeJSONMapStringString(m.Params)
	return nil
}

func (m *SysUserInfo) BeforeSave(tx *gorm.DB) error {
	if m == nil {
		return nil
	}
	m.Tags = normalizeJSONMapStringString(m.Tags)
	m.PubTags = normalizeJSONMapStringString(m.PubTags)
	return nil
}

func (m *SysProjectInfo) BeforeSave(tx *gorm.DB) error {
	if m == nil {
		return nil
	}
	m.Tags = normalizeJSONMapStringString(m.Tags)
	return nil
}

func (m *SysProjectCrud) BeforeSave(tx *gorm.DB) error {
	if m == nil {
		return nil
	}
	m.Params = normalizeJSONString(m.Params)
	return nil
}

func (m *SysAreaInfo) BeforeSave(tx *gorm.DB) error {
	if m == nil {
		return nil
	}
	m.Tags = normalizeJSONMapStringString(m.Tags)
	return nil
}

func (m *SysSlotInfo) BeforeSave(tx *gorm.DB) error {
	if m == nil {
		return nil
	}
	m.Hosts = normalizeStringSlice(m.Hosts)
	m.Handler = normalizeJSONMapStringString(m.Handler)
	return nil
}

func (m *SysTenantOpenWebhook) BeforeSave(tx *gorm.DB) error {
	if m == nil {
		return nil
	}
	m.Hosts = normalizeStringSlice(m.Hosts)
	m.Handler = normalizeJSONMapStringString(m.Handler)
	return nil
}
