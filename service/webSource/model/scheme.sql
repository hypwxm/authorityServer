drop table if exists wb_settings_source;

CREATE TABLE if not exists wb_settings_source
(
    id          varchar(128)      not null unique primary key,
    createtime  bigint            not null,
    updatetime  bigint  default 0 not null,
    deletetime  bigint  default 0 not null,
    isdelete    boolean default false,
    disabled    boolean default false,


    name        varchar(10)       not null check ( name <> '' ),
    api_path    varchar(128),
    source_name varchar(50),
    parent_id   varchar(128)
);
Create Index wb_settings_source_createtime_index On wb_settings_source (createtime);



