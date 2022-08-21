package main

import (
	"context"
	"fmt"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"io"
	"net/http"
	"os"
)
func DeleteFromS3(client *minio.Client, ctx context.Context, fileName string) error {
	err := client.RemoveObject(ctx, os.Getenv("S3_BUCKET"), fileName, minio.RemoveObjectOptions{})
	if err != nil {
		return fmt.Errorf("Could not delete file : %w", err)
	}
	return nil
}

func UploadToS3(client *minio.Client, ctx context.Context, filePath, fileName, contentType string) error {
	userMetaData := map[string]string{"x-amz-acl": "public-read"}
	_, err := client.FPutObject(ctx, os.Getenv("S3_BUCKET"), fileName,
	filePath, minio.PutObjectOptions{ContentType: contentType, UserMetadata: userMetaData})
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
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("Error during creating request to %s: %w", url, err)
	}
	req.Header.Set("Host", "www.govdeals.com")
	req.Header.Set("Sec-Ch-Ua-Mobile", "?1")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Linux; Android 10; M2006C3LG) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/102.0.0.0 Mobile Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
	req.Header.Set("Sec-Fetch-Site", "none")
	req.Header.Set("Sec-Fetch-Mode", "navigate")
	req.Header.Set("Sec-Fetch-User", "?1")
	req.Header.Set("Sec-Fetch-Dest", "document")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9,fa;q=0.8")

	resp, err := http.DefaultClient.Do(req)
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
