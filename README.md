
#**golang  任务调度系统**

1. 从mysql读取cron配置,开始任务和结束任务完全配置话

2. 采用crontab表达式支持到秒级

3. 任务运行状态全部透明化

4. 支持*nux 和 windows  脚本入口命令分别是 shell -c  和  cmd /c 

##Tables

###任务配置
    CREATE TABLE `cron` (
      `id` int(11) NOT NULL AUTO_INCREMENT,
      `uid` int(11) NOT NULL COMMENT '用户id',
      `name` varchar(50) NOT NULL DEFAULT '',
      `time` varchar(50) DEFAULT NULL,
      `cmd` varchar(255) DEFAULT NULL,
      `sTime` int(11) NOT NULL,
      `eTime` int(11) NOT NULL,
      `status` tinyint(1) NOT NULL DEFAULT '0',
      `isrunning` tinyint(1) NOT NULL DEFAULT '0' COMMENT '是否在运行中 0不是,1是',
      `modify` tinyint(1) NOT NULL DEFAULT '0',
      `process` tinyint(2) NOT NULL DEFAULT '1' COMMENT '进程数量',
      `ip` varchar(20) CHARACTER SET utf8 COLLATE utf8_bin NOT NULL DEFAULT '',
      `singleton` tinyint(1) NOT NULL DEFAULT '0' COMMENT '是否单例执行0非单例，1单例',
      PRIMARY KEY (`id`)
    ) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8

###任务执行记录

    CREATE TABLE `cron_hist` (
      `id` int(11) NOT NULL AUTO_INCREMENT,
      `cId` int(11) NOT NULL,
      `step` tinyint(3) NOT NULL,
      `time` datetime NOT NULL,
      `ret` varchar(255) DEFAULT NULL,
      PRIMARY KEY (`id`)
    ) ENGINE=InnoD式B  DEFAULT CHARSET=utf8;

###示例:
    INSERT INTO `cron` (`id`, `name`, `time`, `cmd`, `sTIme`, `eTime`, `status`, `isrunning`, `modify`, `process`, `ip`) VALUES
      (1, 'test', '*/1 * * * * ?', '/home/wida/sh.sh', 1427337701, 1437337701, 1, 0, 0, 1, ''),
      (2, 'test2', '*/1 * * * * ?', '/home/wida/sh2.sh', 1427337701, 1447337701, 1, 0, 0, 1, '');

###文件 /home/wida/sh.sh

    #!/bin/sh

    php /home/wida/php.php

###文件 /home/wida/sh2.sh

    #!/bin/sh

    php /home/wida/php.php

###文件 /home/wida/php.php

    <?php
        echo "test";
    ?>