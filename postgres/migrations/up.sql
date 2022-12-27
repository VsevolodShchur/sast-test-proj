CREATE TABLE IF NOT EXISTS items (
    id serial NOT NULL,
    title varchar(64) NOT NULL,
    CONSTRAINT form_pkey PRIMARY KEY (id)
);

INSERT INTO items (id, title) VALUES
    (1, 'Item 1'),
    (2, 'Item 2'),
    (3, 'Item 3');
