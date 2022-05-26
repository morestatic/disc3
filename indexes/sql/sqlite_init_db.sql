DROP TABLE IF EXISTS artists;
CREATE TABLE artists (
  id integer PRIMARY KEY
  , did integer NOT NULL
  , start_pos integer
  , end_pos integer
  , UNIQUE( did )
);

DROP TABLE IF EXISTS releases;
CREATE TABLE releases (
  id integer PRIMARY KEY
  , did integer NOT NULL
  , start_pos bigint
  , end_pos bigint
  , UNIQUE( did )
);

DROP TABLE IF EXISTS labels;
CREATE TABLE labels (
  id integer PRIMARY KEY
  , did integer NOT NULL
  , start_pos bigint
  , end_pos bigint
  , UNIQUE( did )
);

DROP TABLE IF EXISTS masters;
CREATE TABLE masters (
  id integer PRIMARY KEY
  , did integer NOT NULL
  , start_pos bigint
  , end_pos bigint
  , UNIQUE( did )
);

DROP TABLE IF EXISTS artist_releases;
CREATE TABLE artist_releases (
  artist_did integer NOT NULL
  , release_did integer NOT NULL
  , role smallint NOT NULL
);

DROP TABLE IF EXISTS label_releases;
CREATE TABLE label_releases (
  id integer PRIMARY KEY
  , label_did integer NOT NULL
  , release_did integer NOT NULL
);

DROP TABLE IF EXISTS stats;
CREATE TABLE stats (
  property text NOT NULL PRIMARY KEY
  , value text
  , UNIQUE( property )
);

DROP INDEX IF EXISTS idx_artist_releases_artist;
CREATE INDEX idx_artist_releases_artist ON artist_releases(artist_did);

DROP TABLE IF EXISTS idx_artist_releases_artist_role;
CREATE INDEX idx_artist_releases_artist_role ON artist_releases(artist_did, role);

-- CREATE INDEX idx_label_releases_label ON label_releases(label_did);
-- CREATE INDEX idx_label_releases_release ON label_releases(release_did);

drop table if exists artists_search_info;
CREATE TABLE artists_search_info (
  artist_did integer PRIMARY KEY
  , name text
  , realname text
  , is_group integer NOT NULL
);

drop table if exists release_genres;
CREATE TABLE release_genres (
  genre string PRIMARY KEY
);

drop table if exists release_styles;
CREATE TABLE release_styles (
  style string PRIMARY KEY
);

-- select * from artists_search_info limit 10;
-- .header
-- .header on
-- .mode csv
-- .output artists_search_info.csv
-- select * from artists_search_info;
-- .quit
-- .output artists_search_info.sql
-- .dump artists_search_info
-- .quit
