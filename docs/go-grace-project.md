# 优雅的Go项目
从以下几个方面入手：
+ 代码结构
  + 目录结构
  + 按功能拆分模块
+ 代码规范
  + 编码规范
  + 最佳实践
+ 代码质量
  + 编写可测试的代码
  + 高单元测试覆盖率
  + 代码审查
+ 编程哲学
  + 面向接口编程
  + 面向“对象”编程
+ 软件设计方法
  + 设计模式
  + 遵循SOLID原则

## 代码规范
[Uber Go语言编码规范](https://github.com/xxjwxc/uber_go_guide_cn?tab=readme-ov-file)

## 代码质量

使用Mock工具：
- golang/mock，是官方提供的Mock框架。它实现了基于interface的Mock功能，能够与Golang内置的testing包做很好的集成，是最常用的Mock工具。golang/mock提供了mockgen工具用来生成interface对应的Mock源文件。
- sqlmock，可以用来模拟数据库连接。数据库是项目中比较常见的依赖，在遇到数据库依赖时都可以用它。
- httpmock，可以用来Mock HTTP请求。
- bouk/monkey，猴子补丁，能够通过替换函数指针的方式来修改任意函数的实现。如果golang/mock、sqlmock和httpmock这几种方法都不能满足我们的需求，我们可以尝试通过猴子补丁的方式来Mock依赖。可以这么说，猴子补丁提供了单元测试 Mock 依赖的最终解决方案。

使用单元测试
- 使用gotests工具自动生成单元测试代码，减少编写单元测试用例的工作量，将你从重复的劳动中解放出来。
- 定期检查单元测试覆盖率。你可以通过以下方法来检查：

## 编程哲学

