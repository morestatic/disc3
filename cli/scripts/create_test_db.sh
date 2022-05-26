psql --set=db_name="$1" -f indexes/sql/drop_db.sql
psql --set=db_name="$1" -f indexes/sql/init_db.sql
