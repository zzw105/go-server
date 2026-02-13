package router

import (
	"go-server/handler"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func InitRouter() *gin.Engine {
	r := gin.Default()

	// Swagger 文档路由
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	upload := r.Group("/upload")
	{
		upload.POST("", handler.UploadFile)
	}

	classification := r.Group("/classification")
	{
		classification.GET("", handler.GetClassificationTree)
		classification.PUT("", handler.UpdateClassificationTree)
	}

	incomeClassification := r.Group("/income-classification")
	{
		incomeClassification.GET("", handler.GetIncomeClassificationTree)
		incomeClassification.PUT("", handler.UpdateIncomeClassificationTree)
	}

	export := r.Group("/export")
	{
		export.POST("", handler.ExportExcel)
	}

	return r
}
