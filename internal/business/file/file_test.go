package file

import (
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"testing"
)

const baseURL = "http://localhost:8081"

func TestFileUploadAndDownload(t *testing.T) {
	content := []byte("hello gim file test")
	fileName := "test.txt"

	// 上传
	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)
	fw, err := w.CreateFormFile("file", fileName)
	if err != nil {
		t.Fatalf("create form file: %v", err)
	}
	if _, err = fw.Write(content); err != nil {
		t.Fatalf("write form file: %v", err)
	}
	w.Close()

	resp, err := http.Post(baseURL+"/upload", w.FormDataContentType(), &buf)
	if err != nil {
		t.Fatalf("upload request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("upload status = %d, want 200", resp.StatusCode)
	}

	var result Response
	if err = json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("decode upload response: %v", err)
	}
	if result.Code != 0 {
		t.Fatalf("upload code = %d, want 0, message: %s", result.Code, result.Message)
	}

	data, ok := result.Data.(map[string]interface{})
	if !ok {
		t.Fatalf("unexpected data type: %T", result.Data)
	}
	fileURL, ok := data["url"].(string)
	if !ok || fileURL == "" {
		t.Fatalf("url not found in response data")
	}
	t.Logf("uploaded url: %s", fileURL)

	// 下载
	dlResp, err := http.Get(fileURL)
	if err != nil {
		t.Fatalf("download request: %v", err)
	}
	defer dlResp.Body.Close()

	if dlResp.StatusCode != http.StatusOK {
		t.Fatalf("download status = %d, want 200", dlResp.StatusCode)
	}

	got, err := io.ReadAll(dlResp.Body)
	if err != nil {
		t.Fatalf("read download body: %v", err)
	}
	if !bytes.Equal(got, content) {
		t.Fatalf("downloaded content = %q, want %q", got, content)
	}
}
