package setting

import (
	"mall/internal/global"
	"mall/internal/pkg/token"
)

type maker struct {
}

func (maker) Init() {
	var err error
	global.Maker, err = token.NewPasetoMaker([]byte(global.Setting.Token.Key))
	if err != nil {
		panic(err)
	}
}
