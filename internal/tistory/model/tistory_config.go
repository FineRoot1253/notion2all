package model

import tistoryAPIModel "github.com/fineroot1253/tistoryAPI/model"

type TistoryConfiguration struct {
	Tistory struct {
		BlogName         string `json:"blog_name"`
		ChromeDriverPath string `json:"chrome_driver_path"`
		tistoryAPIModel.UserData
		UserId  string `json:"user_id"`
		UserPwd string `json:"user_pwd"`
	} `json:"tistory"`
}

func (tc TistoryConfiguration) GetUserData() tistoryAPIModel.UserData {
	return tistoryAPIModel.UserData{
		ClientId:          tc.Tistory.ClientId,
		SecretKey:         tc.Tistory.SecretKey,
		RedirectUrl:       tc.Tistory.RedirectUrl,
		AuthorizationCode: tc.Tistory.AuthorizationCode,
	}
}
