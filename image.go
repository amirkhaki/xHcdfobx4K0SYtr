package main

import (
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"context"
	"net/http"
	"os"
	"io"
	"bytes"
	"net/textproto"
	"mime/multipart"
	"strings"
	"fmt"
	"path/filepath"

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







var quoteEscaper = strings.NewReplacer("\\", "\\\\", `"`, "\\\"")


func escapeQuotes(s string) string {
	return quoteEscaper.Replace(s)
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

func UploadFile(url string, filename string) (string, error) {
	file, err := os.Open(filename)

	if err != nil {
		return "", err
	}
	defer file.Close()
	defer os.Remove(file.Name())
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition",
	fmt.Sprintf(`form-data; name="%s"; filename="%s"`,
	escapeQuotes("file"), escapeQuotes(filepath.Base(file.Name())+".jpg")))
	h.Set("Content-Type", "image/jpeg")
	part, err := writer.CreatePart(h)

	if err != nil {
		return "", fmt.Errorf("Error during creting form file: %w", err)
	}

	io.Copy(part, file)
	writer.Close()
	request, err := http.NewRequest("POST", url, body)

	if err != nil {
		return "", fmt.Errorf("Error during creating request: %w", err)
	}
	request.Header.Add("Content-Type", writer.FormDataContentType())
	client := &http.Client{}

	response, err := client.Do(request)

	if err != nil {
		return "", fmt.Errorf("Error during sending request: %w", err)
	}
	defer response.Body.Close()

	content, err := io.ReadAll(response.Body)

	if err != nil {
		return "", fmt.Errorf("Error during reading response body: %w", err)
	}

	return string(content), nil
}
