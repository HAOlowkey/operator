# operator

##1、初始化项目
```bash
go mod init github.com/HAOlowkey/operator
```

##2、kubebuilder初始化
```bash
kubebuilder init --domain haolowkey.github.io
```

##3、kubebuilder创建api
```bash
kubebuilder create api --group application --version v1 --kind Unit
```