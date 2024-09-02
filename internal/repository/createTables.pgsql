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
    is_second_price int,	
 	time   varchar(100) not null,
 	name    varchar(100) not null,
 	address varchar(100) not null,
	stats varchar(1500) not null
    -- FOREIGN KEY (id_user) REFERENCES users(id)
);

DROP TABLE azs_button;
create table if not exists azs_button
(
	id_azs  bigint,
    price   int,
    button  int
);
-- insert into azs_button (id_azs, price, button)
-- values (10111999, 4300, 33);
-- DELETE FROM azses WHERE id_user = -1;
-- DELETE FROM azs_button
-- WHERE ctid NOT IN (
--     SELECT MIN(ctid)
--     FROM azs_button
--     GROUP BY id_azs, price, button
-- );
-- DELETE FROM azs_button_v2
-- WHERE ctid NOT IN (
--     SELECT MIN(ctid)
--     FROM azs_button_v2
--     GROUP BY id_azs, value, button
-- );
-- DELETE FROM log_button
-- WHERE ctid NOT IN (
--     SELECT MIN(ctid)
--     FROM log_button
--     GROUP BY id_azs, download
-- );
-- DELETE FROM update_command
-- WHERE ctid NOT IN (
--     SELECT MIN(ctid)
--     FROM update_command
--     GROUP BY id_azs, url
-- );
-- DELETE FROM ya_azs_info
-- WHERE ctid NOT IN (
--     SELECT MIN(ctid)
--     FROM ya_azs_info
--     GROUP BY id_azs, lat, lon, enable
-- );