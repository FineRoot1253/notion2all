package model

import (
	"strings"
)

type NotionConfiguration struct {
	Notion struct {
		Token        string `json:"token"`
		TablePageUrl string `json:"table_page_url"`
		DownloadPath string `json:"download_path"`
	} `json:"notion"`
}

func (nc NotionConfiguration) IsEmpty() bool {
	if nc.Notion.Token == "" || nc.Notion.TablePageUrl == "" || nc.Notion.DownloadPath == "" {
		return true
	}
	return false
}

func (nc NotionConfiguration) GetTableBlockId() string {
	tableUrlArr := strings.Split(nc.Notion.TablePageUrl, "/")

	tableData := tableUrlArr[len(tableUrlArr)-1]

	return strings.Split(tableData, "?")[0]
}

func (nc NotionConfiguration) GetTableViewId() string {
	tableUrlArr := strings.Split(nc.Notion.TablePageUrl, "/")

	tableData := tableUrlArr[len(tableUrlArr)-1]
	splitTableData := strings.Split(tableData, "=")

	return splitTableData[len(splitTableData)-1]
}
