package relationDB

import (
	"strings"

	"gorm.io/gorm"
)

func normalizeJSONString(in string) string {
	if strings.TrimSpace(in) == "" {
		return "{}"
	}
	return in
}

func normalizeMapStringString(in map[string]string) map[string]string {
	if in == nil {
		return map[string]string{}
	}
	return in
}

func (m *TimedTaskLog) BeforeSave(tx *gorm.DB) error {
	if m == nil {
		return nil
	}
	m.Params = normalizeJSONString(m.Params)
	return nil
}

func (m *TimedTaskGroup) BeforeSave(tx *gorm.DB) error {
	if m == nil {
		return nil
	}
	m.Env = normalizeMapStringString(m.Env)
	m.Config = normalizeJSONString(m.Config)
	return nil
}

func (m *TimedTaskInfo) BeforeSave(tx *gorm.DB) error {
	if m == nil {
		return nil
	}
	m.Params = normalizeJSONString(m.Params)
	return nil
}
