package notion

import (
	"archive/zip"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/fineroot1253/notion2all/cmd/config"
	"github.com/fineroot1253/notion2all/internal/common"
	"github.com/fineroot1253/notion2all/internal/common/logger"
	"github.com/fineroot1253/notion2all/internal/common/utils"
	"github.com/fineroot1253/notion2all/internal/notion/model"
	"github.com/kjk/notionapi"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type Service interface {
	/*
		GetDeployPostList 배포 포스트 리스트 받기
		@Return []model.NotionPostItem, error // 다운 받은 html 파일 path 리스트, error
	*/
	GetDeployPostList() ([]model.NotionPostItem, error)

	/*
		GetPostHtmlDataList 배포할 포스트의 html byte 데이터 리스트
		@Param []model.ExportTask // 배포할 포스트데이터 리스트
		@Return []model.NotionPostData, []model.ExportTask, error // 다운 받은 html 파일 data 리스트, 실제 수행된 html Export 요청 리스트, error
	*/
	GetPostHtmlDataList(postItemList []model.NotionPostItem) ([]model.ExportTask, error)

	/*	enqueueExportTask Html파일 export 링크 생성 요청
		@Param string // export할 페이지ID
		@Return string, error // 생성요청 처리 진행중인 TaskID, error
	*/
	enqueueExportTask(pageId string) (model.NotionEnqueueTaskResponse, error)

	/*	getTaskStatusList 생성 작업 현황 리스트 조회
		@Param []string // 생성 대기열 리스트
		@Return string, error //  TaskID, error
	*/
	getTaskStatusList(taskIdList []string) (model.NotionGetTasksResponse, error)

	/*	downloadHtmlFiles 생성된 html 다운로드 url 이용해 파일 다운로드
		@Param []string // html 다운로드 url 리스트
		@Return []string, error //  다운받은 파일 리스트, error
	*/
	downloadHtmlFiles(taskList []model.ExportTask) ([]model.ExportTask, error)
}

type service struct {
	context      context.Context
	notionClient notionapi.Client
	notionConfig model.NotionConfiguration
	logger       logger.LogTemplate
}

func NewService(ctx context.Context, notionConfig model.NotionConfiguration, logger logger.LogTemplate) (Service, error) {
	if ctx == nil || logger.IsEmpty() || notionConfig.IsEmpty() {
		return nil, &common.CommonError{
			Func: "NewService",
			Data: common.Arguments_Error.String(),
			Err:  errors.New("파라미터 개수가 모자릅니다"),
		}
	}
	return service{context: ctx, logger: logger, notionConfig: notionConfig, notionClient: notionapi.Client{AuthToken: notionConfig.Notion.Token}}, nil
}

func (s service) GetDeployPostList() ([]model.NotionPostItem, error) {
	var spaceId string
	var collectionId string
	var collectionViewId string
	var collectionKey model.NotionSchema
	var resultData []model.NotionPostItem
	req := new(notionapi.QueryCollectionRequest)

	blockRecords, err := s.notionClient.DownloadPage(utils.ConvertId(s.notionConfig.GetTableBlockId()))
	if err != nil {
		return nil, err
	}
	for _, record := range blockRecords.SpaceRecords {
		if spaceId == "" {
			spaceId = record.Space.ID
		}
	}
	for _, record := range blockRecords.CollectionViewRecords {
		if collectionViewId == "" {
			collectionViewId = record.CollectionView.ID
		}
	}
	for _, record := range blockRecords.CollectionRecords {
		if collectionId == "" {
			collectionId = record.Collection.ID
		}

		for s2, schema := range record.Collection.Schema {
			collectionKey.SetProperty(schema.Name, s2)
		}
	}

	req.CollectionView.SpaceID = spaceId
	req.CollectionView.ID = collectionViewId
	req.Collection.SpaceID = spaceId
	req.Collection.ID = collectionId

	collection, err := s.notionClient.QueryCollection(*req, nil)
	if err != nil {
		return nil, err
	}

	for _, record := range collection.RecordMap.Blocks {
		var temp model.NotionPostItem
		for s3, i := range record.Block.Properties {
			temp.SetProperty(collectionKey, s3, i)
		}
		resultData = append(resultData, temp)
	}

	return resultData, nil
}

func (s service) GetPostHtmlDataList(postItemList []model.NotionPostItem) ([]model.ExportTask, error) {

	var taskIdList []string
	var exportTaskList []model.ExportTask

	for _, item := range postItemList {
		pageId := utils.ConvertPageUrlToPageId(item.Content)
		task, _ := s.enqueueExportTask(pageId)
		taskIdList = append(taskIdList, task.TaskId)
		exportTaskList = append(exportTaskList, model.ExportTask{
			PageId:     pageId,
			TaskId:     task.TaskId,
			TaskResult: task,
		})
	}

	for true {
		res, err := s.getTaskStatusList(utils.ConvertIdList(taskIdList))
		if err != nil {
			return nil, err
		}
		if res.HasIncompleteTask() {
			time.Sleep(750 * time.Millisecond)
		} else {
			exportTaskList = res.GetExportDataList(exportTaskList)
			break
		}
	}

	exportTaskList, err := s.downloadHtmlFiles(exportTaskList)
	if err != nil {
		return nil, err
	}

	for _, task := range exportTaskList {
		destDirPath := filepath.Join(config.NOTION_HTML_DATA_PATH, task.PageId)
		htmlFilePath, err := unzip(task.FilePath, destDirPath)
		if err != nil {
			return nil, err
		}
		task.FilePath = htmlFilePath
	}

	return exportTaskList, nil

}

func (s service) enqueueExportTask(pageId string) (model.NotionEnqueueTaskResponse, error) {
	var result model.NotionEnqueueTaskResponse
	request := model.ExportRequest{
		Task: model.Task{
			EventName: "exportBlock",
			Request: model.TaskRequest{
				Block: model.BlockStruct{
					Id: pageId,
				},
				ExportOptions: model.ExportOption{
					ExportType: "exportType",
					TimeZone:   "Asia/Seoul",
					Locale:     "en",
				},
				Recursive: false,
			},
		},
	}
	byteData, err := sendPostRequest(s.context, request, config.NOTION_ENQUEUE_TASK_PATH, s.notionConfig.Notion.Token)
	if err != nil {
		return result, err
	}

	if err := json.Unmarshal(byteData, &result); err != nil {
		return result, &common.CommonError{
			Func: "enqueueExportTask",
			Data: common.Marshal_Error.String(),
			Err:  err,
		}
	}

	return result, nil
}

func (s service) getTaskStatusList(taskIdList []string) (model.NotionGetTasksResponse, error) {

	var result model.NotionGetTasksResponse
	request := model.CheckRequest{
		TaskIds: taskIdList,
	}
	byteData, err := sendPostRequest(s.context, request, config.NOTION_GET_TASKS_PATH, s.notionConfig.Notion.Token)
	if err != nil {
		return result, err
	}

	if err := json.Unmarshal(byteData, result); err != nil {
		return result, &common.CommonError{
			Func: "getTaskStatusList",
			Data: common.Marshal_Error.String(),
			Err:  err,
		}
	}

	return result, nil
}

/* sendRequest Export용 Html 파일 다운로드 링크 생성 요청
Task를 하나 등록한다.
	=> 비동기
비동기로써 테이블 로우별 Task등록(다운로드 링크 생성)시 빠짐없이 모든 Task가 완료 할때까지 상당시간 대기 필요
	=> sleep 필요
실행 순서
1. Export된 HTML 파일 다운 링크 생성
	=> 비동기라서 생성까지 시간이 걸린다. 여기에 Sleep으로 대기를 하던지 고루틴을 돌려 백그라운드로 돌리던지 이럴 필요가 존재함
	=> 0.0.1v 기준 단일 고루틴으로 진행 [Sleep으로 대기]
2. Html 로컬 저장
3. Html
*/
func sendPostRequest[NR model.NotionRequest](ctx context.Context, request NR, path string, token string) ([]byte, error) {

	client := http.Client{}
	body := &bytes.Buffer{}

	reqByte, err := json.Marshal(request)
	if err != nil {
		return nil, &common.CommonError{
			Func: "sendPostRequest",
			Data: common.Marshal_Error.String(),
			Err:  err,
		}
	}
	if _, err := body.Write(reqByte); err != nil {
		return nil, &common.CommonError{
			Func: "sendPostRequest",
			Data: common.Unexpected_Error.String(),
			Err:  err,
		}
	}

	reqWithCtx, err := http.NewRequestWithContext(ctx, http.MethodPost, config.NOTION_V3_BASE_URL+path, body)
	reqWithCtx.Header.Set("Cookie", "token_v2="+token)
	response, err := client.Do(reqWithCtx)
	if err != nil {
		return nil, &common.CommonError{
			Func: "sendPostRequest",
			Data: common.NotionAPIRequest_Error.String(),
			Err:  err,
		}
	}

	defer response.Body.Close()
	all, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, &common.CommonError{
			Func: "sendPostRequest",
			Data: common.Unexpected_Error.String(),
			Err:  err,
		}
	}

	return all, nil
}

