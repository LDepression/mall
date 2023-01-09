package times

import (
	"errors"
	"time"
)

const (
	LayoutDate     = "2006-01-02"
	LayoutTime     = "15:04:05"
	LayoutDateTime = "2006-01-02 15:04:05"
)

var cstZone *time.Location

func init() {
	//给定一个地点名,并且再给一个时间偏移量,UTC是以0时区开始算的
	cstZone = time.FixedZone("UTC", 8*3600) //北京时间
}

func Location() *time.Location {
	return cstZone
}

// ParseTimeToStr 格式化时间
func ParseTimeToStr(tm time.Time) string {
	return parseTimeToStr(tm.In(cstZone), LayoutTime)
}

// ParseDataToStr 格式化日期
func ParseDataToStr(tm time.Time) string {
	return parseTimeToStr(tm.In(cstZone), LayoutDate)
}

// ParseDataTimeToStr 格式化日期时间
func ParseDataTimeToStr(tm time.Time) string {
	return parseTimeToStr(tm.In(cstZone), LayoutDateTime)
}

// GetNowDateTimeStr 获取当前时间日期的字符串
func GetNowDateTimeStr() string {
	return parseTimeToStr(time.Now().In(cstZone), LayoutDateTime)
}

// GetNowDateStr 获取当前日期的字符串
func GetNowDateStr() string {
	return parseTimeToStr(time.Now().In(cstZone), LayoutDate)
}

// GetNowTimeStr 获取当前时间的字符串
func GetNowTimeStr() string {
	return parseTimeToStr(time.Now().In(cstZone), LayoutTime)
}

// ParseDateTime 解析日期时间
func ParseDateTime(dateTimeStr string) (time.Time, error) {
	return parseStrToTime(dateTimeStr, LayoutDateTime)
}

// ParseDate 解析日期
func ParseDate(dateStr string) (time.Time, error) {
	return parseStrToTime(dateStr, LayoutDate)
}

// ParseTime ...解析时间
func ParseTime(timeStr string) (time.Time, error) {
	return parseStrToTime(timeStr, LayoutTime)
}

func parseTimeToStr(t time.Time, layout string) string {
	return t.Format(layout)
}

func parseStrToTime(str, layout string) (time.Time, error) {
	if str == "" {
		return time.Now(), errors.New("args It can't be empty")
	}
	return time.ParseInLocation(layout, str, cstZone)
}

func GetNowTime() time.Time {
	return time.Now().In(cstZone)
}

// FuncTiming 程序运行时间
func FuncTiming(fn func()) time.Duration {
	startT := GetNowTime()
	fn()
	return GetNowTime().Sub(startT)
}

// 1970-01-01 08:00:00 +0800 CST
var zeroTime = time.Unix(0, 0)

// IsZero reports whether t represents the zero time instant
func IsZero(t time.Time) bool {
	return t.IsZero() || zeroTime.Equal(t)
}
