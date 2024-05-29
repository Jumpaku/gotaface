CREATE TABLE I
(
    PK INT64 NOT NULL,
    C1 INT64 NOT NULL,
    C2 INT64 NOT NULL,
) PRIMARY KEY (PK);

CREATE INDEX IDX_I_C1Desc_C2Asc ON I (C1 DESC, C2 ASC);
CREATE INDEX IDX_I_C1Asc_C2Desc ON I (C1 ASC, C2 DESC);
CREATE INDEX IDX_I_Storing ON I (C1) STORING (C2);
