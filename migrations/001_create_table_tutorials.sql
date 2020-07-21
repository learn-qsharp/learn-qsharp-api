CREATE TABLE tutorials(
  id serial PRIMARY KEY,
  title VARCHAR (50),
  credits VARCHAR (50),
  description text,
  body text,
  difficulty VARCHAR (10),
  tags VARCHAR (50)[]
);

---- create above / drop below ----

DROP TABLE tutorials;
