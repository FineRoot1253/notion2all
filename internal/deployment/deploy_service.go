package deployment

import (
	"encoding/json"
	"github.com/fineroot1253/notion2all/internal/common/logger"
	"github.com/fineroot1253/notion2all/internal/deployment/model"
	"github.com/fineroot1253/notion2all/internal/notion"
	"github.com/fineroot1253/tistoryAPI"
	tistoryModel "github.com/fineroot1253/tistoryAPI/model"
)

type Service interface {
	/*	Deploy NotionService를 토대로 각 플랫폼에 노션 포스트를 배포한다.
	 */
	Deploy(notionTablePath string) error
}

// service deployment 구조체
// logger를 통해 각 레이어별 실행 과정을 로깅한다.
type service struct {
	logger         logger.LogTemplate
	notionService  notion.Service
	tistoryService tistoryAPI.Service
	//... TODO 추후 서비스들을 추가해 나갈 예정
}

func NewService(logger logger.LogTemplate, notionService notion.Service, tistoryService tistoryAPI.Service) Service {
	return service{logger: logger, notionService: notionService, tistoryService: tistoryService}
}

func (s service) Deploy(notionTablePath string) error {
	// output record setting
	var postData []model.DeployResult

	// execution
	list, err := s.notionService.GetDeployPostList(notionTablePath)
	if err != nil {
		return err
	}

	// record initialize
	postData = make([]model.DeployResult, len(list))

	// parsing
	for i, data := range list {
		parsedPostData := tistoryModel.PostData{}
		if err := json.Unmarshal(data.PostData, &parsedPostData); err != nil {
			return err
		}
		post, _ := s.tistoryService.WritePost(parsedPostData)
		if post.Status != "200" {

		}
		postData[i].DeployList = append(postData[i].DeployList, post)

	}

	return nil
}
