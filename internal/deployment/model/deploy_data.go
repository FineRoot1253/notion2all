package model

import "github.com/fineroot1253/tistoryAPI/model"

// DeployRequestContent 배포 컨텐츠 타입
// TODO 일단 any로 열어뒀지만 추후 any를 지우고 다른 플랫폼 컨텐츠 타입을 넣을 것
type DeployRequestContent interface {
	model.PostData | any
}

type DeployPostData[DC DeployRequestContent] struct {
	PageId      string
	PostContent DC
}
