# 快速开始
滴滴云MC主要是基于Minio Client (https://github.com/minio/mc)开发，支持Minio Client (mc)所有功能同时增加了超大对象检查功能。

## 如何申请滴滴云S3的Bucket？
先注册**滴滴云账号**，进入：https://app.didiyun.com/#/s3/add 申请Bucket，如下图：

![didyun s3 bucket](http://img-ys011.didistatic.com/static/doc/S3_bucket_01.png)

填写名称和设置访问权限，点立即创建即可。

## 如何申请AK和SK？
![didyun s3 AK SK](http://img-ys011.didistatic.com/static/doc/S3_AK_SK_01.png)

```
操作步骤：
（1）点击“API”按钮。
（2）选择“对象存储密钥”。
（3）点击“创建API密钥”。
即可得到的SecretID和SecretKey值.
```

## 如何配置滴滴云Minio Client？
获取到S3 API密钥后，得到了SecretID和SecretKey值，通过这两个值来配置滴滴云S3。

### 公共配置信息如下：

```
mc config host add didiyuns3 https://s3.didiyunapi.com AKDD002DYS7H379X1YQKZFSCGOFNX1 V7f1CwQqAcwo80UEIJEjc5gVQUSSx5ohQ9GSrr12

```

### DC2配置信息如下：

```
mc config host add didiyuns3 https://s3-internal.didiyunapi.com AKDD002DYS7H379X1YQKZFSCGOFNX1 V7f1CwQqAcwo80UEIJEjc5gVQUSSx5ohQ9GSrr12

```

配置成功后，在用户目录下.mc/config.json会生成新的配置信息。

```
{
	"version": "9",
	"hosts": {
		"didiyuns3": {
			"url": "https://s3-gz.didiyunapi.com",
			"accessKey": "AKDD002DYS7H379X1YQKZFSCGOFNX1",
			"secretKey": "V7f1CwQqAcwo80UEIJEjc5gVQUSSx5ohQ9GSrr12",
			"api": "s3v4",
			"lookup": "auto"
		}
	}
}
```

## 如何使用的Minio Client？
### 查询滴滴云S3上的所有bucket
```
➜  ~ mc ls didiyuns3
[2018-02-09 15:08:04 CST]     0B didiyun/
```
### 查询滴滴云S3上某bucket的文件列表
```
➜  ~ mc ls didiyuns3/didiyun
[2018-11-04 10:57:03 CST] 107KiB 6a6f178b009847dca.jpg
[2018-10-31 10:24:09 CST]    40B test
[2018-08-30 15:50:07 CST]  13MiB test.mp4
[2018-08-24 09:59:25 CST] 107KiB test_6a6f178b009847163649c7cb9s
[2018-12-10 17:49:36 CST]     0B test/
```

### 上传文件到滴滴云S3上
```
➜  ~ mc cp ./test1 didiyuns3/didiyun/
./test1:  40 B / 40 B ▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓  100.00% 296 B/s 0s
➜  ~ mc ls didiyuns3/didiyun
[2018-11-04 10:57:03 CST] 107KiB 6a6f178b009847163649c7cb96a9e4ca.jpg
[2018-11-13 17:56:44 CST] 3.1KiB das.graffle
[2018-10-31 10:24:09 CST]    40B test
[2018-08-30 15:50:07 CST]  13MiB test.mp4
[2018-12-10 17:52:30 CST]    40B test1
[2018-08-24 09:59:25 CST] 107KiB test_6a6f178b009847163649c7cb96a9e4ca
[2018-12-10 17:53:39 CST]     0B test/
```
使用MC CP上传成功后，再重新获取到列表就会多出test1文件。

### 下载滴滴云S3上的文件到本地
```
➜  ~ mc cp didiyuns3/didiyun/test1 ./
...gz.didiyunapi.com/didiyun/test1:  40 B / 40 B  ▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓  100.00% 109 B/s 0s
```
### 检查分片对象是否正确
```

➜ mc check didiyun/why-test/image_1543334400_1.log.tar.gz ./image_1543334400_1.log.tar.gz.part.minio
Right: localMd5[91ac858d19e436c7972931831efdf914] remoteMd5[91ac858d19e436c7972931831efdf914] partNumber[1]
Wrong: localMd5[633be061e11b37d209c6d61105b39dd1] remoteMd5[adcf84c9f80612e098b922a900686b85] partNumber[2]
Wrong: localMd5[d41d8cd98f00b204e9800998ecf8427e] remoteMd5[bb1e534133eaec67258c2b7eeb8d5b24] partNumber[3]
Wrong: localMd5[d41d8cd98f00b204e9800998ecf8427e] remoteMd5[181edd35eea4569b7a7e9700b572f892] partNumber[4]
Wrong: localMd5[d41d8cd98f00b204e9800998ecf8427e] remoteMd5[98e0246164ce67ce53bce9d2c77f589d] partNumber[5]
Wrong: localMd5[d41d8cd98f00b204e9800998ecf8427e] remoteMd5[a19106f58867dbfd60edd471d22faa58] partNumber[6]
Wrong: localMd5[d41d8cd98f00b204e9800998ecf8427e] remoteMd5[0892e7849f1c53dace7e4e4ad5cc1279] partNumber[7]
Wrong: localMd5[d41d8cd98f00b204e9800998ecf8427e] remoteMd5[98929c633e6bee560f5d11109c37d6cc] partNumber[8]

All Wrong Info
Wrong: PartNumber[2] LocalMD5[633be061e11b37d209c6d61105b39dd1] RemoteMd5[adcf84c9f80612e098b922a900686b85] Size[67108864] LastModified[2019-06-13 11:46:04 +0000 UTC]
Wrong: PartNumber[3] LocalMD5[d41d8cd98f00b204e9800998ecf8427e] RemoteMd5[bb1e534133eaec67258c2b7eeb8d5b24] Size[67108864] LastModified[2019-06-13 11:46:04 +0000 UTC]
Wrong: PartNumber[4] LocalMD5[d41d8cd98f00b204e9800998ecf8427e] RemoteMd5[181edd35eea4569b7a7e9700b572f892] Size[67108864] LastModified[2019-06-13 11:46:04 +0000 UTC]
Wrong: PartNumber[5] LocalMD5[d41d8cd98f00b204e9800998ecf8427e] RemoteMd5[98e0246164ce67ce53bce9d2c77f589d] Size[67108864] LastModified[2019-06-13 11:46:05 +0000 UTC]
Wrong: PartNumber[6] LocalMD5[d41d8cd98f00b204e9800998ecf8427e] RemoteMd5[a19106f58867dbfd60edd471d22faa58] Size[67108864] LastModified[2019-06-13 11:46:05 +0000 UTC]
Wrong: PartNumber[7] LocalMD5[d41d8cd98f00b204e9800998ecf8427e] RemoteMd5[0892e7849f1c53dace7e4e4ad5cc1279] Size[67108864] LastModified[2019-06-13 11:46:06 +0000 UTC]
Wrong: PartNumber[8] LocalMD5[d41d8cd98f00b204e9800998ecf8427e] RemoteMd5[98929c633e6bee560f5d11109c37d6cc] Size[59743711] LastModified[2019-06-13 11:46:06 +0000 UTC]

ALL Total: Right[1] Wrong[7]

```

### 修复分片对象
```
➜ mc check didiyun/why-test/image_1543334400_1.log.tar.gz ./image_1543334400_1.log.tar.gz.part.minio --repair
Right: localMd5[91ac858d19e436c7972931831efdf914] remoteMd5[91ac858d19e436c7972931831efdf914] partNumber[1]
Wrong: localMd5[633be061e11b37d209c6d61105b39dd1] remoteMd5[adcf84c9f80612e098b922a900686b85] partNumber[2]
Wrong: localMd5[d41d8cd98f00b204e9800998ecf8427e] remoteMd5[bb1e534133eaec67258c2b7eeb8d5b24] partNumber[3]
Wrong: localMd5[d41d8cd98f00b204e9800998ecf8427e] remoteMd5[181edd35eea4569b7a7e9700b572f892] partNumber[4]
Wrong: localMd5[d41d8cd98f00b204e9800998ecf8427e] remoteMd5[98e0246164ce67ce53bce9d2c77f589d] partNumber[5]
Wrong: localMd5[d41d8cd98f00b204e9800998ecf8427e] remoteMd5[a19106f58867dbfd60edd471d22faa58] partNumber[6]
Wrong: localMd5[d41d8cd98f00b204e9800998ecf8427e] remoteMd5[0892e7849f1c53dace7e4e4ad5cc1279] partNumber[7]
Wrong: localMd5[d41d8cd98f00b204e9800998ecf8427e] remoteMd5[98929c633e6bee560f5d11109c37d6cc] partNumber[8]

All Wrong Info
Wrong: PartNumber[2] LocalMD5[633be061e11b37d209c6d61105b39dd1] RemoteMd5[adcf84c9f80612e098b922a900686b85] Size[67108864] LastModified[2019-06-13 11:46:04 +0000 UTC]
Wrong: PartNumber[3] LocalMD5[d41d8cd98f00b204e9800998ecf8427e] RemoteMd5[bb1e534133eaec67258c2b7eeb8d5b24] Size[67108864] LastModified[2019-06-13 11:46:04 +0000 UTC]
Wrong: PartNumber[4] LocalMD5[d41d8cd98f00b204e9800998ecf8427e] RemoteMd5[181edd35eea4569b7a7e9700b572f892] Size[67108864] LastModified[2019-06-13 11:46:04 +0000 UTC]
Wrong: PartNumber[5] LocalMD5[d41d8cd98f00b204e9800998ecf8427e] RemoteMd5[98e0246164ce67ce53bce9d2c77f589d] Size[67108864] LastModified[2019-06-13 11:46:05 +0000 UTC]
Wrong: PartNumber[6] LocalMD5[d41d8cd98f00b204e9800998ecf8427e] RemoteMd5[a19106f58867dbfd60edd471d22faa58] Size[67108864] LastModified[2019-06-13 11:46:05 +0000 UTC]
Wrong: PartNumber[7] LocalMD5[d41d8cd98f00b204e9800998ecf8427e] RemoteMd5[0892e7849f1c53dace7e4e4ad5cc1279] Size[67108864] LastModified[2019-06-13 11:46:06 +0000 UTC]
Wrong: PartNumber[8] LocalMD5[d41d8cd98f00b204e9800998ecf8427e] RemoteMd5[98929c633e6bee560f5d11109c37d6cc] Size[59743711] LastModified[2019-06-13 11:46:06 +0000 UTC]

ALL Total: Right[1] Wrong[7]

Start repail size

Start repail part
Right: localMd5[adcf84c9f80612e098b922a900686b85] remoteMd5[adcf84c9f80612e098b922a900686b85] partNumber[2]
Repair succeed: partNumber[2]
^@Right: localMd5[bb1e534133eaec67258c2b7eeb8d5b24] remoteMd5[bb1e534133eaec67258c2b7eeb8d5b24] partNumber[3]
Repair succeed: partNumber[3]
^@Right: localMd5[181edd35eea4569b7a7e9700b572f892] remoteMd5[181edd35eea4569b7a7e9700b572f892] partNumber[4]
Repair succeed: partNumber[4]
^@Right: localMd5[98e0246164ce67ce53bce9d2c77f589d] remoteMd5[98e0246164ce67ce53bce9d2c77f589d] partNumber[5]
Repair succeed: partNumber[5]
^@Right: localMd5[a19106f58867dbfd60edd471d22faa58] remoteMd5[a19106f58867dbfd60edd471d22faa58] partNumber[6]
Repair succeed: partNumber[6]
^@Right: localMd5[0892e7849f1c53dace7e4e4ad5cc1279] remoteMd5[0892e7849f1c53dace7e4e4ad5cc1279] partNumber[7]
Repair succeed: partNumber[7]
Right: localMd5[98929c633e6bee560f5d11109c37d6cc] remoteMd5[98929c633e6bee560f5d11109c37d6cc] partNumber[8]
Repair succeed: partNumber[8]

Repair succedd!

```

### 更多指令请使用Help
```
➜  ~ mc help
NAME:
  mc - Minio Client for cloud storage and filesystems.

USAGE:
  mc [FLAGS] COMMAND [COMMAND FLAGS | -h] [ARGUMENTS...]

COMMANDS:
  ls       list buckets and objects
  mb       make a bucket
  cat      display object contents
  pipe     stream STDIN to an object
  share    generate URL for temporary access to an object
  cp       copy objects
  mirror   synchronize object(s) to a remote site
  find     search for objects
  sql      run sql queries on objects
  stat     show object metadata
  diff     list differences in object name, size, and date between two buckets
  rm       remove objects
  event    configure object notifications
  watch    listen for object notification events
  policy   manage anonymous access to buckets and objects
  admin    manage minio servers
  session  resume interrupted operations
  config   configure minio client
  update   update mc to latest release
  check    check big object content
  version  show version info

GLOBAL FLAGS:
  --config-dir value, -C value  path to configuration folder (default: "/Users/didi/.mc")
  --quiet, -q                   disable progress bar display
  --no-color                    disable color theme
  --json                        enable JSON formatted output
  --debug                       enable debug output
  --insecure                    disable SSL certificate verification
  --help, -h                    show help
  --version, -v                 print the version

VERSION:
  DEVELOPMENT.2019-07-12T03-55-00Z
```

