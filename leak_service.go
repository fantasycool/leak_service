package main

import (
	"database/sql"
	"flag"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"gzlog"
	"leak_service/leak_service"
	"log"
	"strconv"
)

var (
	MysqlHost   = flag.String("db_host", "10.174.33.90:3306", "db host to connect to ")
	ListenPort  = flag.String("serv_port", "8001", "server port to listen")
	LogFile     = flag.String("file_name", "leak_service.log", "leak_service log to write")
	RollingSize = flag.String("rolling_size", "50", "log file rolling size (M)")
	RollingNum  = flag.String("rolling_num", "5", "log file rolling num")
)

const (
	MysqlUser   = "yom"
	MysqlPasswd = ""
)

func main() {
	flag.Parse()
	var rollingSize int
	var rollingNum int
	var err error
	if rollingSize, err = strconv.Atoi(*RollingSize); err != nil {
		log.Printf("rolling size is not valid! \n")
		return
	}
	if rollingNum, err = strconv.Atoi(*RollingNum); err != nil {
		log.Printf("rolling num is not valid! \n")
		return
	}

	leak_service.Port = *ListenPort
	gzlog.InitGZLogger(*LogFile, rollingSize*1000*1000, rollingNum)

	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@(%s)/jwlwl?charset=utf8&parseTime=True", MysqlUser, MysqlPasswd, *MysqlHost))
	defer db.Close()
	if err != nil {
		log.Printf("mysql db connect failed !errMessage:%s \n", err)
		return
	}
	leak_service.Db = db
	log.Printf("Start to ServeRequest!")
	leak_service.ServeRequest()
}
