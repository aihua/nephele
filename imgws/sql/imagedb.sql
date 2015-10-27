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
