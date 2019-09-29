DELIMITER $$

DROP FUNCTION IF EXISTS parseRecipe$$

CREATE FUNCTION parseRecipe (v1 VARCHAR (255), material INT)
RETURNS TEXT
BEGIN
	DECLARE i INT;
	SET i = 0;
    
    DROP TEMPORARY TABLE IF EXISTS loc_table;
    
    CREATE TEMPORARY TABLE loc_table(
        id int,
        json text
    );
    
	WHILE i < JSON_LENGTH(v1) DO
    	INSERT INTO loc_table (id, json) VALUES (i, JSON_EXTRACT(v1, CONCAT_WS( '', '$[', i, ']' )) );
        SET i = i + 1;
    END WHILE;

    RETURN (SELECT SUM(`json`->>'$.recipe_amount') FROM `loc_table` WHERE `json`->>'$.material_id' = material);
END$$

DELIMITER ;


SELECT
    B.`name`,
    @saldo_start :=(
        SELECT
            CAST( SUM(
                `weight_with_load` - `weight_without_load`
            ) AS DECIMAL(38, 2))
        FROM
            `warehouse`
        WHERE
            `zombie` <> TRUE 
        	AND date_create >= '1348891200' 
        	AND date_create <= '1569729600' AND `id` = A.`id`
        GROUP BY
            `id`
    ) AS saldo_start,
    (
    	SELECT
        	SUM( parseRecipe(`calculated_data`->>'$.recipe', A.`materials`) ) as 'count'
        FROM `shipment`
    )
    
FROM
    `warehouse` AS A
INNER JOIN `materials` AS B
ON
    B.`id` = A.`materials`
WHERE
    A.`zombie` <> TRUE
