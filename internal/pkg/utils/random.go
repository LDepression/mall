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

var (
	pictures = []string{"http://humraja.xyz/picture/wallhaven-13pv13.jpg", "http://humraja.xyz/picture/wallhaven-6ozkzl.jpg", "http://humraja.xyz/picture/wallhaven-dp3ogo.jpg", "http://humraja.xyz/picture/wallhaven-q2ep7l.jpg"}
)

func RandomAvatar() string {
	return pictures[rand.Intn(len(pictures))]
}

func RandomEmail() string {
	return RandomString(10) + "@" + RandomString(rand.Intn(2)+2) + ".com"
}

var areas = []string{"大陆",
	"美国",
	"韩国",
	"日本",
	"中国香港",
	"中国台湾",
	"泰国",
	"印度",
	"法国",
	"英国",
	"俄罗斯",
	"意大利",
	"西班牙",
	"德国",
	"波兰",
	"澳大利亚",
	"伊朗",
	"其他"}

func RandomArea() string {
	return areas[rand.Intn(len(areas))]
}

var tags = []string{"爱情", "喜剧", "动画", "剧情", "恐怖", "惊悚", "科幻", "动作", "悬疑", "犯罪", "冒险", "战争", "奇幻", "运动", "家庭", "古装", "武侠", "西部", "历史", "传记", "歌舞", "黑色电影", "短片", "纪录片", "戏曲", "音乐", "灾难", "青春", "儿童", "其他"}

func RandomTag() string {
	return tags[rand.Intn(len(tags))]
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
