use drivers;

create table driver
(
id int auto_increment primary key,
first_name varchar(90),
last_name varchar(90),
mobile_number int,
email varchar(90),
ic_number varchar(20),
license_number varchar(30),
password char(64),
salt char(32),
available bool default true
);