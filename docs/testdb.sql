drop table if exists customers;
drop table if exists addresses;
drop table if exists invoices;

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

create table invoices (
    id int not null,
    customer_id int not null,
    number varchar not null,
    date date not null,
    amount money not null default 0.00,
    primary key (id)
);

insert into invoices values (1, 1, 'A0001', '05/21/18', 1000.00);
insert into invoices values (2, 1, 'C0072', '07/22/18', 1010.00);
insert into invoices values (3, 1, 'D0084', '09/04/18', 2050.00);
insert into invoices values (4, 1, 'D0092', '09/07/18', 750.00);
insert into invoices values (5, 2, 'B0050', '05/30/18', 3432.45);
insert into invoices values (6, 2, 'C0007', '07/10/18', 492.72);
insert into invoices values (7, 2, 'D0192', '09/22/18', 4234.98);