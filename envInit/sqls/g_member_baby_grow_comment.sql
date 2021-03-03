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
    baby_id varchar(128) not null check (user_id <> ''),
    dairy_id varchar(128) not null check (dairy_id <> '')
);
Create g_member_baby_grow_comment_createtime_index On g_member_baby_grow_comment (createtime);
Create g_member_baby_grow_comment_dairy_id_index On g_member_baby_grow_comment (dairy_id);
Create g_member_baby_grow_comment_user_id_index On g_member_baby_grow_comment (user_id);