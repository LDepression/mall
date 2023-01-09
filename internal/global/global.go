package global

import (
	ut "github.com/go-playground/universal-translator"
	"mall/internal/config"
	"mall/internal/pkg/goroutine/work"
	"mall/internal/pkg/logger"
	"mall/internal/pkg/token"
)

var (
	Trans   ut.Translator
	Setting config.All
	Logger  *logger.Log
	Maker   token.Maker
	Worker  = new(work.Worker)
)
