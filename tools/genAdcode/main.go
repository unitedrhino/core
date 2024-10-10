package main

import (
	"fmt"
	"gitee.com/unitedrhino/share/utils"
	"github.com/zeromicro/go-zero/core/logx"
	"os"
	"strings"
)

var file = "./Township_Area_A_20240719.xlsx"
var idStart int64 = 100
var dictCode = "adcode"
var dictMap = map[string]*SysDictDetail{}
var dictSlice []*SysDictDetail

// 获取地址: https://lbsyun.baidu.com/faq/api?title=webapi/download
func main() {
	f, err := os.Open(file)
	logx.Must(err)
	exe, err := utils.ReadExcel(f, file)
	logx.Must(err)
	for _, v := range exe {
		ConvertOne(v)
	}
	rst := GenResult()
	err = os.WriteFile("out/modelMigrateAdcode.go", []byte(rst), 0666)
	logx.Must(err)
	fmt.Println("结束转换")
}

type SysDictDetail struct {
	ID       int64            `gorm:"column:id;type:BIGINT;primary_key;AUTO_INCREMENT"`      // id编号
	DictCode string           `gorm:"column:dict_code;type:VARCHAR(50);default:'';NOT NULL"` // 关联标记
	Label    string           `gorm:"column:label;comment:展示值"`                              // 展示值
	Value    string           `gorm:"column:value;comment:字典值"`                              // 字典值
	IDPath   string           `gorm:"column:id_path;type:varchar(100);NOT NULL"`             // 1-2-3-的格式记录顶级区域到当前区域的路径
	ParentID int64            `gorm:"column:parent_id;type:BIGINT"`                          // id编号
	Children []*SysDictDetail `gorm:"foreignKey:ParentID;references:ID"`
	Parent   *SysDictDetail   `gorm:"foreignKey:ID;references:ParentID"`
}

func GetID() int64 {
	idStart++
	return idStart
}

var tmp = `package relationDB

import "gitee.com/unitedrhino/share/def"

var (
	MigrateDictDetailAdcode = []SysDictDetail{
		%s,
	}
)`

func GenResult() string {
	var dicts []string
	for _, v := range dictSlice {
		dicts = append(dicts, fmt.Sprintf(`{ID: %d,DictCode: "%s",Label:    "%s",Value:    "%s",Status:   def.True,ParentID: %d,IDPath:   "%s",Sort:     1}`, v.ID, v.DictCode, v.Label, v.Value, v.ParentID, v.IDPath))
	}
	return fmt.Sprintf(tmp, strings.Join(dicts, ",\n"))

}

func ConvertOne(in []string) {
	l1 := dictMap[in[1]]
	if l1 == nil {
		id := GetID()
		dictMap[in[1]] = &SysDictDetail{
			ID:       id,
			DictCode: dictCode,
			Label:    in[0],
			Value:    in[1],
			IDPath:   utils.GenIDPath(id),
			ParentID: 1,
		}
		dictSlice = append(dictSlice, dictMap[in[1]])
		l1 = dictMap[in[1]]
	}
	l2 := dictMap[in[4]]
	if l2 == nil {
		id := GetID()
		dictMap[in[4]] = &SysDictDetail{
			ID:       id,
			DictCode: dictCode,
			Label:    in[2],
			Value:    in[4],
			IDPath:   utils.GenIDPath(l1.ID, id),
			ParentID: l1.ID,
		}
		dictSlice = append(dictSlice, dictMap[in[4]])
		l2 = dictMap[in[4]]
	}
	l3 := dictMap[in[7]]
	if l3 == nil {
		id := GetID()
		dictMap[in[7]] = &SysDictDetail{
			ID:       id,
			DictCode: dictCode,
			Label:    in[5],
			Value:    in[7],
			IDPath:   utils.GenIDPath(l1.ID, l2.ID, id),
			ParentID: l2.ID,
		}
		l3 = dictMap[in[7]]
		dictSlice = append(dictSlice, dictMap[in[7]])
	}
	id := GetID()
	dictMap[in[10]] = &SysDictDetail{
		ID:       id,
		DictCode: dictCode,
		Label:    in[8],
		Value:    in[10],
		IDPath:   utils.GenIDPath(l1.ID, l2.ID, l3.ID, id),
		ParentID: l3.ID,
	}
	dictSlice = append(dictSlice, dictMap[in[10]])
	return
}
