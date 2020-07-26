git-merge
=====
一个 drone 插件

能力：在 source_branchs 中的分支推送时，将会自动通过 ssh 连接到 server，并进入 projects_path 中，将推送的分支自动合并到 target_branch。

## 开发计划
- [x] 自动合并
- [x] 测试代码
- [ ] 自动识别端口
- [ ] 自动钉钉通知

## 原理：
伪代码

```shell
if commit_branch in source_branchs
    connect server
        each path in projects_path
            cd path
            if current_branch == target_branch
                git fetch origin commit_branch:temp-commit_branch
                git merge temp-commit_branch
                git branch -D temp-commit_branch 
```


## 配置
```yaml
kind: pipeline
name: merge-barnch

clone:
    disable: true

steps:
    -   name: merge-barnch
        image: alextechs/drone-git-merge
        settings:
            server:
                host: host
                port: 22
                user: root
                password: "password"
            target_branch: develops
            source_branchs: 
                - release
                - develop
            projects_path:
                - /app
```

## 编译
```shell
GOOS=linux CGO_ENABLED=0 go build -o merge ./src/...
```

## 运行

export 这些都是drone在流水线中设置的变量，我们只是模拟

```shell
export CI_BUILD_NUMBER="61"
export CI_BUILD_STARTED="1599572953"
export CI_BUILD_STATUS="success"
export CI_COMMIT_AUTHOR_EMAIL="im@println.org"
export CI_COMMIT_AUTHOR_NAME="Alex"
export CI_COMMIT_BRANCH="branch"
export CI_COMMIT_MESSAGE="1"
export DRONE_REPO_LINK="https://host/"
export PLUGIN_PROJECTS_PATH='/app./go-app'
export PLUGIN_SERVER='{"host":"host","port" : "port","user":"user","password":"password"}'
export PLUGIN_SOURCE_BRANCHS="master,release,physical-server"
export PLUGIN_TARGET_BRANCH="develops"
./merge
```

## 测试
```shell
chmod +x test.sh
./test.sh
```
