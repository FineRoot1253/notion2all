package utils

import (
	"github.com/fineroot1253/notion2all/internal/common"
	"regexp"
	"strings"
	"time"
)

func ConvertPageUrlToPageId(contentUrl string) string {
	urlArr := strings.Split(contentUrl, "-")
	return ConvertId(urlArr[len(urlArr)-1])
}

/*	ConvertIdList pageid 리스트를 uuid-v4 형식으로 변환
	@Param	[]string // pageId 리스트
	@Return	[]string // 변환된 pageId 리스트
*/
func ConvertIdList(idList []string) []string {
	for _, id := range idList {
		id = ConvertId(id)
	}
	return idList
}

/*	ConvertId pageid를 uuid-v4 형식으로 변환

	재사용을 위해 추출함

	@Param	string // pageId

	@Return	string // 변환된 pageId
*/
func ConvertId(id string) string {
	return regexp.MustCompile(`(.{8})(.{4})(.{4})(.{4})(.+)`).ReplaceAllString(id, `$1-$2-$3-$4-$5`)
}

func ConvertStringToTimeStamp(timeString string) (string, error) {
	parse, err := time.Parse("yyyy/mm/dd", timeString)
	if err != nil {
		return "", &common.CommonError{
			Func: "ConvertStringToTimeStamp",
			Data: common.Unexpected_Error.String(),
			Err:  err,
		}
	}
	return parse.Format("yyyy-mm-dd hh:mm:ss"), nil
}
