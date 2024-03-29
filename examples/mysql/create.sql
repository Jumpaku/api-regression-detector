SET FOREIGN_KEY_CHECKS = 0;
DROP TABLE IF EXISTS example_table;
CREATE TABLE example_table (
    id integer auto_increment,
    c0 text,
    c1 integer,
    c2 boolean,
    c3 datetime,
    PRIMARY KEY(id)
);
DROP TABLE IF EXISTS child_example_table_1;
CREATE TABLE child_example_table_1 (
    id integer auto_increment,
    example_table_id integer,
    PRIMARY KEY(id),
    FOREIGN KEY (example_table_id) REFERENCES example_table (id)
);
DROP TABLE IF EXISTS child_example_table_2;
CREATE TABLE child_example_table_2 (
    id integer auto_increment,
    example_table_id integer,
    PRIMARY KEY(id),
    FOREIGN KEY (example_table_id) REFERENCES example_table (id)
);
SET FOREIGN_KEY_CHECKS = 1;