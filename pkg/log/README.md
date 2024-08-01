
# IAM 日志库

参考

- [江湖十年：使用go第三方日志库zap](https://jianghushinian.cn/2023/03/19/use-of-zap-in-go-third-party-log-library/)
- [江湖十年：如何优雅地封装一个更友好的日志库](https://jianghushinian.cn/2023/04/16/how-to-wrap-a-more-user-friendly-logging-package-based-on-zap/]https://jianghushinian.cn/2023/04/16/how-to-wrap-a-more-user-friendly-logging-package-based-on-zap/)
- [errors 包](https://pkg.go.dev/errors)
- [zap](https://pkg.go.dev/go.uber.org/zap)
- [marmotedu/log](https://github.com/marmotedu/iam/pkg/log)

## 日志处理
> 基础功能
+ 支持基础的日志信息：
+ 支持不同的日志级别
+ 支持自定义配置
+ 支持输出到标准输出和保存到文件
> 高级功能
+ 支持多种日志格式
+ 支持按级别分类输出
+ 支持结构化日志
+ 支持日志轮转
+ 具备Hook能力
> 额外可选功能
+ 支持颜色输出
+ 兼容标准库log
+ 支持输出到不同位置
