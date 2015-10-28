CREATE TABLE `channel` (
  `name` varchar(45) NOT NULL,
  `code` char(2) NOT NULL,
  PRIMARY KEY (`code`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

CREATE TABLE `config` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `channelCode` char(2) NOT NULL,
  `key` varchar(45) NOT NULL,
  `value` varchar(2048) NOT NULL,
  `recordTime` varchar(19) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `i_channelkey` (`channelCode`,`key`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

CREATE TABLE `imageindex_1` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `channelCode` char(2) NOT NULL,
  `storagePath` varchar(256) NOT NULL,
  `storageType` varchar(45) NOT NULL,
  `profile` varchar(2048) NOT NULL,
  `createTime` datetime NOT NULL,
  `partitionKey` smallint(6) NOT NULL,
  PRIMARY KEY (`id`,`partitionKey`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=latin1
/*!50100 PARTITION BY HASH (partitionKey)
PARTITIONS 256 */;

CREATE TABLE `imageplan_1` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `imgId` bigint(20) NOT NULL,
  `plan` varchar(4000) NOT NULL,
  `partitionKey` smallint(6) NOT NULL,
  PRIMARY KEY (`id`,`partitionKey`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1
/*!50100 PARTITION BY HASH (partitionKey)
PARTITIONS 256 */;

CREATE TABLE `imageindex_2` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `channelCode` char(2) NOT NULL,
  `storagePath` varchar(256) NOT NULL,
  `storageType` varchar(45) NOT NULL,
  `profile` varchar(2048) NOT NULL,
  `createTime` datetime NOT NULL,
  `partitionKey` smallint(6) NOT NULL,
  PRIMARY KEY (`id`,`partitionKey`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=latin1
/*!50100 PARTITION BY HASH (partitionKey)
PARTITIONS 256 */;

CREATE TABLE `imageplan_2` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `imgId` bigint(20) NOT NULL,
  `plan` varchar(4000) NOT NULL,
  `partitionKey` smallint(6) NOT NULL,
  PRIMARY KEY (`id`,`partitionKey`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1
/*!50100 PARTITION BY HASH (partitionKey)
PARTITIONS 256 */;

CREATE TABLE `imageindex_3` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `channelCode` char(2) NOT NULL,
  `storagePath` varchar(256) NOT NULL,
  `storageType` varchar(45) NOT NULL,
  `profile` varchar(2048) NOT NULL,
  `createTime` datetime NOT NULL,
  `partitionKey` smallint(6) NOT NULL,
  PRIMARY KEY (`id`,`partitionKey`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=latin1
/*!50100 PARTITION BY HASH (partitionKey)
PARTITIONS 256 */;

CREATE TABLE `imageplan_3` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `imgId` bigint(20) NOT NULL,
  `plan` varchar(4000) NOT NULL,
  `partitionKey` smallint(6) NOT NULL,
  PRIMARY KEY (`id`,`partitionKey`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1
/*!50100 PARTITION BY HASH (partitionKey)
PARTITIONS 256 */;

CREATE TABLE `imageindex_4` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `channelCode` char(2) NOT NULL,
  `storagePath` varchar(256) NOT NULL,
  `storageType` varchar(45) NOT NULL,
  `profile` varchar(2048) NOT NULL,
  `createTime` datetime NOT NULL,
  `partitionKey` smallint(6) NOT NULL,
  PRIMARY KEY (`id`,`partitionKey`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=latin1
/*!50100 PARTITION BY HASH (partitionKey)
PARTITIONS 256 */;

CREATE TABLE `imageplan_4` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `imgId` bigint(20) NOT NULL,
  `plan` varchar(4000) NOT NULL,
  `partitionKey` smallint(6) NOT NULL,
  PRIMARY KEY (`id`,`partitionKey`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1
/*!50100 PARTITION BY HASH (partitionKey)
PARTITIONS 256 */;

CREATE TABLE `imageindex_5` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `channelCode` char(2) NOT NULL,
  `storagePath` varchar(256) NOT NULL,
  `storageType` varchar(45) NOT NULL,
  `profile` varchar(2048) NOT NULL,
  `createTime` datetime NOT NULL,
  `partitionKey` smallint(6) NOT NULL,
  PRIMARY KEY (`id`,`partitionKey`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=latin1
/*!50100 PARTITION BY HASH (partitionKey)
PARTITIONS 256 */;

CREATE TABLE `imageplan_5` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `imgId` bigint(20) NOT NULL,
  `plan` varchar(4000) NOT NULL,
  `partitionKey` smallint(6) NOT NULL,
  PRIMARY KEY (`id`,`partitionKey`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1
/*!50100 PARTITION BY HASH (partitionKey)
PARTITIONS 256 */;

CREATE TABLE `imageindex_6` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `channelCode` char(2) NOT NULL,
  `storagePath` varchar(256) NOT NULL,
  `storageType` varchar(45) NOT NULL,
  `profile` varchar(2048) NOT NULL,
  `createTime` datetime NOT NULL,
  `partitionKey` smallint(6) NOT NULL,
  PRIMARY KEY (`id`,`partitionKey`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=latin1
/*!50100 PARTITION BY HASH (partitionKey)
PARTITIONS 256 */;

CREATE TABLE `imageplan_6` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `imgId` bigint(20) NOT NULL,
  `plan` varchar(4000) NOT NULL,
  `partitionKey` smallint(6) NOT NULL,
  PRIMARY KEY (`id`,`partitionKey`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1
/*!50100 PARTITION BY HASH (partitionKey)
PARTITIONS 256 */;


trucate table channel;
trucate table config;
INSERT INTO channel(name,code) VALUES('tg','10');
INSERT INTO channel(name,code) VALUES('hotel','20');

INSERT INTO config(channelCode,`key`,value,recordTime)VALUES('00','fdfsdomain','fdfs.tracker.fx.uat.qa.nt.ctripcorp.com','1445936136683331362');
INSERT INTO config(channelCode,`key`,value,recordTime)VALUES('00','fdfsport','22122','1445936136683331362');
INSERT INTO config(channelCode,`key`,value,recordTime)VALUES('00','nfs1','http://10.2.25.0:8081/target/','1445936136683331362');
INSERT INTO config(channelCode,`key`,value,recordTime)VALUES('00','nfs2','http://10.2.25.0:8082/target/','1445936136683331362');
INSERT INTO config(channelCode,`key`,value,recordTime)VALUES('00','sizes',',100x100,120x120,130x130,248x186,250x250,600x400,1000x1000,300x10000,10000x200,192x192,500x100000,290x170,','1445936136683331362');
INSERT INTO config(channelCode,`key`,value,recordTime)VALUES('00','fdfsgroups','group1','1445936136683331362');
