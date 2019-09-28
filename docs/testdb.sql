drop table if exists customers;
drop table if exists addresses;

create table customers (
    id int not null,
    name varchar not null,
    created_on date,
    updated_on date,
    primary key (id)
);

insert into customers values (1, 'Shoe Shine Ltd', '09/30/2018', '07/14/2019');
insert into customers values (2, 'Fractured Inc',  '04/19/2019', '09/15/2019');
insert into customers values (3, 'Corner Store, Inc', '08/09/2018', '07/21/2019');

create table addresses (
    id int not null,
    customer_id int not null,
    address1 varchar,
    address2 varchar,
    city varchar,
    state varchar(2),
    zip varchar(5),
    primary key (id)
);

insert into addresses (id, customer_id, address1, city, state, zip)
    values (1, 1, '19 Hampshire Pl', 'Portland', 'ME', '12345');
insert into addresses (id, customer_id, address1, address2, city, state, zip)
    values (2, 1, '4768 Bobshire St', 'Suite 500', 'Portsmouth', 'NH', '19999');
insert into addresses (id, customer_id, address1, city, state, zip)
    values (3, 2, '99 Marshal St', 'Bloomington', 'IN', '55555');
