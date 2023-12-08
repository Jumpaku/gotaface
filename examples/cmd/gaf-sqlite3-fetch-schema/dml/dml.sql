CREATE TABLE A(

    PK INT64 NOT NULL,

    Col_01 BOOL ,

    Col_02 BOOL NOT NULL,

    Col_03 BYTES(50) ,

    Col_04 BYTES(50) NOT NULL,

    Col_05 DATE ,

    Col_06 DATE NOT NULL,

    Col_07 FLOAT64 ,

    Col_08 FLOAT64 NOT NULL,

    Col_09 INT64 ,

    Col_10 INT64 NOT NULL,

    Col_11 JSON ,

    Col_12 JSON NOT NULL,

    Col_13 NUMERIC ,

    Col_14 NUMERIC NOT NULL,

    Col_15 STRING(50) ,

    Col_16 STRING(50) NOT NULL,

    Col_17 TIMESTAMP ,

    Col_18 TIMESTAMP NOT NULL,



    PRIMARY KEY (
        
        PK
        
    )
);CREATE TABLE C_1(

    PK_11 INT64 NOT NULL,

    PK_12 INT64 NOT NULL,



    PRIMARY KEY (
        
        PK_11
        
        , PK_12
        
    )
);CREATE TABLE C_2(

    PK_21 INT64 NOT NULL,

    PK_22 INT64 NOT NULL,


    FOREIGN KEY (
        
            PK_21
        
            , PK_22
        
    ) C_1 (
        
            PK_11
        
            , PK_12
        
    ),


    PRIMARY KEY (
        
        PK_21
        
        , PK_22
        
    )
);CREATE TABLE C_3(

    PK_31 INT64 NOT NULL,

    PK_32 INT64 NOT NULL,


    FOREIGN KEY (
        
            PK_31
        
            , PK_32
        
    ) C_2 (
        
            PK_21
        
            , PK_22
        
    ),


    PRIMARY KEY (
        
        PK_31
        
        , PK_32
        
    )
);CREATE TABLE C_4(

    PK_41 INT64 NOT NULL,

    PK_42 INT64 NOT NULL,


    FOREIGN KEY (
        
            PK_41
        
            , PK_42
        
    ) C_2 (
        
            PK_21
        
            , PK_22
        
    ),


    PRIMARY KEY (
        
        PK_41
        
        , PK_42
        
    )
);CREATE TABLE C_5(

    PK_51 INT64 NOT NULL,

    PK_52 INT64 NOT NULL,


    FOREIGN KEY (
        
            PK_51
        
            , PK_52
        
    ) C_4 (
        
            PK_41
        
            , PK_42
        
    ),

    FOREIGN KEY (
        
            PK_51
        
            , PK_52
        
    ) C_3 (
        
            PK_31
        
            , PK_32
        
    ),


    PRIMARY KEY (
        
        PK_51
        
        , PK_52
        
    )
);CREATE TABLE D_1(

    PK_11 INT64 NOT NULL,

    PK_12 INT64 NOT NULL,


    FOREIGN KEY (
        
            PK_11
        
    ) D_1 (
        
            PK_12
        
    ),


    PRIMARY KEY (
        
        PK_11
        
        , PK_12
        
    )
);CREATE TABLE E_1(

    PK_11 INT64 NOT NULL,

    PK_12 INT64 NOT NULL,


    FOREIGN KEY (
        
            PK_11
        
            , PK_12
        
    ) E_2 (
        
            PK_21
        
            , PK_22
        
    ),


    PRIMARY KEY (
        
        PK_11
        
        , PK_12
        
    )
);CREATE TABLE E_2(

    PK_21 INT64 NOT NULL,

    PK_22 INT64 NOT NULL,


    FOREIGN KEY (
        
            PK_21
        
            , PK_22
        
    ) E_1 (
        
            PK_11
        
            , PK_12
        
    ),


    PRIMARY KEY (
        
        PK_21
        
        , PK_22
        
    )
);CREATE TABLE F_1(

    PK_11 INT64 NOT NULL,

    PK_12 INT64 NOT NULL,


    FOREIGN KEY (
        
            PK_11
        
            , PK_12
        
    ) F_3 (
        
            PK_31
        
            , PK_32
        
    ),


    PRIMARY KEY (
        
        PK_11
        
        , PK_12
        
    )
);CREATE TABLE F_2(

    PK_21 INT64 NOT NULL,

    PK_22 INT64 NOT NULL,


    FOREIGN KEY (
        
            PK_21
        
            , PK_22
        
    ) F_1 (
        
            PK_11
        
            , PK_12
        
    ),


    PRIMARY KEY (
        
        PK_21
        
        , PK_22
        
    )
);CREATE TABLE F_3(

    PK_31 INT64 NOT NULL,

    PK_32 INT64 NOT NULL,


    FOREIGN KEY (
        
            PK_31
        
            , PK_32
        
    ) F_2 (
        
            PK_21
        
            , PK_22
        
    ),


    PRIMARY KEY (
        
        PK_31
        
        , PK_32
        
    )
);CREATE TABLE G(

    PK INT64 NOT NULL,

    C1 INT64 NOT NULL,

    C2 INT64 NOT NULL,

    C3 INT64 NOT NULL,



    

    

    

    

    

    

    

    

    

    

    

    

    

    

    

    PRIMARY KEY (
        
        PK
        
    )
);

