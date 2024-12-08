package image_processing

import (
	"context"
	"fmt"
	"image"
	"image/jpeg"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"cloud.google.com/go/storage"
	"github.com/nfnt/resize"
	"google.golang.org/api/option"
)

// DownloadImage downloads an image from a URL and saves it to the specified destination
func DownloadImage(url, destination string) error {
    response, err := http.Get(url)
    if err != nil {
        return err
    }
    defer response.Body.Close()

    err = os.MkdirAll(filepath.Dir(destination), os.ModePerm)
    if err != nil {
        return err
    }

    file, err := os.Create(destination)
    if err != nil {
        return err
    }
    defer file.Close()

    _, err = io.Copy(file, response.Body)
    if err != nil {
        return err
    }
    return nil
}

// CompressImage compresses and resizes an image
func CompressImage(source, destination string, maxWidth, maxHeight uint) error {
    err := os.MkdirAll(filepath.Dir(destination), os.ModePerm)
    if err != nil {
        return err
    }

    file, err := os.Open(source)
    if err != nil {
        return err
    }
    defer file.Close()

    img, _, err := image.Decode(file)
    if err != nil {
        return err
    }

    resizedImg := resize.Resize(maxWidth, maxHeight, img, resize.Lanczos3)

    out, err := os.Create(destination)
    if err != nil {
        return err
    }
    defer out.Close()

    err = jpeg.Encode(out, resizedImg, nil)
    if err != nil {
        return err
    }

    return nil
}

// UploadToFirebase uploads a file to Firebase Storage and returns the URL
func UploadToFirebase(filePath, bucketName, objectName string) (string, error) {
	ctx := context.Background()
	client, err := storage.NewClient(ctx, option.WithCredentialsFile("./serviceAccountKey.json"))
	if err != nil {
			return "", err
	}
	defer client.Close()

	bucket := client.Bucket(bucketName)
	object := bucket.Object(objectName)
	writer := object.NewWriter(ctx)
	writer.ACL = []storage.ACLRule{{Entity: storage.AllUsers, Role: storage.RoleReader}}
	defer writer.Close()

	file, err := os.Open(filePath)
	if err != nil {
			return "", err
	}
	defer file.Close()

	if _, err := io.Copy(writer, file); err != nil {
			return "", err
	}

	if err := writer.Close(); err != nil {
			return "", err
	}

	url := fmt.Sprintf("https://storage.googleapis.com/%s/%s", bucketName, objectName)
	return url, nil
}

// ProcessImage downloads, compresses, and uploads an image, then returns the URL
func ProcessImage(url, downloadPath, compressedPath, bucketName, objectName string) (string, error) {
    err := DownloadImage(url, downloadPath)
    if err != nil {
        return "", err
    }

    err = CompressImage(downloadPath, compressedPath, 800, 600)
    if err != nil {
        return "", err
    }

    
    objectName = "image_processing/" + objectName

    uploadedURL, err := UploadToFirebase(compressedPath, bucketName, objectName)
    if err != nil {
        return "", err
    }

    return uploadedURL, nil
}