package tistory

import (
	"github.com/fineroot1253/notion2all/internal/common"
	"github.com/fineroot1253/notion2all/internal/tistory/model"
	"github.com/fineroot1253/tistoryAPI"
	tistoryAPImodel "github.com/fineroot1253/tistoryAPI/model"
)

type service struct {
	apiService tistoryAPI.Service
	config     model.TistoryConfiguration
}

type Service interface {
	GetBlogName() string
	SendPost(data tistoryAPImodel.PostData) (tistoryAPImodel.PostWriteResult, error)
	GetCategoryList() ([]tistoryAPImodel.CategoryData, error)
}

func NewService(config model.TistoryConfiguration, apiService tistoryAPI.Service) Service {
	return service{config: config, apiService: apiService}
}

func (s service) GetBlogName() string {
	return s.config.Tistory.BlogName
}

func (s service) SendPost(data tistoryAPImodel.PostData) (tistoryAPImodel.PostWriteResult, error) {
	result, err := s.apiService.WritePost(data)
	if err != nil {
		return tistoryAPImodel.PostWriteResult{}, &common.CommonError{
			Func: "SendPost",
			Data: common.TistoryAPIRequest_Error.String(),
			Err:  err,
		}
	}

	return result, nil
}

func (s service) GetCategoryList() ([]tistoryAPImodel.CategoryData, error) {

	list, err := s.apiService.GetCategoryList(s.GetBlogName())
	if err != nil {
		return nil, &common.CommonError{
			Func: "GetCategoryList",
			Data: common.TistoryAPIRequest_Error.String(),
			Err:  err,
		}
	}
	return list.Item.Categories, nil

}
