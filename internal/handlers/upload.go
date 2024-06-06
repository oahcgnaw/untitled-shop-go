package handlers

import (
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
)

func UploadFile(c *gin.Context) {
	file, err := c.FormFile("product")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// Create the filename with a timestamp
	filename := filepath.Join("upload/images", "product_" + time.Now().Format("20060102150405") + filepath.Ext(file.Filename))
	// Save the file
	if err := c.SaveUploadedFile(file, filename);  err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to save file"})
		return
	}
	var imageURL string
	if os.Getenv("SERVE_MODE") == "dev" {
		imageURL = "http://localhost:" + os.Getenv("PORT") + "/" + filename
	} else {
		imageURL = "https://ushop.cws-project.site/" + filename
	}

	c.JSON(http.StatusOK, gin.H{
		"success":true,
		"imageURL": imageURL,
	})


}