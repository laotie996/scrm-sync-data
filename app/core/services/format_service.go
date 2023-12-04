package services

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"context"
	"database/sql"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	uuid "github.com/google/uuid"
	jsoniter "github.com/json-iterator/go"
	"io"
	"math/rand"
	"os"
	"scrm-sync-data/app/config"
	"strconv"
	"strings"
	"time"
)

type FormatService struct {
	context context.Context
	cancel  context.CancelFunc
	config  *config.Config
	logger  *LoggerService
	State   bool
}

func (formatService *FormatService) Init(parentContext context.Context, config *config.Config, logger *LoggerService) {
	formatService.config = config
	formatService.logger = logger
	formatService.State = false
	formatService.context, formatService.cancel = context.WithCancel(parentContext)
	formatService.Start()
}

func (formatService *FormatService) Start() {
	fmt.Println("start format service...", time.Now())
	formatService.logger.Debug(fmt.Sprintf("%s,%v", "start format service...", time.Now()))
	formatService.State = true
}

func (formatService *FormatService) Stop() {
	fmt.Println("stop format service...", time.Now())
	formatService.logger.Debug(fmt.Sprintf("%s,%v", "stop format service...", time.Now()))
	formatService.cancel()
	formatService.State = false
}

func (formatService *FormatService) ToByte(v interface{}) []byte {
	switch v.(type) {
	case string:
		return []byte(v.(string))
	case []byte:
		return v.([]byte)
	case int:
		bytesBuffer := bytes.NewBuffer([]byte{})
		_ = binary.Write(bytesBuffer, binary.BigEndian, int64(v.(int))) //没有int型转byte方法
		return bytesBuffer.Bytes()
	case int8:
		bytesBuffer := bytes.NewBuffer([]byte{})
		_ = binary.Write(bytesBuffer, binary.BigEndian, v.(int8))
		return bytesBuffer.Bytes()
	case int16:
		bytesBuffer := bytes.NewBuffer([]byte{})
		_ = binary.Write(bytesBuffer, binary.BigEndian, v.(int16))
		return bytesBuffer.Bytes()
	case int32:
		bytesBuffer := bytes.NewBuffer([]byte{})
		_ = binary.Write(bytesBuffer, binary.BigEndian, v.(int32))
		return bytesBuffer.Bytes()
	case int64:
		bytesBuffer := bytes.NewBuffer([]byte{})
		_ = binary.Write(bytesBuffer, binary.BigEndian, v.(int64))
		return bytesBuffer.Bytes()
	case float32:
		bytesBuffer := bytes.NewBuffer([]byte{})
		_ = binary.Write(bytesBuffer, binary.BigEndian, v.(float32))
		return bytesBuffer.Bytes()
	case float64:
		bytesBuffer := bytes.NewBuffer([]byte{})
		_ = binary.Write(bytesBuffer, binary.BigEndian, v.(float64))
		return bytesBuffer.Bytes()
	default:
		return nil
	}
}

func (formatService *FormatService) ToInt(v interface{}) int {
	switch v.(type) {
	case string:
		i, _ := strconv.Atoi(v.(string))
		return i
	case []byte:
		var i int
		bytesBuffer := bytes.NewBuffer(v.([]byte))
		_ = binary.Read(bytesBuffer, binary.BigEndian, &i)
		return i
	case int:
		return v.(int)
	case int8:
		return int(v.(int8))
	case int16:
		return int(v.(int16))
	case int32:
		return int(v.(int32))
	case int64:
		return int(v.(int64))
	case float32:
		return int(v.(float32))
	case float64:
		return int(v.(float64))
	default:
		return 0
	}
}

