package utils

import (
	"bufio"
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"ultagic.com/pkg/flags"
)

var (
	BasePath  string
	loggerMap map[string]*log.Logger
)

func init() {
	BasePath, _ = os.Getwd()
	loggerMap = make(map[string]*log.Logger, 10)
}

func Logger(name string) *log.Logger {
	if loggerMap[name] == nil {
		file := filepath.Join(flags.LogPath, "main.log")
		logFile, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0766)
		if err != nil {
			panic(err)
		}
		loggerMap[name] = log.New(logFile, "["+name+"] ", log.LstdFlags|log.LUTC)
	}

	return loggerMap[name]
}

func ParseIP(s string) (net.IP, int) {
	ip := net.ParseIP(s)
	if ip == nil {
		return nil, 0
	}

	for i := 0; i < len(s); i++ {
		switch s[i] {
		case '.':
			return ip, 4
		case ':':
			return ip, 6
		}
	}

	return nil, 0
}

func MD5(s string) string {
	sum := md5.Sum([]byte(s))
	return hex.EncodeToString(sum[:])
}

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func getStoragePath(app string, fileName string) string {
	// 拼接路径
	var buffer bytes.Buffer
	buffer.WriteString(BasePath)
	buffer.WriteString("/storage/")
	buffer.WriteString(app)
	buffer.WriteString("/")
	buffer.WriteString(fileName)

	return buffer.String()
}

func GetStorageContent(app string, fileName string) (string, error) {
	// 读取文件
	content, err := ioutil.ReadFile(getStoragePath(app, fileName))

	return string(content), err
}

func GetStorageJSON(app string, fileName string, v interface{}) error {
	// 读取文件
	content, err := ioutil.ReadFile(getStoragePath(app, fileName))
	if err != nil {
		return err
	}

	if err := json.Unmarshal(content, &v); err != nil {
		return err
	}

	return err
}

func SetStorageContent(app string, fileName string, content string) error {
	// 写入文件
	return ioutil.WriteFile(getStoragePath(app, fileName), []byte(content), 0666)
}

func AppendStorageContent(app string, fileName string, content string) error {
	file, err := os.OpenFile(getStoragePath(app, fileName), os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		return err
	}

	defer file.Close()

	writer := bufio.NewWriter(file)

	_, err = writer.WriteString(content + "\n")
	if err != nil {
		return err
	}

	writer.Flush()

	return nil
}

func structToMap(obj interface{}, lowerKey bool) map[string]interface{} {
	t := reflect.TypeOf(obj)
	v := reflect.ValueOf(obj)

	var data = make(map[string]interface{})
	for i := 0; i < t.NumField(); i++ {
		name := t.Field(i).Name
		if lowerKey {
			name = strings.ToLower(name)
		}
		data[name] = v.Field(i).Interface()
	}

	return data
}

func StructToMap(obj interface{}) map[string]interface{} {
	return structToMap(obj, false)
}

func StructToMapWithLowerKey(obj interface{}) map[string]interface{} {
	return structToMap(obj, true)
}

func StructCopyFields(a interface{}, b interface{}, fields ...string) error {
	at := reflect.TypeOf(a)
	av := reflect.ValueOf(a)
	bt := reflect.TypeOf(b)
	bv := reflect.ValueOf(b)

	if at.Kind() != reflect.Ptr {
		return errors.New("a must be a struct pointer")
	}
	av = reflect.ValueOf(av.Interface())

	_fields := make([]string, 0)
	if len(fields) > 0 {
		_fields = fields
	} else {
		for i := 0; i < bv.NumField(); i++ {
			_fields = append(_fields, bt.Field(i).Name)
		}
	}

	for i := 0; i < len(_fields); i++ {
		name := _fields[i]
		f := av.Elem().FieldByName(name)
		bValue := bv.FieldByName(name)

		// a 中有同名的字段并且类型一致才复制
		if f.IsValid() && f.Type() == bValue.Type() {
			f.Set(bValue)
		}
	}
	return nil
}

func GetStructFieldName(s interface{}) ([]string, error) {
	t := reflect.TypeOf(s)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	if t.Kind() != reflect.Struct {
		return nil, errors.New("not Struct")
	}

	fieldNum := t.NumField()
	result := make([]string, 0, fieldNum)
	for i := 0; i < fieldNum; i++ {
		result = append(result, t.Field(i).Name)
	}

	return result, nil
}

func GetStructFieldNameToSnake(s interface{}) ([]string, error) {
	t := reflect.TypeOf(s)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	if t.Kind() != reflect.Struct {
		return nil, errors.New("not Struct")
	}

	fieldNum := t.NumField()
	result := make([]string, 0, fieldNum)
	for i := 0; i < fieldNum; i++ {
		result = append(result, CamelToSnake(t.Field(i).Name))
	}

	return result, nil
}

/**
 * 驼峰转蛇形
 * @description XxYy to xx_yy , XxYY to xx_y_y
 * @date 2020/7/30
 * @param s 需要转换的字符串
 * @return string
 **/
func CamelToSnake(s string) string {
	data := make([]byte, 0, len(s)*2)
	j := false
	num := len(s)
	for i := 0; i < num; i++ {
		d := s[i]
		// or通过ASCII码进行大小写的转化
		// 65-90（A-Z），97-122（a-z）
		//判断如果字母为大写的A-Z就在前面拼接一个_
		if i > 0 && d >= 'A' && d <= 'Z' && j {
			data = append(data, '_')
		}
		if d != '_' {
			j = true
		}
		data = append(data, d)
	}

	//ToLower把大写字母统一转小写
	return strings.ToLower(string(data[:]))
}

/**
 * 蛇形转驼峰
 * @description xx_yy to XxYx  xx_y_y to XxYY
 * @date 2020/7/30
 * @param s要转换的字符串
 * @return string
 **/
func SnakeToCamel(s string) string {
	data := make([]byte, 0, len(s))
	j := false
	k := false
	num := len(s) - 1
	for i := 0; i <= num; i++ {
		d := s[i]
		if !k && d >= 'A' && d <= 'Z' {
			k = true
		}
		if d >= 'a' && d <= 'z' && (j || !k) {
			d = d - 32
			j = false
			k = true
		}
		if k && d == '_' && num > i && s[i+1] >= 'a' && s[i+1] <= 'z' {
			j = true
			continue
		}
		data = append(data, d)
	}

	return string(data[:])
}
