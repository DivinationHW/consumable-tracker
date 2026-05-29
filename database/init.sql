CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) NOT NULL UNIQUE,
    password_hash VARCHAR(200) NOT NULL,
    role VARCHAR(20) NOT NULL CHECK (role IN ('admin', 'readonly')),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE consumable_models (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    unit VARCHAR(20) DEFAULT '个',
    is_default BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE offices (
    id SERIAL PRIMARY KEY,
    room_number VARCHAR(20) NOT NULL UNIQUE,
    device_type VARCHAR(50),
    device_model VARCHAR(100),
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE usage_records (
    id SERIAL PRIMARY KEY,
    office_id INTEGER REFERENCES offices(id),
    consumable_id INTEGER REFERENCES consumable_models(id),
    quantity INTEGER NOT NULL CHECK (quantity > 0),
    usage_date DATE NOT NULL,
    note TEXT,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE notes (
    id SERIAL PRIMARY KEY,
    title VARCHAR(100) NOT NULL,
    content TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE problem_types (
    id SERIAL PRIMARY KEY,
    device_type VARCHAR(50) NOT NULL,
    name VARCHAR(50) NOT NULL,
    sort_order INTEGER DEFAULT 0,
    is_default BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(device_type, name)
);

CREATE TABLE qr_codes (
    id SERIAL PRIMARY KEY,
    code VARCHAR(100) NOT NULL UNIQUE,
    office_id INTEGER REFERENCES offices(id),
    device_type VARCHAR(50),
    device_model VARCHAR(100),
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE tickets (
    id VARCHAR(36) PRIMARY KEY DEFAULT gen_random_uuid()::text,
    office_id INTEGER REFERENCES offices(id),
    device_type VARCHAR(50),
    device_model VARCHAR(100),
    problem_type VARCHAR(50) NOT NULL,
    description TEXT,
    contact VARCHAR(100),
    status VARCHAR(20) DEFAULT 'pending' CHECK (status IN ('pending', 'processing', 'completed')),
    consumable_used VARCHAR(100),
    consumable_quantity INTEGER,
    handled_by_user_id INTEGER REFERENCES users(id),
    handle_note TEXT,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

INSERT INTO consumable_models (name, unit, is_default)
VALUES ('联想LJ2600D/2605D/2400/2600通用墨盒', '个', TRUE);

INSERT INTO problem_types (device_type, name, sort_order, is_default) VALUES
('printer', '卡纸', 1, FALSE),
('printer', '无墨粉', 2, FALSE),
('printer', '不响应', 3, FALSE),
('printer', '其它', 99, TRUE),
('scanner', '不响应', 1, FALSE),
('scanner', '扫描模糊', 2, FALSE),
('scanner', '其它', 99, TRUE),
('computer', '不响应', 1, FALSE),
('computer', '运行缓慢', 2, FALSE),
('computer', '其它', 99, TRUE),
('copier', '卡纸', 1, FALSE),
('copier', '无墨粉', 2, FALSE),
('copier', '不响应', 3, FALSE),
('copier', '其它', 99, TRUE);