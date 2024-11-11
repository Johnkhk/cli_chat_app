# Database


## Starting the database:

1. Start the database in a shell:
`mysql -u root`

2. Create DB
`CREATE DATABASE cli_chat_app;`

3. See users:
```sql
mysql> SELECT User, Host FROM mysql.user;
+------------------+-----------+
| User             | Host      |
+------------------+-----------+
| mysql.infoschema | localhost |
| mysql.session    | localhost |
| mysql.sys        | localhost |
| root             | localhost |
+------------------+-----------+
```

4. Create a user for the app:
```sql
CREATE USER 'cli_chat_dev'@'localhost' IDENTIFIED BY 'your_password';
```
5. Grant privileges to the user:
```sql
mysql> GRANT ALL PRIVILEGES ON cli_chat_app.* TO 'cli_chat_dev'@'localhost';
Query OK, 0 rows affected (0.01 sec)

mysql> FLUSH PRIVILEGES;
Query OK, 0 rows affected (0.00 sec)
```

6. Log in as the new user:
```sql
mysql -u cli_chat_dev -p
```

## setup:


# Apply migrations

1. Load the environment variables
`export $(cat .env | xargs)`

2. Apply the migration
`mysql -u cli_chat_dev -p cli_chat_app < db/migrations/02_up.sql;`

`mysql -u cli_chat_dev -p cli_chat_app < db/migrations/01_down.sql`


```
CREATE TABLE messages (
    id INT AUTO_INCREMENT PRIMARY KEY,
    sender_id INT NOT NULL, -- User ID of the sender
    recipient_id INT NOT NULL, -- User ID of the recipient
    content TEXT, -- The content of the message, optional for text messages
    sent_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, -- When the message was sent
    is_read BOOLEAN DEFAULT FALSE, -- Whether the message has been read by the recipient
    attachment_id INT, -- Optional: Reference to the attachment if this message has one
    FOREIGN KEY (sender_id) REFERENCES users(id),
    FOREIGN KEY (recipient_id) REFERENCES users(id),
    FOREIGN KEY (attachment_id) REFERENCES attachments(id)
);

CREATE TABLE attachments (
    id INT AUTO_INCREMENT PRIMARY KEY,
    message_id INT NOT NULL, -- The ID of the associated message
    file_type VARCHAR(50) NOT NULL, -- Type of file (e.g., 'image', 'video')
    file_url VARCHAR(255) NOT NULL, -- URL to the file storage location
    file_size INT, -- Size of the file in bytes
    uploaded_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, -- When the file was uploaded
    FOREIGN KEY (message_id) REFERENCES messages(id)
);
```