CREATE DATABASE :db_name;
\c :db_name;

CREATE TABLE artists (
  id integer GENERATED ALWAYS AS IDENTITY
  , did integer NOT NULL
  , start_pos bigint
  , end_pos bigint
  , PRIMARY KEY (id)
  , UNIQUE ( did )
);

CREATE TABLE releases (
  id integer GENERATED ALWAYS AS IDENTITY
  , did integer NOT NULL
  , start_pos bigint
  , end_pos bigint
  , PRIMARY KEY (id)
  , UNIQUE ( did )
);

CREATE TABLE labels (
  id integer GENERATED ALWAYS AS IDENTITY
  , did integer NOT NULL
  , start_pos bigint
  , end_pos bigint
  , PRIMARY KEY (id)
  , UNIQUE ( did )
);

CREATE TABLE masters (
  id integer GENERATED ALWAYS AS IDENTITY
  , did integer NOT NULL
  , start_pos bigint
  , end_pos bigint
  , PRIMARY KEY (id)
  , UNIQUE ( did )
);

CREATE TABLE artist_releases (
  id integer GENERATED ALWAYS AS IDENTITY
  , artist_did integer NOT NULL
  , release_did integer NOT NULL
  , role smallint NOT NULL
  , PRIMARY KEY (id)
);

CREATE TABLE label_releases (
  id integer GENERATED ALWAYS AS IDENTITY
  , label_did integer NOT NULL
  , release_did integer NOT NULL
  , PRIMARY KEY (id)
);

-- CREATE INDEX idx_artist_releases_artist ON artist_releases(artist_did);
-- CREATE INDEX idx_artist_releases_release ON artist_releases(release_did);

-- CREATE INDEX idx_label_releases_label ON label_releases(label_did);
-- CREATE INDEX idx_label_releases_release ON label_releases(release_did);
