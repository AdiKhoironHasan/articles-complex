CREATE TABLE articles
(
  id BIGSERIAL,
  author varchar(60) NOT NULL,
  title varchar(150) NOT NULL,
  body text NOT NULL,
  created timestamp NOT NULL,
  PRIMARY KEY (id)
);