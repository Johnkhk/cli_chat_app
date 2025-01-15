CREATE TABLE users (
    id INT AUTO_INCREMENT PRIMARY KEY,
    username VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);


CREATE TABLE friend_requests (
    id SERIAL PRIMARY KEY,
    requester_id INT NOT NULL, -- User ID of the requester
    recipient_id INT NOT NULL, -- User ID of the recipient
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, -- When the request was made
    status VARCHAR(20) DEFAULT 'PENDING', -- Status of the request (pending, accepted, declined)
    response_at TIMESTAMP NULL DEFAULT NULL,
    FOREIGN KEY (requester_id) REFERENCES users(id),
    FOREIGN KEY (recipient_id) REFERENCES users(id)
);

CREATE TABLE friends (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL, -- ID of the user
    friend_id INT NOT NULL, -- ID of the friend
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, -- When the friendship was established
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (friend_id) REFERENCES users(id),
    UNIQUE(user_id, friend_id) -- Ensure that each friendship is unique
);


CREATE TABLE prekey_bundle (
    id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY, 
    user_id INT NOT NULL UNIQUE,  -- Make user_id unique
    registration_id INT UNSIGNED NOT NULL,      
    device_id INT UNSIGNED NOT NULL,            
    identity_key BLOB NOT NULL,                 
    pre_key_id INT UNSIGNED NOT NULL,           
    pre_key BLOB NOT NULL,                      
    signed_pre_key_id INT UNSIGNED NOT NULL,    
    signed_pre_key BLOB NOT NULL,               
    signed_pre_key_signature BLOB NOT NULL,     
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);


CREATE TABLE IF NOT EXISTS onetime_prekeys (
    user_id INT NOT NULL,                    -- Foreign key to reference prekey_bundles
    prekey_id INT PRIMARY KEY,               -- ID of the one-time prekey
    prekey BLOB NOT NULL,                    -- The one-time prekey
    FOREIGN KEY (user_id) REFERENCES prekey_bundle(user_id)
);