
-- +migrate Up
CREATE TABLE `exchange_holiday` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `exchange_code` VARCHAR(32) NOT NULL '代號',
    `name` VARCHAR(32) NOT NULL COMMENT '名稱',
    `date` DATE NOT NULL COMMENT '開始日期',
    `end_date` DATE NULL DEFAULT NULL COMMENT '結束日期',
    `type` ENUM('none', 'fullday', 'halfday') NOT NULL COMMENT '交易所休假日類別',
    `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新時間',
    `half_day_close_time` TIMESTAMP NULL DEFAULT NULL COMMENT '半天假日結束交易時間',
    `memo` VARCHAR(128) NULL DEFAULT NULL COMMENT '備註',
    PRIMARY KEY (`id`),
    UNIQUE INDEX (`exchange_code`, `date`),
    FOREIGN KEY (`exchange_code`) REFERENCES exchange(`code`) ON DELETE CASCADE
) AUTO_INCREMENT=1 CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='交易所休假日';


-- +migrate Down
SET FOREIGN_KEY_CHECKS=0;
DROP TABLE IF EXISTS `exchange`;
