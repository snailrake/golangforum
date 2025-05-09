INSERT INTO topics (id, title, description, created_at) VALUES
  (1, 'Topic 1', 'Desc 1', NOW()),
  (2, 'Topic 2', 'Desc 2', NOW()),
  (3, 'Topic 3', 'Desc 3', NOW()),
  (4, 'Topic 4', 'Desc 4', NOW()),
  (5, 'Topic 5', 'Desc 5', NOW());

INSERT INTO posts (id, topic_id, title, content, user_id, username, timestamp) VALUES
  (1, 1, 'Post 1', 'Content 1', 1, 'user1', NOW()),
  (2, 2, 'Post 2', 'Content 2', 2, 'user2', NOW()),
  (3, 3, 'Post 3', 'Content 3', 3, 'user3', NOW()),
  (4, 4, 'Post 4', 'Content 4', 4, 'user4', NOW()),
  (5, 5, 'Post 5', 'Content 5', 5, 'user5', NOW());

INSERT INTO messages (id, user_id, username, content, timestamp) VALUES
  (1, 1, 'user1', 'Message 1', NOW()),
  (2, 2, 'user2', 'Message 2', NOW()),
  (3, 3, 'user3', 'Message 3', NOW()),
  (4, 4, 'user4', 'Message 4', NOW()),
  (5, 5, 'user5', 'Message 5', NOW());
