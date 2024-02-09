package controller

import (
	"OpenTeens/model"
	"OpenTeens/utils"
	"bytes"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// AttachmentUploadHandler handles the upload of attachments.
func AttachmentUploadHandler(c *gin.Context) {
	// Parse multipart form
	err := c.Request.ParseMultipartForm(30 << 20) // 10 MB
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "解析表单数据失败"})
		return
	}

	// Get the file from the multipart form
	file, fileHeader, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "获取文件失败"})
		return
	}
	defer file.Close()

	// Now you can use fileHeader.Filename to get the original file name
	_, valid := determineFileType(fileHeader.Filename)
	if !valid {
		c.JSON(http.StatusBadRequest, gin.H{"error": "不支持的文件类型"})
		return
	}

	// Get userID from context
	userID, ok := c.Get("user")
	userID = userID.(uint)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无法获取用户ID"})
		return
	}

	// Read binary data from the request body
	data, err := io.ReadAll(file) // 修改这里，从file读取数据
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "读取文件数据失败"})
		return
	}
	defer c.Request.Body.Close()

	filePath := createTempFile(data)
	if filePath == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建临时文件失败"})
		return
	}

	//fileHeader := createFileHeader(filePath, data)
	ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
	// Process upload
	fileURL, fileHash, err := processUpload(fileHeader, filePath, userID.(uint), ext)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "文件上传失败", "details": err.Error()})
		return
	}
	fmt.Println(fileURL)

	c.JSON(http.StatusOK, gin.H{"message": "文件上传成功", "url": "Cant Show for Security", "fileHash": fileHash})
}

// processUpload processes the file upload.
func processUpload(fileHeader *multipart.FileHeader, filePath string, userID uint, fileExt string) (string, string, error) {
	// Upload file to server
	fileURL, fileHash, err := uploadFileToServer(filePath, fileExt, utils.GenerateToken()+fileExt, userID)
	if err != nil {
		return "", "", err
	}

	fileType, _ := determineFileType(fileURL)

	// Store attachment info in database
	attachment := model.Attachment{
		Type:       fileType,
		FileHash:   fileHash,
		Nickname:   fileHeader.Filename,
		Url:        fileURL,
		UploaderID: userID,
	}

	if err := model.DBAttachmentCreate(&attachment); err != nil {
		return "", "", err
	}

	return fileURL, fileHash, nil
}

// createTempFile creates a temporary file and returns its path.
func createTempFile(data []byte) string {
	uploadDir := "uploads"
	createDirectory(uploadDir)

	fileName := generateFileName()
	filePath := filepath.Join(uploadDir, fileName)

	if err := os.WriteFile(filePath, data, 0644); err != nil {
		fmt.Println("无法写入临时文件:", err)
		return ""
	}

	return filePath
}

// createFileHeader creates a multipart.FileHeader.
func createFileHeader(filePath string, data []byte) *multipart.FileHeader {
	return &multipart.FileHeader{
		Filename: filepath.Base(filePath),
		Size:     int64(len(data)),
	}
}

// determineFileType determines the file type based on its extension.
func determineFileType(fileName string) (string, bool) {
	//fmt.Println(fileName)
	ext := strings.ToLower(filepath.Ext(fileName))
	// Update with desired extensions
	allowedExtensions := map[string]string{
		".jpg":  "Image",
		".jpeg": "Image",
		".png":  "Image",
		".gif":  "Image",
		".bmp":  "Image",
		".zip":  "Zip",
		".rar":  "Zip",
		".tar":  "Zip",
		".gz":   "Zip",
		".pdf":  "PDF",
		// Add other extensions and types
	}

	fileType, ok := allowedExtensions[ext]
	return fileType, ok
}

// createDirectory creates a directory if it does not exist.
func createDirectory(dir string) {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		os.Mkdir(dir, 0755)
	}
}

// generateFileName generates a unique file name.
func generateFileName() string {
	timestamp := strconv.FormatInt(time.Now().UnixNano(), 10)
	hash := md5.Sum([]byte(timestamp))
	return fmt.Sprintf("%x", hash) + ".tmp"
}

// uploadFileToServer uploads a file to the server.
func uploadFileToServer(filePath, fileType, filename string, uid uint) (string, string, error) {
	fileURL := "http://119.28.78.54/upload.php" // PHP上传接口地址

	// 打开文件
	file, err := os.Open(filePath)
	if err != nil {
		return "", "", err
	}
	defer file.Close()

	// 创建表单文件
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", filename)
	if err != nil {
		return "", "", err
	}
	_, err = io.Copy(part, file)
	if err != nil {
		return "", "", err
	}
	writer.Close()

	// 发送POST请求
	request, err := http.NewRequest("POST", fileURL, body)
	if err != nil {
		return "", "", err
	}
	request.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return "", "", err
	}
	defer response.Body.Close()

	// 读取响应
	responseBytes, err := io.ReadAll(response.Body)
	if err != nil {
		return "", "", err
	}

	// 解析JSON响应
	var respData struct {
		Status   string `json:"status"`
		Message  string `json:"message"`
		FileHash string `json:"fileHash"`
		Url      string `json:"url"`
	}
	err = json.Unmarshal(responseBytes, &respData)
	if err != nil {
		return "", "", err
	}

	if respData.Status != "success" {
		return "", "", fmt.Errorf("error uploading file: %s", respData.Message)
	}

	return respData.Url, respData.FileHash, nil
}

func AttachmentGetHandler(c *gin.Context) {
	fileHash := c.Query("hash") // 从查询参数获取URL
	if fileHash == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Hash is required"})
		return
	}

	// 查询数据库
	attachment, err := model.DBAttachmentGetByHash(fileHash)
	url := attachment.Url
	typer := attachment.Type

	// 调用函数来获取文件内容
	content, err := fetchFileContent(url)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to fetch file content", "details": err.Error()})
		return
	}

	// 如果Type是Image，返回图片
	if typer == "Image" {
		c.Data(http.StatusOK, "image/jpeg", content)
		return
	} else {
		// 设置响应的Content-Type为"application/octet-stream"，表示是一个二进制文件
		c.Data(http.StatusOK, "application/octet-stream", content)
	}
}

// fetchFileContent takes a URL to a file and returns the content of the file.
func fetchFileContent(fileURL string) ([]byte, error) {
	resp, err := http.Get(fileURL)
	if err != nil {
		return nil, fmt.Errorf("error fetching file: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("server returned non-200 status code: %d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %v", err)
	}

	return body, nil
}
