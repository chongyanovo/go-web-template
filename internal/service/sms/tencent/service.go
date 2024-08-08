package tencent

import (
	"context"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
)

type Service struct {
	appId     *string
	signature *string
	client    *common.Client
}

func (s Service) Send(ctx context.Context, tpl string, args []string, numbers ...string) error {
	//TODO implement me
	panic("implement me")
}
