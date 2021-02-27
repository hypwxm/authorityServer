drop table if exists g_member_baby_relation;

CREATE TABLE if not exists g_member_baby_relation
(
    id             varchar(128) not null unique primary key,
    createtime     bigint       not null,
    updatetime     bigint                default 0 not null,
    deletetime     bigint                default 0 not null,
    isdelete       boolean               default false,
    disabled       boolean               default false,


    user_id        varchar(128) not null check ( user_id <> '' ),
    baby_id        varchar(128) not null check ( baby_id <> '' ),
    role_name      varchar(50) default '' not null

);
Create Index g_member_baby_relation_createtime_index On g_member_baby_relation (createtime);


