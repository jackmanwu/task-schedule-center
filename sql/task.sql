CREATE TABLE `task`
(
    `id`          bigint       NOT NULL AUTO_INCREMENT,
    `gid`         int          NOT NULL DEFAULT '0',
    `name`        varchar(30)  NOT NULL DEFAULT '',
    `cron`        varchar(10)  NOT NULL DEFAULT '',
    `state`       tinyint(1) NOT NULL DEFAULT '1' COMMENT '状态：0-禁用；1-启用',
    `path`        varchar(100) NOT NULL DEFAULT '',
    `uid`         bigint       NOT NULL DEFAULT '0',
    `create_time` int          NOT NULL DEFAULT '0',
    `update_time` int          NOT NULL DEFAULT '0',
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;