# 3pScan

## Features

 - Simple tool for port scan

## Installation

```sh
go install -v github.com/tangxiaofeng7/3pScan@latest
```

## Usage
```console

Simple:
3pScan -h baidu.com -top 1000

Scan with disable ping:
3pScan -h baidu.com -Pn -top 1000

Scan with read file:
3pScan -hf url.txt -top 1000

Scan with full ports:
3pScan -h baidu.com -top full
or
3pScan -h baidu.com -p 1-65535

Scan with show stats:
3pScan -h baidu.com -p 1-65535 -t 10

Scan with other tools:

httpx:
3pScan -h baidu.com -top 1000 ｜ httpx -silent -title -status-code

nuclei:
3pScan -h baidu.com -top 1000 ｜ nuclei

subfinder:
echo baidu.com | subfinder -silent | 3pScan -top 1000

```