func sendGetRequest(ctx context.Context, url string) ([]byte, error) {
	client := http.Client{}

	reqWithCtx, err := http.NewRequestWithContext(ctx, http.MethodPost, url, nil)
	response, err := client.Do(reqWithCtx)
	if err != nil {
		return nil, &common.CommonError{
			Func: "sendPostRequest",
			Data: common.NotionAPIRequest_Error.String(),
			Err:  err,
		}
	}

	defer response.Body.Close()
	all, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, &common.CommonError{
			Func: "sendPostRequest",
			Data: common.Unexpected_Error.String(),
			Err:  err,
		}
	}

	return all, nil
}

func (s service) downloadHtmlFiles(taskList []model.ExportTask) ([]model.ExportTask, error) {

	resultList := taskList

	if err := os.MkdirAll(config.NOTION_HTML_DATA_PATH, 0755); err != nil {
		return nil, &common.CommonError{
			Func: "downloadHtmlFiles",
			Data: common.Unexpected_Error.String(),
			Err:  err,
		}
	}

	for _, data := range resultList {
		byteData, _ := sendGetRequest(s.context, data.GetExportUrl())
		filePath := filepath.Join(config.NOTION_HTML_DATA_PATH, data.PageId, ".zip")
		file, err := os.Create(filePath)
		if err != nil {
			return nil, &common.CommonError{
				Func: "downloadHtmlFiles",
				Data: common.Unexpected_Error.String(),
				Err:  err,
			}
		}

		data.FilePath = filePath

		defer func() {
			file.Close()
		}()

		_, err = file.Write(byteData)
		if err != nil {
			return nil, &common.CommonError{
				Func: "downloadHtmlFiles",
				Data: common.Unexpected_Error.String(),
				Err:  err,
			}
		}
	}

	return resultList, nil

}

