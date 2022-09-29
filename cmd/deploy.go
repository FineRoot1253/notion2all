/*
Copyright © 2022 Jun-Geun Hong : dev.fineroot1253@gmail.com

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"context"
	"github.com/fineroot1253/notion2all/internal/common/logger"
	"github.com/fineroot1253/notion2all/internal/common/utils"
	"github.com/fineroot1253/notion2all/internal/deployment"
	"github.com/fineroot1253/notion2all/internal/notion"
	notionModel "github.com/fineroot1253/notion2all/internal/notion/model"
	"github.com/fineroot1253/notion2all/internal/tistory"
	tistoryModel "github.com/fineroot1253/notion2all/internal/tistory/model"
	"github.com/fineroot1253/tistoryAPI"
	"github.com/spf13/cobra"
	"log"
)

var deployService deployment.Service

// deployCmd represents the deployment command
// 동작 순서
// 1. notion 서비스 세팅 [필수]
// 2. tistory 서비스 세팅 TODO 플랫폼별로 이 주입될 서비스를 늘려야 함
// 3. [notion 서비스] notion 배포용 테이블 탐색
// 4. [notion 서비스] notion 테이블에서 파싱 -> 일종의 배포 실행 계획
// 5. 배포 시작
// 		5-1. 블로그 카테고리를 먼저 탐색하여 카테고리 id를 onmemory에 세팅 TODO 카테고리 ID 구조체 필요
// 		5-2. 모든 테이블을 순회하여 블로그 컨텐츠 page id : {블로그 포스트 명, 블로그 컨텐츠(html 포멧), 카테고리명, 공개 비공개 여부, 배포 시간} 구조체 리스트 세팅 TODO notion -> platform 파싱 데이터 구조체 필요
//			5-2-1. html 다운로드후 zip 파일 해제
//			5-2-2. html 파일 열고 on-memory 로드
//			5-2-3-a. 파일이 1개가 아닐시 이미지 파일 포함 버전 => 이미지 파일 base64 변환후 img 태그 src에 replace
//			5-2-3-b. 파일이 html파일 딱 1개일시 이미지 미포함 버전 => 그냥 무시
//		5-4. 5-2 구조체 리스트를 토대로 tistoryAPI를 통해 배포 시작
//
var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Start distributing notion content to other platforms.",
	Long: `Start to deployment notion contents to other platform.
Depending on the config.json content, you can specify the platform on which it is deployed.
`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := deployService.Deploy(); err != nil {
			log.Panicln(err)
		}
	},
}

/*
	설정 파일 바인딩
	현 실행 파일로부터 위치를 찾으며 따로 플래그 설정을 하게 되면 그 설정을 토대로 탐색한다.
	만약 이 바인딩이 실패하면 실행을 중단한다.
	TODO 추후 업데이트는 셀레니움 도입 예정, 늘려갈 플랫폼들 대부분 OPEN api 지원 불가한 플랫폼이기 때문
*/
func init() {
	var cfgFile string
	var notionConfig notionModel.NotionConfiguration
	var tistoryConfig tistoryModel.TistoryConfiguration
	ctx := context.Background()
	rootCmd.AddCommand(deployCmd)

	rootCmd.Flags().StringVar(&cfgFile, "config", "./config.json", "config file (default is ./config.json)")
	//rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	if err := utils.ParseConfigFile(cfgFile, &notionConfig, &tistoryConfig); err != nil {
		log.Panicln(err)
	}

	// 1) 로그 템플릿 생성
	logTemplate := logger.NewLogTemplate(logger.CommonLogRunner{})
	// 2) notion service 생성
	notionService, err := notion.NewService(ctx, notionConfig, logTemplate)
	if err != nil {
		log.Panicln(err)
	}
	// 3) tistory service 생성
	tistoryAPIService, err := tistoryAPI.NewService(ctx, tistoryConfig.GetUserData())
	if err != nil {
		log.Panicln(err)
	}

	tistoryService := tistory.NewService(tistoryConfig, tistoryAPIService)

	deployService = deployment.NewService(logTemplate, notionService, tistoryService)

}
