sqlite3 -header -csv $1 "select * from artists_search_info;" > $2/artists_search_info.csv
gzip -k $2/artists_search_info.csv