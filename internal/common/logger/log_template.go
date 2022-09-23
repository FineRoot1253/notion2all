package logger

import (
	"errors"
	"fmt"
	"github.com/fineroot1253/notion2all/internal/common"
	"log"
)

type LogTemplate struct {
	runner LogRunner
}

func NewLogTemplate(runner LogRunner) LogTemplate {
	return LogTemplate{runner: runner}
}

func (lt LogTemplate) IsEmpty() bool {
	if lt.runner == nil {
		return true
	}
	return false
}

func (l LogTemplate) Execute(msg string, run Run[any]) TY {
	// 예외처리
	defer func() {
		if r := recover(); r != nil {
			log.Panicln(&common.CommonError{
				Func: "Execute",
				Data: fmt.Sprint(r),
				Err:  errors.New(common.Unexpected_Error.String()),
			})
		}
	}()
	logStatus := l.runner.PreRun(msg)
	run()
	l.runner.PostRun(logStatus)
	return nil
}
