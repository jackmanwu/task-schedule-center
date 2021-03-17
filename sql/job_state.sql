CREATE TABLE `job_state`
(
    `id`          bigint      NOT NULL AUTO_INCREMENT,
    `job_id`      bigint      NOT NULL DEFAULT '0',
    `state`       varchar(45) NOT NULL DEFAULT '',
    `time`        bigint      NOT NULL DEFAULT '0',
    `create_time` int         NOT NULL DEFAULT '0',
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;