-- db schema


CREATE TABLE users (
	id INT PRIMARY KEY AUTO_INCREMENT,
	oauth_id VARCHAR(255) DEFAULT NULL,
	name VARCHAR(255) NOT NULL,
	email VARCHAR(255) NOT NULL,
	picture VARCHAR(255) DEFAULT NULL,
	oauth VARCHAR(255) DEFAULT "google",
	created_at DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  
  
CREATE TABLE users_session (
	id INT PRIMARY KEY AUTO_INCREMENT,
	user_id INT UNSIGNED NOT NULL,
	login_token VARCHAR(255) NOT NULL,
	ip_addr VARCHAR(255) NOT NULL,
	superseded TINYINT(1) DEFAULT 0,
	created_at DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  


CREATE TABLE notes (
	id INT PRIMARY KEY AUTO_INCREMENT,
	user_id INT UNSIGNED NOT NULL,
	title VARCHAR(255) DEFAULT NULL,
	content VARCHAR(255) DEFAULT NULL,
	type INT UNSIGNED DEFAULT 0,
	editable TINYINT(1) UNSIGNED DEFAULT 0,
	is_visible TINYINT(1) UNSIGNED DEFAULT 0,
	is_deleted TINYINT(1) UNSIGNED DEFAULT 0,
	version INT UNSIGNED DEFAULT 1,
	created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
	deleted_at DATETIME DEFAULT NULL,
	updated_at DATETIME DEFAULT NULL
  );
  
    
CREATE TABLE notes_alert (
	id INT PRIMARY KEY AUTO_INCREMENT,
	user_id INT UNSIGNED NOT NULL,
	notes_id INT UNSIGNED NOT NULL,
	version INT UNSIGNED DEFAULT 1,
	alert_type VARCHAR(255) DEFAULT "ONCE",
	alert_time datetime NOT NULL,
	superseded TINYINT(1) DEFAULT 0,
	created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
	updated_at DATETIME DEFAULT NULL
  );