#!/usr/bin/env bash

# 使用方法：
# ./genModel.sh usercenter user
# ./genModel.sh usercenter user_auth
# 再将./genModel下的文件剪切到对应服务的model目录里面，记得改package


# 数据库配置
host=127.0.0.1
port=3306
dbname=$1
username=root
passwd='0}{;Pp1W"@HRMxTOSNvhl$Yw'

#生成的表名
tables=$2
#表生成的genmodel目录
modeldir=$1/dao/model
pname=$3 #包名

echo "开始创建库：$dbname 的表：$3"
#goctl model mysql datasource -url="${username}:${passwd}@tcp(${host}:${port})/${dbname}" -table="${tables}"  -dir="${modeldir}" -cache=false --style=goZero
gentool -dsn "${username}:${passwd}@tcp(${host}:${port})/${dbname}?charset=utf8mb4&parseTime=true&loc=Local" -tables="${tables}" 
#gentool -dsn 'root:toor@tcp(localhost)/main?charset=utf8mb4&parseTime=true&loc=Local' -table="${tables}"  -dir="${modeldir}"

# 
go run main.go -dsn "${username}:${passwd}@tcp(${host}:${port})/${dbname}?charset=utf8mb4&parseTime=true&loc=Local" -path "${path}"./dao/query -pname "${pname}" \
-tables "${tables}"
