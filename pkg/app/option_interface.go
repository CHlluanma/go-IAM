package cli

// FlagInterface 命令行读取配置
type FlagIntf interface {
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
