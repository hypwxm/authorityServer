drop table if exists wb_admin_role;

CREATE TABLE if not exists wb_admin_role
(
    id               varchar(128) not null unique primary key,
    createtime       bigint       not null,
    updatetime       bigint                default 0 not null,
    deletetime       bigint                default 0 not null,
    isdelete         boolean               default false,
    disabled         boolean               default false,


    name             varchar(20)  not null check ( name <> '' ),
    intro            varchar(128) not null default '',
    parent_role_id   varchar(128) not null default '',
    parent_role_link text         not null default ''
);
Create Index wb_admin_role_createtime_index On wb_admin_role (createtime);
comment on column wb_admin_role.parent_role_id is '父级的角色id，角色可以属于另一个角色，相当于部门管理，存储';
comment on column wb_admin_role.parent_role_link is '父级的角色id的链路，a=>b=>c=>d；为了方便查询和递归查找上级，当中间某一层级进行上级切换的时候需要对下级进行操作，所以这一步可能比较繁琐';




