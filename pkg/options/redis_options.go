package options

import "github.com/spf13/pflag"

// RedisOptions 定义 Redis 配置选项.
type RedisOptions struct {
	Addr     string `json:"addr" mapstructure:"addr"`
	Password string `json:"password" mapstructure:"password"`
	DB       int    `json:"db" mapstructure:"db"`
}

// NewRedisOptions 创建带有默认值的 RedisOptions 实例.
func NewRedisOptions() *RedisOptions {
	return &RedisOptions{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	}
}

// Validate 验证选项是否合法.
func (o *RedisOptions) Validate() []error {
	var errs []error
	return errs
}

// AddFlags 将选项绑定到命令行标志.
func (o *RedisOptions) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&o.Addr, "redis.addr", o.Addr, "Redis server address")
	fs.StringVar(&o.Password, "redis.password", o.Password, "Redis password")
	fs.IntVar(&o.DB, "redis.db", o.DB, "Redis database number")
}
