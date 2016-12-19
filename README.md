# goTDmd
分析项目中的 `TODO`/`FIXME`，生成 `TODOs.md` 文件，方便查阅。  

# 用法
```
# 下载
$ go get -v -u github.com/ohko/goTDmd
$ go install github.com/ohko/goTDmd

# 搜集当前目录下的所有TODO/FIXME，写到TODOs.md
$ cd {you path}
$ goTDmd

# 自定义目录和输出文件
$ goTDmd -p ./ -d TODOs.md
```