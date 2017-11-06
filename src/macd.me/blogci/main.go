package main

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
)

var logger *log.Logger

func main() {

	logfile, err := os.OpenFile("/home/macd/blogci/bin/blogci.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)

	if err != nil {
		log.Fatalf("%s\n", err.Error())
	}

	defer logfile.Close()

	logger = log.New(logfile, "", log.Ldate|log.Ltime|log.Llongfile)

	http.HandleFunc("/githooks", gitHooks)
	err = http.ListenAndServe(":8888", nil)

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
	os.Chdir("/home/macd/goblog/content")
	b, err := cmd.Output()
	if err != nil {
		logger.Println("execute cmd [git pull origin master] error: ", err)
		io.WriteString(w, "execute cmd [git pull origin master] error")
		return
	}
	logger.Println(string(b))

	// 2. 调用hugo命令生成静态页面
	os.Chdir("/home/macd/goblog")
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
	os.Chdir("/home/macd/goblog/public")
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
	os.Chdir("/home/macd/goblog/public")
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
	os.Chdir("/home/macd/goblog/public")
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
