CREATE TABLE IF NOT EXISTS `Links` (
	`id` INT AUTO_INCREMENT NOT NULL,
	`path` VARCHAR(255) NOT NULL,
	`flux` VARCHAR(255) NOT NULL,
	`parent` VARCHAR(255) DEFAULT NULL,
	`checked` BOOLEAN NOT NULL DEFAULT FALSE,
	`date` DATETIME DEFAULT NULL,
	UNIQUE `path` (`path`),
	PRIMARY KEY (`id`)
);
CREATE TABLE IF NOT EXISTS `PhotoLinks` (
	`id` INT AUTO_INCREMENT NOT NULL,
	`path` VARCHAR(255) NOT NULL,
	`dir` VARCHAR(255) DEFAULT NULL,
	`flux` VARCHAR(255) NOT NULL,
	`checked` BOOLEAN NOT NULL DEFAULT FALSE,
	`linksId` INT(11) DEFAULT NULL,
	`date` DATETIME DEFAULT NULL,
	UNIQUE `path` (`path`),
	PRIMARY KEY (`id`)
);
