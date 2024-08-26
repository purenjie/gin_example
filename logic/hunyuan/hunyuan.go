package hunyuan

import (
	"log"
	"sync"

	"gin.example.com/entity/config"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	hunyuan "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/hunyuan/v20230901"
)

type HunyuanClient struct {
	Client *hunyuan.Client
}

var client *HunyuanClient
var hunyuanOnce sync.Once

func GetHunYuanClient() *HunyuanClient {
	hunyuanOnce.Do(func() {
		config := config.GetHunyuanConfig()
		credential := common.NewCredential(config.SecretID, config.SecretKey)
		cpf := profile.NewClientProfile()
		cpf.HttpProfile.Endpoint = "hunyuan.tencentcloudapi.com"
		hunyuanClient, _ := hunyuan.NewClient(credential, "", cpf)
		client = &HunyuanClient{Client: hunyuanClient}
	})
	return client
}

func (c *HunyuanClient) SendMessage(user, message, model string) (string, error) {
	request := hunyuan.NewChatCompletionsRequest()
	request.Model = common.StringPtr(model)
	request.Messages = []*hunyuan.Message{
		{
			Role:    common.StringPtr(user),
			Content: common.StringPtr(message),
		},
	}
	request.Stream = common.BoolPtr(false)
	response, err := c.Client.ChatCompletions(request)
	log.Printf("req: %s, rsp: %s", entity.ObjectToStr(request), entity.ObjectToStr(response))
	if err != nil {
		return "", err
	}
	var msg string
	for _, choice := range response.Response.Choices {
		msg += *choice.Message.Content
	}
	return msg, nil
}
