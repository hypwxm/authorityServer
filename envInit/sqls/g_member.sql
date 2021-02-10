drop table if exists g_member;

CREATE TABLE if not exists g_member
(
    id         varchar(128) not null unique primary key,
    createtime bigint       not null,
    updatetime bigint       not null default 0,
    deletetime bigint       not null default 0,
    isdelete   boolean               default false,
    disabled   boolean               default false,
    nickname   varchar(32)  not null default '',
    realname   varchar(32)  not null default '',
    firstname  varchar(32)  not null default '',
    lastname   varchar(32)  not null default '',

    avatar     varchar(256) not null default '',
    phone      varchar(20) not null default '',

    account    varchar(32)  not null unique,
    password   varchar(32)  not null,
    salt       varchar(32)  not null
);

create index g_member_createtime_index on g_member (createtime);








