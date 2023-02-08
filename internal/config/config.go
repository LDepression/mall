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
	OSS      OSS      `mapstructure:"OSS"`
	Auto     Auto     `mapstructure:"Auto"`
	EsInfo   EsInfo   `mapstructure:"EsInfo"`
	AliPay   AliPay   `mapstructure:"AliPay"`
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
	Addr      string        `json:"addr" mapstructure:"addr"`
	Password  string        `json:"password" mapstructure:"password"`
	PoolSize  int           `json:"poolSize" mapstructure:"poolSize"`
	CacheTime time.Duration `json:"cacheTime" mapstructure:"CacheTime"`
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

type OSS struct {
	BucketName      string `mapstructure:"BucketName"`
	ObjectName      string `mapstructure:"ObjectName"`
	LocalFileName   string `mapstructure:"LocalFileName"`
	Endpoint        string `mapstructure:"Endpoint"`
	AccessKeyId     string `mapstructure:"AccessKeyId"`
	AccessKeySecret string `mapstructure:"AccessKeySecret"`
	BasePath        string `mapstructure:"BasePath"`
	BucketUrl       string `mapstructure:"BucketUrl"`
}

type Auto struct {
	SendEmailTime time.Duration `mapstructure:"SendEmailTime"`
	CodeValidTime time.Duration `mapstructure:"CodeValidTime"`
}

type EsInfo struct {
	Host string `mapstructure:"Host"`
	Port int    `mapstructure:"Port"`
}

type AliPay struct {
	AppID        string `mapstructure:"AppID"`
	PrivateKey   string `mapstructure:"PrivateKey"`
	AliPublicKey string `mapstrvcucture:"AliPublicKey"`
	NotifyURL    string `mapstructure:"NotifyURL"`
	ReturnURL    string `mapstructure:"ReturnURL"`
	ProductCode  string `mapstructure:"ProductCode"`
	IsProduction bool   `mapstructure:"IsProduction"`
}
