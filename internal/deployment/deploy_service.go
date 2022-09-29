package deployment

import (
	"encoding/base64"
	"github.com/PuerkitoBio/goquery"
	"github.com/fineroot1253/notion2all/cmd/config"
	"github.com/fineroot1253/notion2all/internal/common"
	"github.com/fineroot1253/notion2all/internal/common/logger"
	"github.com/fineroot1253/notion2all/internal/deployment/model"
	"github.com/fineroot1253/notion2all/internal/notion"
	"github.com/fineroot1253/notion2all/internal/tistory"
	"golang.org/x/net/html"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

type Service interface {
	/*	Deploy NotionService를 토대로 각 플랫폼에 노션 포스트를 배포한다.
	 */
	Deploy() (model.DeployResult, error)

	parseHtml(pageHtmlPath string, codeBlockLang string) (string, error)
}

// service deployment 구조체
// logger를 통해 각 레이어별 실행 과정을 로깅한다.
type service struct {
	logger         logger.LogTemplate
	notionService  notion.Service
	tistoryService tistory.Service
	//... TODO 추후 서비스들을 추가해 나갈 예정
}

func NewService(logger logger.LogTemplate, notionService notion.Service, tistoryService tistory.Service) Service {
	return service{logger: logger, notionService: notionService, tistoryService: tistoryService}
}

func (s service) Deploy() (model.DeployResult, error) {
	// output record setting
	var deployResult model.DeployResult
	var categoryMap map[string]string

	// record initialize
	categoryList, err := s.tistoryService.GetCategoryList()
	if err != nil {
		return deployResult, err
	}
	for _, item := range categoryList {
		categoryMap[item.Name] = item.Id
	}

	// execution
	list, err := s.notionService.GetDeployPostList()
	if err != nil {
		return deployResult, err
	}
	deployResult.SuccessTaskCount = len(list)

	dataList, err := s.notionService.GetPostHtmlDataList(list)
	if err != nil {
		return deployResult, err
	}
	deployResult.TaskList = dataList
	deployResult.SuccessTaskCount = len(dataList)
	deployResult.FailureTaskCount = len(list) - len(dataList)

	for _, data := range dataList {

		parseHtml, err := s.parseHtml(data.FilePath, config.NOTION_CODE_BLOCK_THEME)
		if err != nil {
			return deployResult, err
		}
		bytes, err := ioutil.ReadFile(parseHtml)
		if err != nil {
			return deployResult, &common.CommonError{
				Func: "Deploy",
				Data: common.Unexpected_Error.String(),
				Err:  err,
			}
		}

		for _, item := range list {
			if item.GetPageId() == data.PageId {

				tistoryPostData, err := item.ToTistoryPostData(s.tistoryService.GetBlogName(), string(bytes), categoryMap[item.Category])
				if err != nil {
					return deployResult, err
				}
				postResult, err := s.tistoryService.SendPost(tistoryPostData)
				if err != nil {
					return deployResult, err
				}
				if postResult.Status != "200" {
					deployResult.FailureTaskCount++
					deployResult.DeployList = append(deployResult.DeployList, postResult)
				} else {
					deployResult.SuccessTaskCount++
					deployResult.DeployList = append(deployResult.DeployList, postResult)
				}
			}
		}

	}

	return deployResult, err
}

func (s service) parseHtml(pageHtmlDirPath string, codeBlockLang string) (string, error) {
	// 1) html 로드
	file, err := os.Open(pageHtmlDirPath)

	if err != nil {
		return "", &common.CommonError{
			Func: "parseHtml",
			Data: common.CanNotFoundFile_Error.String(),
			Err:  err,
		}
	}
	defer func() {
		file.Close()
	}()

	doc, err := goquery.NewDocumentFromReader(file)
	if err != nil {
		return "", &common.CommonError{
			Func: "parseHtml",
			Data: common.Unexpected_Error.String(),
			Err:  err,
		}
	}
	// 2) 기존 해더 정보 제거
	doc.Find("meta").Remove()
	doc.Find("title").Remove()
	doc.Find("style").Remove()

	// 3) 코드 블럭 존재시 코드블럭 수정
	doc.Find("pre").AddClass(codeBlockLang)
	//preTags.Find("class").AddClass(codeBlockLang)

	// 4) 본문 내용 로드 및 class 추가
	article := doc.Find("article").AddClass("Notion")
	//article.Find("class").AddClass("Notion")

	// 5) page-body 태그 가져오기, class 추가
	doc.Find("article").AddClass("Tistory")

	// 6) 테이블 표 제거
	article.Find("table").Remove()

	// 7) 제목 제거
	doc.Find("h1").Find("page-title").Remove()

	// 8) 이미지 태그 수정
	if err := modifyImgTag(doc, pageHtmlDirPath); err != nil {
		return "", err
	}

	// 9) body 태그에 notion style css link 추가
	newStyleTag := new(html.Node)
	newStyleTag.Type = html.ElementNode
	newStyleTag.Data = "link"
	newStyleTag.Attr = append(newStyleTag.Attr, html.Attribute{Key: "rel", Val: "stylesheet"}, html.Attribute{Key: "href", Val: config.NOTION_CSS})

	doc.Find("body").PrependNodes(newStyleTag)

	// 10) body 태그에 code style css 적용
	cssDoc, err := goquery.NewDocumentFromReader(strings.NewReader(config.NOTION_STYLE_TAG))
	if err != nil {
		return "", &common.CommonError{
			Func: "parseHtml",
			Data: common.Unexpected_Error.String(),
			Err:  err,
		}
	}
	cssDocStr, err := cssDoc.Html()
	if err != nil {
		return "", &common.CommonError{
			Func: "parseHtml",
			Data: common.Unexpected_Error.String(),
			Err:  err,
		}
	}

	doc.Find("body").PrependHtml(cssDocStr)

	// 11) html resave
	saveFile := strings.Replace(pageHtmlDirPath, ".html", "_output.html", -1)
	create, err := os.Create(saveFile)
	if err != nil {
		return "", &common.CommonError{
			Func: "parseHtml",
			Data: common.Unexpected_Error.String(),
			Err:  err,
		}
	}
	ret, err := doc.Html()
	if err != nil {
		return "", &common.CommonError{
			Func: "parseHtml",
			Data: common.Unexpected_Error.String(),
			Err:  err,
		}
	}

	_, err = create.Write([]byte(ret))
	if err != nil {
		return "", &common.CommonError{
			Func: "parseHtml",
			Data: common.Unexpected_Error.String(),
			Err:  err,
		}
	}

	return saveFile, nil
}

func modifyImgTag(doc *goquery.Document, path string) error {

	doc.Find("img").Each(func(i int, selection *goquery.Selection) {
		val, _ := selection.Attr("src")
		if !(strings.HasPrefix(val, "http") || strings.HasPrefix(val, "data:image/")) {
			imageFilePath, err := getImageFilePath(path, val)
			if err != nil {
				panic(err)
			}
			encodedImgData, err := getImageBase64String(imageFilePath)
			if err != nil {
				panic(err)
			}
			selection.SetAttr("src", encodedImgData)
		}
	})

	return nil

}

/*	getImageBase64String 이미지 파일 base64 반환
	@Param	imgFilePath string // pageId 리스트
	@Return	[]string,error // 변환된 pageId 리스트, error
*/
func getImageBase64String(imgFilePath string) (string, error) {
	bytes, err := ioutil.ReadFile(imgFilePath)
	if err != nil {
		return "", &common.CommonError{
			Func: "getImageBase64String",
			Data: common.Unexpected_Error.String(),
			Err:  err,
		}
	}
	return getBase64MimeTypeFromImage(bytes) + base64.StdEncoding.EncodeToString(bytes), nil
}

/*	getImageFilePath 이미지 파일 경로 추출 메서드
	@Param	path string // html 파일 path
	@Return	string, error // imgfile 경로, error
*/
func getImageFilePath(path string, srcVal string) (string, error) {
	imgDir := filepath.Dir(srcVal)
	ancestorPath := filepath.Dir(path)
	srcArr := strings.Split(srcVal, "/")
	unquotedImgFileName, err := url.QueryUnescape(srcArr[(len(srcArr) - 1)])
	if err != nil {
		return "", &common.CommonError{
			Func: "getImageFilePath",
			Data: common.Unexpected_Error.String(),
			Err:  err,
		}
	}
	unquotedImgDirName, err := url.QueryUnescape(imgDir)
	if err != nil {
		return "", &common.CommonError{
			Func: "getImageFilePath",
			Data: common.Unexpected_Error.String(),
			Err:  err,
		}
	}
	imgFilePath := filepath.Join(ancestorPath, unquotedImgDirName, unquotedImgFileName)
	return imgFilePath, nil
}

/*	getBase64MimeTypeFromImage 이미지 타입 추출 메서드
	@Param	bytes []byte // 이미지 []byte 데이터
	@Return	string // base64타입 이미지타입 string
*/
func getBase64MimeTypeFromImage(bytes []byte) string {
	switch http.DetectContentType(bytes) {
	case "image/jpeg":
		return "data:image/jpeg;base64,"
	case "image/png":
		return "data:image/png;base64,"
	default:
		return ""
	}
}
