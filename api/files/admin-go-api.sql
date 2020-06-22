/*
 Navicat Premium Data Transfer

 Source Server         : localhost
 Source Server Type    : MySQL
 Source Server Version : 50725
 Source Host           : 127.0.0.1:3306
 Source Schema         : admin-go-api

 Target Server Type    : MySQL
 Target Server Version : 50725
 File Encoding         : 65001

 Date: 22/06/2020 08:16:19
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for app
-- ----------------------------
DROP TABLE IF EXISTS `app`;
CREATE TABLE `app` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(64) NOT NULL DEFAULT '',
  `deploy_type` tinyint(4) NOT NULL,
  `tid` int(11) NOT NULL DEFAULT '0' COMMENT '项目类型id',
  `active` tinyint(4) NOT NULL DEFAULT '0' COMMENT '是否启用项目',
  `enable_sync` tinyint(4) NOT NULL DEFAULT '0' COMMENT '是否需要初始化',
  `desc` varchar(128) NOT NULL DEFAULT '' COMMENT '说明',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `env_id` tinyint(4) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=utf8mb4 COMMENT='项目列表';


-- ----------------------------
-- Table structure for app_deploy
-- ----------------------------
DROP TABLE IF EXISTS `app_deploy`;
CREATE TABLE `app_deploy` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(64) NOT NULL DEFAULT '',
  `tid` int(11) NOT NULL DEFAULT '0' COMMENT '发布模板id',
  `repo_branch` varchar(128) NOT NULL DEFAULT '' COMMENT '项目分支',
  `repo_commit` varchar(255) NOT NULL DEFAULT '' COMMENT '项目commit id号',
  `status` tinyint(4) NOT NULL DEFAULT '0' COMMENT '状态 1: 新建提单， 2: 审核通过, 3:审核失败, 4:上线失败 ,5:上线成功',
  `desc` varchar(128) NOT NULL DEFAULT '' COMMENT '说明',
  `operator` int(11) NOT NULL DEFAULT '0' COMMENT '申请者',
  `review` int(11) NOT NULL DEFAULT '0' COMMENT '审核者',
  `deploy` int(11) NOT NULL DEFAULT '0' COMMENT '执行者',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=12 DEFAULT CHARSET=utf8mb4 COMMENT='应用发布信息表';


-- ----------------------------
-- Table structure for app_sync_value
-- ----------------------------
DROP TABLE IF EXISTS `app_sync_value`;
CREATE TABLE `app_sync_value` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(128) NOT NULL DEFAULT '' COMMENT '变量名',
  `value` varchar(1024) NOT NULL DEFAULT '' COMMENT '值',
  `aid` int(11) NOT NULL DEFAULT '0' COMMENT '业务 id',
  `desc` varchar(128) NOT NULL DEFAULT '' COMMENT '说明',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COMMENT='业务初始变量表';


-- ----------------------------
-- Table structure for app_type
-- ----------------------------
DROP TABLE IF EXISTS `app_type`;
CREATE TABLE `app_type` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(64) NOT NULL DEFAULT '' COMMENT '类型名',
  `desc` varchar(128) NOT NULL DEFAULT '' COMMENT '说明',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=7 DEFAULT CHARSET=utf8mb4 COMMENT='业务类型表';


-- ----------------------------
-- Table structure for app_value
-- ----------------------------
DROP TABLE IF EXISTS `app_value`;
CREATE TABLE `app_value` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(128) NOT NULL DEFAULT '' COMMENT '变量名',
  `value` varchar(1024) NOT NULL DEFAULT '' COMMENT '值',
  `aid` int(11) NOT NULL DEFAULT '0' COMMENT '业务 id',
  `desc` varchar(128) NOT NULL DEFAULT '' COMMENT '说明',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COMMENT='业务变量关联表';


-- ----------------------------
-- Table structure for certificate_info
-- ----------------------------
DROP TABLE IF EXISTS `certificate_info`;
CREATE TABLE `certificate_info` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `did` int(11) NOT NULL DEFAULT '0' COMMENT '域名id',
  `name` varchar(64) NOT NULL COMMENT '证书',
  `channel` varchar(255) NOT NULL COMMENT '申请渠道',
  `start_time` datetime NOT NULL COMMENT '证书申请时间',
  `end_time` datetime NOT NULL COMMENT '证书到期时间',
  `status` tinyint(4) NOT NULL DEFAULT '0' COMMENT '是否生效',
  `desc` varchar(128) NOT NULL DEFAULT '' COMMENT '说明',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `update_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=9 DEFAULT CHARSET=utf8mb4 COMMENT='证书表';


-- ----------------------------
-- Table structure for config_env
-- ----------------------------
DROP TABLE IF EXISTS `config_env`;
CREATE TABLE `config_env` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(64) NOT NULL DEFAULT '',
  `desc` varchar(128) NOT NULL DEFAULT '' COMMENT '说明',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=utf8mb4 COMMENT='环境类型';



-- ----------------------------
-- Table structure for deploy_extend
-- ----------------------------
DROP TABLE IF EXISTS `deploy_extend`;
CREATE TABLE `deploy_extend` (
  `dtid` int(11) NOT NULL AUTO_INCREMENT,
  `aid` int(11) NOT NULL DEFAULT '0' COMMENT '项目id',
  `tag` varchar(64) NOT NULL COMMENT '多功能预留字段，jar包作为打包版本号',
  `template_name` varchar(64) NOT NULL DEFAULT '',
  `dst_dir` varchar(64) NOT NULL DEFAULT '' COMMENT '目标主机部署路径',
  `dst_repo` varchar(64) NOT NULL DEFAULT '' COMMENT '目标主机仓库路径',
  `host_ids` varchar(128) NOT NULL DEFAULT '' COMMENT '发布模板绑定主机，多台以逗号分隔',
  `repo_url` varchar(128) NOT NULL DEFAULT '' COMMENT '项目git url',
  `versions` int(11) NOT NULL DEFAULT '0' COMMENT '项目保留历史版本数',
  `pre_code` varchar(255) NOT NULL DEFAULT '' COMMENT '代码clone前执行的命令',
  `post_code` varchar(255) NOT NULL DEFAULT '' COMMENT '代码clone后执行的命令',
  `pre_deploy` varchar(255) NOT NULL DEFAULT '' COMMENT '发布前执行命令',
  `post_deploy` varchar(255) NOT NULL DEFAULT '' COMMENT '发布后执行命令',
  `enable_check` tinyint(4) NOT NULL DEFAULT '0' COMMENT '是否开启审核',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`dtid`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=15 DEFAULT CHARSET=utf8mb4 COMMENT='项目发布模板表';



-- ----------------------------
-- Table structure for domain_info
-- ----------------------------
DROP TABLE IF EXISTS `domain_info`;
CREATE TABLE `domain_info` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(64) NOT NULL COMMENT '域名',
  `channel` varchar(255) NOT NULL COMMENT '申请渠道',
  `start_time` datetime NOT NULL COMMENT '域名申请时间',
  `end_time` datetime NOT NULL COMMENT '域名到期时间',
  `status` tinyint(4) NOT NULL DEFAULT '0' COMMENT '是否生效',
  `desc` varchar(128) NOT NULL DEFAULT '' COMMENT '说明',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `update_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=10 DEFAULT CHARSET=utf8mb4 COMMENT='域名表';


-- ----------------------------
-- Table structure for host
-- ----------------------------
DROP TABLE IF EXISTS `host`;
CREATE TABLE `host` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(64) NOT NULL DEFAULT '',
  `rid` int(11) NOT NULL DEFAULT '0' COMMENT '主机类型id',
  `env_id` int(11) NOT NULL DEFAULT '0' COMMENT '主机环境id',
  `zone_id` int(11) NOT NULL DEFAULT '0' COMMENT '主机区域Id',
  `status` tinyint(4) NOT NULL DEFAULT '0',
  `enable` tinyint(4) NOT NULL DEFAULT '0' COMMENT '是否在跳板机上显示',
  `username` varchar(32) NOT NULL DEFAULT '',
  `addres` varchar(64) NOT NULL DEFAULT '',
  `port` int(11) NOT NULL DEFAULT '22',
  `operator` int(11) NOT NULL DEFAULT '0',
  `desc` varchar(128) NOT NULL DEFAULT '' COMMENT '说明',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=6 DEFAULT CHARSET=utf8mb4 COMMENT='主机表';


-- ----------------------------
-- Table structure for host_app
-- ----------------------------
DROP TABLE IF EXISTS `host_app`;
CREATE TABLE `host_app` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `hid` int(11) NOT NULL DEFAULT '0' COMMENT '主机 id',
  `aid` int(11) NOT NULL DEFAULT '0' COMMENT '业务 id',
  `status` int(11) NOT NULL COMMENT '绑定项目是否初始化',
  `desc` varchar(128) NOT NULL DEFAULT '' COMMENT '说明',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=14 DEFAULT CHARSET=utf8mb4 COMMENT='主机业务关联表';



-- ----------------------------
-- Table structure for host_role
-- ----------------------------
DROP TABLE IF EXISTS `host_role`;
CREATE TABLE `host_role` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(64) NOT NULL DEFAULT '',
  `desc` varchar(128) NOT NULL DEFAULT '' COMMENT '说明',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=10 DEFAULT CHARSET=utf8mb4 COMMENT='主机类型';



-- ----------------------------
-- Table structure for menu_permissions
-- ----------------------------
DROP TABLE IF EXISTS `menu_permissions`;
CREATE TABLE `menu_permissions` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(100) NOT NULL DEFAULT '' COMMENT '权限名',
  `pid` int(11) NOT NULL DEFAULT '0' COMMENT '父级id',
  `type` tinyint(1) NOT NULL DEFAULT '0' COMMENT '1:菜单项 2: 权限项',
  `permission` varchar(120) NOT NULL DEFAULT '' COMMENT '权限项唯一标识',
  `url` varchar(120) NOT NULL DEFAULT '' COMMENT '菜单url',
  `icon` varchar(50) NOT NULL DEFAULT '' COMMENT '菜单图标',
  `desc` varchar(128) NOT NULL DEFAULT '' COMMENT '简介',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=129 DEFAULT CHARSET=utf8;

-- ----------------------------
-- Records of menu_permissions
-- ----------------------------
BEGIN;
INSERT INTO `menu_permissions` VALUES (1, '系统管理', 0, 1, '', '', 'setting', '');
INSERT INTO `menu_permissions` VALUES (2, '菜单管理', 0, 1, '', '', 'menu', '');
INSERT INTO `menu_permissions` VALUES (9, '用户列表', 1, 1, '', '/user/list', 'team', '用户列表');
INSERT INTO `menu_permissions` VALUES (16, '角色列表', 1, 1, '', '/user/role', 'lock', '角色列表。');
INSERT INTO `menu_permissions` VALUES (17, '权限列表', 1, 1, '', '/user/perm', 'security-scan', '权限列表');
INSERT INTO `menu_permissions` VALUES (18, '一级菜单', 2, 1, '', '/menu/menu', 'tag', '一级菜单');
INSERT INTO `menu_permissions` VALUES (19, '二级菜单', 2, 1, '', '/menu/submenu', 'tags', '二级菜单');
INSERT INTO `menu_permissions` VALUES (24, '用户添加', 9, 2, 'user-add', '', '', '添加用户');
INSERT INTO `menu_permissions` VALUES (31, '用户修改', 9, 2, 'user-edit', '', '', '用户修改');
INSERT INTO `menu_permissions` VALUES (32, '用户删除', 9, 2, 'user-del', '', '', '用户删除');
INSERT INTO `menu_permissions` VALUES (33, '角色添加', 16, 2, 'role-add', '', '', '角色添加');
INSERT INTO `menu_permissions` VALUES (34, '角色编辑', 16, 2, 'role-edit', '', '', '角色编辑');
INSERT INTO `menu_permissions` VALUES (35, '角色删除', 16, 2, 'role-del', '', '', '角色删除');
INSERT INTO `menu_permissions` VALUES (36, '权限项添加', 17, 2, 'perm-add', '', '', '权限项添加');
INSERT INTO `menu_permissions` VALUES (37, '权限项修改', 17, 2, 'perm-edit', '', '', '权限项修改');
INSERT INTO `menu_permissions` VALUES (38, '权限项删除', 17, 2, 'perm-del', '', '', '权限项删除');
INSERT INTO `menu_permissions` VALUES (39, '一级菜单添加', 18, 2, 'menu-add', '', '', '一级菜单添加');
INSERT INTO `menu_permissions` VALUES (40, '一级菜单修改', 18, 2, 'menu-edit', '', '', '一级菜单修改');
INSERT INTO `menu_permissions` VALUES (41, '一级菜单删除', 18, 2, 'menu-del', '', '', '一级菜单删除');
INSERT INTO `menu_permissions` VALUES (58, '主机管理', 0, 1, '', '', 'desktop', '');
INSERT INTO `menu_permissions` VALUES (59, '主机列表', 58, 1, '', '/host/list', 'cloud-server', '主机列表');
INSERT INTO `menu_permissions` VALUES (60, '主机类型', 58, 1, '', '/host/role', 'code-sandbox', '主机类型');
INSERT INTO `menu_permissions` VALUES (62, '应用配置', 0, 1, '', '', 'tool', '');
INSERT INTO `menu_permissions` VALUES (63, '应用发布', 0, 1, '', '', 'deployment-unit', '');
INSERT INTO `menu_permissions` VALUES (64, '环境管理', 62, 1, '', '/config/environment', 'environment', '环境管理');
INSERT INTO `menu_permissions` VALUES (65, '应用配置', 62, 1, '', '/config/app', 'project', '应用配置');
INSERT INTO `menu_permissions` VALUES (67, '应该发布', 63, 1, '', '/deploy/app', 'cloud-sync', '应该发布列表页');
INSERT INTO `menu_permissions` VALUES (68, '用户列表', 9, 2, 'user-list', '', '', '获取用户列表页');
INSERT INTO `menu_permissions` VALUES (69, '发布列表页', 67, 2, 'deploy-app-list', '', '', '应用发布列表页');
INSERT INTO `menu_permissions` VALUES (70, '发布提单', 67, 2, 'deploy-app-add', '', '', '应用发布提单');
INSERT INTO `menu_permissions` VALUES (71, '发布修改', 67, 2, 'deploy-app-edit', '', '', '应用发布修改');
INSERT INTO `menu_permissions` VALUES (72, '发布删除', 67, 2, 'deploy-app-del', '', '', '应用发布删除');
INSERT INTO `menu_permissions` VALUES (73, '发布审核', 67, 2, 'deploy-app-review', '', '', '应用发布审核');
INSERT INTO `menu_permissions` VALUES (74, '发布上线', 67, 2, 'deploy-app-redo', '', '', '应用发布上线');
INSERT INTO `menu_permissions` VALUES (75, '发布回滚', 67, 2, 'deploy-app-undo', '', '', '应用发布回滚');
INSERT INTO `menu_permissions` VALUES (76, '发布版本信息', 67, 2, 'config-app-git', '', '', '发布请求git版本信息');
INSERT INTO `menu_permissions` VALUES (77, '环境列表', 64, 2, 'config-env-list', '', '', '配置中心环境列表');
INSERT INTO `menu_permissions` VALUES (78, '新增环境类型', 64, 2, 'config-env-add', '', '', '新增环境类型');
INSERT INTO `menu_permissions` VALUES (79, '环境类型修改', 64, 2, 'config-env-edit', '', '', '环境类型信息修改');
INSERT INTO `menu_permissions` VALUES (80, '删除环境类型', 64, 2, 'config-env-del', '', '', '删除环境类型');
INSERT INTO `menu_permissions` VALUES (83, '二级菜单列表', 19, 2, 'submenu-list', '', '', '二级菜单列表页');
INSERT INTO `menu_permissions` VALUES (84, '二级菜单添加', 19, 2, 'submenu-add', '', '', '二级菜单添加');
INSERT INTO `menu_permissions` VALUES (85, '二级菜单修改', 19, 2, 'submenu-edit', '', '', '二级菜单添加');
INSERT INTO `menu_permissions` VALUES (86, '二级菜单删除', 19, 2, 'submenu-del', '', '', '二级菜单删除');
INSERT INTO `menu_permissions` VALUES (87, '主机类型列表', 60, 2, 'host-role-list', '', '', '主机类型列表');
INSERT INTO `menu_permissions` VALUES (88, '主机类型添加', 60, 2, 'host-role-add', '', '', '主机类型添加');
INSERT INTO `menu_permissions` VALUES (89, '主机类型修改', 60, 2, 'host-role-edit', '', '', '主机类型修改');
INSERT INTO `menu_permissions` VALUES (90, '主机类型删除', 60, 2, 'host-role-del', '', '', '主机类型删除');
INSERT INTO `menu_permissions` VALUES (91, '主机列表', 59, 2, 'host-list', '', '', '主机列表');
INSERT INTO `menu_permissions` VALUES (92, '添加主机', 59, 2, 'host-add', '', '', '添加主机');
INSERT INTO `menu_permissions` VALUES (93, '修改主机', 59, 2, 'host-edit', '', '', '修改主机');
INSERT INTO `menu_permissions` VALUES (94, '删除主机', 59, 2, 'host-del', '', '', '删除主机');
INSERT INTO `menu_permissions` VALUES (95, '主机业务查看', 59, 2, 'host-app-list', '', '', '主机业务查看');
INSERT INTO `menu_permissions` VALUES (96, '主机业务添加', 59, 2, 'host-app-add', '', '', '主机业务添加');
INSERT INTO `menu_permissions` VALUES (97, '主机业务删除', 59, 2, 'host-app-del', '', '', '主机业务删除');
INSERT INTO `menu_permissions` VALUES (98, '主机业务修改', 59, 2, 'host-app-edit', '', '', '主机业务修改');
INSERT INTO `menu_permissions` VALUES (99, '主机console', 59, 2, 'host-console', '', '', '主机console');
INSERT INTO `menu_permissions` VALUES (100, '角色权限项查看', 16, 2, 'role-perm-list', '', '', '角色权限项查看');
INSERT INTO `menu_permissions` VALUES (101, '角色权限项添加', 16, 2, 'role-perm-add', '', '', '角色权限项添加');
INSERT INTO `menu_permissions` VALUES (102, '应用列表', 65, 2, 'config-app-list', '', '', '应用列表');
INSERT INTO `menu_permissions` VALUES (103, '应用添加', 65, 2, 'config-app-add', '', '', '应用添加');
INSERT INTO `menu_permissions` VALUES (104, '应用修改', 65, 2, 'config-app-edit', '', '', '应用修改');
INSERT INTO `menu_permissions` VALUES (105, '应用删除', 65, 2, 'config-app-del', '', '', '应用删除');
INSERT INTO `menu_permissions` VALUES (106, '应用初始化', 65, 2, 'config-app-init', '', '', '应用初始化');
INSERT INTO `menu_permissions` VALUES (107, '应用变量设置', 65, 2, 'config-app-set', '', '', '应用变量设置');
INSERT INTO `menu_permissions` VALUES (108, '应用类型', 62, 1, '', '/config/appType', 'flag', '应用类型');
INSERT INTO `menu_permissions` VALUES (109, '应用类型列表', 108, 2, 'app-type-list', '', '', '应用类型列表');
INSERT INTO `menu_permissions` VALUES (110, '新增应用类型', 108, 2, 'app-type-add', '', '', '新增应用类型');
INSERT INTO `menu_permissions` VALUES (111, '修改应用类型', 108, 2, 'app-type-edit', '', '', '修改应用类型');
INSERT INTO `menu_permissions` VALUES (112, '删除应用类型', 108, 2, 'app-type-del', '', '', '删除应用类型');
INSERT INTO `menu_permissions` VALUES (113, '域名管理', 0, 1, '', '', 'google', '');
INSERT INTO `menu_permissions` VALUES (114, '域名列表', 113, 1, '', '/domain/list', 'chrome', '域名信息汇总页');
INSERT INTO `menu_permissions` VALUES (115, '证书列表', 113, 1, '', '/domain/certificate', 'security-scan', '证书信息页');
INSERT INTO `menu_permissions` VALUES (116, '域名列表', 114, 2, 'domain-info-list', '', '', '域名列表');
INSERT INTO `menu_permissions` VALUES (117, '添加域名', 114, 2, 'domain-info-add', '', '', '添加域名');
INSERT INTO `menu_permissions` VALUES (118, '修改域名信息', 114, 2, 'domain-info-edit', '', '', '修改域名信息');
INSERT INTO `menu_permissions` VALUES (119, '删除域名', 114, 2, 'domain-info-del', '', '', '删除域名');
INSERT INTO `menu_permissions` VALUES (120, '证书列表', 115, 2, 'domain-cert-list', '', '', '证书列表');
INSERT INTO `menu_permissions` VALUES (121, '添加证书', 115, 2, 'domain-cert-add', '', '', '添加证书');
INSERT INTO `menu_permissions` VALUES (122, '修改证书', 115, 2, 'domain-cert-edit', '', '', '修改证书');
INSERT INTO `menu_permissions` VALUES (123, '删除证书', 115, 2, 'domain-cert-del', '', '', '删除证书');
INSERT INTO `menu_permissions` VALUES (124, '任务计划', 0, 1, '', '', 'schedule', '');
INSERT INTO `menu_permissions` VALUES (125, '任务列表', 124, 1, '', '/schedule/list', 'bars', '任务列表');
INSERT INTO `menu_permissions` VALUES (126, '新增任务 ', 125, 2, 'schedule-job-add', '', '', '新增Job');
INSERT INTO `menu_permissions` VALUES (127, '任务修改', 125, 2, 'schedule-job-edit', '', '', '修改job');
INSERT INTO `menu_permissions` VALUES (128, '删除任务', 125, 2, 'schedule-job-del', '', '', 'Job删除');
COMMIT;

-- ----------------------------
-- Table structure for notify
-- ----------------------------
DROP TABLE IF EXISTS `notify`;
CREATE TABLE `notify` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `title` varchar(255) NOT NULL DEFAULT '' COMMENT '通知标题',
  `type` tinyint(4) NOT NULL DEFAULT '0' COMMENT '1：通知， 2：代办',
  `source` tinyint(4) NOT NULL DEFAULT '0' COMMENT '1：monitor 监控中心， 2：schedule 任务计划',
  `content` varchar(255) NOT NULL DEFAULT '' COMMENT '通知内容',
  `unread` tinyint(1) NOT NULL DEFAULT '0' COMMENT '0：已经查看处理， 1：未处理',
  `link` varchar(128) NOT NULL DEFAULT '' COMMENT '通知附加链接',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=21 DEFAULT CHARSET=utf8mb4 COMMENT='通知信息表';



-- ----------------------------
-- Table structure for role
-- ----------------------------
DROP TABLE IF EXISTS `role`;
CREATE TABLE `role` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(50) NOT NULL DEFAULT '' COMMENT '角色名',
  `desc` varchar(255) NOT NULL DEFAULT '' COMMENT '角色介绍',
  PRIMARY KEY (`id`),
  UNIQUE KEY `name` (`name`)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8 COMMENT='角色表';



-- ----------------------------
-- Table structure for role_permission_rel
-- ----------------------------
DROP TABLE IF EXISTS `role_permission_rel`;
CREATE TABLE `role_permission_rel` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `rid` int(11) NOT NULL DEFAULT '0' COMMENT '角色id',
  `pid` int(11) NOT NULL DEFAULT '0' COMMENT '权限id',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=56 DEFAULT CHARSET=utf8;


-- ----------------------------
-- Table structure for task
-- ----------------------------
DROP TABLE IF EXISTS `task`;
CREATE TABLE `task` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(64) NOT NULL DEFAULT '' COMMENT '任务名字',
  `is_more` tinyint(4) NOT NULL DEFAULT '0' COMMENT '是否支持多实例同时运行0：否，1：是',
  `host_ids` varchar(128) NOT NULL DEFAULT '' COMMENT '任务绑定主机，多台以逗号分隔',
  `active` tinyint(4) NOT NULL DEFAULT '0' COMMENT '是否启用任务',
  `command` varchar(2048) NOT NULL DEFAULT '' COMMENT '任务执行命令',
  `spec` varchar(32) NOT NULL DEFAULT '' COMMENT 'UNIX cron',
  `trigger_type` tinyint(4) NOT NULL DEFAULT '0' COMMENT '任务类型1：一次性任务， 2： UNIX cron',
  `start_time` datetime DEFAULT NULL COMMENT '任务执行开始时间',
  `end_time` datetime DEFAULT NULL COMMENT '任务执行结束时间',
  `operator` int(11) NOT NULL DEFAULT '0' COMMENT '添加任务用户',
  `desc` varchar(128) NOT NULL DEFAULT '' COMMENT '说明',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `task_name` (`name`)
) ENGINE=InnoDB AUTO_INCREMENT=18 DEFAULT CHARSET=utf8mb4 COMMENT='任务列表';



-- ----------------------------
-- Table structure for task_history
-- ----------------------------
DROP TABLE IF EXISTS `task_history`;
CREATE TABLE `task_history` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `task_id` int(11) NOT NULL DEFAULT '0' COMMENT '任务id',
  `host_id` int(11) NOT NULL DEFAULT '0' COMMENT '任务执行主机id',
  `status` tinyint(4) NOT NULL DEFAULT '0' COMMENT '上次执行状态0：成功，1：异常',
  `run_time` varchar(20) NOT NULL DEFAULT '' COMMENT '任务执行时长',
  `output` mediumtext NOT NULL COMMENT '任务执行输出信息',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=198 DEFAULT CHARSET=utf8mb4 COMMENT='任务执行历史信息表';


-- ----------------------------
-- Table structure for user
-- ----------------------------
DROP TABLE IF EXISTS `user`;
CREATE TABLE `user` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `rid` int(11) NOT NULL DEFAULT '0' COMMENT '角色id',
  `name` varchar(50) NOT NULL DEFAULT '' COMMENT '用户名',
  `nickname` varchar(50) NOT NULL DEFAULT '' COMMENT '昵称',
  `password_hash` varchar(100) NOT NULL DEFAULT '' COMMENT 'hash密码',
  `email` varchar(120) NOT NULL DEFAULT '' COMMENT '邮箱',
  `mobile` varchar(30) NOT NULL DEFAULT '' COMMENT '电话',
  `is_supper` tinyint(1) NOT NULL DEFAULT '0' COMMENT '是否为超级用户',
  `is_active` tinyint(1) NOT NULL DEFAULT '0' COMMENT '用户是否激活',
  `access_token` varchar(120) NOT NULL DEFAULT '' COMMENT '用户token',
  `token_expired` int(11) NOT NULL DEFAULT '0' COMMENT 'token过期时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `name` (`name`)
) ENGINE=InnoDB AUTO_INCREMENT=8 DEFAULT CHARSET=utf8 COMMENT='用户表';

SET FOREIGN_KEY_CHECKS = 1;
