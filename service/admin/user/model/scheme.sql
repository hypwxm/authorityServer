drop table if exists wb_admin_user;

CREATE TABLE if not exists wb_admin_user
(
    id         varchar(128)       not null unique primary key,
    createtime bigint             not null,
    updatetime bigint   default 0 not null,
    deletetime bigint   default 0 not null,
    isdelete   boolean  default false,
    disabled   boolean  default false,


    account    varchar(20)        not null check ( account <> '' ),
    password   varchar(128)       not null check ( password <> '' ),
    username   varchar(20)        not null check ( username <> '' ),
    salt       varchar(128)       not null check ( salt <> '' ),
    avatar     varchar(512)       not null default '',
    type       smallint default 1 not null,
    role_id    varchar(128)       not null default ''
);
Create Index wb_admin_user_createtime_index On wb_admin_user (createtime);
comment on column wb_admin_user.type is '1：超级管理员，2：一般角色';



