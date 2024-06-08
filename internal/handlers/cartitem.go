package handlers

import (
	"backend-go/db"
	"backend-go/internal/utils"
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

// POST /api/v1/cartitems
func GetCartItems(c *gin.Context) {
	token := c.GetHeader("auth-token")
	user, err := utils.GetUserByToken(token)
	if err != nil || user == nil {
		c.JSON(http.StatusNotFound, gin.H{"error":"User not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"cartData": user.CartData})
}

// PATCH /api/v1/cartitems
func UpdateCartItem(c *gin.Context) {
    var req struct {
        ItemID string `json:"itemID"`
    }
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	token := c.GetHeader("auth-token")
	user, err := utils.GetUserByToken(token)
	if err != nil || user == nil {
		c.JSON(http.StatusNotFound, gin.H{"error":"User not found"})
		return
	}
	cartData := user.CartData
	if cartData == nil {
		cartData = make(map[string]int)
	}
	cartData[req.ItemID]++
	collection := db.GetCollection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()
	_, err = collection.UpdateOne(ctx, bson.M{"_id": user.ID}, bson.M{"$set": bson.M{"cartData": cartData}})
	if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update cart"})
        return
    }
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Item added to cart"})
}

// DELETE /api/v1/cartitems/:id
func DeleteCartItem(c *gin.Context) {
	id := c.Param("id")
	token := c.GetHeader("auth-token")
    user, err := utils.GetUserByToken(token)
    if err != nil || user == nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
        return
    }
	cartData := user.CartData
    if cartData == nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Item not found in cart"})
        return
    }
	cartData[id]--
    if cartData[id] <= 0 {
        delete(cartData, id)
    }
	collection := db.GetCollection("users")
	_, err = collection.UpdateOne(context.Background(), bson.M{"_id": user.ID}, bson.M{"$set": bson.M{"cartData": cartData}})
	if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove item from cart"})
        return
    }
    c.JSON(http.StatusOK, gin.H{"success": true, "message": "Item removed from cart"})
}