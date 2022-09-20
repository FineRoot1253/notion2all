package config

const (
	// 릴리즈 버전
	RELEASE_VERSION = "v0.0.1"

	COLUMN_TITLE    = "제목"
	COLUMN_CATEGORY = "카테고리"
	COLUMN_TAG      = "태그"
	COLUMN_STATUS   = "상태"
	COLUMN_URL      = "링크"

	POST_STATUS_UPLOAD   = "발행 요청"
	POST_STATUS_MODIFY   = "수정 요청"
	POST_STATUS_COMPLETE = "발행 완료"

	NOTION_V3_BASE_URL = "https://www.notion.so/api/v3"

	NOTION_ENQUEUE_TASK_PATH = "/enqueueTask"
	NOTION_GET_TASKS_PATH    = "/getTasks"
)
