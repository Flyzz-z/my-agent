-- 创建知识库主数据库
CREATE DATABASE IF NOT EXISTS knowledge_base DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

USE knowledge_base;

-- ==========================================
-- 秒杀模块表
-- ==========================================

-- 优惠券表
CREATE TABLE IF NOT EXISTS coupons (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(255) NOT NULL COMMENT '优惠券名称',
    description TEXT COMMENT '优惠券描述',
    total_stock BIGINT NOT NULL DEFAULT 0 COMMENT '总库存',
    remain_stock BIGINT NOT NULL DEFAULT 0 COMMENT '剩余库存',
    start_time TIMESTAMP NOT NULL COMMENT '开始时间',
    end_time TIMESTAMP NOT NULL COMMENT '结束时间',
    status TINYINT NOT NULL DEFAULT 0 COMMENT '状态: 0-未开始, 1-进行中, 2-已结束',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_status (status),
    INDEX idx_time (start_time, end_time)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='秒杀-优惠券表';

-- 订单表
CREATE TABLE IF NOT EXISTS orders (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    user_id BIGINT NOT NULL COMMENT '用户ID',
    coupon_id BIGINT NOT NULL COMMENT '优惠券ID',
    status TINYINT NOT NULL DEFAULT 0 COMMENT '状态: 0-待支付, 1-已支付, 2-已取消',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_user_id (user_id),
    INDEX idx_coupon_id (coupon_id),
    INDEX idx_status (status),
    UNIQUE KEY uk_user_coupon (user_id, coupon_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='秒杀-订单表';

-- 补偿任务表（MQ失败补偿，预留）
CREATE TABLE IF NOT EXISTS compensation_tasks (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    user_id BIGINT NOT NULL COMMENT '用户ID',
    coupon_id BIGINT NOT NULL COMMENT '优惠券ID',
    order_id BIGINT NOT NULL COMMENT '订单ID',
    status TINYINT NOT NULL DEFAULT 0 COMMENT '状态: 0-待处理, 1-处理中, 2-已完成, -1-失败',
    retry_count INT NOT NULL DEFAULT 0 COMMENT '重试次数',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_status (status),
    INDEX idx_created_at (created_at),
    INDEX idx_retry_count (retry_count)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='秒杀-补偿任务表';

-- ==========================================
-- 其他模块表（预留）
-- ==========================================

-- AI搜索模块可能需要的表（如果不用 ES，可以用 MySQL 做元数据存储）
-- CREATE TABLE IF NOT EXISTS documents (...);

-- 用户模块表（预留）
-- CREATE TABLE IF NOT EXISTS users (...);

-- ==========================================
-- 测试数据（仅开发环境）
-- ==========================================

-- 插入测试优惠券
INSERT INTO coupons (name, description, total_stock, remain_stock, start_time, end_time, status)
VALUES
('双十一优惠券', '满100减50', 1000, 1000, NOW(), DATE_ADD(NOW(), INTERVAL 7 DAY), 1),
('新用户专享', '满50减20', 500, 500, NOW(), DATE_ADD(NOW(), INTERVAL 30 DAY), 1),
('限时秒杀', '全场5折', 100, 100, NOW(), DATE_ADD(NOW(), INTERVAL 1 DAY), 1);
