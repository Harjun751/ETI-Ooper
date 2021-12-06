use passengers;

create table passenger
(
id int auto_increment primary key,
first_name varchar(90),
last_name varchar(90),
mobile_number int,
email varchar(90),
password char(64),
salt char(32)
);