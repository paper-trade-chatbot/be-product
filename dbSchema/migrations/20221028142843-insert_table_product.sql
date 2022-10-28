
-- +migrate Up
INSERT INTO 
	`product`(
        `id`,
        `type`,
        `exchange_id`,
        `exchange_code`,
        `code`,
        `name`,
        `status`,
        `display`,
        `currency_code`,
        `tick_unit`)
VALUES
	(
        '1',
        '1',
        '2',
        'TWSE',
        '2330',
        '台灣積體電路製造股份有限公司',
        '1',
        '1',
        'NTD',
        '5'
    );


-- +migrate Down
SET SQL_SAFE_UPDATES = 0;
DELETE FROM `product`;
SET SQL_SAFE_UPDATES = 1;
