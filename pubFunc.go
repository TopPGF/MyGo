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
	"math"
	mathRand "math/rand"
	"reflect"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync/atomic"
	"time"

	"github.com/henrylee2cn/faygo"
)

//GetMapValue ...
func GetMapValue(key int64, _map map[int64]string) string {
	value := "--"
	val1, isPresent := _map[key]
	if isPresent {
		value = val1
	}
	return value
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

//CheckNum 判断字符串是否为数字(整数+小数)
func CheckNum(str string) bool {
	rs := true
	var numMinByte, numMaxByte byte
	byteSlice := []byte(str)
	numMinByte = 46 //46 ASCII对应 .
	numMaxByte = 57
	for _, v := range byteSlice {
		if v == 47 {
			//47 ASCII对应 /
			rs = false
			break
		}
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

//CheckInt 判断字符串是否为整数
func CheckInt(str string) bool {
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

func AccountMd5Pw(ps string) string {
	data := []byte("ppfs-" + ps + "-gym")
	has := md5.Sum(data)
	md5Password := fmt.Sprintf("%x", has)
	return md5Password
}

//UniqueId 生成Guid字串
func UniqueId() string {
	b := make([]byte, 48)

	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return ""
	}
	return GetMd5String(base64.URLEncoding.EncodeToString(b))
}

//RemoveRepByLoop slice过滤重复元素（去重）
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

//GetOrderNo 生成订单号(11:会员卡开卡；12：课程预约扣费；13：消费签到扣费；14：课程取消预约；15：消费签到取消)
func GetOrderNo(OrderType string) string {
	year := fmt.Sprintf("%02v", time.Now().Year())
	month := fmt.Sprintf("%02v", int(time.Now().Month()))
	day := fmt.Sprintf("%02v", time.Now().Day())
	hour := fmt.Sprintf("%02v", time.Now().Hour())
	minute := fmt.Sprintf("%02v", time.Now().Minute())
	second := fmt.Sprintf("%02v", time.Now().Second())
	rnd := mathRand.New(mathRand.NewSource(time.Now().UnixNano()))
	rndString := fmt.Sprintf("%05v", rnd.Int31n(100000))
	return OrderType + year[2:] + month + day + hour + minute + second + rndString
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

//字符串日期转换为时间戳
func StrToTime(toBeCharge string) int64 {

	var timeLayout string
	if len(toBeCharge) > 10 {
		timeLayout = "2006-01-02 15:04:05" //转化所需模板
	} else {
		timeLayout = "2006-01-02" //转化所需模板
	}
	loc, _ := time.LoadLocation("Local")                            //重要：获取时区
	theTime, _ := time.ParseInLocation(timeLayout, toBeCharge, loc) //使用模板在对应时区转化为time.time类型
	unix_time := theTime.Unix()                                     //转化为时间戳 类型是int64

	return unix_time
}

//UTC 25/Jul/2018:14:14:41 +0800  转为 时间戳
//需把25/Jul/2018:14:14:41 +0800 格式转换为2006/Jan/01:15:04:05  +080格式
func UtcToTime(toBeCharge string) int64 {
	var timeLayout string
	timeLayout = "2006/Jan/02:15:04:05  +0800" //转化所需模板
	data_arr := strings.Split(toBeCharge[0:11], "/")
	toBeCharge = data_arr[2] + "/" + data_arr[1] + "/" + data_arr[0] + ":" + toBeCharge[12:]
	loc, _ := time.LoadLocation("Local")                            //重要：获取时区
	theTime, _ := time.ParseInLocation(timeLayout, toBeCharge, loc) //使用模板在对应时区转化为time.time类型
	unix_time := theTime.Unix()
	return unix_time
}

//时间戳转换为日期格式
func TimeToStr(unix_time int64, ishms bool) string {
	if ishms {
		return time.Unix(unix_time, 0).Format("2006-01-02 15:04:05")
	} else {
		return time.Unix(unix_time, 0).Format("2006-01-02")
	}

}

//NowData 今日日期
func NowData() int64 {
	year := fmt.Sprintf("%02v", time.Now().Year())
	month := fmt.Sprintf("%02v", int(time.Now().Month()))
	day := fmt.Sprintf("%02v", time.Now().Day())

	return ToInt64(year + month + day)
}

//DataCount 日期天数增减（t:日期； a：天数）
func DataCount(t, a int64, action string) int64 {
	tm2, _ := time.Parse("20060102", ToString(t))
	_n := int64(0)
	if action == "+" {
		_n = tm2.Unix() + a*3600*24
	} else {
		_n = tm2.Unix() - a*3600*24
	}

	tm := time.Unix(_n, 0)
	return ToInt64(tm.Format("20060102"))
}

//DataDifference 日期相差天数
func DataDifference(end, start int64) int64 {
	_e, _ := time.Parse("20060102", ToString(end))
	_s, _ := time.Parse("20060102", ToString(start))
	_n := _e.Unix() - _s.Unix()
	return int64(_n / 3600 / 24)
}

//WeeksList 星期列表
func WeeksList() map[int8]string {
	return map[int8]string{
		1: "星期一",
		2: "星期二",
		3: "星期三",
		4: "星期四",
		5: "星期五",
		6: "星期六",
		0: "星期日",
	}
}

//GetWeekStr ...
func GetWeekStr(weekId int8) string {
	weekStr := ""
	_week := WeeksList()
	for _key, _value := range _week {
		if weekId == _key {
			weekStr = _value
		}
	}
	return weekStr
}

//GetWeekDates 获得一个时间段之间内，是周几的日期列表
func GetWeekDates(startDate string, endDate string, weekId int8) []string {
	_dateStrSlice := make([]string, 0)

	_ttStartDate, _ := time.Parse("2006-01-02", startDate)
	_startWeek := int8(_ttStartDate.Weekday())
	_duration, _ := time.ParseDuration("24h")

	if _startWeek != weekId {
		_m := _startWeek - weekId
		var _n int8
		if _m > 0 {
			_n = 7 - _m
		} else {
			_n = int8(math.Abs(float64(_m)))
		}
		_s := time.Duration(_n)
		_startDateT := _ttStartDate.Add(_duration * _s)
		startDate = _startDateT.Format("2006-01-02")
	}

	_startDateT, _ := time.Parse("2006-01-02", startDate)
	_endDateT, _ := time.Parse("2006-01-02", endDate)
	if _startDateT.After(_endDateT) {
		return _dateStrSlice
	}

	for i := 0; i < 1000; i++ {

		if i == 0 {
			_dateStrSlice = append(_dateStrSlice, startDate)
			continue
		}
		_startDateT = _startDateT.Add(_duration * 7)
		startDate = _startDateT.Format("2006-01-02")
		_startDateT, _ := time.Parse("2006-01-02", startDate)
		if _startDateT.After(_endDateT) {
			break
		}
		_dateStrSlice = append(_dateStrSlice, startDate)
	}
	return _dateStrSlice
}

//GetIncreaseID 并发环境下生成一个增长的id,按需设置局部变量或者全局变量
func GetIncreaseID(ID *uint64) uint64 {
	var n, v uint64
	for {
		v = atomic.LoadUint64(ID)
		n = v + 1
		if atomic.CompareAndSwapUint64(ID, v, n) {
			break
		}
	}
	return n
}

//TempReplace 模板替换
func TempReplace(temp, param string) string {
	paramMap := strings.Split(param, ",")
	reg := regexp.MustCompile(`({{[^}}]+}})`)
	tempReg := reg.FindAllString(temp, -1)
	le := len(paramMap)
	for key, value := range tempReg {
		src := "***"
		if key < le {
			src = paramMap[key]
		}
		temp = strings.Replace(temp, value, src, -1)
	}
	return temp
}

//ValidatorValve 验证 rules："must|int|2-10,string|1-10,date,age,phone,email,cn,en""
func ValidatorValve(v, rules string) (bool, string) {
	rule := strings.Split(rules, ",")
	for _, ru := range rule {
		r := strings.Split(ru, "|")
		bt := []string{"0", "0"}
		if len(r) == 2 {
			bt = strings.Split(r[1], "-")
		}
		//保证bt最少2个元素
		bt = append(bt, "0")
		min := ToInt64(bt[0])
		max := ToInt64(bt[1])
		switch r[0] {
		case "must":
			if len(v) == 0 {
				return false, "不能为空"
			}
		case "price":
			if m, _ := regexp.MatchString(`^(\d+\.\d{2})$`, v); !m && len(v) > 0 {
				return false, "必须为两位小数的价格"
			}
		case "int":
			if !CheckInt(v) && len(v) > 0 {
				return false, "必须为整数"
			}
			if ToInt64(v) < min && min > 0 && len(v) > 0 {
				return false, "必须大于" + bt[0]
			}
			if ToInt64(v) > max && max > 0 && len(v) > 0 {
				return false, "必须小于" + bt[1]
			}
		case "num": //可以带小数点
			if !CheckNum(v) && len(v) > 0 {
				return false, "必须为数值类型"
			}
			v1, err := strconv.ParseFloat(v, 64)
			if err != nil {
				return false, "必须为数值类型"
			}
			min, err := strconv.ParseFloat(bt[0], 64)
			if err != nil {
				min = float64(0)
			}
			max, err := strconv.ParseFloat(bt[1], 64)
			if err != nil {
				max = float64(0)
			}
			if v1 < min && min > 0 && len(v) > 0 {
				return false, "必须大于" + bt[0]
			}
			if v1 > max && max > 0 && len(v) > 0 {
				return false, "必须小于" + bt[1]
			}
		case "date":
			//兼容elex 单元格格式
			if CheckInt(v) {
				return true, ""
			}
			_, err := String2Time(v)
			if err != nil && len(v) > 0 {
				return false, "格式错误"
			}
		case "string":
			if ToInt64(len(v)) < min && min > 0 && len(v) > 0 {
				return false, "字符数必须大于" + bt[0]
			}
			if ToInt64(len(v)) > max && max > 0 && len(v) > 0 {
				return false, "字符数必须小于" + bt[1]
			}
		case "phone":
			if m, _ := regexp.MatchString(`^(1[3|4|5|8][0-9]\d{4,8})$`, v); !m && len(v) > 0 {
				return false, "手机号码不正确"
			}
		case "email":
			if m, _ := regexp.MatchString(`^([\w\.\_]{2,10})@(\w{1,}).([a-z]{2,4})$`, v); !m && len(v) > 0 {
				return false, "电子邮件格式不正确"
			}
		case "age":
			if ToInt64(v) > 140 || ToInt64(v) < 0 && len(v) > 0 {
				return false, "年龄取值范围在0 ~ 140岁之间"
			}
		case "cn":
			if m, _ := regexp.MatchString("^\\p{Han}+$", v); !m && len(v) > 0 {
				return false, "必须为中文汉字"
			}
		case "en":
			if m, _ := regexp.MatchString("^[a-zA-Z]+$", v); !m && len(v) > 0 {
				return false, "必须为英文字母"
			}
		}
	}
	return true, ""
}
func String2Time(in string) (out time.Time, err error) {
	in = strings.Replace(in, "/", "-", -1)
	if len(in) > 10 {
		out, err = time.Parse("2006-01-02 15:04:05", in) //layout使用"2006/01/02 15:04:05"此数据格式转换会出错
	} else {
		out, err = time.Parse("2006-01-02", in) //layout使用"2006/01/02"此数据格式转换会出错
	}
	return out, err
}
