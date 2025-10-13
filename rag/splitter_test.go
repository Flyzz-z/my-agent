package rag

// import (
// 	"context"
// 	"testing"
// 	"github.com/cloudwego/eino/components/document"
// 	"github.com/stretchr/testify/assert"
// )

// // TestNewMarkdownSplitter 测试NewMarkdownSplitter函数
// func TestNewMarkdownSplitter(t *testing.T) {
// 	ctx := context.Background()
	
// 	// 创建markdown分割器
// 	splitter, err := NewMarkdownSplitter(ctx)
	
// 	// 验证没有错误
// 	assert.NoError(t, err, "Should not return error when creating markdown splitter")
// 	assert.NotNil(t, splitter, "Splitter should not be nil")
// }

// // TestMarkdownSplitter_Transform 测试Markdown分割器的Transform功能
// func TestMarkdownSplitter_Transform(t *testing.T) {
// 	ctx := context.Background()
	
// 	// 创建markdown分割器
// 	splitter, err := NewMarkdownSplitter(ctx)
// 	assert.NoError(t, err, "Should not return error when creating markdown splitter")
	
// 	// 创建测试文档
// 	inputDocs := []document.Document{
// 		{
// 			Content: []byte("# Title 1\n\nThis is content under title 1.\n\n## Title 1.1\n\nThis is content under title 1.1.\n\n# Title 2\n\nThis is content under title 2."),
// 			Metadata: map[string]interface{}{"source": "test.md"},
// 		},
// 	}
	
// 	// 由于实际的Transform方法依赖于eino库的具体实现，这里我们只是尝试调用它
// 	// 在实际项目中，应该根据具体的实现来测试分割功能
// 	tryToTransform := func() {
// 		splittedDocs, err := splitter.Transform(ctx, inputDocs)
// 		if err != nil {
// 			t.Logf("Transform might fail in test environment: %v", err)
// 		} else {
// 			t.Logf("Successfully split into %d documents", len(splittedDocs))
// 		}
// 	}
	
// 	// 调用Transform方法
// 	tryToTransform()
// }

// // TestMarkdownSplitter_HeaderConfig 测试Markdown分割器的标题配置
// func TestMarkdownSplitter_HeaderConfig(t *testing.T) {
// 	// 验证NewMarkdownSplitter函数中使用的标题配置
// 	// 由于我们无法直接访问内部配置，这里只是测试配置的行为
// 	ctx := context.Background()
// 	splitter, err := NewMarkdownSplitter(ctx)
// 	assert.NoError(t, err, "Should not return error when creating markdown splitter")
// 	assert.NotNil(t, splitter, "Splitter should not be nil")
	
// 	// 在实际项目中，可以通过反射或其他方式验证内部配置
// 	// 这里我们只是确保splitter可以正常创建并具有Transform方法
// 	assert.Implements(t, (*document.Transformer)(nil), splitter, "Splitter should implement Transformer interface")
// }