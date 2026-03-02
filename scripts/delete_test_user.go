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
	fmt.Println("删除测试用户 18077357550...")

	// 删除用户（软删除）
	result, err := db.Exec(`
		UPDATE sys_user_info
		SET deleted_time = UNIX_TIMESTAMP()
		WHERE phone = '18077357550' AND deleted_time = 0
	`)
	if err != nil {
		log.Fatalf("删除用户失败: %v", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected > 0 {
		fmt.Printf("✓ 成功删除 %d 个用户记录\n", rowsAffected)
		fmt.Println("现在可以重新测试注册流程了")
	} else {
		fmt.Println("未找到该用户或用户已被删除")
	}
}
