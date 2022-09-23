package model

type NotionRequest interface {
	CheckRequest | ExportRequest
}

type CheckRequest struct {
	TaskIds []string `json:"taskIds"`
}

type ExportRequest struct {
	Task Task `json:"task"`
}

type Task struct {
	EventName string      `json:"eventName"`
	Request   TaskRequest `json:"request"`
}

type TaskRequest struct {
	Block         BlockStruct  `json:"block"`
	ExportOptions ExportOption `json:"exportOptions"`
	Recursive     bool         `json:"recursive"`
}

type BlockStruct struct {
	Id string `json:"id"`
	//SpaceId string `json:"spaceId"`
}

type ExportOption struct {
	ExportType string `json:"exportType"`
	TimeZone   string `json:"timeZone"`
	Locale     string `json:"locale"`
}

type ExportTask struct {
	PageId               string
	TaskId               string
	TaskResult           NotionEnqueueTaskResponse
	TaskStatusResultItem NotionGetTaskResultItem
}

func (et ExportTask) GetExportUrl() string {
	return et.TaskStatusResultItem.Status.ExportURL
}
