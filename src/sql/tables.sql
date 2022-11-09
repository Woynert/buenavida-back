
-- drop

DROP TABLE IF EXISTS sales_items;
DROP TABLE IF EXISTS sales;

-- create

CREATE TABLE IF NOT EXISTS sales
(
	id          SERIAL PRIMARY KEY,
	id_client   CHAR(36) NOT NULL, -- uuid lenght in mongodb

	iva         INT    NOT NULL DEFAULT 19,
	subtotal    NUMERIC(8, 4) NOT NULL DEFAULT 0,
	total       NUMERIC(8, 4) NOT NULL DEFAULT 0,
	created     TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS sales_items
(
	id          SERIAL PRIMARY KEY,
	id_sale     INT      NOT NULL,
	id_product  CHAR(36) NOT NULL, -- mongodb uuid

	quantity    INT    NOT NULL,
	discount    INT    NOT NULL DEFAULT 0,
	price_base  NUMERIC(8, 4) NOT NULL,
	price_final NUMERIC(8, 4) NOT NULL,

	FOREIGN KEY (id_sale) REFERENCES sales(id)
);
