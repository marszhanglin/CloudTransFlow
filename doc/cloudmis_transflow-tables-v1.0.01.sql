

-- 流水表
drop TABLE if exists trans_flow; 
CREATE TABLE `trans_flow` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '主键',
  `flow_no` varchar(128) DEFAULT NULL COMMENT '上送流水号',
  `device_sn` varchar(128) DEFAULT NULL COMMENT '设备sn',
  `upload_time` varchar(128) DEFAULT NULL COMMENT '流水上送时间 YYYYMMDDHHmmss',
  `trans_time` varchar(128) DEFAULT NULL COMMENT '流水上送时间 YYYYMMDDHHmmss',
  `trans_type` varchar(128) DEFAULT NULL COMMENT '交易类型',
  `channel_id` varchar(128) DEFAULT NULL COMMENT '支付渠道',
  `merchant_id` varchar(128) DEFAULT NULL COMMENT '商户号',
  `terminal_id` varchar(128) DEFAULT NULL COMMENT '终端号',
  `merchant_name` varchar(128) DEFAULT NULL COMMENT '商户名称',
  `amount` bigint(20) DEFAULT NULL COMMENT '交易金额(分)',
  `trans_amount` bigint(20) DEFAULT NULL COMMENT '实际交易金额(分)',
  `currency_code` varchar(128) DEFAULT '156' COMMENT '货币类型',
  `out_order_no` varchar(128) DEFAULT NULL COMMENT '外部订单号',
  `voucher_no` varchar(128) DEFAULT NULL COMMENT '终端流水号',
  `reference_no` varchar(128) DEFAULT NULL COMMENT '系统流水号',
  `auth_code` varchar(128) DEFAULT NULL COMMENT '授权码',
  `ori_out_order_no` varchar(128) DEFAULT NULL COMMENT '原外部流水号',
  `ori_voucher_no` varchar(128) DEFAULT NULL COMMENT '原终端流水号',
  `ori_reference_no` varchar(128) DEFAULT NULL COMMENT '原系统流水号',
  `ori_auth_code` varchar(128) DEFAULT NULL COMMENT '原授权码',
  `card_no` varchar(128) DEFAULT NULL COMMENT '卡号',
  `operator_no` varchar(128) DEFAULT NULL COMMENT '操作员编号',
  `combination_no` varchar(128) DEFAULT NULL COMMENT '组合支付编号',
  `cardType` varchar(128) DEFAULT NULL COMMENT '银行卡类型',
  `remark` varchar(512) DEFAULT NULL COMMENT '备注',
  `extendParams` varchar(1024) DEFAULT NULL COMMENT '扩展字段',
  `shop_id` bigint(20) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=5 DEFAULT CHARSET=utf8 COMMENT='流水表'; 


-- 用户表
drop TABLE if exists t_user; 
CREATE TABLE `t_user` (
  `user_id` bigint(20) NOT NULL AUTO_INCREMENT,
  `user_name` varchar(128) DEFAULT NULL COMMENT '登陆账号',
  `password` varchar(128) DEFAULT NULL COMMENT '登陆密码',
  `user_status` bigint(20) DEFAULT NULL COMMENT '用户状态:\r\n            1 - 正常\r\n            0 - 锁定\r\n            ',
  `user_type` bigint(20) NOT NULL COMMENT '用户类型:     管理员|运营商|店长|操作员\r\n             1 - 管理员\r\n             2 - 运营商\r\n             3 - 店长\r\n             4 - 操作员\r\n            ',
  `operator_id` bigint(20) DEFAULT NULL COMMENT '所属运营商信息表主键',
  `shop_id` bigint(20) DEFAULT NULL COMMENT '所属店铺信息表主键',
  `cre_time` datetime DEFAULT NULL,
  `upd_time` datetime DEFAULT NULL, 
  `remark` varchar(512) DEFAULT NULL,
  PRIMARY KEY (`user_id`)
) ENGINE=InnoDB AUTO_INCREMENT=0 DEFAULT CHARSET=utf8;


INSERT INTO `t_user` VALUES (1, 'admin',      '88888888', 1, 1, NULL, NULL, '2015-11-3 17:15:41', '2017-1-23 11:25:53', NULL);
INSERT INTO `t_user` VALUES (2, 'newlandyys', '88888888', 1, 2, 1,    NULL, '2016-12-11 11:12:12', '2017-1-15 15:17:22', NULL);
INSERT INTO `t_user` VALUES (3, 'newlanddz',  '88888888', 1, 3, 1,    1,    '2016-12-11 11:12:12', '2017-1-15 15:17:22', NULL);
INSERT INTO `t_user` VALUES (4, 'newlandczy', '88888888', 1, 4, 1,    1,    '2016-12-11 11:12:12', '2017-1-15 15:17:22', NULL);