CREATE UNIQUE INDEX UQ_G_C3_C2_C1 ON G(
    
    C3
    
    , C2
    
    , C1
    
);



CREATE UNIQUE INDEX UQ_G_C3_C1_C2 ON G(
    
    C3
    
    , C1
    
    , C2
    
);



CREATE UNIQUE INDEX UQ_G_C2_C1_C3 ON G(
    
    C2
    
    , C1
    
    , C3
    
);



CREATE UNIQUE INDEX UQ_G_C2_C3_C1 ON G(
    
    C2
    
    , C3
    
    , C1
    
);



CREATE UNIQUE INDEX UQ_G_C1_C3_C2 ON G(
    
    C1
    
    , C3
    
    , C2
    
);



CREATE UNIQUE INDEX UQ_G_C1_C2_C3 ON G(
    
    C1
    
    , C2
    
    , C3
    
);



CREATE UNIQUE INDEX UQ_G_C1_C3 ON G(
    
    C1
    
    , C3
    
);



CREATE UNIQUE INDEX UQ_G_C3_C1 ON G(
    
    C3
    
    , C1
    
);



CREATE UNIQUE INDEX UQ_G_C3_C2 ON G(
    
    C3
    
    , C2
    
);



CREATE UNIQUE INDEX UQ_G_C2_C3 ON G(
    
    C2
    
    , C3
    
);



CREATE UNIQUE INDEX UQ_G_C2_C1 ON G(
    
    C2
    
    , C1
    
);



CREATE UNIQUE INDEX UQ_G_C1_C2 ON G(
    
    C1
    
    , C2
    
);



CREATE UNIQUE INDEX UQ_G_C3 ON G(
    
    C3
    
);



CREATE UNIQUE INDEX UQ_G_C2 ON G(
    
    C2
    
);



CREATE UNIQUE INDEX UQ_G_C1 ON G(
    
    C1
    
);

CREATE TABLE H(

    PK INT64 NOT NULL,

    C1 INT64 NOT NULL,

    C2 INT64 NOT NULL,

    C3 INT64 NOT NULL,



    
    UNIQUE (
        
        C3
        
        , C2
        
        , C1
        
    ),
    

    
    UNIQUE (
        
        C3
        
        , C1
        
        , C2
        
    ),
    

    
    UNIQUE (
        
        C2
        
        , C1
        
        , C3
        
    ),
    

    
    UNIQUE (
        
        C2
        
        , C3
        
        , C1
        
    ),
    

    
    UNIQUE (
        
        C1
        
        , C3
        
        , C2
        
    ),
    

    
    UNIQUE (
        
        C1
        
        , C2
        
        , C3
        
    ),
    

    
    UNIQUE (
        
        C1
        
        , C3
        
    ),
    

    
    UNIQUE (
        
        C3
        
        , C1
        
    ),
    

    
    UNIQUE (
        
        C3
        
        , C2
        
    ),
    

    
    UNIQUE (
        
        C2
        
        , C3
        
    ),
    

    
    UNIQUE (
        
        C2
        
        , C1
        
    ),
    

    
    UNIQUE (
        
        C1
        
        , C2
        
    ),
    

    
    UNIQUE (
        
        C3
        
    ),
    

    
    UNIQUE (
        
        C2
        
    ),
    

    
    UNIQUE (
        
        C1
        
    ),
    

    PRIMARY KEY (
        
        PK
        
    )
);





























CREATE TABLE I(

    PK INT64 NOT NULL,

    C1 INT64 NOT NULL,

    C2 INT64 NOT NULL,

    C3 INT64 NOT NULL,



    
    UNIQUE (
        
        C3
        
    ),
    

    
    UNIQUE (
        
        C2
        
    ),
    

    
    UNIQUE (
        
        C1
        
    ),
    

    PRIMARY KEY (
        
        PK
        
    )
);





