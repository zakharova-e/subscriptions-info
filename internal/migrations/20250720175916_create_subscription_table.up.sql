begin;

create table if not exists subscription(
    id integer generated always as identity primary key,
    service_name varchar(1000) not null,
    price integer not null,
    user_id uuid not null,
    start_date date not null,
    finish_date date
);

create index idx_subscription_date on subscription(start_date,finish_date);

--in case of frequent user_id requests
-- create index idx_subscription_user_date on subscription(user_id,start_date,finish_date);

--in case of frequent service_name requests
--create index idx_subscription_service_name_date on subscription(service_name,start_date,finish_date);

commit;