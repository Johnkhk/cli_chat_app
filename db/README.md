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
`mysql -u cli_chat_dev -p cli_chat_app < db/migrations/up.sql;`

`mysql -u cli_chat_dev -p cli_chat_app < db/migrations/down.sql`
