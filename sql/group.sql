CREATE TABLE `group`
(
    `id`          int         NOT NULL AUTO_INCREMENT,
    `name`        varchar(45) NOT NULL DEFAULT '',
    `create_time` int         NOT NULL DEFAULT '0',
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;