/*	unzip export 파일 결과인 zip파일 해제
	@Param	string, string // zip 파일 위치, 파일 해제 위치
	@Return	string, error // html 파일 위치, error
*/
func unzip(src string, destDirPath string) (string, error) {

	var htmlFilePath string

	r, err := zip.OpenReader(src)
	if err != nil {
		return "", &common.CommonError{
			Func: "unzip",
			Data: common.Unexpected_Error.String(),
			Err:  err,
		}
	}
	defer func() {
		if err := r.Close(); err != nil {
			panic(&common.CommonError{
				Func: "extractAndWriteFile",
				Data: common.Unexpected_Error.String(),
				Err:  err,
			})
		}
	}()
	if err := os.MkdirAll(destDirPath, 0755); err != nil {
		return "", &common.CommonError{
			Func: "unzip",
			Data: common.Unexpected_Error.String(),
			Err:  err,
		}
	}

	for _, f := range r.File {

		err := extractAndWriteFile(f, destDirPath)
		if err != nil {
			return "", err
		}
		if strings.HasSuffix(f.Name, ".html") {
			htmlFilePath = filepath.Join(destDirPath, f.Name)
		}
	}
	return htmlFilePath, nil
}

func extractAndWriteFile(f *zip.File, destDirPath string) error {

	rc, err := f.Open()
	if err != nil {
		return &common.CommonError{
			Func: "extractAndWriteFile",
			Data: common.Unexpected_Error.String(),
			Err:  err,
		}
	}
	defer func() {
		if err := rc.Close(); err != nil {
			panic(&common.CommonError{
				Func: "extractAndWriteFile",
				Data: common.Unexpected_Error.String(),
				Err:  err,
			})
		}
	}()

	path := filepath.Join(destDirPath, f.Name)

	if !strings.HasPrefix(path, filepath.Clean(destDirPath)+string(os.PathSeparator)) {
		return fmt.Errorf("illegal file path: %s", path)
	}

	if f.FileInfo().IsDir() {
		if err := os.MkdirAll(path, f.Mode()); err != nil {
			return &common.CommonError{
				Func: "extractAndWriteFile",
				Data: common.Unexpected_Error.String(),
				Err:  err,
			}
		}

	} else {

		if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			return &common.CommonError{
				Func: "extractAndWriteFile",
				Data: common.Unexpected_Error.String(),
				Err:  err,
			}
		}
		f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return &common.CommonError{
				Func: "extractAndWriteFile",
				Data: common.Unexpected_Error.String(),
				Err:  err,
			}
		}
		defer func() {
			if err := f.Close(); err != nil {
				panic(&common.CommonError{
					Func: "extractAndWriteFile",
					Data: common.Unexpected_Error.String(),
					Err:  err,
				})
			}
		}()

		_, err = io.Copy(f, rc)
		if err != nil {
			return &common.CommonError{
				Func: "extractAndWriteFile",
				Data: common.Unexpected_Error.String(),
				Err:  err,
			}
		}

	}
	return nil
}