func (formatService *FormatService) ToUint(v interface{}) uint {
	switch v.(type) {
	case string:
		i, _ := strconv.Atoi(v.(string))
		return uint(i)
	case []byte:
		var i int
		bytesBuffer := bytes.NewBuffer(v.([]byte))
		_ = binary.Read(bytesBuffer, binary.BigEndian, &i)
		return uint(i)
	case int:
		return uint(v.(int))
	case int8:
		return uint(v.(int8))
	case int16:
		return uint(v.(int16))
	case int32:
		return uint(v.(int32))
	case int64:
		return uint(v.(int64))
	case float32:
		return uint(v.(float32))
	case float64:
		return uint(v.(float64))
	default:
		return 0
	}
}

func (formatService *FormatService) ToInt8(v interface{}) int8 {
	switch v.(type) {
	case string:
		i, _ := strconv.ParseInt(v.(string), 10, 8)
		return int8(i)
	case []byte:
		var i int8
		bytesBuffer := bytes.NewBuffer(v.([]byte))
		_ = binary.Read(bytesBuffer, binary.BigEndian, &i)
		return i
	case int:
		return int8(v.(int))
	case int8:
		return v.(int8)
	case int16:
		return int8(v.(int16))
	case int32:
		return int8(v.(int32))
	case int64:
		return int8(v.(int64))
	case float32:
		return int8(v.(float32))
	case float64:
		return int8(v.(float64))
	default:
		return 0
	}
}
func (formatService *FormatService) ToUint8(v interface{}) uint8 {
	switch v.(type) {
	case string:
		i, _ := strconv.ParseInt(v.(string), 10, 8)
		return uint8(i)
	case []byte:
		var i int8
		bytesBuffer := bytes.NewBuffer(v.([]byte))
		_ = binary.Read(bytesBuffer, binary.BigEndian, &i)
		return uint8(i)
	case int:
		return uint8(v.(int))
	case int8:
		return uint8(v.(int8))
	case int16:
		return uint8(v.(int16))
	case int32:
		return uint8(v.(int32))
	case int64:
		return uint8(v.(int64))
	case float32:
		return uint8(v.(float32))
	case float64:
		return uint8(v.(float64))
	default:
		return 0
	}
}

func (formatService *FormatService) ToInt16(v interface{}) int16 {
	switch v.(type) {
	case string:
		i, _ := strconv.ParseInt(v.(string), 10, 16)
		return int16(i)
	case []byte:
		var i int16
		bytesBuffer := bytes.NewBuffer(v.([]byte))
		_ = binary.Read(bytesBuffer, binary.BigEndian, &i)
		return i
	case int:
		return int16(v.(int))
	case int8:
		return int16(v.(int8))
	case int16:
		return v.(int16)
	case int32:
		return int16(v.(int32))
	case int64:
		return int16(v.(int64))
	case float32:
		return int16(v.(float32))
	case float64:
		return int16(v.(float64))
	default:
		return 0
	}
}
func (formatService *FormatService) ToUint16(v interface{}) uint16 {
	switch v.(type) {
	case string:
		i, _ := strconv.ParseInt(v.(string), 10, 16)
		return uint16(i)
	case []byte:
		var i int16
		bytesBuffer := bytes.NewBuffer(v.([]byte))
		_ = binary.Read(bytesBuffer, binary.BigEndian, &i)
		return uint16(i)
	case int:
		return uint16(v.(int))
	case int8:
		return uint16(v.(int8))
	case int16:
		return uint16(v.(int16))
	case int32:
		return uint16(v.(int32))
	case int64:
		return uint16(v.(int64))
	case float32:
		return uint16(v.(float32))
	case float64:
		return uint16(v.(float64))
	default:
		return 0
	}
}

