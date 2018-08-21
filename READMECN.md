# SWAPIDC_NO
添加随机用户和工单

### 编译方法(你可以直接从Release里面下载)
* go get github.com/json-iterator/go
* go build
* ./SwapNo -log -rate 5 -url https://www.site.com/index/register/

### 功能
* -log                保存输出
* -loglimit           记录文件限制大小
* -proxy              启动代理模式
* -proxyupdate        从Github更新代理列表
* -rate               刷新频率(秒)
* -url                目标注册路径(工单模式下为目标登录路径)
* -saveusers          保存注册成功的用户(只是POST提交成功的，不保证一定注册上了)
* -debug              显示Post和Get的网页输出
* -overclock          启动超频模式(使用毫秒计算刷新频率)
* -tickets            启用工单模式
* -ticketurl          目标工单的提交地址
* -ticketprocess      提交工单的线程数量

### 示例
```./SwapNo -log -proxy -rate 2 -url https://site.com/index/register/```

### 保存用户模式
```./SwapNo -log -proxy -rate 2 -url https://site.com/index/register/ -saveusers```

### 超频模式
```./SwapNo -log -proxy -rate 20 -url https://site.com/index/register/ -overclock```

### 更新代理
```./SwapNo -proxyupdate```

### 工单模式
* 1. 在 users.txt 中添加用户，格式如下:
* ({"username": "xxxxx", "password": "123456"})
* 2. 运行脚本
```./SwapNo -log -rate 2 -url http://site.com/index/login/ -proxy -tickets -ticketurl http://site.com/ticket/submit/```
