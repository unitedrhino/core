package relationDB

import "gorm.io/gorm"

func normalizeFilterMap(in map[string]FilterKeywords) map[string]FilterKeywords {
	if in == nil {
		return map[string]FilterKeywords{}
	}
	return in
}

func (m *DataStatisticsInfo) BeforeSave(tx *gorm.DB) error {
	if m == nil {
		return nil
	}
	m.Filter = normalizeFilterMap(m.Filter)
	return nil
}
