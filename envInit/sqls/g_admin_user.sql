drop table if exists g_admin_user;
CREATE TABLE if not exists g_admin_user (
    id varchar(128) not null unique primary key,
    createtime bigint not null,
    updatetime bigint default 0 not null,
    deletetime bigint default 0 not null,
    isdelete boolean default false,
    disabled boolean default false,
    account varchar(20) not null check (account <> ''),
    password varchar(128) not null check (password <> ''),
    username varchar(20) not null check (username <> ''),
    salt varchar(128) not null check (salt <> ''),
    avatar varchar(512) not null default '',
    post varchar(50) not null default '',
    contact_way varchar(50) not null default '',
    sort int not null default 0
);
Create Index g_admin_user_createtime_index On g_admin_user (createtime);