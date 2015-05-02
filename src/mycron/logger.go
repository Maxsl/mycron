package mycron

import (
    "github.com/donnie4w/go-logger/logger"
)

func init() {
    logger.SetConsole(true)
    logger.SetRollingDaily(Log_Path, Log_filename)
    logger.SetLevel(logger.ERROR)
}


func Log(v ...interface{}) {
    logger.Error(v);
}