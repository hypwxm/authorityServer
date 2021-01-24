drop table if exists g_admin_user_role;
CREATE TABLE if not exists g_admin_user_role (
    user_id varchar(128) not null check (user_id <> ''),
    role_id varchar(128) not null check (role_id <> ''),
    org_id varchar(128) not null default ''
);

Create Unique Index g_admin_user_role_user_role_org On g_admin_user_role (user_id, role_id, org_id);