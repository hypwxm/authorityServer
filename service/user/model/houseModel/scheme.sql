
drop table if exists wb_user_house;

CREATE TABLE if not exists wb_user_house
(
    id        varchar(128) not null unique primary key,
    enums_id  varchar(128) not null check (enums_id <> ''),
    option_id varchar(128) not null check (enums_id <> ''),
    user_id   varchar(128) not null check (user_id <> '')
);

create unique index wb_user_house_createtime_index on wb_user_house (enums_id, option_id, user_id);