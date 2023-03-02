package utils

import (
	"math/rand"
	"strings"
	"time"
)

const alphabetic = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func init() {
	rand.Seed(time.Now().UnixNano())
}

// RandomInt 返回min到max之间的一个随机数
func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

// RandomFloat 返回min到max之间的一个随机小数
func RandomFloat(min, max float64) float64 {
	return min + rand.Float64()*(max-min)
}

// RandomString 生成一个长度为n的随机字符串
func RandomString(n int) string {
	var sb strings.Builder
	k := len(alphabetic)
	for i := 0; i < n; i++ {
		c := alphabetic[rand.Intn(k)]
		sb.WriteByte(c)
	}
	return sb.String()
}

// RandomOwner RandomString(6)
func RandomOwner() string {
	return RandomString(6)
}

// RandomStringSlice 指定最大切片长度和切片中元素长度的随机字符串切片
func RandomStringSlice(maxLength int, eleMaxLength int) []string {
	length := int(RandomInt(1, int64(maxLength)))
	ret := make([]string, length)
	for i := range ret {
		ret[i] = RandomString(eleMaxLength)
	}
	return ret
}

func RandomEmail() string {
	return RandomString(10) + "@" + RandomString(rand.Intn(2)+2) + ".com"
}

/*
2022 2021 2020 2019 2018 2017 2016 2015 2014 2013 2012 2011 2000-2010 90年代 80年代 70年代 更早
*/

var (
	periods []time.Time
)

func init() {
	strs := []string{"2022", "2021", "2020", "2019", "2018", "2017", "2016", "2015", "2014", "2013", "2012", "2011", "2000"}
	layout := "2006"
	for i := range strs {
		t, err := time.Parse(layout, strs[i])
		if err != nil {
			panic(err)
		}
		periods = append(periods, t)
	}
}

// RandomPeriod 随机时间点
func RandomPeriod() time.Time {
	return periods[rand.Intn(len(periods))]
}
