INSERT INTO users 
(name, password) VALUES
('testuser1', 'password1'),
('testuser2', 'password2'),
('testuser3', 'password3');

INSERT INTO mails
(user_id, mail_localpart, mail_address) VALUES
(1,"user1", "user1@example.com"),
(1,"hoge-user1", "hoge-user1@example.com"),
(2,"user2", "user2@example.com");
