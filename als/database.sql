CREATE TABLE licenses (
    id int not null unique auto_increment,
    license_key varchar(18) not null unique,
    product varchar(250) not null,
    email varchar(100) not null
);