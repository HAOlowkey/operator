###1.创建项目目录
```
mkdir mysql-operator
cd mysql-operator
```

###2.初始化项目
```
kubebuilder init --domain bsgchina.com --repo gitlab.bsgchina.com/dbscale-kube/mysql-operator
```

###3.添加mysql api
```
kubebuilder create api --group dbscale --version v1 --kind MysqlCluster
```

###4.更新代码和manifest
```
make && make manifests
```

###5.安装code-generator
```
go get k8s.io/code-generator@0.23.1
go mod vendor
```