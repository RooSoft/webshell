# webshell
webshell golang implementation

run http
```
go get github.com/luoqeng/webshell

go build

./webshell -pass="mypass" -addr=":9090"

curl -X POST -d '{"pass": "mypass", "cmd": "bash", "opt": "-c", "args": "ls -l ~; echo hello"}' http://localhost:9090
total 20
drwxrwxr-x 10 luoqeng luoqeng 4096 Jul  6 17:25 dev
drwxrwxr-x  3 luoqeng luoqeng 4096 Jul  9 15:01 download
drwxrwxr-x  5 luoqeng luoqeng 4096 Jun 25 18:05 go
-rw-rw-r--  1 luoqeng luoqeng 5230 Jun 13 21:00 spf13-vim.sh
hello

```