-- 应用注册表  用于注册运营商  暂时不考虑
drop TABLE if exists t_user_regist; 
CREATE TABLE `t_user_regist` (
  `user_regist_id` bigint(20) NOT NULL AUTO_INCREMENT,
  `user_id` bigint(20) NOT NULL,
  `user_name` varchar(128) NOT NULL,
  `active_code` varchar(128) NOT NULL COMMENT '激活码',
  `email` varchar(128) NOT NULL COMMENT '激活邮件的目标邮箱',
  `regist_status` bigint(20) NOT NULL COMMENT '0 无效  1有效',
  `cre_time` datetime NOT NULL,
  `upd_time` datetime DEFAULT NULL,
  PRIMARY KEY (`user_regist_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='注册激活信息表';


-- 角色表   
drop TABLE if exists t_role;  
CREATE TABLE `t_role` (
  `role_id` bigint(20) NOT NULL AUTO_INCREMENT,
  `role_code` varchar(128) NOT NULL COMMENT '角色标识（英文）',
  `role_name` varchar(128) NOT NULL COMMENT '角色名称',
  `role_remark` varchar(128) DEFAULT NULL,
  `cre_time` datetime NOT NULL,
  `upd_time` datetime DEFAULT NULL,
  PRIMARY KEY (`role_id`)
) ENGINE=InnoDB AUTO_INCREMENT=0 DEFAULT CHARSET=utf8;

INSERT INTO `t_role` VALUES (1, 'sys_admin', '系统管理员', NULL, '2015-11-21 00:00:00', NULL);
INSERT INTO `t_role` VALUES (2, 'operator_admin', '运营商管理员', NULL, '2015-11-21 00:00:00', NULL);
INSERT INTO `t_role` VALUES (3, 'shop_admin', '店铺管理员', NULL, '2015-11-11 00:00:00', NULL);
INSERT INTO `t_role` VALUES (4, 'base_admin', '操作员', NULL, '2015-11-11 00:00:00', NULL);






-- 运营商表
drop TABLE if exists t_operator; 
CREATE TABLE `t_operator` (
  `operator_id` bigint(20) NOT NULL AUTO_INCREMENT,
  `operator_name` varchar(128) DEFAULT NULL COMMENT '客户名称',
  `address` varchar(128) DEFAULT NULL COMMENT '地址',
  `linkman` varchar(128) DEFAULT NULL COMMENT '联系人',
  `mobile` varchar(128) DEFAULT NULL COMMENT '手机号',
  `operator_status` bigint(20) DEFAULT NULL COMMENT '状态：\r\n            0 - 锁定\r\n            1 - 正常\r\n            2 - 待审核',
  `cre_time` datetime DEFAULT NULL,
  `upd_time` datetime DEFAULT NULL, 
  `remark` varchar(512) DEFAULT NULL,
  PRIMARY KEY (`operator_id`)
) ENGINE=InnoDB AUTO_INCREMENT=0 DEFAULT CHARSET=utf8 COMMENT='运营商信息表';
INSERT INTO `t_operator` VALUES (1, '新大陆测试运营商', '福州马尾儒江', '新大陆联系人', '18912345678', 1, '2016-12-11 11:12:34', '2016-12-11 11:12:57', '审核通过');

-- 店铺表
drop TABLE if exists t_shop; 
CREATE TABLE `t_shop` (
  `shop_id` bigint(20) NOT NULL AUTO_INCREMENT,
  `shop_name` varchar(120) DEFAULT NULL COMMENT '店铺名称',
  `address` varchar(120) DEFAULT NULL COMMENT '地址',
  `operator_id` bigint(20) NOT NULL COMMENT '所属运营商信息表主键',
  `cre_time` datetime DEFAULT NULL,
  `upd_time` datetime DEFAULT NULL, 
  `remark` varchar(512) DEFAULT NULL,
  PRIMARY KEY (`shop_id`)
) ENGINE=InnoDB AUTO_INCREMENT=0 DEFAULT CHARSET=utf8 COMMENT='店铺信息表';

INSERT INTO `t_shop` VALUES (1, '新大陆测试门店', '福州马尾儒江',1, '2016-12-11 11:12:34', '2016-12-11 11:12:57', '');


