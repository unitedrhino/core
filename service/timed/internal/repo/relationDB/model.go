package relationDB

import (
	"gitee.com/unitedrhino/core/service/timed/internal/domain"
	"gitee.com/unitedrhino/share/stores"
	"time"
)

// 示例
type TimedExample struct {
	ID int64 `gorm:"column:id;type:bigint;primary_key;AUTO_INCREMENT"` // id编号
}
type TimedTaskLogScript struct {
	ExecLog []*domain.ScriptLog `gorm:"column:exec_log;type:json;serializer:json"` //执行日志
}

type TimedTaskLogSql struct {
	SelectNum int64 `gorm:"column:select_num"` //查询的数量
	ExecNum   int64 `gorm:"column:exec_num"`   //执行的数量
}

type TimedTaskLog struct {
	ID                  int64     `gorm:"column:id;primary_key"`                         //
	GroupCode           string    `gorm:"column:group_code;index"`                       //组编码
	TaskCode            string    `gorm:"column:task_code;index"`                        //任务编码
	Params              string    `gorm:"column:params;type:json;NOT NULL;default:'{}'"` // 任务参数
	ResultCode          int64     `gorm:"column:result_code;index"`                      //结果code
	ResultMsg           string    `gorm:"column:result_msg"`                             //结果消息
	CreatedTime         time.Time `gorm:"column:created_time;index;sort:desc;default:CURRENT_TIMESTAMP;NOT NULL"`
	*TimedTaskLogSql    `gorm:"embedded;embeddedPrefix:sql_"`
	*TimedTaskLogScript `gorm:"embedded;embeddedPrefix:script_"`
}

func (t *TimedTaskLog) TableName() string {
	return "timed_task_log"
}

type TaskGroupDBConfig struct {
	DSN    string `json:"dsn"`    //数据库连接串
	DBType string `json:"dbType"` //数据库类型(默认mysql)
}

type TimedTaskGroup struct {
	ID       int64             `gorm:"column:id;primary_key"`                                      // 任务组ID
	Code     string            `gorm:"column:code;uniqueIndex:idx_group_code"`                     //任务组编码
	Name     string            `gorm:"column:name;uniqueIndex:idx_group_name"`                     // 组名
	Type     string            `gorm:"column:type"`                                                //组类型:queue(消息队列消息发送)  sql(执行sql) email(邮件发送) http(http请求)
	SubType  string            `gorm:"column:sub_type;default:''"`                                 //组子类型 natsJs nats         normal js
	Priority int64             `gorm:"column:priority"`                                            //组优先级: 6:critical 最高优先级  3: default 普通优先级 1:low 低优先级
	Env      map[string]string `gorm:"column:env;type:json;serializer:json;NOT NULL;default:'{}'"` //环境变量
	/*
		组的配置, sql类型配置格式如下,key若为select,则select默认会选择该配置,exec:exec执行sql默认会选择这个,执行sql的函数也可以指定连接
		database: map[string]TaskGroupDBConfig
	*/
	Config string `gorm:"column:config;type:json;NOT NULL;default:'{}'"` //组的配置
	stores.NoDelTime
	DeletedTime stores.DeletedTime `gorm:"column:deleted_time;default:0;uniqueIndex:idx_group_code;uniqueIndex:idx_group_name"`
}

func (t *TimedTaskGroup) TableName() string {
	return "timed_task_group"
}

type TimedTaskInfo struct {
	ID        int64    `gorm:"column:id;primary_key"`                                // 任务ID
	GroupCode string   `gorm:"column:group_code;uniqueIndex:idx_group_task"`         //组编码
	Type      int64    `gorm:"column:type;default:1"`                                //任务类型 1 定时任务 2 延时任务 3 消息队列触发
	Name      string   `gorm:"column:name"`                                          // 任务名称
	Code      string   `gorm:"column:code;uniqueIndex:idx_group_task"`               //任务编码
	Topics    []string `gorm:"column:topics;type:json;serializer:json;default:'[]'"` //触发topic列表
	Params    string   `gorm:"column:params;type:json;NOT NULL;default:'{}'"`        // 任务参数,延时任务如果没有传任务参数会拿数据库的参数来执行
	CronExpr  string   `gorm:"column:cron_expr"`                                     // cron执行表达式
	Status    int64    `gorm:"column:status"`                                        // 状态
	Priority  int64    `gorm:"column:priority;default:3"`                            //优先级: 10:critical 最高优先级  3: default 普通优先级 1:low 低优先级
	stores.NoDelTime
	DeletedTime stores.DeletedTime `gorm:"column:deleted_time;default:0;uniqueIndex:idx_group_task"`
	Group       *TimedTaskGroup    `gorm:"foreignKey:Code;references:GroupCode"`
}

func (t *TimedTaskInfo) TableName() string {
	return "timed_task_info"
}
