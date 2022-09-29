package deployment

import (
	"context"
	"encoding/json"
	"github.com/fineroot1253/notion2all/internal/common/logger"
	"github.com/fineroot1253/notion2all/internal/notion"
	notionModel "github.com/fineroot1253/notion2all/internal/notion/model"
	"github.com/fineroot1253/notion2all/internal/tistory"
	tistoryModel "github.com/fineroot1253/notion2all/internal/tistory/model"
	"github.com/fineroot1253/tistoryAPI"
	tistoryAPIModel "github.com/fineroot1253/tistoryAPI/model"
	"io/ioutil"
	"log"
	"testing"
)

/*
	초기 설정
*/
var (
	notionConfiguration  notionModel.NotionConfiguration
	tistoryConfiguration tistoryModel.TistoryConfiguration
)

func init() {
	bytes, err := ioutil.ReadFile("./testdata/config.json")
	if err != nil {
		log.Panicln(err)
	}

	if err := json.Unmarshal(bytes, &notionConfiguration); err != nil {
		log.Panicln(err)
	}
	if err := json.Unmarshal(bytes, &tistoryConfiguration); err != nil {
		log.Panicln(err)
	}

}

func Test_service_parseHtml(t *testing.T) {

	ctx := context.Background()
	deployLogger := logger.NewLogTemplate(logger.CommonLogRunner{})
	tistoryAPIService, err := tistoryAPI.NewService(ctx, tistoryAPIModel.UserData{
		ClientId:          tistoryConfiguration.Tistory.ClientId,
		SecretKey:         tistoryConfiguration.Tistory.SecretKey,
		RedirectUrl:       tistoryConfiguration.Tistory.RedirectUrl,
		AuthorizationCode: tistoryConfiguration.Tistory.AuthorizationCode,
	})
	tistoryService := tistory.NewService(tistoryConfiguration, tistoryAPIService)
	if err != nil {
		log.Panicln(err)
	}

	notionService, err := notion.NewService(ctx, notionConfiguration, deployLogger)
	if err != nil {
		log.Panicln(err)
	}
	deployService := NewService(deployLogger, notionService, tistoryService)

	type args struct {
		pageHtmlDirPath string
		codeBlockLang   string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Tistory 전용 Parse 테스트:[success]",
			args: args{
				pageHtmlDirPath: "./testdata/testzip/Spring d455075af98e48c79e8229dfb4d6065d.html",
				codeBlockLang:   "atom-one-dark",
			},
			wantErr: false,
		},
		{
			name: "Tistory 전용 Parse 테스트:[failure] (필수 파라미터 제거)",
			args: args{
				pageHtmlDirPath: "",
				codeBlockLang:   "",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if resultPath, err := deployService.parseHtml(tt.args.pageHtmlDirPath, tt.args.codeBlockLang); err != nil {
				if tt.wantErr {
					log.Println(err)
					log.Println("resultPath: ", resultPath)
				} else {
					log.Panicln(err)
				}
			}
		})
	}
}
