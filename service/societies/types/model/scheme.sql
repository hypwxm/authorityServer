drop table if exists wb_societies_type;

CREATE TABLE if not exists wb_societies_type
(
    id         varchar(128)      not null unique primary key,
    createtime bigint            not null,
    updatetime bigint  default 0 not null,
    deletetime bigint  default 0 not null,
    isdelete   boolean default false,
    disabled   boolean default false,


    name       varchar(24)       not null check (name <> '' )

);
Create Index wb_societies_type_createtime_index On wb_societies_type (createtime);



