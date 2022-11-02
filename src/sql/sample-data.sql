
-- sample data

DO $$
DECLARE
    id_sale INT;
BEGIN

    SELECT * INTO id_sale
	FROM sales_new('asd12-5406-4506das');
	CALL sales_add_item (id_sale, '10155-5406-4506das', 2, 0, 50.4 );
	CALL sales_add_item (id_sale, '11155-5406-4506das', 1, 50, 10.3 );
	CALL sales_add_item (id_sale, '12155-5406-4506das', 5, 20, 15 );

    SELECT * INTO id_sale
	FROM sales_new('asd12-5406-4506das');
	CALL sales_add_item (id_sale, '10155-5406-4506das', 2, 0, 20.5);
	CALL sales_add_item (id_sale, '11155-5406-4506das', 1, 20, 5.5);
    
END;
$$;
