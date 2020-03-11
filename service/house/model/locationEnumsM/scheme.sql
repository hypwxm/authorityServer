drop table if exists wb_house_enums;

CREATE TABLE if not exists wb_house_enums
(
    id         varchar(128) not null unique primary key,
    createtime bigint       not null,
    updatetime bigint  default null,
    deletetime bigint  default null,
    isdelete   boolean default false,
    disabled   boolean default false,

    name       varchar(16),
    note       varchar(255),
    sort       SMALLSERIAL  not null
);

comment on column wb_house_enums.name is '房屋枚举';
Create Index createtime_index On wb_house_enums (createtime);



