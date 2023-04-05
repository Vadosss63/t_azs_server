DROP TABLE users;

create table if not exists users
(
    id      bigint primary key generated always as identity,
    login   varchar(200) not null unique,
    hashed_password varchar(200) not null,
    name    varchar(200) not null,
    surname varchar(200) not null
);
-- password=12345
insert into users (login, hashed_password, name, surname)
values ('admin', '827ccb0eea8a706c4c34a16891f84e7b', 'admin', 'admin');

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

DROP TABLE azs_button;
create table if not exists azs_button
(
	id_azs  bigint,
    price1   int,
    price2   int,
    button  int
);
-- insert into azs_button (id_azs, price1, price2, button)
-- values (10111999, 4300, 5300, 33);