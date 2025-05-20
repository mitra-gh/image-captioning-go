package handlers

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UploadHandler struct{}

func (u *UploadHandler) Upload(c *gin.Context) {
	file, err := c.FormFile("image")
	if err != nil {
		log.Println("Upload error:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid image upload"})
		return
	}

	// Save uploaded file with unique name
	id := strings.Replace(uuid.New().String(), "-", "", -1)
	ext := strings.ToLower(strings.Split(file.Filename, ".")[1])
	imageName := fmt.Sprintf("%s.%s", id, ext)
	savePath := fmt.Sprintf("./images/%s", imageName)

	if err := c.SaveUploadedFile(file, savePath); err != nil {
		log.Println("Save error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save image"})
		return
	}

	// Describe the image using OpenAI
	description, err := describeImage(savePath)
	if err != nil {
		log.Println("OpenAI error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to analyze image"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"image":       imageName,
		"description": description,
	})
}

func describeImage(imagePath string) (string, error) {
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		return "", fmt.Errorf("OPENAI_API_KEY environment variable not set")
	}

	imageBytes, err := os.ReadFile(imagePath)
	if err != nil {
		return "", err
	}
	base64Image := base64.StdEncoding.EncodeToString(imageBytes)

	payload := map[string]interface{}{
		"model": "gpt-4o",
		"messages": []map[string]interface{}{
			{
				"role": "user",
				"content": []map[string]interface{}{
					{
						"type": "text",
						"text": "Describe the objects in this image. give me in persian language.",
					},
					{
						"type": "image_url",
						"image_url": map[string]interface{}{
							"url":    fmt.Sprintf("data:image/jpeg;base64,%s", base64Image),
							"detail": "high",
						},
					},
				},
			},
		},
		"max_tokens": 300,
	}

	bodyBytes, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(bodyBytes))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var result map[string]interface{}
	if err := json.Unmarshal(respBody, &result); err != nil {
		return "", err
	}

	// Extract response text
	if choices, ok := result["choices"].([]interface{}); ok && len(choices) > 0 {
		first := choices[0].(map[string]interface{})
		if message, ok := first["message"].(map[string]interface{}); ok {
			if content, ok := message["content"].(string); ok {
				return content, nil
			}
		}
	}

	return "", fmt.Errorf("unexpected API response format")
}
