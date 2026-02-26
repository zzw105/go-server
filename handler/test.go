package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go-server/config"
	"go-server/model"
	"go-server/service"
	"io"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
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

	apiKey := os.Getenv("ARK_API_KEY") // 建议用环境变量
	fmt.Println("apiKey:", apiKey)
	url := "https://ark.cn-beijing.volces.com/api/v3/responses"
	fmt.Println(url)
	// 构造请求体
	reqBody := map[string]interface{}{
		"model": "doubao-seed-1-8-251228",
		"input": []interface{}{
			map[string]interface{}{
				"role": "user",
				"content": []interface{}{

					map[string]interface{}{
						"type": "input_text",
						"text": `我会给你一个树结构，你需要根据这个树结构，这是账单的分类的类型，我需要你根据我给你的支付宝账单的信息，依次给我返回每条账单对应的合适的分类，
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
		"thinking": map[string]interface{}{
			"type": "disabled",
		},
	}
	fmt.Println(reqBody)
	jsonData, _ := json.Marshal(reqBody)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	// type ARK_res struct {
	// 	Name string `json:"name"`
	// 	Age  int    `json:"age"`
	// }

	var arkRes interface{}

	err = json.Unmarshal([]byte(body), &arkRes)
	if err != nil {
		panic(err)
	}

	type ContentItem struct {
		Type string `json:"type"`
		Text string `json:"text"`
	}

	type OutputItem struct {
		Role    string        `json:"role"`
		Content []ContentItem `json:"content"`
	}

	type ArkRes struct {
		Output []OutputItem `json:"output"`
	}

	var arkRes1 ArkRes

	err = json.Unmarshal([]byte(body), &arkRes1)
	if err != nil {
		panic(err)
	}
	arkText := arkRes1.Output[0].Content[0].Text
	var arkTextObj interface{}

	err = json.Unmarshal([]byte(arkText), &arkTextObj)
	if err != nil {
		panic(err)
	}
	// var arkContent interface{}

	// err = json.Unmarshal([]byte(arkRes.output[0].content[0].text), &arkContent)
	// if err != nil {
	// 	panic(err)
	// }

	// output := arkRes["output"].([]interface{}) // JSON 数组 → []interface{}
	// content := output[0].(map[string]interface{})["content"].([]interface{})
	// text := content[0].(map[string]interface{})["text"].(string)

	// 直接把火山返回的数据透传给前端
	// c.Data(resp.StatusCode, "application/json", body)
	type a struct {
		ArkRes     interface{} `json:"arkRes"`
		ArkTextObj interface{} `json:"arkTextObj"`
	}
	c.JSON(200, model.SuccessWithData(a{
		ArkRes:     arkRes,
		ArkTextObj: arkTextObj,
	}))
}
