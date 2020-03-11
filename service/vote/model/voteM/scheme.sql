drop table if exists wb_vote;

CREATE TABLE if not exists wb_vote
(
    id         varchar(128) not null unique primary key,
    createtime bigint       not null,
    updatetime bigint  default null,
    deletetime bigint  default null,
    isdelete   boolean default false,
    disabled   boolean default false,
    title      varchar(512),
    comment    text
);
