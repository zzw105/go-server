package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
)

type ExportItem struct {
	Name   string  `json:"name"`
	Age    int     `json:"age"`
	Amount float64 `json:"amount"`
}

func ExportExcel(c *gin.Context) {
	var list []ExportItem

	// 1. 接收前端 JSON
	if err := c.ShouldBindJSON(&list); err != nil {
		c.JSON(400, gin.H{"msg": err.Error()})
		return
	}

	// 2. 创建 Excel
	f := excelize.NewFile()
	sheet := "Sheet1"
	f.SetSheetName("Sheet1", sheet)

	// 3. 表头
	headers := []string{"姓名", "年龄", "金额"}
	for i, h := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheet, cell, h)
	}

	// 4. 写数据
	for i, item := range list {
		row := i + 2
		f.SetCellValue(sheet, "A"+strconv.Itoa(row), item.Name)
		f.SetCellValue(sheet, "B"+strconv.Itoa(row), item.Age)
		f.SetCellValue(sheet, "C"+strconv.Itoa(row), item.Amount)
	}

	// 5. 写入内存
	buf, err := f.WriteToBuffer()
	if err != nil {
		c.JSON(500, gin.H{"msg": "生成 Excel 失败"})
		return
	}

	// 6. 设置下载头
	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Disposition", "attachment; filename=export.xlsx")
	c.Data(200, "application/octet-stream", buf.Bytes())
}
