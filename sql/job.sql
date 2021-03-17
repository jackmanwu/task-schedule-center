CREATE TABLE `job`
(
    `id`          bigint       NOT NULL AUTO_INCREMENT,
    `gid`         int          NOT NULL DEFAULT '0',
    `tid`         bigint       NOT NULL DEFAULT '0',
    `Ip`          bigint       NOT NULL DEFAULT '0',
    `task_date`   bigint       NOT NULL DEFAULT '0',
    `param`       varchar(150) NOT NULL DEFAULT '',
    `create_time` int          NOT NULL DEFAULT '0',
    `update_time` int          NOT NULL DEFAULT '0',
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;