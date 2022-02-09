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