func (formatService *FormatService) ToInt32(v interface{}) int32 {
	switch v.(type) {
	case string:
		i, _ := strconv.ParseInt(v.(string), 10, 32)
		return int32(i)
	case []byte:
		var i int32
		bytesBuffer := bytes.NewBuffer(v.([]byte))
		_ = binary.Read(bytesBuffer, binary.BigEndian, &i)
		return i
	case int:
		return int32(v.(int))
	case int8:
		return int32(v.(int8))
	case int16:
		return int32(v.(int16))
	case int32:
		return v.(int32)
	case int64:
		return int32(v.(int64))
	case float32:
		return int32(v.(float32))
	case float64:
		return int32(v.(float64))
	case uint64:
		return int32(v.(uint64))
	default:
		return 0
	}
}
func (formatService *FormatService) ToUint32(v interface{}) uint32 {
	switch v.(type) {
	case string:
		i, _ := strconv.ParseInt(v.(string), 10, 32)
		return uint32(i)
	case []byte:
		var i int32
		bytesBuffer := bytes.NewBuffer(v.([]byte))
		_ = binary.Read(bytesBuffer, binary.BigEndian, &i)
		return uint32(i)
	case int:
		return uint32(v.(int))
	case int8:
		return uint32(v.(int8))
	case int16:
		return uint32(v.(int16))
	case int32:
		return uint32(v.(int32))
	case int64:
		return uint32(v.(int64))
	case float32:
		return uint32(v.(float32))
	case float64:
		return uint32(v.(float64))
	default:
		return 0
	}
}

func (formatService *FormatService) ToInt64(v interface{}) int64 {
	switch v.(type) {
	case string:
		i, _ := strconv.ParseInt(v.(string), 10, 64)
		return i
	case []byte:
		var i int64
		bytesBuffer := bytes.NewBuffer(v.([]byte))
		_ = binary.Read(bytesBuffer, binary.BigEndian, &i)
		return i
	case int:
		return int64(v.(int))
	case int8:
		return int64(v.(int8))
	case int16:
		return int64(v.(int16))
	case int32:
		return int64(v.(int32))
	case int64:
		return v.(int64)
	case float32:
		return int64(v.(float32))
	case float64:
		return int64(v.(float64))
	case uint64:
		return int64(v.(uint64))
	default:
		return 0
	}
}

func (formatService *FormatService) ToUint64(v interface{}) uint64 {
	switch v.(type) {
	case string:
		i, _ := strconv.ParseInt(v.(string), 10, 64)
		return uint64(i)
	case []byte:
		var i int64
		bytesBuffer := bytes.NewBuffer(v.([]byte))
		_ = binary.Read(bytesBuffer, binary.BigEndian, &i)
		return uint64(i)
	case int:
		return uint64(v.(int))
	case int8:
		return uint64(v.(int8))
	case int16:
		return uint64(v.(int16))
	case int32:
		return uint64(v.(int32))
	case int64:
		return uint64(v.(int64))
	case float32:
		return uint64(v.(float32))
	case float64:
		return uint64(v.(float64))
	default:
		return 0
	}
}

func (formatService *FormatService) ToFloat32(v interface{}) float32 {
	switch v.(type) {
	case string:
		f, _ := strconv.ParseFloat(v.(string), 32)
		return float32(f)
	case []byte:
		var f float32
		bytesBuffer := bytes.NewBuffer(v.([]byte))
		_ = binary.Read(bytesBuffer, binary.BigEndian, &f)
		return f
	case int:
		return float32(v.(int))
	case int8:
		return float32(v.(int8))
	case int16:
		return float32(v.(int16))
	case int32:
		return float32(v.(int32))
	case int64:
		return float32(v.(int64))
	case float32:
		return v.(float32)
	case float64:
		return float32(v.(float64))
	default:
		return 0
	}
}

func (formatService *FormatService) ToFloat64(v interface{}) float64 {
	switch v.(type) {
	case string:
		f, _ := strconv.ParseFloat(v.(string), 64)
		return f
	case []byte:
		var f float64
		bytesBuffer := bytes.NewBuffer(v.([]byte))
		_ = binary.Read(bytesBuffer, binary.BigEndian, &f)
		return f
	case int:
		return float64(v.(int))
	case int8:
		return float64(v.(int8))
	case int16:
		return float64(v.(int16))
	case int32:
		return float64(v.(int32))
	case int64:
		return float64(v.(int64))
	case float32:
		return float64(v.(float32))
	case float64:
		return v.(float64)
	default:
		return 0
	}
}

