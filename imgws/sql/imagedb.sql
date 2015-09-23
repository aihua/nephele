CREATE TABLE `channel` (
  `name` varchar(45) NOT NULL,
  `code` char(2) NOT NULL,
  PRIMARY KEY (`name`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

CREATE TABLE `config` (
  `idx` int(11) NOT NULL AUTO_INCREMENT,
  `channel` char(2) NOT NULL,
  `key` varchar(45) NOT NULL,
  `value` varchar(2048) NOT NULL,
  `recordTime` varchar(19) NOT NULL,
  PRIMARY KEY (`idx`),
  KEY `i_channelkey` (`channel`,`key`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

CREATE TABLE `imageindex_1` (
  `idx` bigint(20) NOT NULL AUTO_INCREMENT,
  `channel` char(2) NOT NULL,
  `storagePath` varchar(256) NOT NULL,
  `storageType` varchar(45) NOT NULL,
  `profile` varchar(2048) NOT NULL,
  `createTime` datetime NOT NULL,
  `partitionKey` smallint(6) NOT NULL,
  PRIMARY KEY (`idx`,`partitionKey`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=latin1
/*!50100 PARTITION BY HASH (partitionKey)
PARTITIONS 256 */;

CREATE TABLE `imageplan_1` (
  `idx` bigint(20) NOT NULL AUTO_INCREMENT,
  `imgIdx` bigint(20) DEFAULT NULL,
  `plan` varchar(2048) DEFAULT NULL,
  `partitionKey` smallint(6) NOT NULL,
  PRIMARY KEY (`idx`,`partitionKey`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1
/*!50100 PARTITION BY HASH (partitionKey)
PARTITIONS 256 */;
