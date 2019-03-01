# mysql-to-markdown
mysql 数据库详细信息转化成markdown格式
## 使用方法
### 1.0 安装
    程序使用前需要安装扩展库
    go get -u github.com/go-sql-driver/mysql
    和yaml解析的扩展库
    go get gopkg.in/yaml.v2
    
    开始使用：
    go get github.com/shanlongpan/mysql-to-markdown
    cd $GOPATH/src/shanlongpan/mysql-to-markdown
    go build explainMysql.go
    修改  test.yaml 里面对应的数据库信息
    ./explainMysql test.yaml
  
### 生成的文件格式为 mydatabese.md (对应的数据库名字)

### 数据库：cms 详细信息
#### 1.1表信息
| Name| Engine | Rows | Create_time| Comment|
| ---- | ---- | ---- |---- |---- |
|permission\_role|InnoDB|16|2019-02-28|角色与权限对应关系表|
|permissions|InnoDB|22|2019-02-28|权限表|
|role\_user|InnoDB|5|2019-02-28|管理员与角色对应关系表|
|roles|InnoDB|13|2019-02-28||
|users|InnoDB|3|2019-02-28||
#### 2.0 表 permission\_role 详细信息
| TABLE_SCHEMA| TABLE_NAME | COLUMN_NAME | COLUMN_DEFAULT| COLUMN_TYPE|COLUMN_KEY|COLUMN_COMMENT|
| ---- | ---- | ---- |---- |---- |----|---- | 
|cms|permission\_role|permission\_id||int(10) unsigned|PRI|权限ID|
|cms|permission\_role|role\_id||int(11)|PRI|角色ID|
#### 2.1 表 permissions 详细信息
| TABLE_SCHEMA| TABLE_NAME | COLUMN_NAME | COLUMN_DEFAULT| COLUMN_TYPE|COLUMN_KEY|COLUMN_COMMENT|
| ---- | ---- | ---- |---- |---- |----|---- | 
|cms|permissions|id||int(11)|PRI||
|cms|permissions|parent\_id|0|int(11)||父id|
|cms|permissions|name||varchar(255)||权限名称|
|cms|permissions|display\_name||varchar(255)||显示名称|
|cms|permissions|description||varchar(255)||备注|
|cms|permissions|menu\_nav|0|tinyint(1)||是否是菜单项 1 是项0不是|
|cms|permissions|created\_at||timestamp||创建时间|
|cms|permissions|updated\_at||timestamp||更新时间|
|cms|permissions|icon||varchar(50)|||
|cms|permissions|sequence|0|tinyint(4)||由大到小排序|
|cms|permissions|is\_active|0|tinyint(4)||0使用，1删除|
