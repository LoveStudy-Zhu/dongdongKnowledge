# dongdongKnowledge
go



### 部署主从服务器：

主数据库：
docker run -p 3339:3306 --name mysql-master -e MYSQL_ROOT_PASSWORD=1216 -d mysql:5.7
进入主数据库：
docker exec -it mysql-master /bin/bash

创建用户
create user 'user_w'@'%' IDENTIFIED BY '1216';
赋予读取数据库权限
GRANT ALL PRIVILEGES ON * .* TO 'user_w'@'%';

ALTER USER 'user_w'@'%' IDENTIFIED WITH mysql_native_password BY '1216';

程序外连接:mysql -h127.0.0.1 -P3339 -uuser_w -p

配置主数据库：(/etc/mysql/ my.conf中添加)
[mysqld]
###### 同一局域网内注意要唯一
server-id=100  
###### 开启二进制日志功能，可以随便取（关键）
log-bin=mysql-bin

配置完成后先重启数据库服务，再重启容器


从数据库：docker run -p 3340:3306 --name mysql-slave -e MYSQL_ROOT_PASSWORD=1216 -d mysql:5.7
同上
(/etc/mysql/ my.conf中添加)
[mysqld]
###### 设置server_id,注意要唯一
server-id=101  
###### 开启二进制日志功能，以备Slave作为其它Slave的Master时使用
log-bin=mysql-slave-bin 
relay_log=edu-mysql-relay-bin

查看镜像ip：
docker inspect -f '{{.Name}} - {{.NetworkSettings.IPAddress}}' $(docker ps -aq)
/mysql-master - 172.17.0.4
/mysql-slave - 172.17.0.3

master中创建同步用户
create user 'slave'@'%' IDENTIFIED BY '1216';
user mysql;
授权：GRANT REPLICATION SLAVE,REPLICATION CLIENT ON *.* TO 'slave'@'%';
测试：mysql -h172.17.0.4 -P3306 -uslave -p1216


进入主库中：
show master status
得到Position值
622

进入从数据库库同步：
change master to master_host='172.17.0.4',master_user='slave',master_password='1216',master_port=3306,master_log_file='mysql-bin.000001',master_log_pos=622,master_connect_retry=30;


change master to master_host='172.17.0.4', master_user='slave', master_password='1216', master_port=3306, master_log_file='mysql-bin.000001', master_log_pos=781, master_connect_retry=30;

开启主从复制：start slave;






最后测试，可能会碰到问题
Slave_SQL_Running: No

mysql> stop slave ;
mysql> set GLOBAL SQL_SLAVE_SKIP_COUNTER=1;
mysql> start slave ;