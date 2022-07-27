CREATE TABLE IF NOT EXISTS posts (
	uuid uuid DEFAULT uuid_generate_v4() PRIMARY KEY,
	content text,
	likes int DEFAULT 0,
	dislikes int DEFAULT 0
);
