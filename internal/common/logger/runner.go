package logger

type LogRunner interface {
	PreRun(message string) LogStatus
	PostRun(status LogStatus)
}

type Run[T any] func() T
type TY any
