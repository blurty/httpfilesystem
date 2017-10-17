# httpfilesystem
go版本的http文件服务器与其客户端

## 用法

### 服务端

简单地如下操作即可：

    -> go build httpfileserver.go
    -> ./httpfileserver &

### 客户端

首先编译程序：

    -> go build httpfileclient.go
    -> ./httpfileclient
    -> missing server ip
    -> Usage:./httpfileclient -h 127.0.0.1 <pattern> <filename>
    -> 		-h server ip
    -> 		-u upload file
    -> 		-d download file
    -> 		-q result of file transfer

看到输出提示可知程序必须提供一个http文件服务器地址，端口是固定的12345。然后还需要选择一个功能。

上传文件：

    -> ./httpfileclient -h 127.0.0.1 -u test.txt

下载文件：

    -> ./httpfileclient -h 127.0.0.1 -d test.txt

查询文件上传结果：

    -> ./httpfileclient -h 127.0.0.1 -q test.txt