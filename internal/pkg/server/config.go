package server

import (
	"github.com/gin-gonic/gin"
	"net"
	"strconv"
	"time"
)

// Config 用于配置genericServer的配置.
type Config struct {
	SecureServing   *SecureServingInfo
	InsecureServing *InsecureServingInfo
	JWT             *JWTInfo
	Mode            string
	Middlewares     []string

	Healthz         bool
	EnableProfiling bool
	EnableMetrics   bool
}

// CertKey 结构体用于存储证书和密钥文件的路径。
type CertKey struct {
	CertFile string // CertFile 字段指定了证书文件的路径。
	KeyFile  string // KeyFile 字段指定了密钥文件的路径。
}

// SecureServingInfo 进行TLS加密的服务器信息
type SecureServingInfo struct {
	BindAddress string // BindAddress 字段指定了监听地址。
	BindPort    int    // BindPort 字段指定了监听端口。
	CertKey     CertKey
}

// Address 返回监听地址和端口的组合字符串。
func (s SecureServingInfo) Address() string {
	return net.JoinHostPort(s.BindAddress, strconv.Itoa(s.BindPort))
}

// InsecureServingInfo 不加密的服务器信息
type InsecureServingInfo struct {
	BindAddress string
	BindPort    int
}

func (i InsecureServingInfo) Address() string {
	return net.JoinHostPort(i.BindAddress, strconv.Itoa(i.BindPort))
}

// JWTInfo JWT认证信息
type JWTInfo struct {
	Realm      string // Realm 字段指定了认证域。
	Key        string // Key 字段指定了密钥。
	Timeout    time.Duration
	MaxRefresh time.Duration
}

// NewNilConfig 返回一个空的Config对象。
func NewNilConfig() *Config {
	return &Config{
		SecureServing:   nil,
		InsecureServing: nil,
		JWT: &JWTInfo{
			Realm:      "iam-jwt",
			Key:        "",
			Timeout:    1 * time.Hour,
			MaxRefresh: 1 * time.Hour,
		},
		Mode:            gin.ReleaseMode,
		Middlewares:     make([]string, 0),
		Healthz:         false,
		EnableProfiling: false,
		EnableMetrics:   false,
	}
}

func (c *Config) Complete() CompletedConfig {
	return CompletedConfig{c}
}
