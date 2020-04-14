-- 事宜可见用户配置表
drop table if exists wb_matter_visible;

CREATE TABLE if not exists wb_matter_visible
(
    matter_id varchar(128) not null check ( matter_id <> '' ),
    user_id   varchar(128) not null check ( user_id <> '' )
);


Create Index wb_matter_visible_matter_id_user_id On wb_matter_visible (matter_id, user_id);
