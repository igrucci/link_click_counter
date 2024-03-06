create table URL
(
    id    bigserial,
    url   varchar(255)    not null unique check (url ~* '^https?://(?:\w+\.)?\w+\.\w{2,}(?:/\S*)?$'),
    name  varchar(255)     not null ,
    code  UUID             not null ,
    count bigint default 0 not null ,
    constraint URL_PK primary key (id)
);

create table Statistic
(
    id     bigserial,
    url_id bigint                  not null ,
    time   timestamptz default now() not null ,
    constraint STATISTIC_PK primary key (id),
    constraint STATISTIC_URL_FK foreign key (url_id) references URL (id)
);

create index url_code_index on URL (code);
create index url_url_index on URL (url);

