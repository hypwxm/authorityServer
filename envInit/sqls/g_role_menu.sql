drop table if exists g_role_menu;

CREATE TABLE if not exists g_role_menu
(
    id         varchar(128)      not null unique primary key,
    createtime bigint            not null,
    updatetime bigint  default 0 not null,
    deletetime bigint  default 0 not null,
    isdelete   boolean default false,
    disabled   boolean default false,

    role_id    varchar(128)      not null check ( role_id <> '' ),
    menu_id    varchar(128)      not null check ( menu_id <> '' )
);
Create Index g_role_menu_createtime_index On g_role_menu (createtime);




