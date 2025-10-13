package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

// LoadJSONFile 从文件中加载JSON数据
func LoadJSONFile(filePath string, v interface{}) error {
	// 读取文件内容
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("读取文件失败: %v", err)
	}

	// 解析JSON数据
	if err := json.Unmarshal(content, v); err != nil {
		return fmt.Errorf("解析JSON数据失败: %v", err)
	}

	return nil
}

// SaveJSONFile 将JSON数据保存到文件
func SaveJSONFile(filePath string, v interface{}) error {
	// 确保目录存在
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("创建目录失败: %v", err)
	}

	// 序列化JSON数据
	content, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return fmt.Errorf("序列化JSON数据失败: %v", err)
	}

	// 写入文件
	if err := ioutil.WriteFile(filePath, content, 0644); err != nil {
		return fmt.Errorf("写入文件失败: %v", err)
	}

	return nil
}

// SplitTextIntoChunks 将文本分割成多个块
func SplitTextIntoChunks(text string, chunkSize int) []string {
	if chunkSize <= 0 {
		chunkSize = 1000 // 默认块大小
	}

	var chunks []string
	for i := 0; i < len(text); i += chunkSize {
		end := i + chunkSize
		if end > len(text) {
			end = len(text)
		}
		chunks = append(chunks, text[i:end])
	}

	return chunks
}

// CleanText 清理文本中的特殊字符
func CleanText(text string) string {
	// 这里可以根据需要添加更多的文本清理逻辑
	// 例如去除多余的空白字符、特殊符号等

	return text
}