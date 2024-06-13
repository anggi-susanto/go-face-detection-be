package storage

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/gofiber/fiber/v2"
)

func SavePhoto(c *fiber.Ctx) (string, error) {
	// Retrieve the file from the form
	file, err := c.FormFile("photo")
	if err != nil {
		return "", err
	}

	// Generate a unique timestamped file name
	timestamp := time.Now().UnixNano()
	filePath := filepath.Join(os.Getenv("PHOTO_STORAGE_PATH"), fmt.Sprintf("%d-%s", timestamp, file.Filename))

	// Ensure the directory exists
	err = os.MkdirAll(filepath.Dir(filePath), os.ModePerm)
	if err != nil {
		return "", err
	}

	// Save the file to the specified path
	err = c.SaveFile(file, filePath)
	if err != nil {
		return "", err
	}

	return filePath, nil
}
