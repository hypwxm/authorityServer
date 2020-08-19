drop table if exists g_my_babies;

CREATE TABLE if not exists g_my_babies
(
    id             varchar(128) not null unique primary key,
    createtime     bigint       not null,
    updatetime     bigint                default 0 not null,
    deletetime     bigint                default 0 not null,
    isdelete       boolean               default false,
    disabled       boolean               default false,


    weight         real                  default 0 not null,
    height         real                  default 0 not null,
    diary          text                  default '' not null,
    user_id        varchar(128) not null check ( user_id <> '' ),
    baby_id        varchar(128) not null check ( user_id <> '' ),
    name           varchar(20)  not null check (name <> ''),
    birthday       varchar(20)  not null check (birthday <> ''),
    gender         varchar(1)   not null check (gender <> ''),
    avatar         varchar(250) not null default '',

    id_card        varchar(20)  not null default '',
    hobby          varchar(250) not null default '',
    good_at        varchar(250) not null default '',
    favorite_food  varchar(250) not null default '',
    favorite_color varchar(250) not null default '',
    ambition       varchar(250) not null default ''
);
Create Index g_my_babies_createtime_index On g_my_babies (createtime);


