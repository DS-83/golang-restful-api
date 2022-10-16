DROP DATABASE IF EXISTS photogramm;
CREATE DATABASE photogramm;

USE photogramm;

CREATE TABLE users (
  id INTEGER PRIMARY KEY AUTO_INCREMENT,
  username varchar(255) UNIQUE NOT NULL,
  password varchar(255) NOT NULL
  );


CREATE TABLE albums (
  id INTEGER PRIMARY KEY AUTO_INCREMENT,
  album_name varchar(255) NOT NULL,
  user_id INTEGER,
FOREIGN KEY (user_id)
      REFERENCES users(id)
);

CREATE TABLE photos (
    id VARCHAR(255) PRIMARY KEY UNIQUE NOT NULL,
    album_id INTEGER,
    user_id INTEGER,
    FOREIGN KEY (user_id)
      REFERENCES users(id),
    FOREIGN KEY (album_id)
      REFERENCES albums(id)
  );

  CREATE TABLE revoked_tokens (
    jwt_token_id varchar(255) PRIMARY KEY,
    revoketion_date timestamp default now()
  );