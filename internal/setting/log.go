package setting

import (
	"go.uber.org/zap"
	"mall/internal/global"
	"mall/internal/pkg/logger"
)

type log struct{}

func (log) Init() {
	//初始化
	//logger, _ := zap.NewDevelopment()

	/*
				当我们使用生产者模式的时候,我们可以使用Debug
			因为product时候,日志级别是info,要高于debug,所以debug打印不出来
		S和L函数提供了安全访问全局logger的功能(),因为里面加了锁
	*/
	//zap.ReplaceGlobals(logger)
	global.Logger = logger.NewLogger(&logger.InitStruct{
		LogSavePath:   global.Setting.Log.LogSavePath,
		LogFileExt:    global.Setting.Log.LogFileExt,
		MaxSize:       global.Setting.Log.MaxSize,
		MaxBackups:    global.Setting.Log.MaxBackups,
		MaxAge:        global.Setting.Log.MaxAge,
		Compress:      global.Setting.Log.Compress,
		LowLevelFile:  global.Setting.Log.LowLevelFile,
		HighLevelFile: global.Setting.Log.HighLevelFile,
	}, global.Setting.Log.Level)
	zap.ReplaceGlobals(global.Logger.Logger)
}
