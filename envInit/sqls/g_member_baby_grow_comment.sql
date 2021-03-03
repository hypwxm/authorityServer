drop table if exists g_member_baby_grow_comment;
CREATE TABLE if not exists g_member_baby_grow_comment (
    id varchar(128) not null unique primary key,
    createtime bigint not null,
    updatetime bigint default 0 not null,
    deletetime bigint default 0 not null,
    isdelete boolean default false,
    disabled boolean default false,
    content text default '' not null,
    user_id varchar(128) not null check (user_id <> ''),
    baby_id varchar(128) not null check (baby_id <> ''),
    diary_id varchar(128) not null check (diary_id <> ''),
    comment_id varchar(128) not null default ''
);
Create Index g_member_baby_grow_comment_createtime_index On g_member_baby_grow_comment (createtime);
Create Index g_member_baby_grow_comment_diary_id_index On g_member_baby_grow_comment (diary_id);
Create Index g_member_baby_grow_comment_user_id_index On g_member_baby_grow_comment (user_id);