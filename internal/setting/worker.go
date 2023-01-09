package setting

import (
	"mall/internal/global"
	"mall/internal/pkg/goroutine/work"
)

type worker struct {
}

func (worker) Init() {
	global.Worker = work.Init(work.Config{
		TaskChanCapacity:   global.Setting.Worker.TaskChanCapacity,
		WorkerChanCapacity: global.Setting.Worker.WorkerChanCapacity,
		WorkerNum:          global.Setting.Worker.WorkerNum,
	})
}
