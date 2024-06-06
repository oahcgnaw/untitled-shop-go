package main

import (
	"backend-go/db"
	"backend-go/internal/handlers"
	"log"
	"os"

	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
    err := godotenv.Load()
    if err != nil {
        log.Fatalf("Error loading .env file")
    }
	db.ConnectDB()
    router := gin.Default()

    router.Use(cors.New(cors.Config{
        AllowOrigins:     []string{"*"},
        AllowMethods:     []string{"GET", "POST", "PATCH", "DELETE", "OPTIONS"},
        AllowHeaders:     []string{"*"},
        MaxAge:           12 * time.Hour,
    }))

    // Create the upload folder if it doesn't exist
    os.MkdirAll("upload/images", os.ModePerm)

    router.Static("images", "upload/images")

    api := router.Group("/api/v1")
    {
        product := api.Group("/product")
        {
            product.POST("", handlers.CreateProduct)
            product.GET("/:id", handlers.GetProduct)
            product.DELETE("/:id", handlers.DeleteProduct)
            product.GET("", handlers.GetAllProducts)
            product.GET("/category/:category", handlers.GetProductsByCategory)
        }

        cart := api.Group("/cartitems")
        {
            cart.POST("", handlers.GetCartItems)
            cart.PATCH("", handlers.UpdateCartItem)
            cart.DELETE("/:id", handlers.DeleteCartItem)
        }

        api.POST("/signup", handlers.CreateUser)
        api.POST("/login", handlers.LoginUser)
        api.GET("/popular/:category", handlers.GetPopularProducts)

        newcollection := api.Group("/newcollections")
        {
            newcollection.GET("", handlers.GetNewCollections)
            newcollection.GET("/:category", handlers.GetNewProductsByCategory)
        }

        api.POST("/upload", handlers.UploadFile)
        
    }

    router.Run(":"+os.Getenv("PORT"))

}
