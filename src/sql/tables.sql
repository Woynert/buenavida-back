
USE buenavida;

-- drop

DROP TABLE IF EXISTS sales_items;
DROP TABLE IF EXISTS sales;

-- create

CREATE TABLE IF NOT EXISTS sales
(
	id          SERIAL PRIMARY KEY,
	id_client   CHAR(36) NOT NULL, -- uuid lenght in mongodb

	iva         INT    NOT NULL DEFAULT 5,
	subtotal    FLOAT8 NOT NULL DEFAULT 0,
	total       FLOAT8 NOT NULL DEFAULT 0,
	created     TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS sales_items
(
	id          SERIAL PRIMARY KEY,
	id_sale     INT      NOT NULL,
	id_product  CHAR(36) NOT NULL, -- mongodb uuid

	quantity    INT    NOT NULL,
	discount    INT    NOT NULL DEFAULT 0,
	price_base  BIGINT NOT NULL,
	price_final FLOAT8 NOT NULL,

	FOREIGN KEY (id_sale) REFERENCES sales(id)
);
