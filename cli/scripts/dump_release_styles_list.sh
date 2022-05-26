sqlite3 -header -csv $1 "select style from release_styles;" > $2/release_styles.csv
# gzip -k $2/release_styles.csv