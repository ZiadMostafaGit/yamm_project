-- +goose Up
-- +goose StatementBegin



CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    password TEXT NOT NULL,
    role VARCHAR(20) NOT NULL CHECK (role IN ('admin', 'merchant', 'customer')),
    first_name VARCHAR(100),
    last_name VARCHAR(100),
    age INT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);




CREATE TABLE stores (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    user_id INT UNIQUE NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);




CREATE TABLE categories (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) UNIQUE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);




CREATE TABLE faqs (
    id SERIAL PRIMARY KEY,
    category_id INT NOT NULL REFERENCES categories(id) ON DELETE CASCADE,
    store_id INT REFERENCES stores(id) ON DELETE CASCADE, 
    is_global BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    


    CONSTRAINT check_global_store CHECK (
        (is_global = TRUE AND store_id IS NULL) OR 
        (is_global = FALSE AND store_id IS NOT NULL)
    )
);



CREATE TABLE translations (

id SERIAL PRIMARY KEY,
faq_id INT REFERENCES faqs(id) ON DELETE CASCADE,
language_code VARCHAR(10) NOT NULL DEFAULT 'AR',
question TEXT NOT NULL,
answer TEXT NOT NULL,
created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP

);


CREATE INDEX idx_faqs_category_id ON faqs(category_id);
CREATE INDEX idx_faqs_store_id ON faqs(store_id) WHERE is_global = FALSE;
CREATE INDEX idx_faq_translations_faq_id ON translations(faq_id);
-- +goose StatementEnd



-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS translations;
DROP TABLE IF EXISTS faqs;
DROP TABLE IF EXISTS categories;
DROP TABLE IF EXISTS stores;
DROP TABLE IF EXISTS users;


-- +goose StatementEnd
