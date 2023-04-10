package server

import (
	middleware "bcraft/api/middlewares"
	authModule "bcraft/api/modules/auth"
	"bcraft/api/modules/fileManagement"
	recipesModule "bcraft/api/modules/recipes"
	"github.com/gin-gonic/gin"
	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
)

type Router struct{}

func (r *Router) Init() *gin.Engine {
	router := gin.New()
	gin.SetMode(gin.DebugMode)

	authController := &authModule.AuthController{}
	recipeController := &recipesModule.RecipesController{}
	fileController := fileManagement.FileManagementController{}

	router.MaxMultipartMemory = 8 << 20 // 8 MiB
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.Static("/content", fileManagement.ContentDirPath)
	v1 := router.Group("/v1")
	{
		auth := v1.Group("/auth")
		{
			auth.POST("/login", authController.Login)
			auth.POST("/register", authController.Register)
		}

		file := v1.Group("/file").Use(middleware.Auth())
		{
			file.POST("/upload", fileController.Upload)
		}

		recipe := v1.Group("/recipes")
		{
			recipe.GET("", recipeController.Get)
			recipe.POST("/filter", recipeController.GetFiltered)
			recipe.GET("/:id", recipeController.GetById)
		}

		recipeAuth := v1.Group("/recipes").Use(middleware.Auth())
		{
			recipeAuth.PUT("/:id", recipeController.Update)
			recipeAuth.POST("", recipeController.Create)
			recipeAuth.POST("/rate/:id", recipeController.Rate)
			recipeAuth.DELETE("/:id", recipeController.Delete)
		}
	}

	return router
}
