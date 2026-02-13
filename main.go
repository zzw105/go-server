package main

import (
	"go-server/config"
	_ "go-server/docs" // Swagger 生成的文档
	"go-server/router"
)

// @title           账本系统 API
// @version         1.0
// @description     账本管理系统的后端 API 接口文档
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.email  your-email@example.com

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:5200
// @BasePath  /

// @schemes http https
func main() {
	config.InitDB()
	r := router.InitRouter()
	r.Run(":5200")
}
