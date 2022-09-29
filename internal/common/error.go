package common

type ErrorInt interface {
	CurrentError() ErrorStateType
}

func New(state ErrorStateType) ErrorInt {
	return &errorIntModel{state}
}

type errorIntModel struct {
	errorState ErrorStateType
}

func (e *errorIntModel) CurrentError() ErrorStateType {
	return e.errorState
}

type CommonError struct {
	Func string // the failing functions
	Data string // the input
	Err  error  // the reason the conversion failed (e.g. ErrRange, ErrSyntax, etc.)
}

func (e *CommonError) Error() string {
	return e.Func + ": " + "parsing " + e.Data + ": " + e.Err.Error()
}

func (e *CommonError) Unwrap() error { return e.Err }

type ErrorStateType int64

const (
	Unexpected_Error ErrorStateType = iota
	CanNotFoundFile_Error
	Marshal_Error
	NotionAPIRequest_Error
	TistoryAPIRequest_Error
	Arguments_Error
)

func (err ErrorStateType) String() string {
	return []string{"Unexpected_Error", "CanNotFoundFile_Error", "Marshal_Error", "NotionAPIRequest_Error", "TistoryAPIRequest_Error", "Arguments_Error"}[err]
}
