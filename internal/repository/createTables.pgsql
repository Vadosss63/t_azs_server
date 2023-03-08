
create table if not exists users
(
    id      bigint primary key generated always as identity,
    login   varchar(200) not null unique,
    hashed_password varchar(200) not null,
    name    varchar(200) not null,
    surname varchar(200) not null
);

DROP TABLE azses;
create table if not exists azses
(
 	id  bigint primary key generated always as identity,
	id_azs  bigint,
	id_user int,
    is_authorized  int,
    count_colum int,	
 	time   varchar(100) not null,
 	name    varchar(100) not null,
 	address varchar(100) not null,
	stats varchar(1000) not null
    -- FOREIGN KEY (id_user) REFERENCES users(id)
);