drop table if exists wb_user_vote;

CREATE TABLE if not exists wb_user_vote
(
    id         varchar(128) not null unique primary key,
    createtime bigint       not null,
    updatetime bigint  default null,
    deletetime bigint  default null,
    isdelete   boolean default false,
    disabled   boolean default false,
    user_id    varchar(128) not null,
    vote_id    varchar(128) not null,
    option_id  varchar(128) not null
);
