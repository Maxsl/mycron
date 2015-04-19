
CREATE TABLE IF NOT EXISTS `cron` (
  `id` int(11) NOT NULL,
  `name` varchar(50) NOT NULL DEFAULT '',
  `time` varchar(50) DEFAULT '',
  `cmd` varchar(255) DEFAULT '',
  `sTIme` int(11) NOT NULL,
  `eTime` int(11) NOT NULL,
  `status` tinyint(1) NOT NULL DEFAULT '0',
  `isrunning` tinyint(1) NOT NULL DEFAULT '0' COMMENT '是否在运行中 0不是,1是',
  `modify` tinyint(1) NOT NULL DEFAULT '0',
  `process` tinyint(2) NOT NULL DEFAULT '1' COMMENT '进程数量',
  `ip` varchar(20) CHARACTER SET utf8 COLLATE utf8_bin NOT NULL DEFAULT ''
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=latin1;

--
-- 转存表中的数据 `cron`
--

INSERT INTO `cron` (`id`, `name`, `time`, `cmd`, `sTIme`, `eTime`, `status`, `isrunning`, `modify`, `process`, `ip`) VALUES
  (1, 'test', '*/1 * * * * ?', '/home/wida/sh.sh', 1427337701, 1437337701, 1, 0, 0, 1, ''),
  (2, 'test2', '*/1 * * * * ?', '/home/wida/sh2.sh', 1427337701, 1447337701, 1, 0, 0, 1, '');

-- --------------------------------------------------------

--
-- 表的结构 `history`
--

CREATE TABLE IF NOT EXISTS `history` (
  `id` int(11) NOT NULL,
  `jid` int(11) NOT NULL,
  `ip` varchar(20) NOT NULL COMMENT '执行的ip',
  `dotime` int(11) NOT NULL,
  `costtime` int(11) NOT NULL COMMENT '执行时间',
  `ret` int(11) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=latin1 COMMENT='执行历史记录';