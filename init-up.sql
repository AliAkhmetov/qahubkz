CREATE TABLE IF NOT EXISTS users (
    id              SERIAL         primary key,
    email           varchar(50)     not null unique,
    username        varchar(50)     not null unique,
    password_hash   varchar(255)    not null,
    token           varchar(50)     unique,
    user_type       varchar(50)     not null, 
    mod_requested   boolean, 
    expire_at       date            
);
CREATE TABLE IF NOT EXISTS mod_requests (
    id			SERIAL         primary key,
    user_id	    integer         not null,
    created_at  date,            
    updated_at  date,      
    status		varchar(255)    not null,
    foreign key (user_id)       references users(id)
);
CREATE TABLE IF NOT EXISTS posts (
    id			SERIAL         primary key,
    created_by	integer         not null,
    created_at  date,            
    updated_at  date,                 
    title		varchar(30),
    status      varchar(50)     not null,
    content		varchar(305),
    foreign key (created_by)    references users(id)
);
CREATE TABLE IF NOT EXISTS reports (
    id			    SERIAL         primary key,
    created_by	    integer         not null,
    post_id         integer         not null,
    created_at      date,            
    updated_at      date,            
    moderator_msg	varchar(255),
    admin_msg		varchar(255),
    status		    varchar(255)    not null,
    foreign key (created_by)       references users(id),
    foreign key (post_id)       references posts(id)

);
CREATE TABLE IF NOT EXISTS categories (
    id			SERIAL        primary key,
    name		varchar(255)
);
CREATE TABLE IF NOT EXISTS posts_categories (
    id			SERIAL         primary key,
    post_id     integer         not null,
    category_id integer         not null,
    foreign key (post_id)       references posts(id),
    foreign key (category_id)   references categories(id)
);
CREATE TABLE IF NOT EXISTS comments (
    id			SERIAL         primary key,
    created_by	integer,
    created_at  date,            
    updated_at  date,               
    post_id		integer,
    content		varchar(305),
    status      varchar(50)     not null,
    foreign key (created_by)    references users(id),
    foreign key (post_id)       references posts(id)
);
CREATE TABLE IF NOT EXISTS posts_likes (
    id			SERIAL        primary key,
    created_by  integer,
    post_id		integer,
    type        boolean         not null,
    foreign key (created_by)    references users(id),
    foreign key (post_id)       references posts(id),
    unique      (post_id, created_by)

);
CREATE TABLE IF NOT EXISTS comments_likes (
    id			SERIAL           primary key,
    created_by	integer,
    comment_id	integer,
    type        boolean             not null,
    foreign key (created_by)        references users(id),
    foreign key (comment_id)        references comments(id),
    unique      (comment_id, created_by)
);

INSERT INTO categories (name) values ('GO');
INSERT INTO categories (name) values ('JS');
INSERT INTO categories (name) values ('PHP');
INSERT INTO categories (name) values ('HTML');