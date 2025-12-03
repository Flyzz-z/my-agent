package rag

import (
	"context"
	"log"

	"github.com/cloudwego/eino-ext/components/document/transformer/splitter/markdown"
	"github.com/cloudwego/eino/components/document"
)

/*
NewMarkdownSplitter 创建 markdown 分割器
作用: 对 markdown 文档进行标题分割
*/
func NewMarkdownSplitter(ctx context.Context) (document.Transformer, error) {
	splitter, err := markdown.NewHeaderSplitter(ctx, &markdown.HeaderConfig{
		Headers: map[string]string{
			"#": "h1",
			"##": "h2",
			"###": "h3",
		},
		TrimHeaders: false,
	})
	if err != nil {
		log.Printf("NewMarkdownSplitter err: %v", err)
		return nil, err
	}
	return splitter, nil
}