func (formatService *FormatService) ToString(v interface{}) string {
	switch v.(type) {
	case string:
		return v.(string)
	case []byte:
		return string(v.([]byte))
	case int:
		i, _ := v.(int)
		return strconv.Itoa(i)
	case int8:
		i, _ := v.(int8)
		return strconv.Itoa(int(i))
	case int16:
		i, _ := v.(int16)
		return strconv.Itoa(int(i))
	case int32:
		i, _ := v.(int32)
		return strconv.Itoa(int(i))
	case int64:
		i, _ := v.(int64)
		return strconv.FormatInt(i, 10)
	case float32:
		f, _ := v.(float32)
		return strconv.FormatFloat(float64(f), 'E', -1, 32)
	case float64:
		f, _ := v.(float64)
		return strconv.FormatFloat(f, 'E', -1, 64)
	default:
		return ""
	}
}

func (formatService *FormatService) ToSqlNullString(v interface{}) sql.NullString {
	return sql.NullString{String: formatService.ToString(v), Valid: true}
}

func (formatService *FormatService) ToSqlNullInt64(v interface{}) sql.NullInt64 {
	return sql.NullInt64{Int64: formatService.ToInt64(v), Valid: true}
}

func (formatService *FormatService) ToSqlNullFloat64(v interface{}) sql.NullFloat64 {
	return sql.NullFloat64{Float64: formatService.ToFloat64(v), Valid: true}
}

func (formatService *FormatService) Split(buf []byte, lim int) [][]byte {
	var chunk []byte
	chunks := make([][]byte, 0, len(buf)/lim+1)
	for len(buf) >= lim {
		chunk, buf = buf[:lim], buf[lim:]
		chunks = append(chunks, chunk)
	}
	if len(buf) > 0 {
		chunks = append(chunks, buf)
	}
	return chunks
}

func (formatService *FormatService) StringArrayRemove(str string, array []string) []string {
	var newArray []string
	for index := range array {
		if array[index] != str {
			newArray = append(newArray, array[index])
		}
	}
	return newArray
}

func (formatService *FormatService) UUID() string {
	u := uuid.New()
	return fmt.Sprintf("%s", u)
}

func (formatService *FormatService) ShortUUID(length int32) string {
	u := uuid.New()
	id := u.String()
	if length > 36 || length < 0 {
		length = 36
	}
	b := []byte(id[:length])
	return base64.RawURLEncoding.EncodeToString(b)
}

func (formatService *FormatService) ParseFormData(formData string) map[string]string {
	var formMap = make(map[string]string)
	formItems := strings.Split(formData, "&")
	for index := range formItems {
		formItem := strings.Split(formItems[index], "=")
		if len(formItem) == 0 {
			continue
		}
		if len(formItem) > 1 {
			formMap[formItem[0]] = formItem[1]
		} else {
			formMap[formItem[0]] = ""
		}
	}
	return formMap
}

// RandomString
// @Description: 生成随机字符串
// @Param size body integer false "随机字符串长度"
// @Param kind body integer false "随机字符串类型 0 纯数字 1 小写字母 2 大写字母 3 数字和字母组合"
// @Return string
func (formatService *FormatService) RandomString(size int, kind int) string {
	strKind, strKinds, buffer := kind, [][]int{{10, 48}, {26, 97}, {26, 65}}, make([]byte, size)
	isAll := kind > 2 || kind < 0
	for i := 0; i < size; i++ {
		if isAll {
			strKind = rand.Intn(3)
		}
		scope, base := strKinds[strKind][0], strKinds[strKind][1]
		buffer[i] = uint8(base + rand.Intn(scope))
	}
	return string(buffer)
}

