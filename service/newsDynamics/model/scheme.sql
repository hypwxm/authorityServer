drop table if exists wb_news_dynamics;

CREATE TABLE if not exists wb_news_dynamics
(
    id         varchar(128) not null unique primary key,
    createtime bigint       not null,
    updatetime bigint       default null,
    deletetime bigint       default null,
    isdelete   boolean      default false,
    disabled   boolean      default false,


    title      varchar(24)  not null check (title <> '' ),
    intro      varchar(256) default '',
    surface    varchar(512) not null check (surface <> '' ),
    content    text         not null check (content <> '' ),
    publisher  varchar(128) not null check (publisher <> '' ),
    type       smallint,
    sort       SMALLSERIAL  not null
);
Create Index wb_news_dynamics_createtime_index On wb_news_dynamics (createtime);
comment on column wb_news_dynamics.title is '标题';
comment on column wb_news_dynamics.type is '1:管理员，系统发布，2:某一用户发布';



