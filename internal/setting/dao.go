package setting

import (
	"mall/internal/dao"
	"mall/internal/dao/db"
	"mall/internal/dao/redis"
)

type mdao struct {
}

func (mdao) Init() {
	dao.Group.DB = db.Init()
	dao.Group.Redis = redis.Init()
}
