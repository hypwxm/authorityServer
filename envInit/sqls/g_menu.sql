drop table if exists g_menu;
CREATE TABLE if not exists g_menu (
    id varchar(128) not null unique primary key,
    createtime bigint not null,
    updatetime bigint default 0 not null,
    deletetime bigint default 0 not null,
    isdelete boolean default false,
    disabled boolean default false,
    name varchar(50) default '' not null,
    icon varchar(50) default '' not null,
    path varchar(50) default '' not null,
    parent_id varchar(128) default '' not null
);
Create Index g_menu_createtime_index On g_menu (createtime);