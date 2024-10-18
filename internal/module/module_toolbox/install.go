package module_toolbox

import "teamide/internal/install"

func GetInstallStages() []*install.StageModel {

	return []*install.StageModel{

		// 创建工具箱表
		{
			Version: "1.0",
			Module:  ModuleToolbox,
			Stage:   `创建表[` + TableToolbox + `]`,
			Sql: &install.StageSqlModel{
				Mysql: []string{`
CREATE TABLE ` + TableToolbox + ` (
	toolboxId bigint(20) NOT NULL COMMENT '工具箱ID',
	toolboxType varchar(10) NOT NULL COMMENT '工具箱类型',
	name varchar(50) NOT NULL COMMENT '名称',
	option varchar(2000) NOT NULL COMMENT '配置',
	userId bigint(20) DEFAULT NULL COMMENT '用户ID',
	deleted int(1) NOT NULL DEFAULT 2 COMMENT '启用状态:1-删除、2-正常',
	createTime datetime NOT NULL COMMENT '创建时间',
	updateTime datetime DEFAULT NULL COMMENT '修改时间',
	deleteTime datetime DEFAULT NULL COMMENT '删除时间',
	PRIMARY KEY (toolboxId),
	KEY index_userId (userId),
	KEY index_toolboxType (toolboxType),
	KEY index_name (name),
	KEY index_deleted (deleted)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='` + TableToolboxComment + `';
`},
				Sqlite: []string{`
CREATE TABLE ` + TableToolbox + ` (
	toolboxId bigint(20) NOT NULL,
	toolboxType varchar(10) NOT NULL,
	name varchar(50) NOT NULL,
	option varchar(2000) NOT NULL,
	userId bigint(20) DEFAULT NULL,
	deleted int(1) NOT NULL DEFAULT 2,
	createTime datetime NOT NULL,
	updateTime datetime DEFAULT NULL,
	deleteTime datetime DEFAULT NULL,
	PRIMARY KEY (toolboxId)
);
`,
					`CREATE INDEX ` + TableToolbox + `_index_userId on ` + TableToolbox + ` (userId);`,
					`CREATE INDEX ` + TableToolbox + `_index_toolboxType on ` + TableToolbox + ` (toolboxType);`,
					`CREATE INDEX ` + TableToolbox + `_index_name on ` + TableToolbox + ` (name);`,
					`CREATE INDEX ` + TableToolbox + `_index_deleted on ` + TableToolbox + ` (deleted);`,
				},
			},
		},
		// 创建工具箱打开记录表
		{
			Version: "1.0",
			Module:  ModuleToolbox,
			Stage:   `创建表[` + TableToolboxOpen + `]`,
			Sql: &install.StageSqlModel{
				Mysql: []string{`
CREATE TABLE ` + TableToolboxOpen + ` (
	openId bigint(20) NOT NULL COMMENT '开启ID',
	userId bigint(20) NOT NULL COMMENT '用户ID',
	toolboxId bigint(20) NOT NULL COMMENT '工具箱ID',
	extend varchar(4000) DEFAULT NULL COMMENT '扩展',
	createTime datetime NOT NULL COMMENT '创建时间',
	updateTime datetime DEFAULT NULL COMMENT '修改时间',
	openTime datetime DEFAULT NULL COMMENT '打开时间',
	PRIMARY KEY (openId),
	KEY index_userId (userId),
	KEY index_toolboxId (toolboxId)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='` + TableToolboxOpenComment + `';
`},
				Sqlite: []string{`
CREATE TABLE ` + TableToolboxOpen + ` (
	openId bigint(20) NOT NULL,
	userId bigint(20) NOT NULL,
	toolboxId bigint(20) NOT NULL,
	extend varchar(4000) DEFAULT NULL,
	createTime datetime NOT NULL,
	updateTime datetime DEFAULT NULL,
	openTime datetime DEFAULT NULL,
	PRIMARY KEY (openId)
);
`,
					`CREATE INDEX ` + TableToolboxOpen + `_index_userId on ` + TableToolboxOpen + ` (userId);`,
					`CREATE INDEX ` + TableToolboxOpen + `_index_toolboxId on ` + TableToolboxOpen + ` (toolboxId);`,
				},
			},
		},
		// 创建工具箱打开标签页表
		{
			Version: "1.0",
			Module:  ModuleToolbox,
			Stage:   `创建表[` + TableToolboxOpenTab + `]`,
			Sql: &install.StageSqlModel{
				Mysql: []string{`
CREATE TABLE ` + TableToolboxOpenTab + ` (
	tabId bigint(20) NOT NULL COMMENT '标签页ID',
	openId bigint(20) NOT NULL COMMENT '开启ID',
	userId bigint(20) NOT NULL COMMENT '用户ID',
	toolboxId bigint(20) NOT NULL COMMENT '工具箱ID',
	extend varchar(4000) DEFAULT NULL COMMENT '扩展',
	createTime datetime NOT NULL COMMENT '创建时间',
	updateTime datetime DEFAULT NULL COMMENT '修改时间',
	openTime datetime DEFAULT NULL COMMENT '打开时间',
	PRIMARY KEY (tabId),
	KEY index_openId (openId),
	KEY index_userId (userId),
	KEY index_toolboxId (toolboxId)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='` + TableToolboxOpenTabComment + `';
`},
				Sqlite: []string{`
CREATE TABLE ` + TableToolboxOpenTab + ` (
	tabId bigint(20) NOT NULL,
	openId bigint(20) NOT NULL,
	userId bigint(20) NOT NULL,
	toolboxId bigint(20) NOT NULL,
	extend varchar(4000) DEFAULT NULL,
	createTime datetime NOT NULL,
	updateTime datetime DEFAULT NULL,
	openTime datetime DEFAULT NULL,
	PRIMARY KEY (tabId)
);
`,
					`CREATE INDEX ` + TableToolboxOpenTab + `_index_openId on ` + TableToolboxOpenTab + ` (openId);`,
					`CREATE INDEX ` + TableToolboxOpenTab + `_index_userId on ` + TableToolboxOpenTab + ` (userId);`,
					`CREATE INDEX ` + TableToolboxOpenTab + `_index_toolboxId on ` + TableToolboxOpenTab + ` (toolboxId);`,
				},
			},
		},

		/** 给工具箱添加分组 开始 **/

		// 创建工具箱分组表
		{
			Version: "1.0.1",
			Module:  ModuleToolbox,
			Stage:   `创建表[` + TableToolboxGroup + `]`,
			Sql: &install.StageSqlModel{
				Mysql: []string{`
CREATE TABLE ` + TableToolboxGroup + ` (
	groupId bigint(20) NOT NULL COMMENT '分组ID',
	name varchar(50) NOT NULL COMMENT '名称',
	comment varchar(200) DEFAULT NULL COMMENT '说明',
	userId bigint(20) NOT NULL COMMENT '用户ID',
	option varchar(2000) DEFAULT NULL COMMENT '配置',
	createTime datetime NOT NULL COMMENT '创建时间',
	updateTime datetime DEFAULT NULL COMMENT '修改时间',
	PRIMARY KEY (groupId),
	KEY index_userId (userId),
	KEY index_name (name)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='` + TableToolboxGroupComment + `';
`},
				Sqlite: []string{`
CREATE TABLE ` + TableToolboxGroup + ` (
	groupId bigint(20) NOT NULL,
	name varchar(50) NOT NULL,
	comment varchar(200) DEFAULT NULL,
	userId bigint(20) NOT NULL,
	option varchar(2000) DEFAULT NULL,
	createTime datetime NOT NULL,
	updateTime datetime DEFAULT NULL,
	PRIMARY KEY (groupId)
);
`,
					`CREATE INDEX ` + TableToolboxGroup + `_index_userId on ` + TableToolboxGroup + ` (userId);`,
					`CREATE INDEX ` + TableToolboxGroup + `_index_name on ` + TableToolboxGroup + ` (name);`,
				},
			},
		},
		// 工具表添加分组ID
		{
			Version: "1.0.1",
			Module:  ModuleToolbox,
			Stage:   `工具箱[` + TableToolbox + `]添加分组ID[groupId]`,
			Sql: &install.StageSqlModel{
				Mysql: []string{
					`ALTER TABLE ` + TableToolbox + ` ADD COLUMN comment varchar(200) DEFAULT NULL COMMENT '说明' AFTER name;`,
					`ALTER TABLE ` + TableToolbox + ` ADD COLUMN groupId bigint(20) DEFAULT NULL COMMENT '分组ID' AFTER toolboxType;`,
					`ALTER TABLE ` + TableToolbox + ` ADD INDEX ` + TableToolbox + `_index_groupId (groupId);`,
				},
				Sqlite: []string{
					`ALTER TABLE ` + TableToolbox + ` ADD comment varchar(200);`,
					`ALTER TABLE ` + TableToolbox + ` ADD groupId bigint(20);`,
					`CREATE INDEX ` + TableToolbox + `_index_groupId on ` + TableToolbox + ` (groupId);`,
				},
			},
		},

		/** 给工具箱添加分组 结束 **/

		/** 给工具箱添加 快速命令 开始 **/

		// 创建工具箱 快速命令 表
		{
			Version: "1.0.2",
			Module:  ModuleToolbox,
			Stage:   `创建表[` + TableToolboxQuickCommand + `]`,
			Sql: &install.StageSqlModel{
				Mysql: []string{`
CREATE TABLE ` + TableToolboxQuickCommand + ` (
	quickCommandId bigint(20) NOT NULL COMMENT '快速命令ID',
	quickCommandType int(10) NOT NULL COMMENT '快速命令类型',
	name varchar(50) NOT NULL COMMENT '名称',
	comment varchar(200) DEFAULT NULL COMMENT '说明',
	userId bigint(20) NOT NULL COMMENT '用户ID',
	option varchar(2000) DEFAULT NULL COMMENT '配置',
	createTime datetime NOT NULL COMMENT '创建时间',
	updateTime datetime DEFAULT NULL COMMENT '修改时间',
	PRIMARY KEY (quickCommandId),
	KEY index_quickCommandType (quickCommandType),
	KEY index_userId (userId),
	KEY index_name (name)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='` + TableToolboxQuickCommandComment + `';
`},
				Sqlite: []string{`
CREATE TABLE ` + TableToolboxQuickCommand + ` (
	quickCommandId bigint(20) NOT NULL,
	quickCommandType int(10) NOT NULL,
	name varchar(50) NOT NULL,
	comment varchar(200) DEFAULT NULL,
	userId bigint(20) NOT NULL,
	option varchar(2000) DEFAULT NULL,
	createTime datetime NOT NULL,
	updateTime datetime DEFAULT NULL,
	PRIMARY KEY (quickCommandId)
);
`,
					`CREATE INDEX ` + TableToolboxQuickCommand + `_index_quickCommandType on ` + TableToolboxQuickCommand + ` (quickCommandType);`,
					`CREATE INDEX ` + TableToolboxQuickCommand + `_index_userId on ` + TableToolboxQuickCommand + ` (userId);`,
					`CREATE INDEX ` + TableToolboxQuickCommand + `_index_name on ` + TableToolboxQuickCommand + ` (name);`,
				},
			},
		},

		/** 给工具箱添加 快速命令 结束 **/

		// 工具表 打开 添加顺序号
		{
			Version: "1.0.3",
			Module:  ModuleToolbox,
			Stage:   `工具箱[` + TableToolboxOpen + `]添加顺序号[sequence]`,
			Sql: &install.StageSqlModel{
				Mysql: []string{
					`ALTER TABLE ` + TableToolboxOpen + ` ADD COLUMN sequence bigint(20) DEFAULT NULL COMMENT '顺序号' AFTER extend;`,
					`ALTER TABLE ` + TableToolboxOpen + ` ADD INDEX ` + TableToolboxOpen + `_index_sequence (sequence);`,
				},
				Sqlite: []string{
					`ALTER TABLE ` + TableToolboxOpen + ` ADD sequence bigint(20);`,
					`CREATE INDEX ` + TableToolboxOpen + `_index_sequence on ` + TableToolboxOpen + ` (sequence);`,
				},
			},
		},
		/** 工具表添加顺序号 结束 **/

		/** 工具表 添加 可见性 开始**/
		{
			Version: "1.0.4",
			Module:  ModuleToolbox,
			Stage:   `工具箱[` + TableToolbox + `]添加可见性[visibility]`,
			Sql: &install.StageSqlModel{
				Mysql: []string{
					`ALTER TABLE ` + TableToolbox + ` ADD COLUMN visibility int(10) DEFAULT '2' COMMENT '可见性' AFTER sequence;`,
					`ALTER TABLE ` + TableToolbox + ` ADD INDEX ` + TableToolbox + `_index_visibility (visibility);`,
					`UPDATE ` + TableToolbox + ` SET visibility=2 WHERE visibility=0 OR visibility IS NULL;`,
				},
				Sqlite: []string{
					`ALTER TABLE ` + TableToolbox + ` ADD visibility int(10) DEFAULT '2';`,
					`CREATE INDEX ` + TableToolbox + `_index_visibility on ` + TableToolbox + ` (visibility);`,
					`UPDATE ` + TableToolbox + ` SET visibility=2 WHERE visibility=0 OR visibility IS NULL;`,
				},
			},
		},
		/** 工具表 添加 可见性 结束**/

		/** 给工具箱添加 扩展 开始 **/

		// 创建工具箱 扩展 表
		{
			Version: "1.0.3",
			Module:  ModuleToolbox,
			Stage:   `创建表[` + TableToolboxExtend + `]`,
			Sql: &install.StageSqlModel{
				Mysql: []string{`
CREATE TABLE ` + TableToolboxExtend + ` (
	extendId bigint(20) NOT NULL COMMENT 'ID',
	toolboxId bigint(20) NOT NULL COMMENT '工具箱ID',
	extendType varchar(50) NOT NULL COMMENT '扩展类型',
	name varchar(50) NOT NULL COMMENT '名称',
	value text DEFAULT NULL COMMENT '值',
	userId bigint(20) NOT NULL COMMENT '用户ID',
	createTime datetime NOT NULL COMMENT '创建时间',
	updateTime datetime DEFAULT NULL COMMENT '修改时间',
	PRIMARY KEY (extendId),
	KEY index_toolboxId (toolboxId),
	KEY index_extendType (extendType),
	KEY index_userId (userId)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='` + TableToolboxExtendComment + `';
`},
				Sqlite: []string{`
CREATE TABLE ` + TableToolboxExtend + ` (
	extendId bigint(20) NOT NULL,
	toolboxId bigint(20) NOT NULL,
	extendType varchar(50) NOT NULL,
	name varchar(50) NOT NULL,
	value text DEFAULT NULL,
	userId bigint(20) NOT NULL,
	createTime datetime NOT NULL,
	updateTime datetime DEFAULT NULL,
	PRIMARY KEY (extendId)
);
`,
					`CREATE INDEX ` + TableToolboxExtend + `_index_toolboxId on ` + TableToolboxExtend + ` (toolboxId);`,
					`CREATE INDEX ` + TableToolboxExtend + `_index_extendType on ` + TableToolboxExtend + ` (extendType);`,
					`CREATE INDEX ` + TableToolboxExtend + `_index_userId on ` + TableToolboxExtend + ` (userId);`,
				},
			},
		},

		/** 给工具箱添加 扩展 结束 **/

		/** 工具表 添加顺序号 开始 **/
		{
			Version: "1.0.4",
			Module:  ModuleToolbox,
			Stage:   `工具箱[` + TableToolbox + `]添加顺序号[sequence]`,
			Sql: &install.StageSqlModel{
				Mysql: []string{
					`ALTER TABLE ` + TableToolbox + ` ADD COLUMN sequence int(10) DEFAULT NULL COMMENT '顺序号';`,
				},
				Sqlite: []string{
					`ALTER TABLE ` + TableToolbox + ` ADD sequence int(10);`,
				},
			},
		},
		/** 工具表 添加顺序号 结束 **/

		/** 工具表 分组 添加顺序号 开始 **/
		{
			Version: "1.0.4",
			Module:  ModuleToolbox,
			Stage:   `工具箱[` + TableToolboxGroup + `]添加顺序号[sequence]`,
			Sql: &install.StageSqlModel{
				Mysql: []string{
					`ALTER TABLE ` + TableToolboxGroup + ` ADD COLUMN sequence int(10) DEFAULT NULL COMMENT '顺序号';`,
				},
				Sqlite: []string{
					`ALTER TABLE ` + TableToolboxGroup + ` ADD sequence int(10);`,
				},
			},
		},
		/** 工具表 分组 添加顺序号 结束 **/

		/** 工具表 分组 添加 父ID 开始 **/
		{
			Version: "1.0.4",
			Module:  ModuleToolbox,
			Stage:   `工具箱[` + TableToolboxGroup + `]添加顺序号[parentId]`,
			Sql: &install.StageSqlModel{
				Mysql: []string{
					`ALTER TABLE ` + TableToolboxGroup + ` ADD COLUMN parentId int(10) DEFAULT NULL COMMENT '父ID';`,
				},
				Sqlite: []string{
					`ALTER TABLE ` + TableToolboxGroup + ` ADD parentId int(10);`,
				},
			},
		},
		/** 工具表 分组 添加 父ID 结束 **/
	}

}
