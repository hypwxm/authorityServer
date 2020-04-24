drop table if exists wb_settings_menu;

CREATE TABLE if not exists wb_settings_menu
(
    id         varchar(128)      not null unique primary key,
    createtime bigint            not null,
    updatetime bigint  default 0 not null,
    deletetime bigint  default 0 not null,
    isdelete   boolean default false,
    disabled   boolean default false,


    name       varchar(10)       not null check ( name <> '' ),
    path       varchar(128)      not null check ( path <> '' )
);
Create Index wb_news_dynamics_createtime_index On wb_settings_menu (createtime);



