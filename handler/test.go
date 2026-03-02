package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"go-server/config"
	"go-server/model"
	"go-server/service"
	"os"

	"github.com/gin-gonic/gin"

	"github.com/volcengine/volcengine-go-sdk/service/arkruntime"
	"github.com/volcengine/volcengine-go-sdk/service/arkruntime/model/responses"
)

func PostTest(c *gin.Context) {

	var list []model.Classification
	config.DB.Order("sort ASC").Find(&list)
	tree := service.BuildTree(0, list)
	treeBytes, err := json.Marshal(tree)
	if err != nil {
		panic(err)
	}
	treeStr := string(treeBytes)

	client := arkruntime.NewClientWithApiKey(
		//通过 os.Getenv 从环境变量中获取 ARK_API_KEY
		os.Getenv("ARK_API_KEY"),
		arkruntime.WithBaseUrl("https://ark.cn-beijing.volces.com/api/v3"),
	)
	// 创建一个上下文，通常用于传递请求的上下文信息，如超时、取消等
	ctx := context.Background()
	inputMessage := &responses.ItemInputMessage{
		Role: responses.MessageRole_user,
		Content: []*responses.ContentItem{
			{
				Union: &responses.ContentItem_Text{
					Text: &responses.ContentItemText{
						Type: responses.ContentItemType_input_text,
						Text: `我会给你一个树结构，你需要根据这个树结构，这是账单的分类的类型，我需要你根据我给你的支付宝账单的信息，依次给我返回每条账单对应的合适的分类，
							你给我的返回json，json的结构是一个对象，对象的第一个属性是classifications，
							classifications是一个数组，数组的每个元素是一个对象，对象的属性是first_level_classification_id和first_level_classification_name已经second_level_classification_id和second_level_classification_name，这四个字段分别是一级分类的id，一级分类的名称，二级分类的id，二级分类的名称，二级分类必须是一级分类的children，
							如果没有合适的分类，上述四个字段为空就行，
							你返回的列表长度应该是和我给你账单的长度是一样的，
							如果有什么意外情况，你可以在返回json中增加一个message属性，message属性是一个字符串，字符串是对异常情况的描述` +
							treeStr +
							``,
					},
				},
			},
		},
	}
	resp, err := client.CreateResponses(ctx, &responses.ResponsesRequest{
		Model: "doubao-seed-1-8-251228",
		Input: &responses.ResponsesInput{
			Union: &responses.ResponsesInput_ListValue{
				ListValue: &responses.InputItemList{ListValue: []*responses.InputItem{{
					Union: &responses.InputItem_InputMessage{
						InputMessage: inputMessage,
					},
				}}},
			},
		},
		Thinking: &responses.ResponsesThinking{
			Type: responses.ThinkingType_disabled.Enum(),
		},
	})
	if err != nil {
		fmt.Printf("response error: %v\\n", err)
		return
	}
	fmt.Println(resp)
	text := resp.Output[0].GetOutputMessage().GetContent()[0].GetText().Text

	type ArkText struct {
		FirstLevelClassificationId    string `json:"first_level_classification_id"`
		FirstLevelClassificationName  string `json:"first_level_classification_name"`
		SecondLevelClassificationId   string `json:"second_level_classification_id"`
		SecondLevelClassificationName string `json:"second_level_classification_name"`
	}

	type ArkResponse struct {
		Classifications []ArkText `json:"classifications"`
		Message         string    `json:"message,omitempty"`
	}

	var arkResp ArkResponse

	err = json.Unmarshal([]byte(text), &arkResp)
	if err != nil {
		panic(err)
	}

	type testRes struct {
		ArkRes     *responses.ResponseObject `json:"arkRes"`
		ArkTextObj ArkResponse               `json:"arkTextObj"`
	}

	c.JSON(200, model.SuccessWithData(testRes{
		ArkRes:     resp,
		ArkTextObj: arkResp,
	}))
}
