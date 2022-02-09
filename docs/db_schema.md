# Database Schema

## Posts
```sql
create table posts (
    id INTEGER NOT NULL PRIMARY KEY AUTO_INCREMENT,
    title VARCHAR(100) NOT NULL,
    content TEXT NOT NULL,
    image VARCHAR(255) NOT NULL DEFAULT '',
    status TEXT NOT NULL,
    created DATETIME NOT NULL,
    updated DATETIME NOT NULL
);

CREATE INDEX idx_posts_created ON posts(created);
```

## Pages
```sql
create table pages (
    id INTEGER NOT NULL PRIMARY KEY AUTO_INCREMENT,
    title VARCHAR(100) NOT NULL,
    content TEXT NOT NULL,
    status TEXT NOT NULL,
    slug VARCHAR(255) NOT NULL,
    image VARCHAR(255) NOT NULL DEFAULT '',
    created DATETIME NOT NULL,
    updated DATETIME NOT NULL
);

CREATE INDEX idx_pages_created ON pages(created);
```

## Users
```sql
CREATE TABLE users (
    id INTEGER NOT NULL PRIMARY KEY AUTO_INCREMENT,
    username VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    hashed_password CHAR(60) NOT NULL,
    last_login DATETIME NOT NULL,
    active BOOLEAN NOT NULL DEFAULT TRUE,
    role VARCHAR(255) NOT NULL,
    created DATETIME NOT NULL,
    updated DATETIME NOT NULL
);

ALTER TABLE users ADD CONSTRAINT users_uc_email UNIQUE (email);
```

## Menus
```sql
CREATE TABLE menus (
    id INTEGER NOT NULL PRIMARY KEY AUTO_INCREMENT,
    item_ids JSON NOT NULL,
    slug VARCHAR(255) NOT NULL,
    created DATETIME NOT NULL,
    updated DATETIME NOT NULL
);

ALTER TABLE menus ADD CONSTRAINT menus_uc_slug UNIQUE (slug);
```

## Media
```sql
create table media (
    id INTEGER NOT NULL PRIMARY KEY AUTO_INCREMENT,
    title VARCHAR(255) NOT NULL,
    filename VARCHAR(255) NOT NULL,
    path VARCHAR(255) NOT NULL,
    description VARCHAR(255) NOT NULL,
    width INT NOT NULL,
    height INT NOT NULL,
    size INT NOT NULL,
    created DATETIME NOT NULL,
    updated DATETIME NOT NULL
);

CREATE INDEX idx_media_created ON media(created);
```