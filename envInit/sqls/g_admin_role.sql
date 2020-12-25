drop table if exists g_admin_role;
CREATE TABLE if not exists g_admin_role (
    id varchar(128) not null unique primary key,
    createtime bigint not null,
    updatetime bigint default 0 not null,
    deletetime bigint default 0 not null,
    isdelete boolean default false,
    disabled boolean default false,
    name varchar(20) not null check (name <> ''),
    intro varchar(128) not null default '',
    org_id varchar(128) not null default ''
);
Create Index g_admin_role_createtime_index On g_admin_role (createtime);