package rag

import (
	"context"
	"testing"

	"github.com/cloudwego/eino/schema"
	"github.com/stretchr/testify/assert"
)

/*
TestMarkdownSplitter_Transform 测试Markdown分割器的Transform功能
*/
// ... existing code ...

/*
TestMarkdownSplitter_Transform 测试Markdown分割器的Transform功能
*/
func TestMarkdownSplitter_Transform(t *testing.T) {
	// 创建上下文
	ctx := context.Background()

	// 创建Markdown分割器
	splitter, err := NewMarkdownSplitter(ctx)
	assert.NoError(t, err, "Should not return error when creating Markdown splitter")
	assert.NotNil(t, splitter, "Splitter should not be nil")

	// 准备测试用的Markdown文档
	testDoc := &schema.Document{
		Content: `# 标题一
这是标题一的内容

## 标题二
这是标题二的内容

### 标题三
这是标题三的内容

## 标题二-2
这是标题二-2的内容

# 标题一-2
这是标题一-2的内容`,
		MetaData: map[string]any{"source": "test.md"},
	}

	docs := []*schema.Document{testDoc}

	// 调用分割器的Transform方法
	resultDocs, err := splitter.Transform(ctx, docs)
	assert.NoError(t, err, "Should not return error when transforming documents")
	assert.NotNil(t, resultDocs, "Result documents should not be nil")

	// 验证分割结果是否符合预期
	// 预期应该分割成5个文档：标题一、标题二、标题三、标题二-2、标题一-2
	assert.Equal(t, 5, len(resultDocs), "Should split into 5 documents")

	// 验证每个分割后的文档内容是否正确
	expectedContents := []string{
		"# 标题一\n这是标题一的内容",
		"## 标题二\n这是标题二的内容",
		"### 标题三\n这是标题三的内容",
		"## 标题二-2\n这是标题二-2的内容",
		"# 标题一-2\n这是标题一-2的内容",
	}

	for _, doc := range resultDocs {
		assert.Contains(t, expectedContents, doc.Content, "Document content should match one of the expected contents")
		// 验证元数据是否正确保留
		assert.Equal(t, "test.md", doc.MetaData["source"], "Metadata should be preserved")
	}
}
