package config

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
)

type Target struct {
	Name        string `yaml:"name"`
	URL         string `yaml:"url"`
	AllowPublic bool   `yaml:"allow_public"`
}

type TargetConfig struct {
	Targets []Target `yaml:"targets"`
}

// 从指定yaml文件加载检测目标
func LoadTargetsFromFile(path string) (*TargetConfig, error) {
	//	1. 读取配置文件内容
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("读取配置文件失败: %w，请在config目录下创建targets.yaml", err)
	}

	// 2. 反序列化 YAML 到结构体
	var cfg TargetConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("解析 YAML 失败: %w", err)
	}

	// 3. 基础校验防止空配置
	if len(cfg.Targets) == 0 {
		return nil, fmt.Errorf("配置文件中未定义任何 targets")
	}

	return &cfg, nil
}
