-- 创建数据库
CREATE DATABASE IF NOT EXISTS seckill DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

USE seckill;

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
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='优惠券表';

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
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='订单表';

-- 插入测试数据
INSERT INTO coupons (name, description, total_stock, remain_stock, start_time, end_time, status)
VALUES
('双十一优惠券', '满100减50', 1000, 1000, NOW(), DATE_ADD(NOW(), INTERVAL 7 DAY), 1),
('新用户专享', '满50减20', 500, 500, NOW(), DATE_ADD(NOW(), INTERVAL 30 DAY), 1);
