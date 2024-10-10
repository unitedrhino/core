package relationDB

import (
	"context"
	"fmt"
	"gitee.com/unitedrhino/core/service/timed/internal/domain"
	"gitee.com/unitedrhino/share/conf"
	"gitee.com/unitedrhino/share/def"
	"gitee.com/unitedrhino/share/events/topics"
	"gitee.com/unitedrhino/share/stores"
	"gorm.io/gorm/clause"
	"sync"
)

var once sync.Once

func Migrate(c conf.Database) (err error) {
	if c.IsInitTable == false {
		return
	}
	once.Do(func() {
		db := stores.GetCommonConn(context.TODO())
		var needInitColumn bool
		if !db.Migrator().HasTable(&TimedTaskGroup{}) {
			//需要初始化表
			needInitColumn = true
		}
		err = db.AutoMigrate(
			&TimedTaskLog{},
			&TimedTaskGroup{},
			&TimedTaskInfo{},
		)
		if err != nil {
			return
		}
		if needInitColumn {
			err = migrateTableColumn()
		}
	})
	return
}
func migrateTableColumn() error {
	db := stores.GetCommonConn(context.TODO()).Clauses(clause.OnConflict{DoNothing: true})
	if err := db.CreateInBatches(&MigrateTimedTask, 100).Error; err != nil {
		return err
	}
	if err := db.CreateInBatches(&MigrateTimedTaskGroup, 100).Error; err != nil {
		return err
	}
	return nil
}

var (
	MigrateTimedTask = []TimedTaskInfo{
		//{
		//	GroupCode: "queueTest",
		//	Type:      domain.TaskTypeTiming,
		//	Name:      "消息发送",
		//	Code:      "msgSendTest",
		//	Params:    `{"topic":"server.435","payload":"adfgawe"}`,
		//	CronExpr:  "@every 2s",
		//	Status:    def.StatusWaitRun,
		//	Priority:  2,
		//},
		//{
		//	GroupCode: "sqlJsTest",
		//	Type:      domain.TaskTypeTiming,
		//	Name:      "脚本执行",
		//	Code:      "sqlExec",
		//	Params:    `{"execContent": "function SqlJob(){Set('123','sdafawef');let a = Get('123');LogInfo('get value:',a);let code = GetEnv('code');LogInfo('get code env:',code);Exec(\"insert into test_table(name) values('123123')\");let values = Select('select * from test_table limit 10');LogInfo('select get value :',values);return {code:200,msg:'ok'};}"}`,
		//	CronExpr:  "@every 2s",
		//	Status:    def.StatusWaitRun,
		//	Priority:  4,
		//},
		//{
		//	GroupCode: "queueTest",
		//	Type:      domain.TaskTypeDelay,
		//	Name:      "延时测试",
		//	Code:      "delayTest",
		//	Params:    `{"topic":"server.333","payload":"garegawef"}`,
		//	CronExpr:  "",
		//	Status:    def.StatusRunning,
		//	Priority:  3,
		//},
		//{
		//	GroupCode: def.TimedIThingsQueueGroupCode,
		//	Type:      domain.TaskTypeDelay, //定义一个延时任务
		//	Name:      "流服务数据初始化(自动添加docker到数据库)",
		//	Code:      "VidInfoInitDatabase",
		//	Params:    fmt.Sprintf(`{"topic":"%s","payload":""}`, topics.VidInfoInitDatabase),
		//	CronExpr:  "",
		//	Status:    def.StatusWaitRun,
		//	Priority:  3,
		//},
		{
			GroupCode: def.TimedIThingsQueueGroupCode,
			Type:      domain.TaskTypeTiming,
			Name:      "timedJob服务缓存及日志清理",
			Code:      "timedJobClean",
			Params:    fmt.Sprintf(`{"topic":"%s","payload":""}`, topics.TimedJobClean),
			CronExpr:  "1 1 * * *",
			Status:    def.StatusWaitRun,
			Priority:  3,
		},
		//{
		//	GroupCode: def.TimedIThingsQueueGroupCode,
		//	Type:      domain.TaskTypeTiming,
		//	Name:      "流服务状态更新",
		//	Code:      "VidInfoCheckStatus",
		//	Params:    fmt.Sprintf(`{"topic":"%s","payload":""}`, topics.VidInfoCheckStatus),
		//	CronExpr:  "@every 30s",
		//	Status:    def.StatusWaitRun,
		//	Priority:  1, //低优先级任务
		//},
	}
	MigrateTimedTaskGroup = []TimedTaskGroup{
		{
			Code:     "queueTest",
			Name:     "消息队列测试",
			Type:     domain.TaskGroupTypeQueue,
			SubType:  domain.QueueTypeNatsJs,
			Priority: 9,
		},
		{
			Code:     def.TimedIThingsQueueGroupCode,
			Name:     "联犀系统定时消息任务组",
			Type:     domain.TaskGroupTypeQueue,
			SubType:  domain.QueueTypeNatsJs,
			Priority: 9,
		},
		{
			Code:     "sqlJsTest",
			Name:     "sqlJs模式测试",
			Type:     domain.TaskGroupTypeSql,
			SubType:  domain.SqlTypeJs,
			Priority: 7,
			Env:      map[string]string{"code": "66666"},
			Config:   `{"database":{"select":{"dsn":"root:password@tcp(127.0.0.1:3306)/iThings?charset=utf8mb4&collation=utf8mb4_bin&parseTime=true&loc=Asia%2FShanghai","dbType":"mysql"}}}`,
		},
	}
	//var a = "{\"execContent\": \"function SqlJob(){let values=Select('Select * from model_common_hublog limit 10');LogInfo('select get value :',values);return {code:200,msg:'ok'};}\"}"
)
