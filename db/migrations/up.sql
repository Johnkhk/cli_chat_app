-- Create the users table
CREATE TABLE users (
    id INT AUTO_INCREMENT PRIMARY KEY,
    username VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create the friend_requests table
CREATE TABLE friend_requests (
    id SERIAL PRIMARY KEY,
    requester_id INT NOT NULL, -- User ID of the requester
    recipient_id INT NOT NULL, -- User ID of the recipient
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, -- When the request was made
    status VARCHAR(20) DEFAULT 'pending', -- Status of the request (pending, accepted, declined)
    response_at TIMESTAMP, -- When the request was accepted or declined
    FOREIGN KEY (requester_id) REFERENCES users(id),
    FOREIGN KEY (recipient_id) REFERENCES users(id),
    CHECK (status IN ('pending', 'accepted', 'declined')) -- Ensure only valid statuses
);

-- Create the friends table
CREATE TABLE friends (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL, -- ID of the user
    friend_id INT NOT NULL, -- ID of the friend
    added_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, -- When the friendship was established
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (friend_id) REFERENCES users(id),
    UNIQUE(user_id, friend_id) -- Ensure that each friendship is unique
);