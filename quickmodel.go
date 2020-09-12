package quickmodel

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	_ "github.com/jinzhu/gorm/dialects/mysql"
	"gopkg.in/ini.v1"
)

type Database struct {
	Host     string `ini:"host"`
	User     string `ini:"user"`
	Password string `ini:"password"`
	DbName   string `ini:"dbname"`
}

var DBConfig = Database{}
var DB *sql.DB

func LoadFile(path string) bool {
	cfg, err := ini.Load(path)
	if err != nil {
		fmt.Println("数据库配置文件找不到")
		return false
	}
	errs := cfg.Section("database").MapTo(&DBConfig)
	if errs != nil {
		fmt.Println("数据库配置文件字段错误")
		return false
	}
	var error error
	str := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?charset=utf8&parseTime=True&loc=Local",
		DBConfig.User,
		DBConfig.Password,
		DBConfig.Host,
		DBConfig.DbName,
	)
	DB, error = sql.Open("mysql", str)
	if nil != error {
		fmt.Println("数据库连接失败: ", error)
		return false
	}
	return true
}

func CreateModel() error {
	_, err := os.Stat("models")
	if err != nil {
		os.MkdirAll("models", 644)
	}
	tables := make([]string, 0)
	rows, _ := DB.Query("show tables;")
	for rows.Next() {
		var name string
		rows.Scan(&name)
		tables = append(tables, name)
	}
	bytes, err := ioutil.ReadFile("template/demo.go")
	if err != nil {
		return err
	}

	fmt.Println("总共", len(tables), "张表")
	for _, table := range tables {
		file, _ := os.Create(fmt.Sprintf("models/%s.go", table))
		fmt.Println("正在处理", table, "表...")
		sql := "SELECT " +
			"concat(UPPER( SUBSTRING( COLUMN_NAME, 1, 1 ) )," +
			"SUBSTRING( COLUMN_NAME, 2, length( COLUMN_NAME ) )," +
			"( CASE DATA_TYPE WHEN 'char' THEN ' string ' WHEN 'enum' THEN ' int8 ' WHEN 'mediumint' THEN ' int8 ' WHEN 'smallint' THEN ' int8 '  WHEN 'varchar' THEN ' string ' WHEN 'timestamp' THEN ' int ' WHEN 'int' THEN ' int ' WHEN 'decimal' THEN ' float64 ' WHEN 'text' THEN ' string ' WHEN 'tinyint' THEN ' int ' WHEN 'double' THEN ' float32 ' WHEN 'float' THEN ' float64 ' WHEN 'datetime' THEN ' time.Time ' END )," +
			"'`json:\"',COLUMN_NAME,'\" ','\\comment:','\"',COLUMN_COMMENT,'\"`') AS struct FROM" +
			"( SELECT ORDINAL_POSITION, COLUMN_NAME, DATA_TYPE, COLUMN_COMMENT FROM information_schema.COLUMNS WHERE table_name = '%s' AND table_schema = 'test' ) a ORDER BY ORDINAL_POSITION"
		rows, err := DB.Query(fmt.Sprintf(sql, table))
		if err != nil {
			return err
		}
		var params string
		for rows.Next() {
			var name string
			rows.Scan(&name)
			params += name + "\n\t"
		}

		fmt.Println(table, "表处理完成")
		defer file.Close()
		temp := fmt.Sprintf(string(bytes), Capitalize(strFirstToUpper(table)), params)
		file.WriteString(temp)
	}
	return nil
}

func strFirstToUpper(str string) string {
	temp := strings.Split(str, "_")
	var upperStr string
	for y := 0; y < len(temp); y++ {
		vv := []rune(temp[y])
		if y != 0 {
			for i := 0; i < len(vv); i++ {
				if i == 0 {
					vv[i] -= 32
					upperStr += string(vv[i]) // + string(vv[i+1])
				} else {
					upperStr += string(vv[i])
				}
			}
		}
	}
	return temp[0] + upperStr
}
func Capitalize(str string) string {
	var upperStr string
	vv := []rune(str) // 后文有介绍
	for i := 0; i < len(vv); i++ {
		if i == 0 {
			if vv[i] >= 97 && vv[i] <= 122 { // 后文有介绍
				vv[i] -= 32 // string的码表相差32位
				upperStr += string(vv[i])
			} else {
				fmt.Println("Not begins with lowercase letter,")
				return str
			}
		} else {
			upperStr += string(vv[i])
		}
	}
	return upperStr
}
