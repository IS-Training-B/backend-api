CREATE DATABASE IF NOT EXISTS prisma;
USE prisma;

-- DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS mails;

-- CREATE TABLE users
--   (
--     id             INT NOT NULL AUTO_INCREMENT,
--     name           VARCHAR(255),
--     password       VARCHAR(255),
--     email          VARCHAR(255),
--     email_verified TIMESTAMP(6),
--     image          VARCHAR(255),
--     created_at     TIMESTAMP(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
--     updated_at     TIMESTAMP(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
--     PRIMARY KEY (id)
--   );

CREATE TABLE mails 
  (
    id             INT NOT NULL AUTO_INCREMENT,
    user_id        VARCHAR(255) NOT NULL,
    mail_localpart VARCHAR(255),
    mail_address   VARCHAR(255),
    created_at     TIMESTAMP(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
    updated_at     TIMESTAMP(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
    UNIQUE (mail_address),
    PRIMARY KEY (id)
  )
