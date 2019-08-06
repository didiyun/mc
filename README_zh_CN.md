# 快速入门指南
滴滴云MC是基于Minio客户端开发,具体Minio客户端所有的功能并支持超大文件检查与修复功能。Minio Client (mc)为ls，cat，cp，mirror，diff，find等UNIX命令提供了一种替代方案。


```
ls       列出文件和文件夹。
mb       创建一个存储桶或一个文件夹。
cat      显示文件和对象内容。
pipe     将一个STDIN重定向到一个对象或者文件或者STDOUT。
share    生成用于共享的URL。
cp       拷贝文件和对象。
mirror   给存储桶和文件夹做镜像。
find     基于参数查找文件。
diff     对两个文件夹或者存储桶比较差异。
rm       删除文件和对象。
events   管理对象通知。
watch    监听文件和对象的事件。
policy   管理访问策略。
session  为cp命令管理保存的会话。
config   管理mc配置文件。
update   检查软件更新。
check    检查与修复超大文件。
version  输出版本信息。
```
然后使用[`mc config`命令](#add-a-cloud-storage-service)。

## macOS
### Homebrew
使用[Homebrew](http://brew.sh/)安装mc。

```sh
brew install didiyun/stable/mc
mc --help
```

## 添加一个云存储服务
如果你打算仅在POSIX兼容文件系统中使用`mc`,那你可以直接略过本节，跳到[日常使用](#everyday-use)。

添加一个或多个S3兼容的服务，请参考下面说明。`mc`将所有的配置信息都存储在``~/.mc/config.json``文件中。

```sh
mc config host add <ALIAS> <YOUR-S3-ENDPOINT> <YOUR-ACCESS-KEY> <YOUR-SECRET-KEY> <API-SIGNATURE>
```

别名就是给你的云存储服务起了一个短点的外号。S3 endpoint,access key和secret key是你的云存储服务提供的。API签名是可选参数，默认情况下，它被设置为"S3v4"。

### 示例-Minio云存储
从Minio服务获得URL、access key和secret key。

```sh
mc config host add minio http://192.168.1.51 BKIKJAA5BMMU2RHO6IBB V7f1CwQqAcwo80UEIJEjc5gVQUSSx5ohQ9GSrr12 S3v4
```

### 示例-Amazon S3云存储
参考[AWS Credentials指南](http://docs.aws.amazon.com/general/latest/gr/aws-security-credentials.html)获取你的AccessKeyID和SecretAccessKey。

```sh
mc config host add s3 https://s3.amazonaws.com BKIKJAA5BMMU2RHO6IBB V7f1CwQqAcwo80UEIJEjc5gVQUSSx5ohQ9GSrr12 S3v4
```

### 示例-Google云存储
参考[Google Credentials Guide](https://cloud.google.com/storage/docs/migrating?hl=en#keys)获取你的AccessKeyID和SecretAccessKey。

```sh
mc config host add gcs  https://storage.googleapis.com BKIKJAA5BMMU2RHO6IBB V8f1CwQqAcwo80UEIJEjc5gVQUSSx5ohQ9GSrr12 S3v2
```

注意：Google云存储只支持旧版签名版本V2，所以你需要选择S3v2。

## 验证
`mc`预先配置了云存储服务URL：https://play.minio.io:9000，别名“play”。它是一个用于研发和测试的Minio服务。如果想测试Amazon S3,你可以将“play”替换为“s3”。

*示例:*

列出https://play.minio.io:9000上的所有存储桶。

```sh
mc ls play
[2016-03-22 19:47:48 PDT]     0B my-bucketname/
[2016-03-22 22:01:07 PDT]     0B mytestbucket/
[2016-03-22 20:04:39 PDT]     0B mybucketname/
[2016-01-28 17:23:11 PST]     0B newbucket/
[2016-03-20 09:08:36 PDT]     0B s3git-test/
```
<a name="everyday-use"></a>
## 日常使用

### Shell别名
你可以添加shell别名来覆盖默认的Unix工具命令。

```sh
alias ls='mc ls'
alias cp='mc cp'
alias cat='mc cat'
alias mkdir='mc mb'
alias pipe='mc pipe'
alias find='mc find'
```

### Shell自动补全
你也可以下载[`autocomplete/bash_autocomplete`](https://raw.githubusercontent.com/minio/mc/master/autocomplete/bash_autocomplete)到`/etc/bash_completion.d/`，然后将其重命名为`mc`。别忘了在这个文件运行source命令让其在你的当前shell上可用。

```sh
sudo wget https://raw.githubusercontent.com/minio/mc/master/autocomplete/bash_autocomplete -O /etc/bash_completion.d/mc
source /etc/bash_completion.d/mc
```

```sh
mc <TAB>
admin    config   diff     ls       mirror   policy   session  update   watch
cat      cp       events   mb       pipe     rm       share    version
```

## 了解更多
- [Minio Client完全指南](https://docs.minio.io/docs/minio-client-complete-guide)
- [Minio快速入门](https://docs.minio.io/docs/minio-quickstart-guide)
- [Minio官方文档](https://docs.minio.io)

## 贡献
请遵守Minio[贡献者指南](https://github.com/minio/mc/blob/master/docs/zh_CN/CONTRIBUTING.md)
