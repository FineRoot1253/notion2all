package model

type TistoryConfiguration struct {
	Tistory struct {
		BlogName          string `json:"blog_name"`
		ChromeDriverPath  string `json:"chrome_driver_path"`
		ClientId          string `json:"client_id"`
		SecretKey         string `json:"secret_key"`
		RedirectUrl       string `json:"redirect_url"`
		AuthorizationCode string `json:"authorization_code"`
		UserId            string `json:"user_id"`
		UserPwd           string `json:"user_pwd"`
	} `json:"tistory"`
}
