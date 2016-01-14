package leak_service

import (
	"encoding/json"
	"log"
	"net/http"
	"regexp"
	"strconv"
)

func IgnoreRuleValid(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		log.Printf("parse form failed!%s \n", err)
		w.Write([]byte("{}"))
		return
	}
	data := r.Form.Get("data")
	if data == "" {
		log.Printf("request data is nil, do nothing! \n")
		w.Write([]byte("{}"))
		return
	}
	log.Printf("data is %s \n", data)

	dataMap := make(map[string]interface{}, 10)
	err = json.Unmarshal([]byte(data), &dataMap)
	if err != nil {
		log.Printf("data map parameter is not json type %s ;err:%s\n", data, err)
		w.Write([]byte("{}"))
		return
	}
	if dataMap["userid"] == nil {
		w.Write([]byte("{}"))
		return
	}
	i, err := strconv.ParseInt(dataMap["userid"].(string), 10, 64)
	if err != nil {
		log.Printf("Illegal userId %s;err:%s\n", data, err)
		w.Write([]byte("{}"))
		return
	}
	rules, err := QueryIgnoreRules(i, Db)
	if err != nil {
		log.Printf("Query rules failed !%s;err:%s \n", rules, err)
		w.Write([]byte("{}"))
		return
	}
	log.Printf("Queryied rules is %s \n", rules)
	result, err := ValidateIgnore(dataMap, rules)
	if err != nil {
		log.Printf("ValidateIgnore failed !%s;err:%s\n", dataMap, err)
		w.Write([]byte("{}"))
		return
	}
	j, err := json.Marshal(result)
	if err != nil {
		log.Printf("json marshal result:%s failed;err:%s", result, err)
		w.Write([]byte("{}"))
		return
	}
	w.Write(j)

}

func ValidateIgnore(dataMap map[string]interface{}, rules []*SummaryIgnoreRule) ([]int64, error) {
	matchedRuleId := make([]int64, 0)
	for _, rule := range rules {
		yes := true
		filterMap := make(map[string]string, 10)
		err := json.Unmarshal([]byte(rule.Filter), &filterMap)
		if err != nil {
			log.Printf("json unmarshal filter failed,rule.filter: %s;err:%s \n", rule.Filter, err)
			return nil, err
		}
		for key, value := range dataMap {
			if filterMap[key] == "" {
				continue
			}
			switch key {
			case "files":
				if !validateFiles(dataMap["files"].(string), filterMap[key]) {
					yes = false
					break
				}
			case "path":
				if !validatePath(dataMap["path"].(string), filterMap[key]) {
					yes = false
					break
				}
			default:
				r, err := regexp.Compile(filterMap[key])
				if err != nil {
					log.Printf("regexp compile failed! value %s, err:%s \n", filterMap[key], err)
					return nil, err
				}
				if !r.MatchString(value.(string)) {
					yes = false
					break
				}
			}
		}
		if yes {
			matchedRuleId = append(matchedRuleId, rule.Id)
		}
	}
	return matchedRuleId, nil
}

func validateFiles(data string, regexStr string) bool {
	return validatePath(data, regexStr)
}

func validatePath(data string, regexStr string) bool {
	fileSuffixR, _ := regexp.Compile(`(\.[a-zA-Z0-9]+)"?`)
	filterSuffixR, _ := regexp.Compile(regexStr)
	matchedResult := fileSuffixR.FindAllStringSubmatch(data, -1)
	for _, m := range matchedResult {
		if !filterSuffixR.MatchString(m[1]) {
			return false
		}
	}
	return true
}
