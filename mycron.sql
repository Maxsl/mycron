--
-- Database: `mycron`
--

-- --------------------------------------------------------

--
-- 表的结构 `cron`
--

CREATE TABLE `cron` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(50) NOT NULL DEFAULT '',
  `time` varchar(50) DEFAULT '',
  `cmd` varchar(255) DEFAULT '',
  `sTIme` int(11) NOT NULL,
  `eTime` int(11) NOT NULL,
  `status` tinyint(1) NOT NULL DEFAULT '0',
  `isrunning` tinyint(1) NOT NULL DEFAULT '0' COMMENT '是否在运行中 0不是,1是',
  `modify` tinyint(1) NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB;

--
-- 转存表中的数据 `cron`
--

INSERT INTO `cron` (`id`, `name`, `time`, `cmd`, `sTIme`, `eTime`, `status`, `isrunning`, `modify`) VALUES
  (1, 'test', '*/1 * * * * ?', '/home/wida/sh.sh', 1427337701, 1437337701, 1, 0, 0),
  (2, 'test2', '*/1 * * * * ?', '/home/wida/sh2.sh', 1427337701, 1447337701, 1, 0, 0);