drop table if exists g_media;
CREATE TABLE if not exists g_media (
    id varchar(128) not null unique primary key,
    createtime bigint not null,
    updatetime bigint default 0 not null,
    deletetime bigint default 0 not null,
    isdelete boolean default false,
    disabled boolean default false,
    url varchar(500) default '' not null,
    business varchar(50) default '' not null check (business <> ''),
    business_id varchar(128) default '' not null check (business_id <> ''),
    size int not null default 0,
    user_id varchar(128) not null check (user_id <> '')
);
Create Index g_media_createtime_index On g_media (createtime);