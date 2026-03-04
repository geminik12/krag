package options

import "github.com/spf13/pflag"

// QdrantOptions 定义 Qdrant 配置选项.
type QdrantOptions struct {
	URL            string `json:"url" mapstructure:"url"`
	APIKey         string `json:"api-key" mapstructure:"api-key"`
	CollectionName string `json:"collection-name" mapstructure:"collection-name"`
}

// NewQdrantOptions 创建带有默认值的 QdrantOptions 实例.
func NewQdrantOptions() *QdrantOptions {
	return &QdrantOptions{
		URL:            "http://localhost:6333",
		APIKey:         "",
		CollectionName: "krag_knowledge",
	}
}

// Validate 验证选项是否合法.
func (o *QdrantOptions) Validate() []error {
	var errs []error
	return errs
}

// AddFlags 将选项绑定到命令行标志.
func (o *QdrantOptions) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&o.URL, "qdrant.url", o.URL, "Qdrant server URL")
	fs.StringVar(&o.APIKey, "qdrant.api-key", o.APIKey, "Qdrant API Key")
	fs.StringVar(&o.CollectionName, "qdrant.collection-name", o.CollectionName, "Qdrant collection name")
}
