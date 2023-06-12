create table query_info (
    id serial primary key,
    release_date bigint,
    company_name text not null,
    symbol text not null,
    record_name text, 
    record_title text,
    period_date bigint,
    file_path text not null,
    url_link text not null
);