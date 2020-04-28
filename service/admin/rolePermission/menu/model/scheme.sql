drop table if exists wb_admin_role_menu_permission;

CREATE TABLE if not exists wb_admin_role_menu_permission
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
Create Index wb_admin_role_menu_permission_createtime_index On wb_admin_role_menu_permission (createtime);




