package handlers

import (
	"backend-go/db"
	"backend-go/internal/models"
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// POST /api/v1/product
func CreateProduct(c *gin.Context) {
    var product models.Product
    if err := c.ShouldBindJSON(&product); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    product.Date = time.Now()
    product.Available = true
    collection := db.GetCollection("products")
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    result, err := collection.InsertOne(ctx, product)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"id": result.InsertedID})
}

// GET /api/v1/product/:id
func GetProduct(c *gin.Context) {
    id := c.Param("id")
    objID, _ := primitive.ObjectIDFromHex(id)
    var product models.Product

    collection := db.GetCollection("products")
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    err := collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&product)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
        return
    }

    c.JSON(http.StatusOK, product)
}

// DELETE /api/v1/product/:id
func DeleteProduct(c *gin.Context) {
    id := c.Param("id")
    objID, _ := primitive.ObjectIDFromHex(id)

    collection := db.GetCollection("products")
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    result, err := collection.DeleteOne(ctx, bson.M{"_id": objID})
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    if result.DeletedCount == 0 {
        c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"success": true})
}

// GET /api/v1/product
func GetAllProducts(c *gin.Context) {
    var products []models.Product

    collection := db.GetCollection("products")
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    cursor, err := collection.Find(ctx, bson.M{})
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    defer cursor.Close(ctx)

    for cursor.Next(ctx) {
        var product models.Product
        cursor.Decode(&product)
        products = append(products, product)
    }

    c.JSON(http.StatusOK, products)
}

// GET /api/v1/product/category/:category
func GetProductsByCategory(c *gin.Context) {
    category := c.Param("category")
    page := c.Query("page")
    var pageNumber int
    if page == "" {
        pageNumber = 1
    } else {
        pageNumber, _ = strconv.Atoi(page)
    }
    skip := 10 * (pageNumber - 1)

    products := []models.Product{} 
    collection := db.GetCollection("products")
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    cursor, err := collection.Find(ctx, bson.M{"category": category}, options.Find().SetSkip(int64(skip)).SetLimit(10))
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    defer cursor.Close(ctx)

    for cursor.Next(ctx) {
        var product models.Product
        cursor.Decode(&product)
        products = append(products, product)
    }

    total, err := collection.CountDocuments(ctx, bson.M{"category": category})
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "products": products,
        "total":    total,
        "page":     pageNumber,
    })
}

// GET /api/v1/popular/:category
func GetPopularProducts(c *gin.Context) {
    category := c.Param("category")
    // find and limit 4 documents
    var products []models.Product
    collection := db.GetCollection("products")
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    cursor, err := collection.Find(ctx, bson.M{"category": category}, options.Find().SetLimit(4))
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    defer cursor.Close(ctx)

    for cursor.Next(ctx) {
        var product models.Product
        cursor.Decode(&product)
        products = append(products, product)
    }

    c.JSON(http.StatusOK, products)
}

// GET /api/v1/newcollections
func GetNewCollections(c *gin.Context) {
    var products []models.Product
    collection := db.GetCollection("products")
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    // find 8 most recent added document
    cursor, err := collection.Find(ctx, bson.M{}, options.Find().SetSort(bson.D{{"date", -1}}).SetLimit(8))
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    defer cursor.Close(ctx)

    for cursor.Next(ctx) {
        var product models.Product
        cursor.Decode(&product)
        products = append(products, product)
    }

    c.JSON(http.StatusOK, products)
}

// GET /api/v1/newcollections/:category

func GetNewProductsByCategory(c *gin.Context) {
    category := c.Param("category")
    var products []models.Product

    collection := db.GetCollection("products")
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    // find 4 most recent added document of a certain category
    cursor, err := collection.Find(ctx, bson.M{"category": category}, options.Find().SetSort(bson.D{{"date", -1}}).SetLimit(4))
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
    }
    defer cursor.Close(ctx)

    for cursor.Next(ctx) {
        var product models.Product
        cursor.Decode(&product)
        products = append(products, product)
    }

    c.JSON(http.StatusOK, products)
}