package relationDB

import (
	"database/sql"
	"gitee.com/unitedrhino/share/domain/ops"
	"gitee.com/unitedrhino/share/stores"
)

// 设备维护工单 device Maintenance Work Order
type SysOpsWorkOrder struct {
	ID           int64               `gorm:"column:id;type:bigint;primary_key;AUTO_INCREMENT"`
	TenantCode   stores.TenantCode   `gorm:"column:tenant_code;index;type:VARCHAR(50);NOT NULL"` // 租户编码
	RaiseUserID  int64               `gorm:"column:raise_user_id;type:BIGINT;NOT NULL"`          // 问题提出的用户
	ProjectID    stores.ProjectID    `gorm:"column:project_id;type:bigint;default:0;NOT NULL"`   // 项目ID(雪花ID)
	AreaID       stores.AreaID       `gorm:"column:area_id;type:bigint;default:0;NOT NULL"`      // 项目区域ID(雪花ID)
	Number       string              `gorm:"column:number;unique;type:VARCHAR(50);NOT NULL"`     //编号
	Params       map[string]string   `gorm:"column:params;type:json;serializer:json;"`           // 参数 json格式
	Type         string              `gorm:"column:type;type:varchar(100);NOT NULL"`             // 工单类型: deviceMaintenance:设备维修工单
	IssueDesc    string              `gorm:"column:issue_desc;type:varchar(2000);NOT NULL"`
	Status       ops.WorkOrderStatus `gorm:"column:status;type:BIGINT;default:1"` //状态 1:待处理 2:处理中 3:已完成
	HandleTime   sql.NullTime        `gorm:"column:handle_time;default:null"`     //处理时间
	FinishedTime sql.NullTime        `gorm:"column:finished_time;default:null"`   //处理完成时间
	stores.SoftTime
}

func (m *SysOpsWorkOrder) TableName() string {
	return "sys_ops_work_order"
}

// 帮助与反馈
type SysOpsFeedback struct {
	ID                 int64               `gorm:"column:id;type:bigint;primary_key;AUTO_INCREMENT"`
	TenantCode         stores.TenantCode   `gorm:"column:tenant_code;index;type:VARCHAR(50);NOT NULL"`   // 租户编码
	RaiseUserID        int64               `gorm:"column:raise_user_id;type:BIGINT;NOT NULL"`            // 问题提出的用户
	ProjectID          stores.ProjectID    `gorm:"column:project_id;type:bigint;default:0;NOT NULL"`     // 项目ID(雪花ID)
	Type               string              `gorm:"column:type;type:VARCHAR(50);NOT NULL"`                //问题类型 设备问题:thingsDevice 智能场景:thingsScene 体验问题: experience 其他: other
	Status             ops.WorkOrderStatus `gorm:"column:status;type:BIGINT;default:1"`                  //状态 1:待处理 2:处理中 3:已完成
	ContactInformation string              `gorm:"column:contact_information;type:VARCHAR(50);NOT NULL"` //联系信息
	IssueDesc          string              `gorm:"column:issue_desc;type:varchar(2000);NOT NULL"`
	stores.SoftTime
}

func (m *SysOpsFeedback) TableName() string {
	return "sys_ops_feedback"
}
