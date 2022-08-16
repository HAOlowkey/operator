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

###4、通过code-generate生成clientset
```bash
/Users/haolowkey/Dev/golang/src/github.com/kubernetes/code-generator/generate-groups.sh all github.com/HAOlowkey/operator/generated github.com/HAOlowkey/operator/api application.haolowkey.github.io:v1 --go-header-file=./hack/boilerplate.go.txt --output-base  ../../../
```