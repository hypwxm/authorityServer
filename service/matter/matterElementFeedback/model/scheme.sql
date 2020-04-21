drop table if exists wb_matter_element_feedback;

CREATE TABLE if not exists wb_matter_element_feedback
(
    id                 varchar(128)       not null unique primary key,
    createtime         bigint             not null,
    updatetime         bigint  default 0  not null,
    deletetime         bigint  default 0  not null,
    isdelete           boolean default false,
    disabled           boolean default false,


    title              varchar(20)        not null check ( title <> '' ),
    matter_id          varchar(128)       not null check ( matter_id <> '' ),
    element_id         varchar(128)       not null check ( element_id <> '' ),
    element_option_ids text    default '' not null,
    user_id            varchar(128)       not null check ( user_id <> '' )

);
Create Index wb_matter_element_feedback_createtime_index On wb_matter_element_feedback (createtime);
comment on column wb_matter_element_feedback.title is '标题';
comment on column wb_matter_element_feedback.element_option_ids is '反馈，选则的id列表，单选或者多选，文本的情况这里存提交的文本内容';





