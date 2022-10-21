
-- +migrate Up
ALTER TABLE `exchange` AUTO_INCREMENT = 1;
INSERT INTO 
	`exchange`(
        `code`,
        `product_type`,
        `name`,
        `status`,
        `display`,
        `country_code`,
        `timezone_offset`,
        `open_time`,
        `close_time`,
        `exchange_day`,
        `exception_time`,
        `daylight_saving`,
        `location`)
VALUES
	('TEST1',
    1,
    'TEST1',
    1,
    1,
    'TW',
    8.0,
    null,
    null,
    '{"startDay":0,"endDay":6}',
    '{"trade":[],"stopTrade":[]}',
    false,
    'Asia/Taipei'
    );

-- +migrate Down
SET SQL_SAFE_UPDATES = 0;
DELETE FROM `exchange`;
SET SQL_SAFE_UPDATES = 1;
