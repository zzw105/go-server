package handler

import (
	"github.com/gin-gonic/gin"

	"go-server/config"
	"go-server/model"
	"go-server/service"
)

// GetIncomeClassificationTree 获取收入分类树
// @Summary      获取收入分类树
// @Description  获取所有收入分类，按树形结构返回
// @Tags         收入分类
// @Accept       json
// @Produce      json
// @Success      200 {object} model.Response[[]service.TreeNodeDTO] "成功返回分类树"
// @Failure      500 {object} model.Response[any] "服务器错误"
// @Router       /income-classification [get]
func GetIncomeClassificationTree(c *gin.Context) {
	var list []model.IncomeClassification
	config.DB.Order("sort ASC").Find(&list)

	tree := service.BuildTree(0, list)

	c.JSON(200, model.SuccessWithData(tree))
}

// UpdateIncomeClassificationTree 更新收入分类树
// @Summary      更新整个收入分类树
// @Description  替换数据库中的所有收入分类数据，前端传入完整的树结构
// @Tags         收入分类
// @Accept       json
// @Produce      json
// @Param        request  body      service.UpdateTreeRequest    true  "分类树数据"
// @Success      200      {object}  model.Response[bool]         "更新成功"
// @Failure      400      {object}  model.Response[any]          "请求参数错误"
// @Failure      500      {object}  model.Response[any]          "服务器错误"
// @Router       /income-classification [put]
func UpdateIncomeClassificationTree(c *gin.Context) {
	var req service.UpdateTreeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, model.Error(400, err.Error()))
		return
	}

	// 开启事务
	tx := config.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 更新树
	if err := service.UpdateTree(tx, "income_classifications", req.Tree, func() *model.IncomeClassification {
		return &model.IncomeClassification{}
	}); err != nil {
		tx.Rollback()
		c.JSON(500, model.Error(500, "更新失败: "+err.Error()))
		return
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		c.JSON(500, model.Error(500, "提交事务失败"))
		return
	}

	c.JSON(200, model.SuccessWithData(true))
}
