package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"project/configs"
	"project/controllers"
	"project/database"
	"project/models"
)

func main() {
	config := configs.LoadConfig()

	db := database.Connect(config.DatabaseURL)

	database.Migrate(db)

	err := db.AutoMigrate(&models.Product{})
	if err != nil {
		return
	}
	router := gin.Default()
	router.LoadHTMLGlob("templates/*")

	productController := controllers.NewProductController(db)
	productGroup := router.Group("/products")
	{
		productGroup.POST("/", productController.CreateProduct)
		productGroup.GET("/", productController.ShowProductsPage)
		productGroup.GET("/:id", productController.GetProductByID)
		productGroup.PUT("/:id", productController.UpdateProduct)
		productGroup.DELETE("/:id", productController.DeleteProduct)
	}

	log.Fatal(router.Run(":8080"))

}
