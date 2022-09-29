package notion

import (
	"context"
	"github.com/fineroot1253/notion2all/internal/common/logger"
	"github.com/fineroot1253/notion2all/internal/common/utils"
	notionModel "github.com/fineroot1253/notion2all/internal/notion/model"
	tistoryModel "github.com/fineroot1253/notion2all/internal/tistory/model"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"testing"
)

var (
	commonService Service
	notionConfig  notionModel.NotionConfiguration
	tistoryConfig tistoryModel.TistoryConfiguration
)

func init() {
	var cfgFile string
	ctx := context.Background()

	cfgFile = "./testdata/config.json"

	if err := utils.ParseConfigFile(cfgFile, &notionConfig, &tistoryConfig); err != nil {
		log.Panicln(err)
	}

	newService, err := NewService(ctx, notionConfig, logger.NewLogTemplate(logger.CommonLogRunner{}))
	if err != nil {
		log.Panicln(err)
	}
	commonService = newService

}

func TestNewService(t *testing.T) {
	type args struct {
		ctx    context.Context
		token  string
		logger logger.LogTemplate
	}
	ctx := context.Background()
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "notion NewService:[success]",
			args: args{
				ctx:    ctx,
				token:  "./testdata/data/testzip/",
				logger: logger.NewLogTemplate(logger.CommonLogRunner{}),
			},
			wantErr: false,
		},
		{
			name: "notion NewService:[failure](모든 필드 생략)",
			args: args{
				ctx:    ctx,
				token:  "./testdata/data/testzip/",
				logger: logger.NewLogTemplate(logger.CommonLogRunner{}),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewService(tt.args.ctx, notionConfig, tt.args.logger)
			if err != nil {
				if tt.wantErr {
					log.Println(err)
				} else {
					log.Panicln(err)
				}
			} else {
				log.Println("NewService complete")
			}
		})
	}
}

func Test_unzip(t *testing.T) {
	type args struct {
		src         string
		destDirPath string
	}
	tests := []struct {
		name    string
		field   args
		wantErr bool
	}{
		{
			name: "zip 파일 해제 테스트:[success]",
			field: args{
				src:         "./testdata/testzip.zip",
				destDirPath: "./testdata/data/testzip/",
			},
			wantErr: false,
		},
		{
			name: "zip 파일 해제 테스트:[success] (data/testzip 디렉토리가 이미 있는 경우)",
			field: args{
				src:         "./testdata/testzip.zip",
				destDirPath: "./testdata/data/testzip/",
			},
			wantErr: false,
		},
		{
			name: "zip 파일 해제 테스트:[failure] (src 필드 생략)",
			field: args{
				src:         "",
				destDirPath: "./testdata/data/testzip/",
			},
			wantErr: true,
		},
		{
			name: "zip 파일 해제 테스트:[failure] (destDirPath 필드 생략)",
			field: args{
				src:         "./testdata/testzip.zip",
				destDirPath: "",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			err := unzip(tt.field.src, tt.field.destDirPath)
			if err != nil {
				if tt.wantErr {
					log.Println(err)
				} else {
					log.Panicln(err)
				}
			} else {
				fileInfos, err := ioutil.ReadDir(tt.field.destDirPath)
				if err != nil {
					log.Panicln(err)
				}
				for _, info := range fileInfos {
					assert.NotEmpty(t, info)
				}
				log.Println(tt.name, " Complete")
			}
			// 2차 성공 테스트를 위해 삭제
			if strings.Contains(tt.name, "success") {
				if err := os.RemoveAll("./testdata/data"); err != nil {
					log.Panicln(err)
				}
			} else {
				// 2차 성공 테스트 이후 반복 테스트를 위해 생성
				if err := os.MkdirAll("./testdata/data", 0755); err != nil {
					log.Panicln(err)
				}
			}

		})
	}
}

func Test_service_GetDeployPostList(t *testing.T) {

	tests := []struct {
		name    string
		want    []notionModel.NotionPostData
		wantErr bool
	}{
		{
			name:    "배포준비 게시글 리스트 가져오기 테스트:[success]",
			wantErr: false,
		},
		{
			name:    "배포준비 게시글 리스트 가져오기 테스트:[failure]",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			list, err := commonService.GetDeployPostList()
			if err != nil {
				if tt.wantErr {
					log.Println(err)
				} else {
					log.Panicln(err)
				}
			}

			assert.NotEmpty(t, list)
		})
	}
}
