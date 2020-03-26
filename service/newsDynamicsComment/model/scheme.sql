drop table if exists wb_news_dynamics_comment;

CREATE TABLE if not exists wb_news_dynamics_comment
(
    id              varchar(128)      not null unique primary key,
    createtime      bigint            not null,
    updatetime      bigint  default 0 not null,
    deletetime      bigint  default 0 not null,
    isdelete        boolean default false,
    disabled        boolean default false,


    content         text              not null check (content <> '' ),
    publisher       varchar(128)      not null check (publisher <> '' ),
    news_id         varchar(128)      not null check (publisher <> '' ),
    prev_publisher  varchar(128),
    prev_comment_id varchar(128),
    top_comment_id  varchar(128)
);
Create Index wb_news_dynamics_comment_createtime_index On wb_news_dynamics_comment (createtime);
comment on column wb_news_dynamics_comment.content is '评论内容';

comment on column wb_news_dynamics_comment.prev_publisher is '评论别人的评论';
comment on column wb_news_dynamics_comment.prev_comment_id is '回复的别人的评论的id';
comment on column wb_news_dynamics_comment.top_comment_id is '评论所在的最上级评论';




