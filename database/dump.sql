# ************************************************************
# Sequel Pro SQL dump
# Version 4541
#
# http://www.sequelpro.com/
# https://github.com/sequelpro/sequelpro
#
# Host: 52.82.24.155 (MySQL 5.7.24)
# Database: deposit
# Generation Time: 2018-11-27 09:44:01 +0000
# ************************************************************


/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

CREATE SCHEMA IF NOT EXISTS `deposit`;

USE `deposit`;

# Dump of table bases
# ------------------------------------------------------------

DROP TABLE IF EXISTS `bases`;

CREATE TABLE `bases` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `asset_id` char(64) NOT NULL,
  `control_program` text NOT NULL,
  PRIMARY KEY (`id`),
  KEY `asset_id` (`asset_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

LOCK TABLES `bases` WRITE;
UNLOCK TABLES;


# Dump of table balances
# ------------------------------------------------------------

DROP TABLE IF EXISTS `balances`;

CREATE TABLE `balances` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `address` varchar(256) NOT NULL,
  `asset_id` char(64) NOT NULL,
  `balance` bigint(20) unsigned DEFAULT '0',
  PRIMARY KEY (`id`),
  UNIQUE KEY `address` (`address`,`asset_id`),
  KEY `asset_id` (`asset_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

LOCK TABLES `balances` WRITE;
UNLOCK TABLES;


# Dump of table utxos
# ------------------------------------------------------------

DROP TABLE IF EXISTS `utxos`;

CREATE TABLE `utxos` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `hash` char(64) NOT NULL,
  `asset_id` char(64) NOT NULL,
  `amount` bigint(20) unsigned DEFAULT '0',
  `control_program` text NOT NULL,
  `is_spend` tinyint(1) DEFAULT '0',
  `is_locked` tinyint(1) DEFAULT '0',
  `submitTime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `duration` bigint(20) unsigned DEFAULT '0',
  PRIMARY KEY (`id`),
  UNIQUE KEY `hash` (`hash`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

LOCK TABLES `utxos` WRITE;
UNLOCK TABLES;

/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;
/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
