package mycron

import (
    "github.com/widaT/go-logger/logger"
)

func init() {
    logger.SetConsole(true)
    //指定日志文件备份方式为日期的方式
    //第一个参数为日志文件存放目录
    //第二个参数为日志文件命名
    //logger.SetRollingDaily(Log_Path, Log_filename)
    //指定日志文件备份方式为文件大小的方式
    //第一个参数为日志文件存放目录
    //第二个参数为日志文件命名
    //第三个参数为备份文件最大数量
    //第四个参数为备份文件大小
    //第五个参数为文件大小的单位 KB，MB，GB TB
    logger.SetRollingFile(Log_Path,Log_filename, 10, 100, logger.MB)
    logger.SetLevel(logger.INFO)
}

func Log(v ...interface{}) {
    logger.Error(v);
}

func Info (v ...interface{} ){
    logger.Info(v);
}