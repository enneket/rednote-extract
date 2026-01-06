package config

// Loader 配置加载器
type Loader struct {
}

// NewLoader 创建配置加载器
func NewLoader() *Loader {
	return &Loader{}
}

// Load 加载配置
func (l *Loader) Load() (*Config, error) {
	return DefaultConfig(), nil
}
