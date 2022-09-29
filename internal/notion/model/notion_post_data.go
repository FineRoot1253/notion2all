package model

import (
	"fmt"
	"github.com/fineroot1253/notion2all/internal/common/utils"
	tistoryAPIModel "github.com/fineroot1253/tistoryAPI/model"
	"strings"
)

type NotionPostItem struct {
	tistoryAPIModel.PostData
}

func (npi *NotionPostItem) ToTistoryPostData(blogName string, content string, categoryId string) (tistoryAPIModel.PostData, error) {
	stamp, err := utils.ConvertStringToTimeStamp(npi.Published)
	if err != nil {
		return tistoryAPIModel.PostData{}, err
	}
	return tistoryAPIModel.PostData{
		BlogName:      blogName,
		Title:         npi.Title,
		Content:       content,
		Visibility:    npi.Visibility,
		Category:      categoryId,
		Published:     stamp,
		Tag:           npi.Tag,
		AcceptComment: npi.getCommentAvailability(),
		Password:      npi.Password,
	}, nil

}

func (npi *NotionPostItem) GetPageId() string {
	if strings.HasPrefix(npi.Content, "http") {
		return utils.ConvertPageUrlToPageId(npi.Content)
	}
	return npi.Content
}

func (npi *NotionPostItem) SetProperty(collectionKey NotionSchema, key string, value interface{}) {
	switch key {
	case collectionKey.Title:
		npi.Title = fmt.Sprintf("%v", value)
		break
	case collectionKey.ContentUrl:
		npi.Content = fmt.Sprintf("%v", value)
		break
	case collectionKey.Tags:
		switch value := value.(type) {
		case []string:
			npi.Tag = strings.Join(value, ",")
		}
		break
	case collectionKey.Category:
		npi.Category = fmt.Sprintf("%v", value)
		break
	case collectionKey.Visibility:
		npi.Visibility = fmt.Sprintf("%v", value)
		break
	case collectionKey.AcceptComment:
		npi.AcceptComment = fmt.Sprintf("%v", value)
		break
	case collectionKey.PublishDate:
		npi.Published = fmt.Sprintf("%v", value)
		break
	default:
		npi.Password = fmt.Sprintf("%v", value)
		break
	}
}

func (npi *NotionPostItem) getCommentAvailability() string {
	if npi.AcceptComment == "YES" {
		return "1"
	} else {
		return "0"
	}
}

type NotionSchema struct {
	Title         string
	ContentUrl    string
	Visibility    string
	Category      string
	PublishDate   string
	Tags          string
	AcceptComment string
	Password      string
}

func (ns *NotionSchema) SetProperty(key string, value string) {
	switch key {
	case "title":
		ns.Title = value
		break
	case "contentUrl":
		ns.ContentUrl = value
		break
	case "visibility":
		ns.Visibility = value
		break
	case "category":
		ns.Category = value
		break
	case "publishDate":
		ns.PublishDate = value
		break
	case "tags":
		ns.Tags = value
		break
	case "comment":
		ns.AcceptComment = value
		break
	default:
		ns.Password = value
		break
	}
}
