package app

import "github.com/spf13/pflag"

// FlagsIntf 提供命令行接口，定义命令行的具体实现
// 该接口的实现由子配置结构体实现
// 例子：
//
//	type Options struct {
//	    MySQLOpts *MySQLOptions `json:"mysql" mapstructure:"mysql"`
//	}
//	type MySQLOptions struct {
//		...
//	}
//	var _ Flags = &MySQLOptions{} // 实现了 FlagsIntf
//	... 这里省略实现函数
type FlagsIntf interface {
	AddFlags(fs *pflag.FlagSet)
	Validate() []error
}

// FlagOptions 命令行读取配置
type FlagsOptions interface {
	// Flags 添加命令行
	Flags() (fs FlagSet)
	// Validate 验证
	Validate() []error
}

// ConfigurableOptions 抽象用于从配置文件读取参数的配置选项。
type ConfigurableOptions interface {
	ApplyFlags() []error
}

// CompletableOptions  抽象可以完成/编译的options
type CompletableOptions interface {
	Complete() error
}

// PrintableOptions 抽象可以打印的options
type PrintableOptions interface {
	String() string
}

// OptionsIntf 提供Options接口，定义Options的具体实现
type OptionsIntf interface {
	FlagsOptions
	ConfigurableOptions
	CompletableOptions
	PrintableOptions
}
