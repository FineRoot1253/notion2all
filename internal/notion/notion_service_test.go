package notion

import (
	"context"
	"github.com/fineroot1253/notion2all/internal/common/logger"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"testing"
)

var commonService Service

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
			serv, err := NewService(tt.args.ctx, tt.args.token, tt.args.logger)
			if err != nil {
				if tt.wantErr {
					log.Println(err)
				} else {
					log.Panicln(err)
				}
			} else {
				log.Println("NewService complete")
				commonService = serv
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
