-- Extension pour générer des UUID si nécessaire, ou on utilise des BIGSERIAL (IDs auto-incrémentés)
CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    username VARCHAR(30) UNIQUE NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE subreddits (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(21) UNIQUE NOT NULL, -- Reddit limite à 21 caractères
    description TEXT,
    creator_id BIGINT REFERENCES users(id) ON DELETE SET NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE posts (
    id BIGSERIAL PRIMARY KEY,
    title VARCHAR(300) NOT NULL,
    content TEXT, -- Peut être vide si c'est juste un lien ou une image (ici text pour simplifier)
    user_id BIGINT REFERENCES users(id) ON DELETE CASCADE,
    subreddit_id BIGINT REFERENCES subreddits(id) ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE comments (
    id BIGSERIAL PRIMARY KEY,
    content TEXT NOT NULL,
    post_id BIGINT REFERENCES posts(id) ON DELETE CASCADE,
    user_id BIGINT REFERENCES users(id) ON DELETE CASCADE,
    parent_id BIGINT REFERENCES comments(id) ON DELETE CASCADE, -- Pour l'arborescence (réponses aux commentaires)
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Gestion des votes (1 = upvote, -1 = downvote)
CREATE TABLE post_votes (
    user_id BIGINT REFERENCES users(id) ON DELETE CASCADE,
    post_id BIGINT REFERENCES posts(id) ON DELETE CASCADE,
    value SMALLINT CHECK (value IN (-1, 1)),
    PRIMARY KEY (user_id, post_id)
);