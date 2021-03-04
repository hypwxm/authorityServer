drop table if exists g_member_family;
CREATE TABLE if not exists g_member_family (
    id varchar(128) not null unique primary key,
    createtime bigint not null,
    updatetime bigint default 0 not null,
    deletetime bigint default 0 not null,
    isdelete boolean default false,
    disabled boolean default false,
    creator varchar(128) not null check (creator <> ''),
    name varchar(50) not null default '',
    label varchar(500) not null default ''
);
Create Index g_member_family_createtime_index On g_member_family (createtime);