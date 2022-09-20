package model

import (
	notionModel "github.com/fineroot1253/notion2all/internal/notion/model"
	tistoryModel "github.com/fineroot1253/tistoryAPI/model"
)

type DeployResult struct {
	SuccessTaskCount int
	FailureTaskCount int
	TaskList         []notionModel.ExportTask
	DeployList       []tistoryModel.PostWriteResult
}
