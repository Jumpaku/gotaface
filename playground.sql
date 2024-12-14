--sql query primary key information
SELECT
    tc.constraint_name AS "Name",
    kcu.column_name AS "ColumnName"
FROM information_schema.table_constraints AS tc
     JOIN information_schema.key_column_usage AS kcu
          ON kcu.constraint_name = tc.constraint_name
WHERE kcu.table_name = 'G' AND tc.constraint_type = 'UNIQUE'
ORDER BY tc.constraint_name, kcu.ordinal_position;
