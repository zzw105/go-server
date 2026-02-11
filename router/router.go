package router

import (
	"go-server/handler"

	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	r := gin.Default()

	user := r.Group("/users")
	{
		user.POST("", handler.CreateUser)
		user.GET("", handler.GetUserList)
		user.GET("/:id", handler.GetUser)
		user.PUT("/:id", handler.UpdateUser)
		user.DELETE("/:id", handler.DeleteUser)
	}

	upload := r.Group("/upload")
	{
		upload.POST("", handler.UploadFile)
	}

	classification := r.Group("/classification")
	{
		classification.GET("", handler.GetClassificationTree)
		classification.PUT("", handler.UpdateClassificationTree)
	}

	export := r.Group("/export")
	{
		export.POST("", handler.ExportExcel)
	}

	return r
}
