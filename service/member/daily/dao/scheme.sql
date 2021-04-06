drop table if exists g_daily;
CREATE TABLE if not exists g_daily (
    id varchar(128) not null unique primary key,
    createtime bigint not null,
    updatetime bigint default 0 not null,
    deletetime bigint default 0 not null,
    isdelete boolean default false,
    disabled boolean default false,
    weight real default 0 not null,
    height real default 0 not null,
    diary text default '' not null,
    user_id varchar(128) not null check (user_id <> ''),
    baby_id varchar(128) not null check (user_id <> ''),
    weather varchar(128) not null default '',
    mood varchar(128) not null default '',
    health varchar(128) not null default '',
    temperature real default 0 not null
);
Create Index g_daily_createtime_index On g_daily (createtime);