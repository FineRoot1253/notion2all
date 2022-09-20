package notion

import (
	"archive/zip"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/fineroot1253/notion2all/cmd/config"
	"github.com/fineroot1253/notion2all/internal/common/logger"
	"github.com/fineroot1253/notion2all/internal/notion/model"
	"github.com/kjk/notionapi"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

type Service interface {
	/*	GetExportPageListToHtmlFiles 각각 페이지를 html로 export후 Html파일 다운로드
		@Param string // export할 페이지ID
		@Return string, error // 다운 받은 html 파일 path 리스트, error
	*/
	GetExportPageListToHtmlFiles(pageIdList []string) ([]model.NotionPostData, error)

	/*
		GetDeployPostList 최종 배포용 데이터 가져오기
		@Param string // export할 페이지ID
		@Return string, error // 다운 받은 html 파일 path 리스트, error
	*/
	GetDeployPostList(tablePagePath string) ([]model.ExportTask, error)

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
	downloadHtmlFiles(downloadLinkList []string) ([]string, error)
}

type service struct {
	context      context.Context
	tokenV2      string
	notionClient notionapi.Client
	logger       logger.LogTemplate
}

func NewService(ctx context.Context, token string, logger logger.LogTemplate) Service {
	return service{context: ctx, tokenV2: token, logger: logger}
}

func (s service) GetExportPageListToHtmlFiles(pageIdList []string) ([]string, error) {

	// Input 레코드 , Output 레코드 초기화 로직
	pageIdMap := map[string]model.NotionEnqueueTaskResponse{}
	var taskIdList []string
	var exportTask []model.ExportTask

	convertedList := convertIdList(pageIdList)

	for idx, pageId := range convertedList {
		result, err := s.enqueueExportTask(pageId)
		if err != nil {
			return nil, err
		}
		exportTask[idx].PageId = pageId
		exportTask[idx].TaskId = result.TaskId
		exportTask[idx].TaskResult = result

		pageIdMap[pageId] = result

	}

	// html 파일 생성요청 상태 확인시 실패한 요청에 대해서는 확인이 불가함!
	// 		=> taskId가 없기 때문임
	// 이를 위해 에러 없이 성공한 요청 TaskId만 모은 taskIdList를 만들고
	// 이를 이용해 상태확인 및 동기 처리를 시도한다.
	for _, response := range pageIdMap {
		if !response.HasErrors() {
			taskIdList = append(taskIdList, response.TaskId)
		}
	}

	// 완성할때까지 0.5초 대기후 재확인하는 동기화 로직
	// 보통은 한번에 통과하지만 정말 큰 페이지[루트블럭]은 좀 오래걸리긴한다.
	for {
		status, err := s.getTaskStatusList(taskIdList)
		if err != nil {
			return nil, err
		}
		if status.HasIncompleteTask() {
			time.Sleep(time.Microsecond * 500)
		} else {
			break
		}
	}

	// TODO 파일 다운로드

	return []string{}, nil

}

func (s service) GetDeployPostList(tablePagePath string) ([]model.NotionPostData, error) {
	//s.notionClient.
	return nil, nil
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
	byteData, err := sendPostRequest(s.context, request, config.NOTION_ENQUEUE_TASK_PATH, s.tokenV2)
	if err != nil {
		return result, err
	}

	if err := json.Unmarshal(byteData, &result); err != nil {
		return result, err
	}

	return result, nil
}

func (s service) getTaskStatusList(taskIdList []string) (model.NotionGetTasksResponse, error) {

	var result model.NotionGetTasksResponse
	request := model.CheckRequest{
		TaskIds: taskIdList,
	}
	byteData, err := sendPostRequest(s.context, request, config.NOTION_GET_TASKS_PATH, s.tokenV2)
	if err != nil {
		return result, err
	}

	if err := json.Unmarshal(byteData, result); err != nil {
		return result, err
	}

	if result.HasErrors() {
		return result, err
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
		return nil, err
	}
	if _, err := body.Write(reqByte); err != nil {
		return nil, err
	}

	reqWithCtx, err := http.NewRequestWithContext(ctx, http.MethodPost, config.NOTION_V3_BASE_URL+path, body)
	reqWithCtx.Header.Set("Cookie", "token_v2="+token)
	response, err := client.Do(reqWithCtx)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()
	all, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	return all, nil
}

func sendGetRequest(ctx context.Context, url string) ([]byte, error) {
	client := http.Client{}

	reqWithCtx, err := http.NewRequestWithContext(ctx, http.MethodPost, url, nil)
	response, err := client.Do(reqWithCtx)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()
	all, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	return all, nil
}

func (serv service) downloadHtmlFiles(downloadLinkList []model.ExportTask) ([]string, error) {

	var downloadFilePathList []string
	tempDir, err := os.MkdirTemp(".", "notion_htmls")
	if err != nil {
		return nil, err
	}

	for _, data := range downloadLinkList {
		downloadPath, _ := sendGetRequest(serv.context, link)
		os.CreateTemp(tempDir)
	}
}

/*	convertIdList pageid 리스트를 uuid-v4 형식으로 변환
	@Param	[]string // pageId 리스트
	@Return	[]string // 변환된 pageId 리스트
*/
func convertIdList(idList []string) []string {
	for _, id := range idList {
		id = regexp.MustCompile(`(.{8})(.{4})(.{4})(.{4})(.+)`).ReplaceAllString(id, `$1-$2-$3-$4-$5`)
	}
	return idList
}

/*	unzip export 파일 결과인 zip파일 해제
	@Param	[]string // pageId 리스트
	@Return	[]string // 변환된 pageId 리스트
*/
func unzip(src string, destDirPath string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer func() {
		if err := r.Close(); err != nil {
			panic(err)
		}
	}()
	if err := os.MkdirAll(destDirPath, 0755); err != nil {
		return err
	}
	for _, f := range r.File {
		err := extractAndWriteFile(f, destDirPath)
		if err != nil {
			return err
		}
	}
	return nil
}

func extractAndWriteFile(f *zip.File, destDirPath string) error {
	rc, err := f.Open()
	if err != nil {
		return err
	}
	defer func() {
		if err := rc.Close(); err != nil {
			panic(err)
		}
	}()

	path := filepath.Join(destDirPath, f.Name)

	if !strings.HasPrefix(path, filepath.Clean(destDirPath)+string(os.PathSeparator)) {
		return fmt.Errorf("illegal file path: %s", path)
	}

	if f.FileInfo().IsDir() {
		if err := os.MkdirAll(path, f.Mode()); err != nil {
			return err
		}

	} else {
		if err := os.MkdirAll(filepath.Dir(path), f.Mode()); err != nil {
			return err
		}
		f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return err
		}
		defer func() {
			if err := f.Close(); err != nil {
				panic(err)
			}
		}()

		_, err = io.Copy(f, rc)
		if err != nil {
			return err
		}
	}
	return nil
}

/*	replaceImgSrc unzip 결과로 파일이 1개 이상일때(이미지파일)
	해당 img 태그 src를 base64로 치환
	@Param	[]string // pageId 리스트
	@Return	[]string // 변환된 pageId 리스트
*/
func replaceImgSrc() error {

	return nil
}

/*	replaceImgSrc 이미지 파일 base64 반환
	@Param	[]string // pageId 리스트
	@Return	[]string // 변환된 pageId 리스트
*/
func getImageBase64String() (string, error) {
	return "", nil
}
