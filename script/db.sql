CREATE TABLE account (
    id INT(11) NOT NULL AUTO_INCREMENT,
    username VARCHAR(64),
    password VARCHAR(64),
    salt VARCHAR(8),
    PRIMARY KEY(id),
    INDEX(username)
);


insert into account (username,password,salt) value("test", md5(concat(md5("test"),"1Uz23")), "1Uz23");

CREATE TABLE user (
    id INT(11) NOT NULL,
    name VARCHAR(64),
    PRIMARY KEY(id),
    INDEX(name)
);

insert into user (id, name) value(1, "test");

