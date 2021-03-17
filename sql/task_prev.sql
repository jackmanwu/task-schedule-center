CREATE TABLE `task_prev`
(
    `id`          bigint NOT NULL AUTO_INCREMENT,
    `tid`         bigint NOT NULL DEFAULT '0',
    `prev_tid`    bigint NOT NULL DEFAULT '0',
    `create_time` int    NOT NULL DEFAULT '0',
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;