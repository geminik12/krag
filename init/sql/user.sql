CREATE TABLE `user` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `user_id` varchar(36) NOT NULL DEFAULT '' COMMENT '用户唯一 ID',
  `username` varchar(255) NOT NULL DEFAULT '' COMMENT '用户名（唯一）',
  `password` varchar(255) NOT NULL DEFAULT '' COMMENT '用户密码（加密后）',
  `nickname` varchar(30) NOT NULL DEFAULT '' COMMENT '用户昵称',
  `email` varchar(256) NOT NULL DEFAULT '' COMMENT '用户电子邮箱地址',
  `phone` varchar(16) NOT NULL DEFAULT '' COMMENT '用户手机号',
  `createdAt` datetime NOT NULL DEFAULT current_timestamp() COMMENT '用户创建时间',
  `updatedAt` datetime NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp() COMMENT '用户最后修改时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `user.user_id` (`user_id`),
  UNIQUE KEY `user.username` (`username`),
  UNIQUE KEY `user.phone` (`phone`)
) ENGINE=InnoDB AUTO_INCREMENT=5 DEFAULT CHARSET=utf8mb3 COLLATE=utf8mb3_general_ci COMMENT='用户表';