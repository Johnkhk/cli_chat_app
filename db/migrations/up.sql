-- Create the chat_users table
CREATE TABLE chat_users (
    id INT AUTO_INCREMENT PRIMARY KEY,
    username VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create the friends table to store friendships between users
CREATE TABLE friends (
    id INT AUTO_INCREMENT PRIMARY KEY,
    user_id INT NOT NULL,
    friend_id INT NOT NULL,
    most_recent_message_id INT DEFAULT NULL, -- New field to store the most recent message ID
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES chat_users(id),
    FOREIGN KEY (friend_id) REFERENCES chat_users(id),
    UNIQUE (user_id, friend_id)
);
-- Create the messages table
CREATE TABLE messages (
    id INT AUTO_INCREMENT PRIMARY KEY,
    sender_id INT,
    receiver_id INT,
    message TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (sender_id) REFERENCES chat_users(id),
    FOREIGN KEY (receiver_id) REFERENCES chat_users(id)
);
