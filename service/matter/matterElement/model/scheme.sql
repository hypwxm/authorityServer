drop table if exists wb_matter_element;

CREATE TABLE if not exists wb_matter_element
(
    id         varchar(128)            not null unique primary key,
    createtime bigint                  not null,
    updatetime bigint       default 0  not null,
    deletetime bigint       default 0  not null,
    isdelete   boolean      default false,
    disabled   boolean      default false,


    title      varchar(24)             not null check (title <> '' ),
    intro      varchar(256) default '' not null,
    type       smallint                not null,
    min        smallint     default 0  not null,
    max        smallint     default 0  not null,
    matter_id  varchar(128)            not null check ( matter_id <> '' ),
    check ( max >= min )

);
Create Index wb_matter_element_createtime_index On wb_matter_element (createtime);
comment on column wb_matter_element.title is '标题';
comment on column wb_matter_element.type is '1：单选，2：多选，3：文案反馈';
comment on column wb_matter_element.max is '多选情况下最多能选择的数量，0代表没有限制';
comment on column wb_matter_element.min is '多选情况下最少需要选择的数量，0代表没有限制';



