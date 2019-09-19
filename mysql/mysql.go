package mysql2

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

type Db struct {
	Sql 	string
}


func (mysql *Db)Explain(sql string) ([]string,[]map[string]string,error) {
	mysql.Sql = "explain " + sql
	return mysql.exec()

}

func (mysql *Db)Profile(id string) ([]string,[]map[string]string,error) {
	mysql.Sql = "show profile BLOCK IO,CPU for query "+id
	return mysql.exec()
}

func (mysql *Db)Result(sql string) ([]string,[]map[string]string,error) {
	mysql.Sql =  sql
	//return mysql.exec()
	columns,res, err := mysql.exec()
	var result []map[string]string
	if len(res) > 0 {
		result = append(result, res[0])
	}
	result = append(result,res...)
	return columns,result,err
}

func (mysql *Db)exec()([]string,[]map[string]string,error) {
	db, err := sql.Open("mysql", "root:GPXrQ71n@tcp(127.0.0.1)/blog")
	if err != nil {
		return nil,nil,err
	}
	defer db.Close()

	db = mysql.setProfiling(db)
	db = mysql.setProfiling(db)

	result, err := db.Query(mysql.Sql)

	if err != nil {
		return nil,nil,err
	}

	columns, _ := result.Columns()
	scanArgs := make([]interface{}, len(columns))
	values := make([]interface{}, len(columns))

	for i := range values {
		scanArgs[i] = &values[i]
	}

	var ret []map[string]string
	for result.Next() {
		//将行数据保存到record字典
		err = result.Scan(scanArgs...)
		record := make(map[string]string)
		for i, col := range values {
			if col != nil {
				record[columns[i]] = string(col.([]byte))
			}else{
				record[columns[i]] = "null"
			}
		}
		ret = append(ret,record)
	}
	return columns,ret,nil
}

func (mysql *Db)setProfiling(db *sql.DB) *sql.DB {

	result, _ := db.Query("SET profiling = 1")

	columns, _ := result.Columns()
	//fmt.Println(columns)
	scanArgs := make([]interface{}, len(columns))
	values := make([]interface{}, len(columns))

	for i := range values {
		scanArgs[i] = &values[i]
	}

	for result.Next() {
		//将行数据保存到record字典
		//result.Scan(scanArgs...)

	}
	return db
}

