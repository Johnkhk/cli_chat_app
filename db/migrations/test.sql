-- Insert sample users (assuming user1 already exists)
INSERT INTO users (username, password_hash) VALUES
('user2', 'hashed_password_2'), -- User ID 2
('user3', 'hashed_password_3'), -- User ID 3
('user4', 'hashed_password_4'), -- User ID 4
('user5', 'hashed_password_5'); -- User ID 5


-- Add friends for user ID 1
INSERT INTO friends (user_id, friend_id) VALUES
(1, 2),  -- user1 is friends with user2
(1, 3),  -- user1 is friends with user3
(2, 1),  -- user2 is friends with user1
(3, 1);  -- user3 is friends with user1


-- Outgoing friend requests from user ID 1
INSERT INTO friend_requests (requester_id, recipient_id, status) VALUES
(1, 4, 'pending'), -- user1 sends a friend request to user4
(1, 5, 'pending'); -- user1 sends a friend request to user5

-- Incoming friend requests to user ID 1
INSERT INTO friend_requests (requester_id, recipient_id, status) VALUES
(2, 1, 'pending'), -- user2 sends a friend request to user1
(3, 1, 'pending'); -- user3 sends a friend request to user1
