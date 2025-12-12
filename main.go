package main

import (
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log"
	"net/http"
	"encoding/json"
)

// Struct to hold the Terraform state
type TerraformState struct {
	Version            int         `json:"version"`
	TerraformVersion   string      `json:"terraform_version"`
	Serial             int         `json:"serial"`
	Lineage            string      `json:"lineage"`
	Resources          []Resource  `json:"resources"`
	Modules            []Module    `json:"modules"`
	Outputs            interface{} `json:"outputs"`
	Variables          interface{} `json:"variables"`
	DataSources        interface{} `json:"data_sources"`
}

type Resource struct {
	Type     string   `json:"type"`
	Name     string   `json:"name"`
	Instances []Instance `json:"instances"`
}

type Instance struct {
	Attributes map[string]interface{} `json:"attributes"`
}

type Module struct {
	Path     []string `json:"path"`
	Resources []Resource `json:"resources"`
}

func UploadStateFile(c *gin.Context) {
	// Retrieve the uploaded file from the form data
	file, _, err := c.Request.FormFile("file")
	if err != nil {
		log.Printf("Error while retrieving file: %s", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Failed to read the file: %s", err.Error())})
		return
	}
	defer file.Close()

	// Read the file content into memory
	data, err := ioutil.ReadAll(file)
	if err != nil {
		log.Printf("Error while reading file content: %s", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to read file content: %s", err.Error())})
		return
	}

	// Parse the file content into a TerraformState struct
	var tfState TerraformState
	if err := json.Unmarshal(data, &tfState); err != nil {
		log.Printf("Error while parsing file content: %s", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to parse JSON: %s", err.Error())})
		return
	}

	// Log file details for debugging
	log.Printf("Received file with size: %d bytes", len(data))
	log.Printf("File content preview (first 100 characters): %s", string(data[:100]))

	// Return the parsed data back to the frontend
	c.JSON(http.StatusOK, tfState)
}

func main() {
	// Create a new Gin router
	r := gin.Default()

	// Enable CORS to allow requests from your frontend (e.g., React on localhost:3000)
	r.Use(cors.Default())

	// Define the /upload POST endpoint
	r.POST("/upload", UploadStateFile)

	// Optionally, set the maximum file upload size (10 MB in this case)
	r.MaxMultipartMemory = 10 << 20 // 10 MB

	// Start the Gin server on port 8080
	log.Println("Starting server on :8080...")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

