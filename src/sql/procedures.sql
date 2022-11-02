
-- drop

DROP FUNCTION IF EXISTS sales_new;
DROP PROCEDURE IF EXISTS sales_add_item;
DROP PROCEDURE IF EXISTS sales_update;

-- create sale

CREATE FUNCTION sales_new (
	IN  in_id_client CHAR(36)
)
RETURNS INTEGER
LANGUAGE PLPGSQL    
AS $$
DECLARE
	p_id_new_sale INT;
BEGIN

	-- insert and return sale id

	INSERT INTO sales ( id_client )
	VALUES ( in_id_client )
	RETURNING id
	INTO p_id_new_sale;

	RETURN p_id_new_sale;

END;$$;

-- add item to sale

CREATE PROCEDURE sales_add_item (
	IN ar_id_sale INT,
	IN ar_id_product CHAR(36),
	IN ar_quantity INT,
	IN ar_discount INT,
	IN ar_price_base BIGINT
)
LANGUAGE PLPGSQL    
AS $$
DECLARE
	p_total FLOAT8;
BEGIN

	p_total := ar_price_base * (1.0 - ar_discount/100.0) * ar_quantity;

	-- insert

	INSERT INTO sales_items (
		id_sale,
		id_product,
		quantity,
		discount,
		price_base,
		price_final
	)
	VALUES (
		ar_id_sale,
		ar_id_product,
		ar_quantity,
		ar_discount,
		ar_price_base,
		p_total
	);

	-- update total

	CALL sales_update(ar_id_sale);
    COMMIT;
END;$$;

-- sales update total

CREATE PROCEDURE sales_update (
	IN  in_id_sale INT
)
LANGUAGE PLPGSQL    
AS $$
DECLARE
	p_subtotal FLOAT8;
BEGIN

	-- get subtotal
	SELECT SUM(price_final)
	FROM sales_items
	WHERE id_sale = in_id_sale
	INTO p_subtotal
	;

	-- set total
	UPDATE sales
	SET
		subtotal = p_subtotal,
		total    = p_subtotal * (1.0 + (sales.iva / 100.0))
	WHERE id = in_id_sale
	;

    COMMIT;
END;$$;

