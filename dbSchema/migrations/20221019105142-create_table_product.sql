
-- +migrate Up
CREATE TABLE `product` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `type` TINYINT(4) NOT NULL COMMENT '產品種類 1:stock, 2:crypto, 3:forex, 4:futures',
    `exchange_code` VARCHAR(32) NOT NULL COMMENT '交易所代號',
    `code` VARCHAR(32) NOT NULL COMMENT '代號',
    `name` VARCHAR(32) NOT NULL COMMENT '名稱',
    `status` TINYINT(4) NOT NULL COMMENT '狀態 1:enabled, 2:disabled',
    `display` TINYINT(4) NOT NULL COMMENT '顯示 1:enabled, 2:disabled',
    `currency_code` VARCHAR(32) NOT NULL COMMENT '貨幣代號',
    `tick_unit` DECIMAL(4,2) NOT NULL COMMENT '報價間隔',
    `minimum_order` DECIMAL(15,10) NULL DEFAULT NULL COMMENT '最小購買單位',
    `icon_id` VARCHAR(128) NULL DEFAULT NULL COMMENT '標誌編號',
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '創建時間',
    `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新時間',

    PRIMARY KEY (`id`),
    UNIQUE INDEX (`exchange_code`, `code`),
    FOREIGN KEY (`exchange_code`) REFERENCES exchange(`code`) ON DELETE CASCADE
) AUTO_INCREMENT=1 CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='產品';


-- +migrate Down
SET FOREIGN_KEY_CHECKS=0;
DROP TABLE IF EXISTS `product`;
