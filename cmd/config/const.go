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

	NOTION_HTML_DATA_PATH   = "./data"
	NOTION_CODE_BLOCK_THEME = "atom-one-dark"

	NOTION_ENQUEUE_TASK_PATH = "/enqueueTask"
	NOTION_GET_TASKS_PATH    = "/getTasks"

	NOTION_CSS       = "https://rawcdn.githack.com/ppuep94/n2t/5ef4dc01e9d6336341e9ab95bb71672f9d3a3dc9/assets/css/style2.css"
	NOTION_STYLE_TAG = `<link rel="stylesheet" href="//cdnjs.cloudflare.com/ajax/libs/highlight.js/10.1.2/styles/{code_theme}.min.css">
				<style> 
.bookmark {
    text-decoration: none;
    max-height: 8em;
    padding: 0;
    display: flex;
    width: 100%;
    align-items: stretch;
}
.source {
    border: 1px solid #ddd;
    border-radius: 3px;
    padding: 1.5em;
    word-break: break-all;
}
a, a.visited {
    color: inherit;
    text-decoration: underline;
}
.bookmark-image {
    width: 33%;
    flex: 1 1 180px;
    display: block;
    position: relative;
    object-fit: cover;
    border-radius: 1px;
}
.bookmark-info {
    flex: 4 1 180px;
    padding: 12px 14px 14px;
    display: flex;
    flex-direction: column;
    justify-content: space-between;
}
.sans {
    font-family: ui-sans-serif, -apple-system, BlinkMacSystemFont, "Segoe UI", Helvetica, "Apple Color Emoji", Arial, sans-serif, "Segoe UI Emoji", "Segoe UI Symbol";
}
img {
    max-width: 100%;
}
.page-header-icon {
    font-size: 3rem;
    margin-bottom: 1rem;
}
.icon {
    display: inline-block;
    max-width: 1.2em;
    max-height: 1.2em;
    text-decoration: none;
    vertical-align: text-bottom;
    margin-right: 0.5em;
}</style>
                <script src="//cdnjs.cloudflare.com/ajax/libs/highlight.js/10.1.2/highlight.min.js"></script>
                <script>hljs.initHighlightingOnLoad();</script>`
)
