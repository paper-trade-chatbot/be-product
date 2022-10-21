
-- +migrate Up
CREATE TABLE `exchange` (
    `code` VARCHAR(32) NOT NULL COMMENT '代號',
    `product_type` TINYINT(4) NOT NULL COMMENT '產品種類 1:stock, 2:crypto, 3:forex, 4:futures',
    `name` VARCHAR(32) NOT NULL COMMENT '名稱',
    `status` TINYINT(4) NOT NULL COMMENT '狀態 1:enabled, 2:disabled',
    `display` TINYINT(4) NOT NULL COMMENT '顯示 1:enabled, 2:disabled',
    `country_code` VARCHAR(16) NOT NULL COMMENT '國家代碼',
    `timezone_offset` DECIMAL(4,2) NOT NULL COMMENT '時區時差',
    `open_time` TIMESTAMP NULL COMMENT '每日開始交易時間',
    `close_time` TIMESTAMP NULL COMMENT '每日結束交易時間',
    `exchange_day` VARCHAR(128) NOT NULL COMMENT '星期幾交易',
    `exception_time` VARCHAR(256) NOT NULL COMMENT '例外交易時間',
    `daylight_saving` BOOLEAN NOT NULL DEFAULT false COMMENT '是否日光節約',
    `location` VARCHAR(32) NOT NULL COMMENT '地區',
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '創建時間',
    `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新時間',
    PRIMARY KEY (`code`)
) AUTO_INCREMENT=1 CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='交易所';


-- +migrate Down
SET FOREIGN_KEY_CHECKS=0;
DROP TABLE IF EXISTS `exchange`;
