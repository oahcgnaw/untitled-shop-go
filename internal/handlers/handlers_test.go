package handlers

import (
	"backend-go/internal/models"
	"context"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
)

func TestGetProduct(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Initialize mock MongoDB
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("success", func(mt *mtest.T) {
		productID := primitive.NewObjectID()

		// Insert a test product into the mock database
		mt.AddMockResponses(mtest.CreateCursorResponse(1, "products.products", mtest.FirstBatch, bson.D{
			{Key: "_id", Value: productID},
			{Key: "name", Value: "Test Product"},
		}))

		// Initialize router
		r := gin.Default()
		r.GET("/api/v1/product/:id", func(c *gin.Context) {
			id := c.Param("id")
			objID, _ := primitive.ObjectIDFromHex(id)
			var product models.Product

			collection := mt.Coll
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()

			err := collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&product)
			if err != nil {
				c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
				return
			}

			c.JSON(http.StatusOK, product)
		})

		// Create a request to send to the endpoint
		req, _ := http.NewRequest(http.MethodGet, "/api/v1/product/"+productID.Hex(), nil)
		rr := httptest.NewRecorder()

		// Perform the request
		r.ServeHTTP(rr, req)

		// Assert the response
		assert.Equal(t, http.StatusOK, rr.Code)

		var response models.Product
		err := json.Unmarshal(rr.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Test Product", response.Name)
	})
}
