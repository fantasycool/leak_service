package leak_service

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

type SummaryIgnoreRule struct {
	Id     int64
	UserId int64
	Filter string
}

func QueryIgnoreRules(userId int64, db *sql.DB) ([]*SummaryIgnoreRule, error) {
	rows, err := db.Query(`select id,user_id, filter from summary_ignore_rules where user_id=?`, userId)
	if err != nil {
		return nil, err
	}
	rules := make([]*SummaryIgnoreRule, 0)
	for rows.Next() {
		rule := &SummaryIgnoreRule{}
		err := rows.Scan(&rule.Id, &rule.UserId, &rule.Filter)
		if err != nil {
			log.Printf("Query mysql summary ignore rules scan failed! userId:%d \n", userId)
			return nil, err
		}
		rules = append(rules, rule)
	}
	return rules, nil
}
