# 概述
利用GitHub提供的WebHooks功能，在blog提交时触发事件请求服务器地址，服务器接收到请求后更新服务器的blog仓库，然后调用hugo进行静态页面的生成，最后将生成好的静态页面提交至GitHub，实现博客的自动发布功能，即博客的持续集成（BlogCI）。

# 前期准备
1. 拥有自己的VPS，并且安装有Git、Hugo
2. 在GitHub上建立博客Markdown文件仓库、静态页面仓库
3. GitHub中需要配置服务器的公钥，以实现免登陆

# 安装服务
可以在`blogci.sh`脚本基础上进行修改，目前只测试过CentOS

# 下一步计划
做成Docker镜像，方便部署

# Docker相关环境变量
* `BLOG_GIT`

    博客MD文件所在git仓库

* `PUBLISH_GIT`

    需要发布的静态页面git仓库

* `THEME_GIT`

    主题git仓库地址

* `CERTIFICATE`

    用户名密码证书串

* `BASE_URL`

    博客根路径

* `TITLE`

    博客标题

* `THEME`

    主题名称