CREATE TABLE `activity` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '编号',
  `createTime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updateTime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `name` varchar(255) NOT NULL DEFAULT '' COMMENT '活动名称',
  `status` tinyint(4) NOT NULL DEFAULT '0' COMMENT '活动状态：0、正常；1、售罄',
  `startTime` timestamp NULL DEFAULT NULL COMMENT '开始时间',
  `endTime` timestamp NULL DEFAULT NULL COMMENT '结束时间',
  `total` int(11) NOT NULL DEFAULT '0' COMMENT '商品默认数量',
  `secSpeed` int(11) NOT NULL DEFAULT '0' COMMENT '商品默认最大被购买速率（件/每秒）',
  `buyLimit` int(11) NOT NULL DEFAULT '0' COMMENT '商品默认最大购买数量',
  `buyRate` int(11) NOT NULL DEFAULT '50' COMMENT '商品默认买到的概率（80：80%）',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=13 DEFAULT CHARSET=utf8mb4 COMMENT='活动表';

CREATE TABLE `product` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '编号',
  `createTime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updateTime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `name` varchar(255) NOT NULL DEFAULT '' COMMENT '商品名',
  `shortName` varchar(255) NOT NULL DEFAULT '' COMMENT '商品简称',
  `area` varchar(255) NOT NULL DEFAULT '' COMMENT '商品产地',
  `total` int(11) NOT NULL DEFAULT '0' COMMENT '商品数量',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=9 DEFAULT CHARSET=utf8mb4 COMMENT='商品表';

CREATE TABLE `activity_product` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '编号',
  `createTime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updateTime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `activityId` int(11) NOT NULL DEFAULT '0' COMMENT '活动id',
  `productId` int(11) NOT NULL DEFAULT '0' COMMENT '商品id',
  `status` tinyint(4) NOT NULL DEFAULT '0' COMMENT '商品状态：0、正常；1、售罄；2、强制售罄',
  `total` int(11) NOT NULL DEFAULT '0' COMMENT '商品数量',
  `secSpeed` int(11) NOT NULL DEFAULT '0' COMMENT '商品最大被购买速率（件/每秒）0、不限制',
  `buyLimit` int(11) NOT NULL DEFAULT '0' COMMENT '商品最大购买数量 0、不限制',
  `buyRate` int(11) NOT NULL DEFAULT '0' COMMENT '商品最大抢购速率（次/每秒）0、不限制',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=8 DEFAULT CHARSET=utf8mb4 COMMENT='活动商品表';
