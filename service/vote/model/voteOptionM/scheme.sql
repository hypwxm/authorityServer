drop table if exists wb_vote_option;

CREATE TABLE if not exists wb_vote_option
(
    id         varchar(128) not null unique primary key,
    createtime bigint       not null,
    updatetime bigint  default null,
    deletetime bigint  default null,
    isdelete   boolean default false,
    disabled   boolean default false,

    vote_id    varchar(128) not null,
    comment    text,
    title      varchar(512)
);
