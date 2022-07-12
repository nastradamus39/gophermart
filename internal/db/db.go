package db

var Schema = `
create table if not exists users
(
	login     varchar(256),
	password  varchar(256),
	accrual   int,
	withdrawn int,
	balance   int
); 

create unique index if not exists users_login on users (login);
`
