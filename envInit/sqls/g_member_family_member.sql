drop table if exists g_member_family_member;
CREATE TABLE if not exists g_member_family_member (
    id varchar(128) not null unique primary key,
    createtime bigint not null,
    updatetime bigint default 0 not null,
    deletetime bigint default 0 not null,
    isdelete boolean default false,
    disabled boolean default false,
    member_id varchar(128) not null check (member_id <> ''),
    family_id varchar(128) not null check (family_id <> ''),
    creator varchar(128) not null check (creator <> ''),
    can_invite boolean default false,
    can_remove boolean default false,
    can_edit boolean default false,
    nickname varchar(50) not null default '',
    role_name varchar(50) not null default ''
);
Create Index g_member_family_member_createtime_index On g_member_family_member (createtime);