package main

import (
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"context"
	"net/http"
	"os"
	"fmt"
	"io"

)


func UploadToS3(client *minio.Client, ctx context.Context,filePath, fileName, contentType string) error {
	userMetaData := map[string]string{"x-amz-acl": "public-read"}
	_, err := client.FPutObject(ctx,os.Getenv("S3_BUCKET"), fileName, 
		filePath, minio.PutObjectOptions{ContentType:contentType, UserMetadata: userMetaData})
		if err != nil {
			return fmt.Errorf("Error during uploading %s: %w", filePath, err)
		}
		return nil
}
func GetObjectUrl(fileName string) string {
	uri := "https://" + os.Getenv("S3_BUCKET") + "." + os.Getenv("S3_ENDPOINT")
	uri += "/" + fileName
	return uri
}
func GetS3Client() (*minio.Client, error) {
	endpoint := os.Getenv("S3_ENDPOINT")
	accessKeyID := os.Getenv("S3_ACCESS")
	secretAccessKey := os.Getenv("S3_SECRET")
	useSSL := true
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("err creating s3 client: %w", err)
	}
	return minioClient, nil
}

func DownloadFile(url, ext string) (*os.File, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("Error during getting file %s: %w", url, err)
	}
	defer resp.Body.Close()
	f, err := os.CreateTemp("", "tmpimg*."+ext)
	if err != nil {
		return nil, fmt.Errorf("Error during creating temp file: %w", err)
	}
	defer f.Close()
	_, err = io.Copy(f, resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Error during writing to file: %w", err)
	}
	return f, nil
}

