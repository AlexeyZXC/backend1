CREATE TABLE links (
	shorturl serial  PRIMARY KEY,
	longurl varchar
);
CREATE TABLE stat (
	id serial PRIMARY KEY,
	shorturl int,
    userip varchar,
    passtime varchar
);