# Git Flow工作流
定义了五个分支

| 分支名     | 描述                                                                                                                 |
|---------|--------------------------------------------------------------------------------------------------------------------|
| master  | 该分支上的代码为发布状态                                                                                                       |
| develop | 开发中的最新代码，该分支只做合并操作，不能直接在该分支上开发                                                                                     |
| feature | 研发阶段做功能开发，基于develop分支新建，分支命名建议规则为feature/xxx-xxx                                                                   |
| release | 在发布阶段作版本发布的预发布分支，基于develop分支创建，分支命名建议为release/xxx-xxx;通过测试后，将release分支合并到master分支，并打上`tag`的版本标签，最后删除对应版本的release分支 |
| hotfix  | 在维护阶段作紧急bug修复分支，在master分支上创建，修复完合并到master和develop分支。分支命名建议规则为hotfix/xxx-xxx                                        |

## 流程示例
```text
# 创建一个常用分支
git checkout -b develop master 

# 基于develop创建feature分支
git checkout -b feature/print-hello-world develop 

# 在新分支feature/print-hello-world上进行开发

@@ 如果需要进行紧急bug修复，步骤如下：
@@ git stash # 临时保存修改至堆栈区
@@ git checkout -b hotfix/print-error master # 从master创建hotfix分支
@@ vi xxxerror.go # 修复bug
@@ git commit -am 'fix print message error bug' # 提交修复
@@ git checkout develop # 切换到develop分支
@@ git merge --no-ff hotfix/print-error # 把hotfix分支合并到develop分支
@@ git checkout master # 切换到master分支
@@ git merge --no-ff hotfix/print-error # 把hotfix分支合并到master分支
@@ git tag -a vxx.xx.x(版本号) -m 'fix log bug' # master分支打tag
@@ go build -v . # 编译代码，将二进制文件更新到生产环境
@@ git branch -d hotfix/print-error # 删除hotfix分支
@@ git checkout feature/print-hello-world # 切换到开发分支下
@@ git merge --no-ff develop # 拉取develop最新代码
@@ git stash pop # 恢复到修复bug前的工作状态
@@ 至此紧急bug的修复步骤结束

# 继续开发feature/print-hello-world分支

# 提交代码
git commit -am "print hello world"

# 将代码推送到代码托管平台
git push origin feature/print-hello-world
## 在GitHub上基于feature/print-hello-world创建PR，
## 创建PR后进行code review

## 通过后由代码仓库matainer将功能分支合并到develop分支上
git checkout develop
git merge --no-ff feature/print-hello-world

## 基于develop创建分支release，测试代码
git checkout -b release/1.0.0 develop
git build -v . # 构建二进制文件

### 如果测试失败，直接在release/1.0.0分支修改并编译部署

## 测试通过后将新功能分支合并到master和develop分支
git checkout develop
git merge --no-ff release/1.0.0
git checkout master
git merge --no-ff release/1.0.0
git tag -a v1.0.0 -m "add print hello world" # master 分支打tag

# 删除feature/print-hello-world分支，也可以删除release/1.0.0分支
git branch -d feature/print-hello-world
```

至此整个Git Flow流程结束，Git Flow工作流比较**适合开发团队相对固定，规模较大的项目**。