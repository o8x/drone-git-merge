package main

import (
    "fmt"
    log "github.com/sirupsen/logrus"
    "github.com/tidwall/gjson"
    "github.com/urfave/cli"
    "os"
    "strings"
)

func main() {
    app := cli.NewApp()
    app.Flags = []cli.Flag{
        cli.StringFlag{Name: "build.url", EnvVar: "CI_BUILD_LINK"},
        cli.StringFlag{Name: "build.number", EnvVar: "CI_BUILD_NUMBER"},
        cli.StringFlag{Name: "build.started", EnvVar: "CI_BUILD_STARTED"},
        cli.StringFlag{Name: "build.status", EnvVar: "CI_BUILD_STATUS"},
        cli.StringFlag{Name: "commit.author_email", EnvVar: "CI_COMMIT_AUTHOR_EMAIL"},
        cli.StringFlag{Name: "commit.author_name", EnvVar: "CI_COMMIT_AUTHOR_NAME"},
        cli.StringFlag{Name: "commit.branch", EnvVar: "CI_COMMIT_BRANCH"},
        cli.StringFlag{Name: "commit.message", EnvVar: "CI_COMMIT_MESSAGE"},
        cli.StringFlag{Name: "drone.link", EnvVar: "DRONE_REPO_LINK"},

        // 需要进行自动合并的分支
        cli.StringFlag{Name: "config.source_branchs", EnvVar: "PLUGIN_SOURCE_BRANCHS"},
        // 目标分支
        cli.StringFlag{Name: "config.target_branch", EnvVar: "PLUGIN_TARGET_BRANCH"},
        // ssh 服务器参数 json
        cli.StringFlag{Name: "config.ssh_server", EnvVar: "PLUGIN_SERVER"},
        cli.StringFlag{Name: "config.ssh_server.server_password", EnvVar: "PLUGIN_SERVER_PASSWORD"},
        // 需要更新的项目
        cli.StringFlag{Name: "config.projects_path", EnvVar: "PLUGIN_PROJECTS_PATH"},
    }

    app.Action = run
    err := app.Run(os.Args)
    if err != nil {
        log.Fatal(err)
    }
}

var result string
var err error

func run(c *cli.Context) error {
    build_url := c.String("build.url")
    build_number := c.String("build.number")
    build_started := c.String("build.started")
    build_status := c.String("build.status")
    commit_author_email := c.String("commit.author_email")
    commit_author_name := c.String("commit.author_name")
    commit_message := c.String("commit.message")
    drone_link := c.String("drone.link")

    commit_branch := c.String("commit.branch")
    source_branchs := strings.Split(c.String("config.source_branchs"), ",")
    target_branch := c.String("config.target_branch")
    ssh_server := c.String("config.ssh_server")
    projects_path := c.String("config.projects_path")
    server_password := c.String("config.ssh_server.server_password")

    log.Info("build_url : ", build_url)
    log.Info("build_number : ", build_number)
    log.Info("build_started : ", build_started)
    log.Info("build_status : ", build_status)
    log.Info("commit_author_email : ", commit_author_email)
    log.Info("commit_author_name : ", commit_author_name)
    log.Info("commit_message : ", commit_message)
    log.Info("drone_link : ", drone_link)
    log.Info("ssh_server : ", ssh_server)
    log.Info("projects_path : ", projects_path)
    log.Info("----------------------------------")
    log.Info("source_branchs : ", source_branchs)
    log.Info("commit_branch : ", commit_branch)
    log.Info("target_branch : ", target_branch)
    log.Info("server_password : ", server_password)

    if inArray(commit_branch, source_branchs) == false {
        return nil
    }

    // 密文不存在，使用明文
    if server_password == "" {
        server_password = gjson.Get(ssh_server, "password").String()
    }

    ssh := NewSSHClient(
        gjson.Get(ssh_server, "host").String(),
        gjson.Get(ssh_server, "user").String(),
        server_password,
        gjson.Get(ssh_server, "port").Int(),
    )

    for _, path := range strings.Split(projects_path, ",") {
        log.Info("current directory：", path)

        shell := "cd " + path + " && git symbolic-ref --short -q HEAD"
        if result, err = ssh.Run(shell); err != nil {
            panic(err)
        }

        // 当前分支不是 target，不进行处理
        if strings.Trim(strings.TrimSpace(result), "\n") != target_branch {
            continue
        }

        command := fmt.Sprintf(`
cd %s
git branch -D temp-%s || true
git fetch origin %s:temp-%s && git merge temp-%s && git branch -D temp-%s && git pull origin %s 
`,
            path,
            commit_branch,
            commit_branch,
            commit_branch,
            commit_branch,
            commit_branch,
            target_branch,
        )

        if result, err = ssh.Run(command); err != nil {
            panic(err)
        }

        log.Info("result", result)
    }

    return nil
}
