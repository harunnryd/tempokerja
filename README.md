This is a example code of temporal.io.

## Basic usage

**Create Table**
```sql
create table if not exists products (
	id serial primary key,
	name varchar(100) not null default '',
	price numeric(19,2) not null default 0.0,
	quantity integer not null default 0,
	created_at timestamp not null default current_timestamp,
	updated_at timestamp not null default current_timestamp,
	deleted_at timestamp null default null
);

create table if not exists users (
	id SERIAL primary key,
	full_name varchar(150) not null default '',
	email varchar(100) not null default '',
	password varchar(100) not null default '',
	created_at timestamp not null default current_timestamp,
	updated_at timestamp not null default current_timestamp,
	deleted_at timestamp null default null
);


create table if not exists balances (
	id SERIAL primary key,
	user_id integer null,
	nominal numeric(19,2) not null default 0.0,
	created_at timestamp not null default current_timestamp,
	updated_at timestamp not null default current_timestamp,
	deleted_at timestamp null default null
);

create table if not exists mutations (
	id SERIAL primary key,
	destination_id integer null,
	origin_id integer null,
	amount  numeric(19,2) not null default 0.0,
	status integer null,
	pocket integer null,
	channel integer null,
	created_at timestamp not null default current_timestamp,
	updated_at timestamp not null default current_timestamp,
	deleted_at timestamp null default null
);

create table if not exists transactions (
	id SERIAL primary key,
	product_id integer null,
	quantity integer not null default 1,
	created_at timestamp not null default current_timestamp,
	updated_at timestamp not null default current_timestamp,
	deleted_at timestamp null default null
);

create table if not exists mutation_transaction(
	id SERIAL primary key,
	mutation_id integer null,
	transaction_id integer null,
	created_at timestamp not null default current_timestamp,
	updated_at timestamp not null default current_timestamp,
	deleted_at timestamp null default null
);

insert into products (name, price, quantity) values ('Kancut', 1000, 100), ('Sempak', 2000, 100), ('Kolor', 2500, 100), ('Boxer', 4000, 100);

insert into users (full_name, email, password) values ('Admin', 'admin@o2o.com', 'p4s5w0rd'), ('User Test', 'usertest@gmail.com', 'p4s5w0rd');

insert into balances (user_id, nominal) values (1, 0), (2, 100000);
```


**Run Temporal Service**
```bash
foo@bar:~$ git clone https://github.com/temporalio/docker-compose.git
foo@bar:~$ cd docker-compose
foo@bar:~$ docker-compose -f docker-compose-postgres.yml up -d
```

**Run Service (1)**
```bash
foo@bar:~$ git clone https://github.com/harunnryd/tempolalu.git
foo@bar:~$ cd tempolalu
foo@bar:~$ go run main.go
```

**Run Service (2)**
```bash
foo@bar:~$ git clone https://github.com/harunnryd/tempokini.git
foo@bar:~$ cd tempokini
foo@bar:~$ go run main.go
```

**Run Core Service**
```bash
foo@bar:~$ git clone https://github.com/harunnryd/tempokerja.git
foo@bar:~$ cd tempokerja
foo@bar:~$ go run main.go
```

## Auxiliary packages

|  Package | Description  |
| ------------ | ------------ |
| [go-chi](https://github.com/go-chi/chi)  | Router (lightweight, idiomatic and composable router for building Go HTTP services) |
| [viper](https://github.com/spf13/viper)  | Go configuration with fangs  |
| [cobra](https://github.com/spf13/cobra)  | Go CLI  |
| [gorm](https://github.com/go-gorm/gorm)  | The fantastic ORM library for Golang, aims to be developer friendly  |

## Contributors
1. [Harun Nur Rasyid](https://github.com/harunnryd)

## Contributing
Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

Please make sure to update tests as appropriate.

## License
Copyright (c) 2021-present [Harun Nur Rasyid](https://github.com/harunnryd)

Licensed under [MIT License](./LICENSE)