drop table if exists wb_user;

CREATE TABLE if not exists wb_user
(
    id         varchar(128) not null unique primary key,
    createtime bigint       not null,
    updatetime bigint  default null,
    deletetime bigint  default null,
    isdelete   boolean default false,
    disabled   boolean default false,
    nickname   varchar(32),
    realname   varchar(32),
    firstname  varchar(32),
    lastname   varchar(32),

    avatar     varchar(256),

    account    varchar(32)  not null unique,
    password   varchar(32)  not null,
    salt       varchar(32)  not null,
    type       varchar(32),
    house      varchar(512)
);

create index wb_user_createtime_index on wb_user (createtime);
comment on column wb_user.house is '用户所在房屋坐落位置';
comment on column wb_user.type is '1:居民用户';



