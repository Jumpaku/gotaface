-- classDiagram
--     E_1 <|-- E_2
--     E_2 <|-- E_1
CREATE TABLE "E_1" (
    "PK_11" integer NOT NULL,
    "PK_12" integer NOT NULL,
    PRIMARY KEY ("PK_11", "PK_12")
);

CREATE TABLE "E_2" (
    "PK_21" integer NOT NULL,
    "PK_22" integer NOT NULL,
    CONSTRAINT "FK_E_2_1" FOREIGN KEY ("PK_21", "PK_22") REFERENCES "E_1" ("PK_11", "PK_12"),
    PRIMARY KEY ("PK_21", "PK_22")
);

ALTER TABLE "E_1" ADD CONSTRAINT "FK_E_1_2" FOREIGN KEY ("PK_11", "PK_12") REFERENCES "E_2" ("PK_21", "PK_22");
