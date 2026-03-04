-- 会话表 --

Create TABLE `conversations` (
    `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
    `conversation_id` varchar(36) NOT NULL COMMENT '会话唯一ID（对外暴露）',
    `user_id` varchar(36) NOT NULL COMMENT '用户ID',
    `title` varchar(255) NOT NULL DEFAULT '' COMMENT '会话标题（可根据首轮对话生成）',
    `model_name` varchar(100) NOT NULL DEFAULT '' COMMENT '使用的模型名称',
    `system_prompt` text COMMENT '系统提示词',
    `status` tinyint(4) NOT NULL DEFAULT 1 COMMENT '状态：1-活跃，2-归档，3-删除',
    `total_tokens` int(11) NOT NULL DEFAULT 0 COMMENT '累计消耗token数',
    `message_count` int(11) NOT NULL DEFAULT 0 COMMENT '消息总数',
    `last_message_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '最后消息时间',
    `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    UNIQUE KEY `idx_conversation_id` (`conversation_id`),
    KEY `idx_user_id_status` (`user_id`, `status`, `last_message_at`),
    KEY `idx_last_message` (`last_message_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='对话会话表';

-- 消息表（每条对话内容）
CREATE TABLE `messages` (
    `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
    `message_id` varchar(36) NOT NULL COMMENT '消息唯一ID',
    `conversation_id` varchar(36) NOT NULL COMMENT '所属会话ID',
    `parent_id` varchar(36) DEFAULT NULL COMMENT '父消息ID（用于分支对话）',
    `role` varchar(20) NOT NULL COMMENT '角色：user/assistant/system/tool',
    `content` longtext NOT NULL COMMENT '消息内容',
    `content_type` varchar(20) NOT NULL DEFAULT 'text' COMMENT '内容类型：text/markdown/image/tool_call',
    `metadata` json DEFAULT NULL COMMENT '元数据（如工具调用信息、附件信息等）',
    `tokens` int(11) NOT NULL DEFAULT 0 COMMENT '本条消息消耗token数',
    `status` tinyint(4) NOT NULL DEFAULT 1 COMMENT '状态：1-正常，2-撤回，3-已编辑',
    `sequence` int(11) NOT NULL DEFAULT 0 COMMENT '会话内序号（用于排序）',
    `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    UNIQUE KEY `idx_message_id` (`message_id`),
    KEY `idx_conversation_role` (`conversation_id`, `role`, `sequence`),
    KEY `idx_conversation_created` (`conversation_id`, `created_at`),
    KEY `idx_parent_id` (`parent_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='对话消息表';

ALTER TABLE `messages` ADD FULLTEXT INDEX `ft_content` (`content`) COMMENT 'engine "mroonga"';