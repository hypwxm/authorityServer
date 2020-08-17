drop table if exists wb_societies_base;

CREATE TABLE if not exists wb_societies_base
(
    id            varchar(128)            not null unique primary key,
    createtime    bigint                  not null,
    updatetime    bigint       default 0  not null,
    deletetime    bigint       default 0  not null,
    isdelete      boolean      default false,
    disabled      boolean      default false,


    name          varchar(24)             not null check (name <> '' ),
    intro         varchar(256) default '' not null,
    surface       varchar(512)            not null check (surface <> '' ),
    creator       varchar(128)            not null check (creator <> '' ),
    creator_id    varchar(128)            not null check (creator_id <> '' ),
    people_max    int                     not null default 0,
    type_id       varchar(128)            not null check (type_id <> '' ),
    type_name     varchar(128)            not null check (type_name <> '' ),
    status        int                     not null default 1,
    status_reason varchar(255)            not null default '',
    publish_time  bigint       default 0  not null

);
Create Index wb_societies_base_createtime_index On wb_societies_base (createtime);



