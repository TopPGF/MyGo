/**
 * 函数
 * ---------------------------------------------------------------------
 * @author		PGF <pgf.@fealive.cn>
 * @date		2017-12-02
 * ---------------------------------------------------------------------
 */

package logic

import (
	"bytes"
	"crypto/md5"
	"crypto/rand"
	"encoding/base64"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"io"
	mathRand "math/rand"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/henrylee2cn/faygo"
)

//GetSQLID ...
func GetSQLID(key string) int64 {
	_ = key
	return time.Now().Unix()
}

//GetMapValue ...
func GetMapValue(key int64, _map map[int64]string) string {
	value := "--"
	val1, isPresent := _map[key]
	if isPresent {
		value = val1
	}
	return value
}

//ErrorJSON 统一的错误码返回接口格式
func ErrorJSON(errorCode interface{}, data ...interface{}) (int, map[string]interface{}) {
	var errorMsg string
	var data2 interface{}
	//fmt.Println(data)
	if len(data) > 0 {
		//fmt.Println(data[0])
		errorMsg = data[0].(string)

	} else {
		errorMsg = ErrCode[ToInt(errorCode)]
	}
	//
	if len(data) > 1 {
		data2 = data[1]
	} else {
		data2 = nil
	}
	statusString := strings.Join(strings.Split(ToString(errorCode), "")[0:3], "")
	status, _ := strconv.Atoi(statusString)
	return status, map[string]interface{}{
		"request":     "",
		"return_code": "SUCCESS",
		"is_error":    true,

		"error_code": statusString,
		"api_code":   errorCode,
		"error":      errorMsg,
		"data":       data2,
	}
}

//RsData ...
func RsData(ctx *faygo.Context, rsData interface{}) (int, interface{}) {
	return 200, map[string]interface{}{
		"request":     ctx.URI(),
		"return_code": "SUCCESS",
		"is_error":    false,

		"data": rsData,
	}
}

//RsDataDoc ...
func RsDataDoc(note string, rsData interface{}) faygo.Doc {
	return faygo.Doc{
		// 向API文档声明接口注意事项
		Note: note,
		// 向API文档声明响应内容格式
		Return: map[string]interface{}{
			"request":     "",
			"return_code": "SUCCESS",
			"is_error":    false,

			"data": rsData,
		},
	}
}

//DeepFields ...
func DeepFields(ifaceType reflect.Type) []reflect.StructField {
	var fields []reflect.StructField

	for i := 0; i < ifaceType.NumField(); i++ {
		v := ifaceType.Field(i)
		if v.Anonymous && v.Type.Kind() == reflect.Struct {
			fields = append(fields, DeepFields(v.Type)...)
		} else {
			fields = append(fields, v)
		}
	}

	return fields
}

//StructCopy ...
func StructCopy(DstStructPtr interface{}, SrcStructPtr interface{}) {
	srcv := reflect.ValueOf(SrcStructPtr)
	dstv := reflect.ValueOf(DstStructPtr)
	srct := reflect.TypeOf(SrcStructPtr)
	dstt := reflect.TypeOf(DstStructPtr)
	if srct.Kind() != reflect.Ptr || dstt.Kind() != reflect.Ptr ||
		srct.Elem().Kind() == reflect.Ptr || dstt.Elem().Kind() == reflect.Ptr {
		panic("Fatal error:type of parameters must be Ptr of value")
	}
	if srcv.IsNil() || dstv.IsNil() {
		panic("Fatal error:value of parameters should not be nil")
	}
	srcV := srcv.Elem()
	dstV := dstv.Elem()
	srcfields := DeepFields(reflect.ValueOf(SrcStructPtr).Elem().Type())
	for _, v := range srcfields {
		if v.Anonymous {
			continue
		}
		dst := dstV.FieldByName(v.Name)
		src := srcV.FieldByName(v.Name)
		if !dst.IsValid() {
			continue
		}
		if src.Type() == dst.Type() && dst.CanSet() {
			dst.Set(src)
			continue
		}
		if src.Kind() == reflect.Ptr && !src.IsNil() && src.Type().Elem() == dst.Type() {
			dst.Set(src.Elem())
			continue
		}
		if dst.Kind() == reflect.Ptr && dst.Type().Elem() == src.Type() {
			dst.Set(reflect.New(src.Type()))
			dst.Elem().Set(src)
			continue
		}
	}
	return
}

// DeepCopy 结构体复制
func DeepCopy(dst, src interface{}) error {
	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(src); err != nil {
		return err
	}
	return gob.NewDecoder(bytes.NewBuffer(buf.Bytes())).Decode(dst)
}

//ToInt ...
func ToInt(d interface{}) int {
	switch n := d.(type) {
	case int:
		return n
	case string:
		i, err := strconv.Atoi(n)
		if err != nil {
			panic(n + err.Error())
		}
		return i
	case int64:
		return int(n)
	case float32:
		return int(n)
	case float64:
		return int(n)
	case bool:
		if n {
			return 1
		}
		return 0
	case []byte:
		i, err := strconv.Atoi(string(n))
		if err != nil {
			panic(err.Error())
		}
		return i
	}
	return 0
}

//ToInt64 ...
func ToInt64(d interface{}) int64 {
	switch n := d.(type) {
	case int64:
		return n
	case string:
		if n == "" {
			return 0
		}
		i, err := strconv.Atoi(n)
		if err != nil {
			panic(n + err.Error())
		}
		return int64(i)
	case int:
		return int64(n)
	case float32:
		return int64(n)
	case float64:
		return int64(n)
	case bool:
		if n {
			return 1
		}
		return 0
	case []byte:
		i, err := strconv.Atoi(string(n))
		if err != nil {
			panic(err.Error())
		}
		return int64(i)
	}
	return 0
}

