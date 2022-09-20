package logger

import (
	"github.com/google/uuid"
	"log"
	"time"
)

type CommonLogRunner struct {
}

func (l CommonLogRunner) PreRun(message string) LogStatus {
	log.Print(message + " start... ")
	return NewLogStatus(uuid.New().String(), time.Now(), message)
}

func (l CommonLogRunner) PostRun(status LogStatus) {
	status.done()
	log.Println("done [time cost = ", status.count(), "ms]")
}
