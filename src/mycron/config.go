package mycron

import (
    "git.oschina.net/wida/goUdpLog/src/conf"
)

const (
    ConfigPath = "./etc/mycron.conf"
)

var(
    Mysql_host,Mysql_user,Mysql_pwd,Mysql_dbname,Log_Path ,Log_filename string
    Mysql_prot int
)

func init() {
    c := conf.NewConfig(ConfigPath)
    //mysql
    Mysql_host =c.Read("mysql", "host")
    Mysql_prot = c.GetInt("mysql", "port")
    Mysql_user = c.Read("mysql", "user")
    Mysql_pwd = c.Read("mysql", "pwd")
    Mysql_dbname =c.Read("mysql", "dbname")

    //logger
    Log_Path = c.Read("log","logPath")
    Log_filename = c.Read("log","filename")
}

