package service

import (
	"github.com/volcengine/volcengine-go-sdk/service/arkruntime/model/responses"
)

type PostAiClassificationRequest struct {
	Data string `json:"data" binding:"required"`
}

type AiClassificationItem struct {
	FirstLevelClassificationId    int64  `json:"firstLevelClassificationId"`
	FirstLevelClassificationName  string `json:"firstLevelClassificationName"`
	SecondLevelClassificationId   int64  `json:"secondLevelClassificationId"`
	SecondLevelClassificationName string `json:"secondLevelClassificationName"`
}

type PostAiClassificationResponse struct {
	Classifications []AiClassificationItem `json:"classifications"`
	Message         string                 `json:"message,omitempty"`
}

type PostAiClassificationResult struct {
	ArkRes     *responses.ResponseObject    `json:"arkRes" swaggertype:"object"`
	ArkTextObj PostAiClassificationResponse `json:"arkTextObj"`
}
