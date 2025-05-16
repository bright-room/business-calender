create table if not exists calender.closed_days (
    date date not null primary key,
    summary varchar(50) not null
);
