create database ooper;
use ooper;

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

create table trip
(
id int auto_increment primary key,
pickup varchar(150),
dropoff varchar(150),
passenger_id int,
driver_id int,
FOREIGN KEY (passenger_id) REFERENCES passenger(id),
FOREIGN KEY (driver_id) REFERENCES driver(id),
requested datetime,
start datetime,
end datetime
);

-- CREATE USER 'user'@'localhost' IDENTIFIED BY 'password';
-- GRANT ALL ON *.* TO 'user'@'localhost'