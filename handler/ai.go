package handler

import (
	"github.com/gin-gonic/gin"

	"context"
	"encoding/json"
	"fmt"
	"go-server/config"
	"go-server/model"
	"go-server/service"
	"os"

	"github.com/volcengine/volcengine-go-sdk/service/arkruntime"
	"github.com/volcengine/volcengine-go-sdk/service/arkruntime/model/responses"
)

// PostAiClassification AI智能分类账单
// @Summary      AI智能分类账单
// @Description  根据支付宝账单数据，调用AI自动匹配支出/收入分类
// @Tags         AI
// @Accept       json
// @Produce      json
// @Param        request  body      service.PostAiClassificationRequest  true  "账单数据"
// @Success      200      {object}  model.Response[service.PostAiClassificationResult]  "分类结果"
// @Failure      400      {object}  model.Response[any]                  "请求参数错误"
// @Failure      500      {object}  model.Response[any]                  "服务器错误"
// @Router       /ai/ai-classification [post]
func PostAiClassification(c *gin.Context) {
	var req service.PostAiClassificationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, model.Error(400, err.Error()))
		return
	}

	// 支出分类
	var list []model.Classification
	config.DB.Order("sort ASC").Find(&list)
	tree := service.BuildTree(0, list)
	treeBytes, err := json.Marshal(tree)
	if err != nil {
		panic(err)
	}
	treeStr := string(treeBytes)

	// 收入分类
	var list1 []model.IncomeClassification
	config.DB.Order("sort ASC").Find(&list1)
	tree1 := service.BuildTree(0, list1)
	treeBytes1, err := json.Marshal(tree1)
	if err != nil {
		panic(err)
	}
	treeStr1 := string(treeBytes1)

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
						Text: `我会给你两个树结构，你需要根据这个树结构，这是账单的分类的类型，我需要你根据我给你的支付宝账单的信息，依次给我返回每条账单对应的合适的分类，
							你给我的返回json，json的结构是一个对象，对象的第一个属性是classifications，
							classifications是一个数组，数组的每个元素是一个对象，对象的属性是firstLevelClassificationId和firstLevelClassificationName已经secondLevelClassificationId和secondLevelClassificationName，这四个字段分别是一级分类的id，一级分类的名称，二级分类的id，二级分类的名称，二级分类必须是一级分类的children，
							尽量给每条账单寻找分类，如果没有合适的分类，上述四个字段为设置为null，
							你返回的列表长度应该是和我给你账单的长度是一样的，
							如果有什么意外情况，你可以在返回json中增加一个message属性，message属性是一个字符串，字符串是对异常情况的描述，有些数据如果你没有给到对应的类型，也可以把说明放在message中` +
							`数据里面如果没有收支类型，也就说没有定义这条记录是收入还是支出，也设置为null` +
							`这个是支出分类树` +
							treeStr +
							`这个是收入分类树` +
							treeStr1 +
							`这个是数据` +
							req.Data,
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

	var arkResp service.PostAiClassificationResponse

	err = json.Unmarshal([]byte(text), &arkResp)
	if err != nil {
		panic(err)
	}

	c.JSON(200, model.SuccessWithData(service.PostAiClassificationResult{
		ArkRes:     resp,
		ArkTextObj: arkResp,
	}))

	// var list []model.IncomeClassification
	// config.DB.Order("sort ASC").Find(&list)

	// tree := service.BuildTree(0, list)

	// c.JSON(200, model.SuccessWithData(tree))
}
