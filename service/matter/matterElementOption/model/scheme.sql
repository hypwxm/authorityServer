drop table if exists wb_matter_element_option;

CREATE TABLE if not exists wb_matter_element_option
(
    id         varchar(128)      not null unique primary key,
    createtime bigint            not null,
    updatetime bigint  default 0 not null,
    deletetime bigint  default 0 not null,
    isdelete   boolean default false,
    disabled   boolean default false,


    title      varchar(20)       not null check ( title <> '' ),
    matter_id  varchar(128)      not null check ( matter_id <> '' ),
    element_id varchar(128)      not null check ( element_id <> '' )

);
Create Index wb_matter_element_option_createtime_index On wb_matter_element_option (createtime);
comment on column wb_matter_element_option.title is '标题';




