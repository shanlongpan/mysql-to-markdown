package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

type Conf struct {
	DB []string `yaml:"DB"`
}

func main() {
	confInputDir := os.Args

	if len(confInputDir) == 1 {
		fmt.Println("请输入配置文件")
		return
	}
	confDir := os.Args[1]

	err, dbInfo := GetMysqlInfo(confDir)
	if err != nil {
		fmt.Println(err)
		return
	}
	if len(dbInfo) < 1 {
		fmt.Println("没有检查到配置文件对应的数据信息")
	}

	for _, mysqlInfo := range dbInfo {
		WriteIntoMarkDown(mysqlInfo)
	}
	fmt.Println("程序执行结束")
}

func GetMysqlInfo(dir string) (error, []string) {
	data, err := ioutil.ReadFile(dir)
	dbConfi := Conf{}
	if err != nil {
		return err, dbConfi.DB
	}
	//把yaml形式的字符串解析成struct类型
	err = yaml.Unmarshal(data, &dbConfi)

	return err, dbConfi.DB
}

func WriteIntoMarkDown(mysqlConInfo string) {
	db, err := sql.Open("mysql", mysqlConInfo)
	if err != nil {
		fmt.Println(err, "联接数据库信息有误--"+mysqlConInfo)
		return
	}

	rows, err := db.Query("SHOW TABLE STATUS;")
	if err != nil {
		fmt.Println(err, "查询数据库有误--"+mysqlConInfo)
		return
	}
	startDatabase := strings.Index(mysqlConInfo, "/")

	endDatabase := strings.Index(mysqlConInfo, "?")

	databaseName := Substr(mysqlConInfo, startDatabase+1, endDatabase)
	fileName := "./" + databaseName + ".md"

	PutInFile("###数据库："+databaseName+" 详细信息\n####1.1表信息\n| Name| Engine | Rows | Create_time| Comment|\n| ---- | ---- | ---- |---- |---- |\n", fileName)
	tableNames := []string{}

	tableInfo := GetRetMap(rows)

	for _, tables := range tableInfo {
		//  | ---- | ---- | ---- |---- |---- |
		tableView := make([]string, 5)
		tableView[0] = MarkdownChang(tables["Name"].(string))
		tableView[1] = tables["Engine"].(string)
		tableView[2] = tables["Rows"].(string)
		tableView[3] = Substr(tables["Create_time"].(string), 0, 10)
		tableView[4] = tables["Comment"].(string)
		formatStr := strings.Join(tableView, "|")
		PutInFile("|"+formatStr+"|\n", fileName)
		tableNames = append(tableNames, tables["Name"].(string))
	}
	for id, tableName := range tableNames {
		rows, err := db.Query("select * from information_schema.columns where table_schema = '" + databaseName + "' and table_name = '" + tableName + "';")

		if err != nil {
			fmt.Println(err, "查询表有误--"+tableName)
			continue
		}
		PutInFile("####2."+strconv.Itoa(id)+" 表 "+MarkdownChang(tableName)+" 详细信息\n| TABLE_SCHEMA| TABLE_NAME | COLUMN_NAME | COLUMN_DEFAULT| COLUMN_TYPE|COLUMN_KEY|COLUMN_COMMENT|\n| ---- | ---- | ---- |---- |---- |----|---- | \n", fileName)
		tableMessage := GetRetMap(rows)
		for _, columnInfo := range tableMessage {
			columnView := make([]string, 7)
			columnView[0] = MarkdownChang(columnInfo["TABLE_SCHEMA"].(string))
			columnView[1] = MarkdownChang(columnInfo["TABLE_NAME"].(string))
			columnView[2] = MarkdownChang(columnInfo["COLUMN_NAME"].(string))
			columnView[3] = MarkdownChang(columnInfo["COLUMN_DEFAULT"].(string))
			columnView[4] = columnInfo["COLUMN_TYPE"].(string)
			columnView[5] = columnInfo["COLUMN_KEY"].(string)
			columnView[6] = columnInfo["COLUMN_COMMENT"].(string)

			formatStr := strings.Join(columnView, "|")
			PutInFile("|"+formatStr+"|\n", fileName)
		}
	}
}

func GetRetMap(result *sql.Rows) []map[string]interface{} {

	//Scan需要的容器
	keys, _ := result.Columns()             //字段名
	values := make([][]byte, len(keys))     //字段值
	scans := make([]interface{}, len(keys)) //Scan时的容器
	for i := range scans {
		scans[i] = &values[i] //对容器的处理，以便能在Scan之后收集到每个字段值
	}

	//循环取出每一行的数据
	datas := make([]map[string]interface{}, 0)
	for result.Next() {
		result.Scan(scans...)
		data := make(map[string]interface{})
		for k, v := range values {
			key := keys[k]
			data[key] = string(v)
		}
		datas = append(datas, data)
	}
	return datas
}

func Substr(str string, start int, end int) string {
	rs := []rune(str)
	length := len(rs)

	if start < 0 || start > length {
		return ""
	}

	if end < 0 || end > length {
		return ""
	}
	return string(rs[start:end])
}

func PutInFile(strContent string, fileName string) {
	fd, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		fmt.Println(err)
		return
	}
	buf := []byte(strContent)
	fd.Write(buf)
	fd.Close()
}

func MarkdownChang(str string) string {
	return strings.ReplaceAll(str, "_", "\\_")
}
