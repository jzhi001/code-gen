create table `jczhi_book_shelf`
(
    id          bigint(20)  not null comment '开课记录id, 开课订单id',
    uid         bigint(20)  not null comment '用户ID',
    account_id  bigint(20)  not null comment '账户id',
    created_at  bigint(20)  not null comment '创建时间',
    updated_at  bigint(20)  not null comment '更新时间',
    is_finished tinyint(1)  not null comment '是否完成 1: 成功，2: 失败',
    op_user     varchar(24) not null comment '操作人',
    book_count  bigint(8)   not null comment '数量',
    PRIMARY KEY (`id`),
    KEY `idx_uid_status` (uid, is_finished) comment '根据用户ID和状态查询'
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_bin comment '开课记录';