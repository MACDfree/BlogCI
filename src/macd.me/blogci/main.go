package main

import (
	"bytes"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

var logger *log.Logger

var (
	// 变量初始化
	blogGit    = os.Getenv("BLOG_GIT")
	publishGit = os.Getenv("PUBLISH_GIT")
	themeGit   = os.Getenv("THEME_GIT")
	//certificate = os.Getenv("CERTIFICATE")
	userName = os.Getenv("USER_NAME")
	password = os.Getenv("PASSWORD")
	host     = os.Getenv("HOST")
	email    = os.Getenv("EMAIL")
	baseURL  = os.Getenv("BASE_URL")
	title    = os.Getenv("TITLE")
	theme    = os.Getenv("THEME")
)

const (
	// BlogPath 博客所在路径
	BlogPath = "/opt/goblog/"
	// CredentialPath 凭证存储路径
	CredentialPath = "/opt/.my-key"
)

func main() {

	logfile, err := os.OpenFile("/opt/blogci/bin/blogci.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("%s\n", err.Error())
	}
	defer logfile.Close()
	logger = log.New(logfile, "", log.Ldate|log.Ltime|log.Llongfile)

	// 修改config内容
	bytes, err := ioutil.ReadFile(BlogPath + "config.yaml")
	if err != nil {
		log.Fatal("Open config file error:", err)
	}
	// 暂时做简单地替换
	config := string(bytes)
	config = strings.Replace(config, "||baseurl||", baseURL, -1)
	config = strings.Replace(config, "||title||", title, -1)
	config = strings.Replace(config, "||theme||", theme, -1)
	err = ioutil.WriteFile(BlogPath+"config.yaml", []byte(config), os.ModePerm)

	// 下载主题
	// 切换至主题文件夹
	os.Chdir(BlogPath + "themes")
	cmd := exec.Command("git", "clone", themeGit)
	b, err := cmd.Output()
	if err != nil {
		logger.Println("clone theme error: ", err)
		return
	}
	logger.Println(string(b))

	// 下载md文件
	os.Chdir(BlogPath + "content")
	cmd = exec.Command("git", "clone", blogGit, "./")
	b, err = cmd.Output()
	if err != nil {
		logger.Println("clone md error: ", err)
		return
	}
	logger.Println(string(b))

	// 下载静态页面
	os.Chdir(BlogPath + "publish")
	cmd = exec.Command("git", "clone", publishGit, "./")
	b, err = cmd.Output()
	if err != nil {
		logger.Println("clone md error: ", err)
		return
	}
	logger.Println(string(b))

	// 设置username
	cmd = exec.Command("git", "config", "user.name", userName)
	b, err = cmd.Output()
	if err != nil {
		logger.Println("git config username error: ", err)
		return
	}
	logger.Println(string(b))

	// 设置email
	cmd = exec.Command("git", "config", "user.email", email)
	b, err = cmd.Output()
	if err != nil {
		logger.Println("git config email error: ", err)
		return
	}
	logger.Println(string(b))

	// 修改publish仓库的gitconfig增加git认证
	// 进入publish目录
	os.Chdir(BlogPath + "publish")
	// 执行git配置凭据命令
	cmd = exec.Command("git", "config", "credential.helper", "store", "--file", CredentialPath)
	b, err = cmd.Output()
	if err != nil {
		logger.Println("config error: ", err)
		return
	}
	logger.Println(string(b))
	// 填充用户名密码
	cmd = exec.Command("/bin/sh", "-c", "echo -e 'protocol=https\nhost="+host+"\nusername="+userName+"\npassword="+password+"\n\n' | git credential-store --file "+CredentialPath+" store")
	b, err = cmd.Output()
	if err != nil {
		logger.Println("credential-store error: ", err)
		return
	}
	logger.Println(string(b))

	http.HandleFunc("/githooks", gitHooks)
	err = http.ListenAndServe(":8080", nil)

	if err != nil {
		log.Fatal("ListenAndServer: ", err)
	}
}

// gitHooks 用于处理github的webhook请求
// 需要设置blog的webhooks，监听push事件
func gitHooks(w http.ResponseWriter, r *http.Request) {
	// 其实只要进了这个方法就可以认为是触发了push事件
	// 此处简单验证一下是否是github发起的请求

	// 接下来需要做的事情：
	// 1. 更新blog本地仓库
	cmd := exec.Command("git", "pull", "origin", "master")
	os.Chdir(BlogPath + "content")
	b, err := cmd.Output()
	if err != nil {
		logger.Println("execute cmd [git pull origin master] error: ", err)
		io.WriteString(w, "execute cmd [git pull origin master] error")
		return
	}
	logger.Println(string(b))

	// 2. 调用hugo命令生成静态页面
	os.Chdir(BlogPath)
	cmd = exec.Command("hugo")
	b, err = cmd.Output()
	if err != nil {
		logger.Println("execute cmd [hugo] error: ", err)
		io.WriteString(w, "execute cmd [hugo] error")
		return
	}
	logger.Println(string(b))

	// 3. push静态页面至github
	// 3.1 添加文件
	os.Chdir(BlogPath + "public")
	//cmd = exec.Command("git", "add", ".")
	cmd = exec.Command("/bin/sh", "-c", "git add .")
	b, err = cmd.Output()
	if err != nil {
		logger.Println("execute cmd [git add .] error: ", err)
		io.WriteString(w, "execute cmd [git add .] error")
		return
	}
	logger.Println(string(b))

	// 3.2 本地提交
	os.Chdir(BlogPath + "public")
	//cmd = exec.Command("git", "commit", "-m", "\"blogci automatic commit\"")
	cmd = exec.Command("/bin/sh", "-c", "git commit -m \"blogci automatic commit\"")
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err = cmd.Run()
	if err != nil {
		logger.Println("execute cmd [git commit] error: ", string(stderr.Bytes()))
		io.WriteString(w, "execute cmd [git commit] error")
		return
	}
	logger.Println(string(stdout.Bytes()))

	// 3.3 推送github
	os.Chdir(BlogPath + "public")
	cmd = exec.Command("/bin/sh", "-c", "git push origin master")
	b, err = cmd.Output()
	if err != nil {
		logger.Println("execute cmd [git push origin master] error: ", err.Error())
		io.WriteString(w, "execute cmd [git push origin master] error")
		return
	}
	logger.Println(string(b))

	io.WriteString(w, "OK")
}
