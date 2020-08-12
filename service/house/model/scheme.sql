drop table if exists wb_house;

CREATE TABLE if not exists wb_house
(
    id         varchar(128) not null unique primary key,
    createtime bigint       not null,
    updatetime bigint  default null,
    deletetime bigint  default null,
    isdelete   boolean default false,
    disabled   boolean default false,

    name       varchar(16),
    note       varchar(255),
    sort       SMALLSERIAL  not null,
    icon       varchar(255) not null,
    parent_id  varchar(255) not null
);

comment on column wb_house.name is '房屋枚举';
Create Index wb_house_createtime_index On wb_house (createtime);



