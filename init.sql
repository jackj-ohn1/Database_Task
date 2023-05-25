IF NOT EXISTS(
		SELECT name
		FROM sys.databases
		WHERE name = 'book_control'
	)
	BEGIN
		CREATE DATABASE book_control;
	END;

USE book_control;

CREATE TABLE book(
	book_id INT IDENTITY(1000000,1) PRIMARY KEY,
	
	book_name VARCHAR(100) NOT NULL,
	book_publish_time DATETIME2 DEFAULT NULL,
	book_price FLOAT NOT NULL,
	book_content VARCHAR(MAX),
	book_used BIT CONSTRAINT book_status_check CHECK (book_used IN (0,1)) DEFAULT 0,
	book_author VARCHAR(200) DEFAULT 'unknown'
);

-- CREATE TABLE author(
-- 	author_id INT IDENTITY(1000000,1) PRIMARY KEY,
-- 	author_name VARCHAR(100) NOT NULL,
-- 	author_age TINYINT DEFAULT -1,
-- 	author_description VARCHAR(250) DEFAULT ''
-- );
--
-- CREATE TABLE book_author(
-- 	author_id INT NOT NULL,
-- 	book_id INT NOT NULL,
--
-- 	PRIMARY KEY(author_id,book_id),
-- 	CONSTRAINT fk_book_author_author FOREIGN KEY(author_id)
-- 	REFERENCES author(author_id),
-- 	CONSTRAINT fk_book_author_book FOREIGN KEY(book_id)
-- 	REFERENCES book(book_id)
--
-- 	ON DELETE CASCADE
-- 	ON UPDATE CASCADE
-- );

CREATE TABLE library_user(
	user_id VARCHAR(10) PRIMARY KEY,
	user_password VARCHAR(20) NOT NULL,
	user_name VARCHAR(100) DEFAULT 'The Unknown Ranger',
	is_admin BIT CONSTRAINT admin_check CHECK (is_admin IN (0,1)) DEFAULT 0,
	borrow_max INT DEFAULT 3,
	borrowed_book INT DEFAULT 0
);

CREATE TABLE borrow(
	borrow_id INT IDENTITY(1000000,1),
	book_id INT NOT NULL,
	user_id VARCHAR(10) NOT NULL,
	borrow_time DATETIME2 DEFAULT GETDATE(),
	should_return_time DATETIME2 DEFAULT DATEADD(day,7,GETDATE()), --默认借七天--
	return_time DATETIME2 DEFAULT NULL,
	
	PRIMARY KEY(book_id,user_id,borrow_id),
	CONSTRAINT fk_star_user FOREIGN KEY(user_id)
	REFERENCES library_user(user_id),
	CONSTRAINT fk_star_book FOREIGN KEY(book_id)
	REFERENCES book(book_id)
	
	ON UPDATE CASCADE
);

-- 修改图书状态,增加用户已借图书数量 --
CREATE TRIGGER book_status_update ON borrow
	AFTER INSERT AS
BEGIN
	UPDATE book SET book_used = 1
	WHERE (SELECT book_id FROM inserted) = book_id;
	
	UPDATE library_user SET borrowed_book = borrowed_book + 1
    WHERE (SELECT user_id FROM inserted) = user_id
    AND borrowed_book < borrow_max;
END;

CREATE TRIGGER delete_changed_update ON borrow
	INSTEAD OF DELETE AS
BEGIN
	UPDATE borrow SET return_time = GETDATE()
	WHERE (SELECT borrow_id FROM deleted) = borrow_id;
	
	UPDATE book SET book_used = 0
	WHERE (SELECT book_id FROM deleted) = book_id;
	
	-- 还书时间超时,借书最大数量限制-1 --
	UPDATE library_user SET borrow_max = borrow_max - 1
	WHERE (SELECT user_id FROM deleted) = user_id
	AND borrow_max > 1
	AND (SELECT should_return_time FROM deleted) < GETDATE();
	
	-- 按时还书,借书最大数量限制+1 --
	UPDATE library_user SET borrow_max = borrow_max + 1
	WHERE (SELECT user_id FROM deleted) = user_id
	AND borrow_max < 10
	AND (SELECT should_return_time FROM deleted) > GETDATE();
END;

---TEST CASE---
INSERT INTO
	author(author_name,author_age,author_description)
VALUES
	('john',18,'i am john'),
	('tom',35,'i am tom');
	
INSERT INTO
	book(book_name, book_price, book_content)
VALUES
	('book_1',18.6,'this is book_1'),
	('book_2',39,'this is book_2');
	
INSERT INTO
	book_author(author_id, book_id)
VALUES
	(1000000,1000000),
	(1000000,1000001),
	(1000001,1000000);

INSERT INTO
	library_user(user_id, user_password)
VALUES
	('2021','2021'),
	('2022','2022');
	
INSERT INTO
	borrow(book_id, user_id)
VALUES
	(1000000,'2021'),
	(1000001,'2022');
	
DELETE FROM borrow WHERE book_id = 1000000 AND user_id = '2021';