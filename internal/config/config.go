package config

import "time"

type All struct {
	Serve    Serve    `mapstructure:"serve"`
	Mysql    Mysql    `mapstructure:"mysql"`
	Redis    Redis    `mapstructure:"redis"`
	SMTPInfo SMTPInfo `mapstructure:"SMTPInfo"`
	Log      Log      `mapstructure:"Log"`
	Captcha  Captcha  `mapstructure:"Captcha"`
	Worker   Worker   `mapstructure:"Worker"`
	Token    Token    `mapstructure:"Token"`
}
type Serve struct {
	Addr           string        `mapstructure:"addr" json:"addr"`
	ReadTimeout    time.Duration `mapstructure:"readTimeout" json:"readTimeout"`
	WriteTimeout   time.Duration `mapstructure:"writeTimeout" json:"writeTimeout"`
	DefaultTimeout time.Duration `mapstructure:"defaultTimeout" json:"defaultTimeout"`
}

type Mysql struct {
	User     string `json:"user" mapstructure:"user"`
	Host     string `json:"host" mapstructure:"host"`
	Port     int    `json:"port" mapstructure:"port"`
	Password string `json:"password" mapstructure:"password"`
	DbName   string `json:"dbName" mapstructure:"dbname"`
}

type Redis struct {
	Addr     string `json:"addr" mapstructure:"addr"`
	Password string `json:"password" mapstructure:"password"`
	PoolSize int    `json:"poolSize" mapstructure:"poolSize"`
}

type SMTPInfo struct {
	Host     string   `json:"host" mapstructure:"host"`
	Port     int      `json:"port" mapstructure:"port"`
	IsSSL    bool     `json:"isSSL" mapstructure:"isSSL"`
	UserName string   `json:"userName" mapstructure:"userName"`
	Password string   `json:"password" mapstructure:"password"`
	From     string   `json:"from" mapstructure:"from"`
	To       []string `json:"to" mapstructure:"to"`
}

type Log struct {
	Level         string `json:"Level" mapstructure:"Level"`
	LogSavePath   string `json:"LogSavePath" mapstructure:"LogSavePath"`     // 保存路径
	LogFileExt    string `json:"LogFileExt" mapstructure:"LogFileExt"`       // 日志文件后缀
	MaxSize       int    `json:"MaxSize" mapstructure:"MaxSize"`             // 备份的大小(M)
	MaxBackups    int    `json:"MaxBackups" mapstructure:"MaxBackups"`       // 最大备份数
	MaxAge        int    `json:"MaxAge" mapstructure:"MaxAge"`               // 最大备份天数
	Compress      bool   `json:"Compress" mapstructure:"Compress"`           // 是否压缩过期日志
	LowLevelFile  string `json:"LowLevelFile" mapstructure:"LowLevelFile"`   // 低级别文件名
	HighLevelFile string `json:"HighLevelFile" mapstructure:"HighLevelFile"` // 高级别文件名
}

type Captcha struct {
	Width    int           `mapstructure:"width" `
	Height   int           `mapstructure:"height"`
	Length   int           `mapstructure:"length"`
	MaxSkew  float64       `mapstructure:"maxSkew"`
	DotCount int           `mapstructure:"dotCount"`
	Expired  time.Duration `mapstructure:"expired"`
}

type Worker struct {
	TaskChanCapacity   int `mapstructure:"TaskChanCapacity"`
	WorkerChanCapacity int `mapstructure:"WorkerChanCapacity"`
	WorkerNum          int `mapstructure:"WorkerNum"`
}

type Token struct {
	Key                string        `mapstructure:"Key"`
	AccessTokenExpire  time.Duration `mapstructure:"AccessTokenExpire"`
	RefreshTokenExpire time.Duration `mapstructure:"RefreshTokenExpire"`
	AuthType           string        `mapstructure:"AuthType"`
	AuthKey            string        `mapstructure:"AuthKey"`
}
