package domain

const (
	TaskGroupTypeQueue  = "queue"
	TaskGroupTypeSql    = "sql"
	TaskGroupTypeScript = "script"
)
const (
	TaskTypeTiming = 1 //定时任务
	TaskTypeDelay  = 2 //延时任务
	TaskTypeQueue  = 3 //消息触发任务
)

const (
	SqlTypeNormal = "normal"
	SqlTypeJs     = "js"
)

const (
	QueueTypeNats   = "nats"
	QueueTypeNatsJs = "natsJs"
)

const (
	SqlEnvDsn    = "dsn"
	SqlEnvDBType = "dbType"
	SqlEnvDriver = "driver"
)

const (
	PriorityCritical = "critical" //最高优先级
	PriorityDefault  = "default"  //默认优先级
	PriorityLow      = "low"      //最低优先级
)

var (
	Priorities = []string{PriorityDefault, PriorityCritical, PriorityLow}
)

type TaskInfo struct {
	ID           int64  `json:"id"`               // 任务ID
	Params       string `json:"params,omitempty"` // 任务参数,延时任务如果没有传任务参数会拿数据库的参数来执行
	Code         string `json:"code"`             //任务编码
	GroupType    string `json:"-"`                //组类型:queue(消息队列消息发送)  sql(执行sql) email(邮件发送) http(http请求)
	GroupSubType string `json:"-"`                //组子类型 natsJs nats         normal js
	GroupCode    string `json:"groupCode"`        //组编码
	//需要使用的环境变量 sql类型中 dsn:如果填写,默认使用该dsn来连接,dbType:mysql|pgsql 默认是mysql
	Env    map[string]string `json:"env"`
	Queue  *ParamQueue       `json:"queue"`  //队列消息
	Sql    *Sql              `json:"sql"`    //sql执行类型
	Script *Script           `json:"script"` //脚本执行类型
}

type ParamQueue struct {
	Topic   string `json:"topic"`
	Payload string `json:"payload"`
}

type ParamSql struct {
	Sql string `json:"sql"`
}

type ParamScript struct {
	Param       map[string]string `json:"param"` //脚本参数,会通过函数入参传进去
	ExecContent string            `json:"execContent"`
}
type SqlDBConfig struct {
	DSN    string `json:"dsn"`    //数据库连接串
	DBType string `json:"dbType"` //数据库类型(默认mysql)
	Driver string `json:"driver"` //连接的驱动
}

type ConfigSql struct {
	Database map[string]SqlDBConfig `json:"database"`
}

type Sql struct {
	Param  ParamSql
	Config ConfigSql
}

type Script struct {
	Param  ParamScript
	Config ConfigSql
}

func ToPriority(level int64) string {
	if level >= 6 {
		return PriorityCritical
	}
	if level >= 3 {
		return PriorityDefault
	}
	return PriorityLow
}
