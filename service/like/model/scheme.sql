-- noinspection SqlNoDataSourceInspectionForFile

drop table if exists wb_like;

CREATE TABLE if not exists wb_like
(
    id          varchar(128)       not null unique primary key,
    createtime  bigint             not null,
    updatetime  bigint   default 0 not null,
    deletetime  bigint   default 0 not null,
    isdelete    boolean  default false,
    disabled    boolean  default false,

    user_id     varchar(128)       not null check (user_id <> ''),
    source_type smallint default 0 not null,
    source_id   varchar(128)       not null check (source_id <> '')
);
Create INDEX wb_like_createtime_index On wb_like (createtime);
Create UNIQUE INDEX wb_like_user_id_source_id_source_type_index On wb_like (user_id, source_id, source_type);

comment on column wb_like.user_id is '点赞的用户';
comment on column wb_like.source_type is '点赞的类型，1：新闻动态，2：用户，';
comment on column wb_like.source_id is '点赞的对象的id';



