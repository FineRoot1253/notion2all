package model

type NotionConfiguration struct {
	Notion struct {
		Token        string `json:"token"`
		TablePageUrl string `json:"table_page_url"`
		DownloadPath string `json:"download_path"`
	} `json:"notion"`
}
