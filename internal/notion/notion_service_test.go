package notion

import (
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"log"
	"os"
	"testing"
)

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
			wantErr: false,
		},
		{
			name: "zip 파일 해제 테스트:[failure] (destDirPath 필드 생략)",
			field: args{
				src:         "./testdata/testzip.zip",
				destDirPath: "",
			},
			wantErr: false,
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
			if tt.name == "zip 파일 해제 테스트:[success]" {
				if err := os.RemoveAll("/testdata/data"); err != nil {
					log.Panicln(err)
				}
			} else {
				// 2차 성공 테스트 이후 반복 테스트를 위해 생성
				if err := os.MkdirAll("/testdata/data", 0755); err != nil {
					log.Panicln(err)
				}
			}

		})
	}
}
