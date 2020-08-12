CREATE TABLE problems
(
    id         serial PRIMARY KEY,
    name       VARCHAR(50),
    credits    VARCHAR(50),
    body       text,
    difficulty VARCHAR(10),
    tags       VARCHAR(50)[]
);

---- create above / drop below ----

DROP TABLE problems;
