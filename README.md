# 3pScan

## 工具目标

· perfect

· precise

· Painless


## 使用方法
```
# 简单扫描
./3pScan -h baidu.com -top 1000

# 禁用ping扫描
./3pScan -h baidu.com -Pn -top 1000

# 文件读取方式扫描
./3pScan -hf url.txt -top 1000

# 终端传入扫描
echo 127.0.0.1 | ./3pScan
cat url.txt | ./3pScan

# 全端口扫描
./3pScan -h baidu.com -top full
或
./3pScan -h baidu.com -p 1-65535

# 端口扫描进度提示，默认为5秒
./3pScan -h baidu.com -p 1-65535 -t 10

```



