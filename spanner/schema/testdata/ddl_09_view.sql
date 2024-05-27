CREATE TABLE J (
    PK INT64 NOT NULL,
    Col1 INT64 NOT NULL,
    Col2 INT64 NOT NULL,
) PRIMARY KEY (PK);

CREATE VIEW K SQL SECURITY INVOKER AS (
    SELECT
        J.PK + 1 AS PK_2,
        J.Col1 * 2 AS Col1_2,
        J.Col2 -3 AS Col2_2
    FROM J
);