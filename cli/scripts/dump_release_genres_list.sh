sqlite3 -header -csv $1 "select genre from release_genres;" > $2/release_genres.csv
# gzip -k $2/release_styles.csv