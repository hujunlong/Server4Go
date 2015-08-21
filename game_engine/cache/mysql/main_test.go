package main

import (
	"database/sql"
	"fmt"
	_ "game_engine/cache/mysql"
)

func CheckError(err error) {
	if err != nil {
		fmt.Println("err:", err)
	}
}

func insert(db *sql.DB) {
	stmt, err := db.Prepare("insert into user(username,password) values(?,?)")
	defer stmt.Close()

	CheckError(err)
	stmt.Exec("guotie", "123456")
	stmt.Exec("abcdef", "222222")
}

func main() {
	db, err := sql.Open("mysql", "root:game9z@/test")
	CheckError(err)
	defer db.Close()

	err = db.Ping()
	CheckError(err)
	insert(db)

	rows, err := db.Query("select username, password from user where password =?", "123456")
	CheckError(err)
	defer rows.Close()
	var name string
	var password int
	for rows.Next() {
		err := rows.Scan(&name, &password)
		fmt.Println("name = ", name, "password = ", password)
		CheckError(err)
	}
}
