drop table if exists g_org;
CREATE TABLE if not exists g_org (
    id varchar(128) not null unique primary key,
    createtime bigint not null,
    updatetime bigint default 0 not null,
    deletetime bigint default 0 not null,
    isdelete boolean default false,
    disabled boolean default false,
    name varchar(50) default '' not null,
    parent_id varchar(128) default '' not null
);
Create Index g_org_createtime_index On g_org (createtime);