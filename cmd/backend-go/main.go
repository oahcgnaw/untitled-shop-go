package main

import (
	"backend-go/db"
	"backend-go/internal/handlers"

	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	db.ConnectDB()
    router := gin.Default()

    router.Use(cors.New(cors.Config{
        AllowOrigins:     []string{"*"},
        AllowMethods:     []string{"GET", "POST", "PATCH", "DELETE", "OPTIONS"},
        AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
        MaxAge:           12 * time.Hour,
    }))

    api := router.Group("/api/v1")
    {
        product := api.Group("/product")
        {
            product.POST("/", handlers.CreateProduct)
            product.GET("/:id", handlers.GetProduct)
            product.DELETE("/:id", handlers.DeleteProduct)
            product.GET("/", handlers.GetAllProducts)
            product.GET("/category/:category", handlers.GetProductsByCategory)
        }

        cart := api.Group("/cartitems")
        {
            cart.GET("/", handlers.GetCartItems)
            cart.PATCH("/", handlers.UpdateCartItem)
            cart.DELETE("/:id", handlers.DeleteCartItem)
        }

        api.POST("/signup", handlers.CreateUser)
        api.POST("/login", handlers.LoginUser)
        api.GET("/popular/:category", handlers.GetPopularProducts)

        newcollection := api.Group("/newcollections")
        {
            newcollection.GET("/", handlers.GetNewCollections)
            newcollection.GET("/:category", handlers.GetNewProductsByCategory)
        }
        
    }

    router.Run(":4000")
}
