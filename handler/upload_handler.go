package handler

import (
	"encoding/csv"
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"

	"go-server/model"

	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

// UploadFile 上传 CSV 文件
// @Summary      上传 CSV 文件（GB2312 编码）
// @Description  上传 GB2312/GBK 编码的 CSV 文件，自动转码后解析并返回数据行列表
// @Tags         上传
// @Accept       multipart/form-data
// @Produce      json
// @Param        file  formData  file    true  "CSV 文件（GB2312 编码）"
// @Success      200   {object}  model.Response[[][]string]  "解析后的数据行列表"
// @Failure      400   {object}  model.Response[bool]        "文件获取或解析失败"
// @Failure      500   {object}  model.Response[bool]        "服务器错误"
// @Router       /upload [post]
func UploadFile(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, model.Error(400, "获取文件失败: "+err.Error()))
		return
	}

	f, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.Error(500, "打开文件失败: "+err.Error()))
		return
	}
	defer f.Close()

	// 1️⃣ GB2312 -> UTF-8 转码
	reader := csv.NewReader(transform.NewReader(f, simplifiedchinese.GBK.NewDecoder()))
	reader.FieldsPerRecord = -1

	// 2️⃣ 跳过表头（可选）
	_, err = reader.Read()
	if err != nil {
		c.JSON(http.StatusBadRequest, model.Error(400, "读取表头失败: "+err.Error()))
		return
	}
	var recordsList [][]string
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			c.JSON(http.StatusBadRequest, model.Error(400, "读取 CSV 行失败: "+err.Error()))
			return
		}

		if len(record) < 5 {
			continue
		}

		if record[0] == "交易时间" {
			continue
		}
		// [交易时间 交易分类 交易对方 对方账号 商品说明 收/支 金额 收/付款方式 交易状态 交易订单号 商家订单号 备注 ]
		fmt.Println(record)
		recordsList = append(recordsList, record)
	}

	c.JSON(http.StatusOK, model.SuccessWithData(recordsList))
}
