package utils

import (
	"encoding/json"
	"github.com/fineroot1253/notion2all/internal/common"
	notionModel "github.com/fineroot1253/notion2all/internal/notion/model"
	tistoryModel "github.com/fineroot1253/notion2all/internal/tistory/model"
	"io/ioutil"
	"os"
)

func ParseConfigFile(cfgFile string, notionConfig *notionModel.NotionConfiguration, tistoryConfig *tistoryModel.TistoryConfiguration) error {
	openFile, err := os.Open(cfgFile)
	if err != nil {
		return &common.CommonError{Func: "parseConfigFile", Data: common.CanNotFoundFile_Error.String(), Err: err}
	}

	defer openFile.Close()
	byteData, err := ioutil.ReadAll(openFile)
	if err != nil {
		return &common.CommonError{Func: "parseConfigFile", Data: common.Marshal_Error.String(), Err: err}
	}

	if err := json.Unmarshal(byteData, notionConfig); err != nil {
		return &common.CommonError{Func: "parseConfigFile", Data: common.Marshal_Error.String(), Err: err}
	}

	if err := json.Unmarshal(byteData, tistoryConfig); err != nil {
		return &common.CommonError{Func: "parseConfigFile", Data: common.Marshal_Error.String(), Err: err}
	}
	return nil
}
