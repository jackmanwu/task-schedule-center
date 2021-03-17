CREATE TABLE `node`
(
    `id`          int    NOT NULL AUTO_INCREMENT,
    `gid`         int    NOT NULL DEFAULT '0' COMMENT '														',
    `ip`          bigint NOT NULL DEFAULT '0',
    `state`       tinyint(1) NOT NULL DEFAULT '1' COMMENT '0-下线；1-在线',
    `update_time` int    NOT NULL DEFAULT '0',
    PRIMARY KEY (`id`),
    UNIQUE KEY `idx_gid_ip` (`gid`,`ip`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;