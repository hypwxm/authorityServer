drop table if exists wb_news_dynamics;

CREATE TABLE if not exists wb_news_dynamics
(
    id            varchar(128)            not null unique primary key,
    createtime    bigint                  not null,
    updatetime    bigint       default 0  not null,
    deletetime    bigint       default 0  not null,
    isdelete      boolean      default false,
    disabled      boolean      default false,


    title         varchar(24)             not null check (title <> '' ),
    intro         varchar(256) default '' not null,
    surface       varchar(512)            not null check (surface <> '' ),
    content       text                    not null check (content <> '' ),
    publisher     varchar(128)            not null check (publisher <> '' ),
    type          smallint,
    sort          SMALLSERIAL             not null,
    status        smallint     default 2  not null,
    status_reason varchar(255) default '' not null,
    publish_time  bigint       default 0  not null
);
Create Index wb_news_dynamics_createtime_index On wb_news_dynamics (createtime);
comment on column wb_news_dynamics.title is '标题';
comment on column wb_news_dynamics.type is '1:管理员，系统发布，2:某一用户发布';
comment on column wb_news_dynamics.status is '1:已发布，2:未完成，未提交审核，3:提交审核，4:审核通过，5：审核不通过，6:取消审核，7:撤回审批通过';



