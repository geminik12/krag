package options

import "github.com/spf13/pflag"

// MinIOOptions 定义 MinIO 配置选项.
type MinIOOptions struct {
	Endpoint        string `json:"endpoint" mapstructure:"endpoint"`
	AccessKeyID     string `json:"access-key-id" mapstructure:"access-key-id"`
	SecretAccessKey string `json:"secret-access-key" mapstructure:"secret-access-key"`
	Bucket          string `json:"bucket" mapstructure:"bucket"`
	UseSSL          bool   `json:"use-ssl" mapstructure:"use-ssl"`
}

// NewMinIOOptions 创建带有默认值的 MinIOOptions 实例.
func NewMinIOOptions() *MinIOOptions {
	return &MinIOOptions{
		Endpoint:        "localhost:9000",
		AccessKeyID:     "minioadmin",
		SecretAccessKey: "minioadmin",
		Bucket:          "krag-files",
		UseSSL:          false,
	}
}

// Validate 验证选项是否合法.
func (o *MinIOOptions) Validate() []error {
	var errs []error
	return errs
}

// AddFlags 将选项绑定到命令行标志.
func (o *MinIOOptions) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&o.Endpoint, "minio.endpoint", o.Endpoint, "MinIO endpoint")
	fs.StringVar(&o.AccessKeyID, "minio.access-key-id", o.AccessKeyID, "MinIO access key ID")
	fs.StringVar(&o.SecretAccessKey, "minio.secret-access-key", o.SecretAccessKey, "MinIO secret access key")
	fs.StringVar(&o.Bucket, "minio.bucket", o.Bucket, "MinIO bucket name")
	fs.BoolVar(&o.UseSSL, "minio.use-ssl", o.UseSSL, "Use SSL for MinIO connection")
}
