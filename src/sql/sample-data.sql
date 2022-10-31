
-- sample data

DO $$
DECLARE
    id_sale INT;
BEGIN

    CALL sales_new('asd12-5406-4506das', id_sale);
    CALL sales_add_item (id_sale, '10155-5406-4506das', 2, 0, 1050 );
    CALL sales_add_item (id_sale, '11155-5406-4506das', 1, 50, 1000 );
    CALL sales_add_item (id_sale, '12155-5406-4506das', 5, 20, 1000 );

    CALL sales_new('abc12-5406-4506das', id_sale);
    CALL sales_add_item (id_sale, '10155-5406-4506das', 2, 0, 1000 );
    CALL sales_add_item (id_sale, '11155-5406-4506das', 1, 20, 1000 );
    
END;
$$;
