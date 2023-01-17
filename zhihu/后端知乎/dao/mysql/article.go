package mysql

import (
	"fmt"
	"log"
	"main/app/global"
	"main/model"
	"time"
)

func AddArticle(context, userID, topicID string) (bool, string) {

	sqlStr := "insert into article(id,user_id,topic_id,create_time,update_time,context) values (?,?,?,?,?,?)"
	_, err := global.MysqlDB.Exec(sqlStr, FindArticleID(), userID, topicID, time.Now(), time.Now(), context)
	if err != nil {
		fmt.Printf("insert failed, err:%v\n", err)
		return false, "another error"
	}
	log.Println("insert success")
	return true, ""
}

func FindArticleID() int {
	sqlStr := "select id from article where id >=?"
	rows, err := global.MysqlDB.Query(sqlStr, 0)
	if err != nil {
		fmt.Printf("query failed, err:%v\n", err)
		return 0
	}
	// 非常重要：关闭rows释放持有的数据库链接
	defer rows.Close()

	// 循环读取结果集中的数据
	i := 1
	for rows.Next() {
		var u model.User
		err := rows.Scan(&u.ID)
		if err != nil {
			fmt.Printf("scan failed, err:%v\n", err)
			return 0
		}
		i++
	}
	return i
}
