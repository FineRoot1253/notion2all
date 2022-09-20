package logger

import "time"

type LogStatus struct {
	logId     string
	startTime time.Time
	endTime   time.Time
	message   string
}

func NewLogStatus(logId string, startTime time.Time, message string) LogStatus {
	return LogStatus{logId: logId, startTime: startTime, message: message}
}

func (logStatus *LogStatus) done() {
	logStatus.endTime = time.Now()
}

func (logStatus *LogStatus) count() int {

	return int(logStatus.endTime.UnixMilli() - logStatus.startTime.UnixMilli())
}
