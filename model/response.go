package model

// Response 统一响应结构
// swagger:model Response
type Response[T any] struct {
	Code    int    `json:"code" binding:"required"`
	Msg     string `json:"msg" binding:"required"`
	Success bool   `json:"success" binding:"required"`
	Data    T      `json:"data" binding:"required"`
}

// SuccessWithData 成功响应（带数据）
func SuccessWithData[T any](data T) Response[T] {
	return Response[T]{
		Code:    200,
		Msg:     "操作成功",
		Success: true,
		Data:    data,
	}
}

// Error 错误响应
func Error(code int, msg string) Response[bool] {
	return Response[bool]{
		Code:    code,
		Msg:     msg,
		Success: false,
		Data:    false,
	}
}