//ToString ...
func ToString(d interface{}) string {
	switch n := d.(type) {
	case int64:
		return strconv.FormatInt(n, 10)
	case string:
		return n
	case int:
		return strconv.Itoa(n)
	case float32:
		return strconv.Itoa(int(n))
	case float64:
		return strconv.Itoa(int(n))
	case []byte:
		return string(n)
	case bool:
		if n {
			return "1"
		}
		return "0"
	}
	return "0"
}

//CheckNum 判断字符串是否为数字
func CheckNum(str string) bool {
	rs := true
	var numMinByte, numMaxByte byte
	byteSlice := []byte(str)
	numMinByte = 48
	numMaxByte = 57
	for _, v := range byteSlice {
		if v < numMinByte {
			rs = false
			break
		}
		if v > numMaxByte {
			rs = false
			break
		}
	}
	return rs
}

//SetFiledToModel 通过反射，"设置"某个“数据表模型”结构体中单个字段的值
func SetFiledToModel(table interface{}, filed string, value interface{}) {
	f := reflect.ValueOf(table).Elem().FieldByName(filed)
	if f.IsValid() {
		f.Set(reflect.Value(reflect.ValueOf(value)))
	}
}

//GetFiledFromModel 通过反射，“获取”某个“数据表模型”结构体中单个字段的值
func GetFiledFromModel(filed string, table interface{}) interface{} {
	switch c := table.(type) {
	case map[string]interface{}:
		return c[filed]
	default:
		f := reflect.ValueOf(table).Elem().FieldByName(filed)
		if f.IsValid() {
			return f.Interface()
		}
		return nil
	}
}

//SnakeCasedName ...
func SnakeCasedName(name string) string {
	newstr := make([]rune, 0)
	for idx, chr := range name {
		if isUpper := 'A' <= chr && chr <= 'Z'; isUpper {
			if idx > 0 {
				newstr = append(newstr, '_')
			}
			chr -= ('A' - 'a')
		}
		newstr = append(newstr, chr)
	}

	return string(newstr)
}

//TitleCasedName ...
func TitleCasedName(name string) string {
	newstr := make([]rune, 0)
	upNextChar := true

	name = strings.ToLower(name)

	for _, chr := range name {
		switch {
		case upNextChar:
			upNextChar = false
			if 'a' <= chr && chr <= 'z' {
				chr -= ('a' - 'A')
			}
		case chr == '_':
			upNextChar = true
			continue
		}

		newstr = append(newstr, chr)
	}

	return string(newstr)
}

//MapToSlice map 转 slice
func MapToSlice(m map[int64]int64) []int64 {
	s := make([]int64, 0, len(m))
	for _, v := range m {
		s = append(s, v)
	}
	return s
}

//Md5Encode ...
func Md5Encode(s string) string {
	h := md5.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}

func XlsxCellIndex(col int, row int) string {
	colMap := strings.Split("ABCDEFGHIJKLMNOPQRSTUVWXYZ", "")
	return colMap[col] + ToString(row)
}

func AppendSliceNoRepeat(s []int, x int) []int {
	leng := len(s)
	if leng < 1 {
		return append(s, x)
	}
	sort.Sort(sort.IntSlice(s))
	index := sort.Search(leng, func(i int) bool {
		return s[i] >= x
	})
	//fmt.Println(index, s)

	if len(s) > index {
		if s[index] == x {
			return s
		}
	}
	s = append(s, x)
	sort.Sort(sort.IntSlice(s))
	return s
}

//GetMd5String 生成32位md5字串
func GetMd5String(s string) string {
	h := md5.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}

//UniqueId 生成Guid字串
func UniqueId() string {
	b := make([]byte, 48)

	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return ""
	}
	return GetMd5String(base64.URLEncoding.EncodeToString(b))
}

//RemoveRepByLoop slice过滤重复元素
func RemoveRepByLoop(slc []int64) []int64 {
	result := []int64{} // 存放结果
	for i := range slc {
		flag := true
		for j := range result {
			if slc[i] == result[j] {
				flag = false // 存在重复元素，标识为false
				break
			}
		}
		if flag { // 标识为false，不添加进结果
			result = append(result, slc[i])
		}
	}
	return result
}

//GetOrderNo 生成订单号
func GetOrderNo(OrderType int64) string {
	year := fmt.Sprintf("%02v", time.Now().Year())
	month := fmt.Sprintf("%02v", int(time.Now().Month()))
	day := fmt.Sprintf("%02v", time.Now().Day())
	hour := fmt.Sprintf("%02v", time.Now().Hour())
	minute := fmt.Sprintf("%02v", time.Now().Minute())
	second := fmt.Sprintf("%02v", time.Now().Second())
	rnd := mathRand.New(mathRand.NewSource(time.Now().UnixNano()))
	rndString := fmt.Sprintf("%05v", rnd.Int31n(100000))
	return fmt.Sprintf("%02v", OrderType) + year[2:] + month + day + hour + minute + second + rndString
}

//ToData 格式化日期
func ToData(d int64) string {
	data := fmt.Sprintf("%08v", d)
	data = data[0:4] + "-" + data[4:6] + "-" + data[6:]
	if d == 0 {
		data = ""
	}
	return data
}
