drop table if exists wb_house_option;

CREATE TABLE if not exists wb_house_option
(
    id             varchar(128) not null unique primary key,
    createtime     bigint       not null,
    updatetime     bigint  default null,
    deletetime     bigint  default null,
    isdelete       boolean default false,
    disabled       boolean default false,

    name           varchar(16),
    house_enums_id varchar(128) not null,
    note           varchar(255),
    sort           SMALLSERIAL  not null
);

comment on column wb_house_option.house_enums_id is '对应的房屋枚举';
Create Index createtime_index On wb_house_option (createtime);


drop table if exists wb_house_option_associate;

CREATE TABLE if not exists wb_house_option_associate
(
    super_option_id varchar(128) not null,
    sub_option_id   varchar(128) not null
);
Create Unique Index super_sub_index On wb_house_option_associate (super_option_id, sub_option_id);