// RandomTimestamp
// @Description: 生成随机时间戳
// @Return string
func (formatService *FormatService) RandomTimestamp(t time.Time) int64 {
	//return time.Now().Unix() + 5
	return t.Unix() + formatService.ToInt64(rand.Intn(68)*900+108000)
}

func (formatService *FormatService) DateStrToTime(dateStr string) time.Time {
	loc, _ := time.LoadLocation("Asia/Shanghai")
	t, _ := time.ParseInLocation(time.DateTime, dateStr, loc)
	return t
}

func (formatService *FormatService) DateStrToDate(dateStr string) string {
	loc, _ := time.LoadLocation("Asia/Shanghai")
	t, _ := time.ParseInLocation(time.RFC3339, dateStr, loc)
	return t.Format(time.DateTime)
}

func (formatService *FormatService) CompareTwoDate(d1, d2 string) bool {
	if len(d1) == 10 {
		d1 = d1 + " 00:00:00"
	}
	if len(d2) == 10 {
		d2 = d2 + " 00:00:00"
	}
	return formatService.DateStrToTime(d1).Sub(formatService.DateStrToTime(d2)) > 0
}

func (formatService *FormatService) GetAge(dateStr string) (age int) {
	birthdayTime := formatService.DateStrToTime(dateStr)
	yearGap := time.Now().Year() - birthdayTime.Year()
	monthGap := int(time.Now().Month() - birthdayTime.Month())
	dayGap := time.Now().Day() - birthdayTime.Day()
	fmt.Println(yearGap, monthGap, dayGap)
	if yearGap > 0 {
		age = yearGap
	}
	if monthGap > 0 {
		age = yearGap + 1
	} else if monthGap == 0 {
		if dayGap > 0 {
			age = yearGap + 1
		}
	}
	return age
}

var quoteEscape = strings.NewReplacer("\\", "\\\\", `"`, "\\\"")

func (formatService *FormatService) EscapeQuotes(s string) string {
	return quoteEscape.Replace(s)
}

const compress = false

func (formatService *FormatService) MustReadStdin() string {
	r := bufio.NewReader(os.Stdin)

	var in string
	for {
		var err error
		in, err = r.ReadString('\n')
		if err != io.EOF {
			if err != nil {
				panic(err)
			}
		}
		in = strings.TrimSpace(in)
		if len(in) > 0 {
			break
		}
	}

	fmt.Println("")

	return in
}

// Encode encodes the input in base64
// It can optionally zip the input before encoding
func (formatService *FormatService) Encode(obj interface{}) string {
	b, err := jsoniter.Marshal(obj)
	if err != nil {
		panic(err)
	}

	if compress {
		b = formatService.Zip(b)
	}

	return base64.StdEncoding.EncodeToString(b)
}

// Decode decodes the input from base64
// It can optionally unzip the input after decoding
func (formatService *FormatService) Decode(in string, obj interface{}) {
	b, err := base64.StdEncoding.DecodeString(in)
	if err != nil {
		panic(err)
	}

	if compress {
		b = formatService.Unzip(b)
	}

	err = jsoniter.Unmarshal(b, obj)
	if err != nil {
		panic(err)
	}
}

func (formatService *FormatService) Zip(in []byte) []byte {
	var b bytes.Buffer
	gz := gzip.NewWriter(&b)
	_, err := gz.Write(in)
	if err != nil {
		panic(err)
	}
	err = gz.Flush()
	if err != nil {
		panic(err)
	}
	err = gz.Close()
	if err != nil {
		panic(err)
	}
	return b.Bytes()
}

func (formatService *FormatService) Unzip(in []byte) []byte {
	var b bytes.Buffer
	_, err := b.Write(in)
	if err != nil {
		panic(err)
	}
	r, err := gzip.NewReader(&b)
	if err != nil {
		panic(err)
	}
	res, err := io.ReadAll(r)
	if err != nil {
		panic(err)
	}
	return res
}
