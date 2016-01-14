package leak_service

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"net/http"
)

var Db *sql.DB

var (
	Port = "8001"
)

func ServeRequest() error {
	http.HandleFunc("/ignore_valid", IgnoreRuleValid)
	err := http.ListenAndServe(":"+Port, nil)
	if err != nil {
		log.Printf("err:%s \n", err)
	}
	return nil
}
