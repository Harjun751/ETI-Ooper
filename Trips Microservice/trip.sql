use trips;

create table trip
(
id int auto_increment primary key,
pickup varchar(150),
dropoff varchar(150),
passenger_id int,
driver_id int,
requested datetime,
start datetime,
end datetime
);