package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	// 数据库连接字符串
	dsn := "root:dcDlqT67IAq7V1eeGD1hDlBDlIHP@tcp(gz-cynosdbmysql-grp-7cqua4ep.sql.tencentcdb.com:25601)/iThings?charset=utf8mb4&parseTime=true&loc=Asia%2FShanghai"

	// 连接数据库
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("连接数据库失败: %v", err)
	}
	defer db.Close()

	// 测试连接
	if err := db.Ping(); err != nil {
		log.Fatalf("数据库连接测试失败: %v", err)
	}

	fmt.Println("数据库连接成功！")
	fmt.Println("查询手机号 18077357550 的用户信息...")
	fmt.Println()

	// 查询用户信息
	query := `
		SELECT user_id, user_name, phone,
		       COALESCE(huawei_union_id, 'NULL') as huawei_union_id,
		       COALESCE(huawei_open_id, 'NULL') as huawei_open_id,
		       COALESCE(wechat_union_id, 'NULL') as wechat_union_id,
		       COALESCE(wechat_open_id, 'NULL') as wechat_open_id
		FROM sys_user_info
		WHERE phone = '18077357550' AND deleted_time = 0
	`

	rows, err := db.Query(query)
	if err != nil {
		log.Fatalf("查询失败: %v", err)
	}
	defer rows.Close()

	found := false
	for rows.Next() {
		found = true
		var userID int64
		var userName, phone, huaweiUnionID, huaweiOpenID, wechatUnionID, wechatOpenID string

		err := rows.Scan(&userID, &userName, &phone, &huaweiUnionID, &huaweiOpenID, &wechatUnionID, &wechatOpenID)
		if err != nil {
			log.Fatalf("读取数据失败: %v", err)
		}

		fmt.Printf("用户ID: %d\n", userID)
		fmt.Printf("用户名: %s\n", userName)
		fmt.Printf("手机号: %s\n", phone)
		fmt.Printf("华为UnionID: %s\n", huaweiUnionID)
		fmt.Printf("华为OpenID: %s\n", huaweiOpenID)
		fmt.Printf("微信UnionID: %s\n", wechatUnionID)
		fmt.Printf("微信OpenID: %s\n", wechatOpenID)
		fmt.Println()

		// 判断华为账号是否绑定
		if huaweiUnionID == "NULL" && huaweiOpenID == "NULL" {
			fmt.Println("❌ 问题确认：华为账号信息未保存到数据库！")
			fmt.Println("   这就是为什么再次登录时返回'未注册'的原因。")
		} else {
			fmt.Println("✓ 华为账号信息已保存")
		}
	}

	if !found {
		fmt.Println("未找到该用户记录")
	}

	if err := rows.Err(); err != nil {
		log.Fatalf("遍历结果失败: %v", err)
	}
}
