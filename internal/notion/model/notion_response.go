package model

type NotionEnqueueTaskResponse struct {
	TaskId     string `json:"taskId"`
	ErrorId    string `json:"errorId"`
	Name       string `json:"name"`
	Message    string `json:"message"`
	ClientData struct {
		Type string `json:"type"`
	} `json:"clientData"`
}

type NotionGetTasksResponse struct {
	Results []NotionGetTaskResultItem `json:"results"`
}

type NotionGetTaskResultItem struct {
	Id        string `json:"id"`
	EventName string `json:"eventName"`
	Request   struct {
		Block struct {
			Id string `json:"id"`
		} `json:"block"`
		Recursive     bool `json:"recursive"`
		ExportOptions struct {
			ExportType string `json:"exportType"`
			TimeZone   string `json:"timeZone"`
		} `json:"exportOptions"`
	} `json:"request"`
	Actor struct {
		Table string `json:"table"`
		Id    string `json:"id"`
	} `json:"actor"`
	State       string `json:"state"`
	RootRequest struct {
		EventName string `json:"eventName"`
		RequestId string `json:"requestId"`
	} `json:"rootRequest"`
	Headers struct {
		Ip string `json:"ip"`
	} `json:"headers"`
	Status struct {
		Type          string `json:"type"`
		PagesExported int    `json:"pagesExported"`
		ExportURL     string `json:"exportURL"`
	} `json:"status"`
	Error string `json:"error"`
}

func (res NotionEnqueueTaskResponse) HasErrors() bool {
	if res.TaskId != "" {
		return false
	}
	return true
}

func (res NotionGetTasksResponse) HasErrors() bool {
	for _, result := range res.Results {
		if result.Error != "" {
			return true
		}
	}
	return false
}

func (res NotionGetTasksResponse) HasIncompleteTask() bool {
	for _, result := range res.Results {
		if result.Status.Type != "complete" {
			return true
		}
	}
	return false
}
