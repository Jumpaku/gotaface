-- classDiagram
--     D_1 <|-- D_1
CREATE TABLE "D_1" (
    "PK_11" integer NOT NULL,
    "PK_12" integer unique NOT NULL,
    CONSTRAINT "FK_D_1_1" FOREIGN KEY ("PK_11") REFERENCES "D_1" ("PK_12"),
    PRIMARY KEY ("PK_11", "PK_12")
);
