CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    username TEXT NOT NULL UNIQUE COLLATE NOCASE,
    password_hash TEXT NOT NULL,
    role TEXT NOT NULL CHECK (role IN ('admin', 'readonly')),
    created_at DATETIME DEFAULT (datetime('now')),
    updated_at DATETIME DEFAULT (datetime('now'))
);

CREATE TABLE IF NOT EXISTS consumable_models (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    unit TEXT DEFAULT '个',
    is_default INTEGER DEFAULT 0,
    created_at DATETIME DEFAULT (datetime('now'))
);

CREATE TABLE IF NOT EXISTS offices (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    room_number TEXT NOT NULL UNIQUE,
    device_type TEXT DEFAULT '',
    device_model TEXT DEFAULT '',
    created_at DATETIME DEFAULT (datetime('now'))
);

CREATE TABLE IF NOT EXISTS usage_records (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    office_id INTEGER NOT NULL REFERENCES offices(id),
    consumable_id INTEGER NOT NULL REFERENCES consumable_models(id),
    quantity INTEGER NOT NULL CHECK (quantity > 0),
    usage_date TEXT NOT NULL,
    note TEXT DEFAULT '',
    created_at DATETIME DEFAULT (datetime('now')),
    updated_at DATETIME DEFAULT (datetime('now'))
);

CREATE TABLE IF NOT EXISTS notes (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    title TEXT NOT NULL,
    content TEXT NOT NULL,
    created_at DATETIME DEFAULT (datetime('now')),
    updated_at DATETIME DEFAULT (datetime('now'))
);

CREATE TABLE IF NOT EXISTS problem_types (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    device_type TEXT NOT NULL,
    name TEXT NOT NULL,
    sort_order INTEGER DEFAULT 0,
    is_default INTEGER DEFAULT 0,
    created_at DATETIME DEFAULT (datetime('now')),
    UNIQUE(device_type, name)
);

CREATE TABLE IF NOT EXISTS qr_codes (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    code TEXT NOT NULL UNIQUE,
    office_id INTEGER REFERENCES offices(id),
    device_type TEXT DEFAULT '',
    device_model TEXT DEFAULT '',
    created_at DATETIME DEFAULT (datetime('now'))
);

CREATE TABLE IF NOT EXISTS tickets (
    id TEXT PRIMARY KEY,
    office_id INTEGER NOT NULL REFERENCES offices(id),
    device_type TEXT DEFAULT '',
    device_model TEXT DEFAULT '',
    problem_type TEXT NOT NULL,
    description TEXT DEFAULT '',
    contact TEXT DEFAULT '',
    status TEXT DEFAULT 'pending' CHECK (status IN ('pending', 'processing', 'completed')),
    consumable_used TEXT DEFAULT '',
    consumable_quantity INTEGER DEFAULT 0,
    handled_by_user_id INTEGER REFERENCES users(id),
    handle_note TEXT DEFAULT '',
    created_at DATETIME DEFAULT (datetime('now')),
    updated_at DATETIME DEFAULT (datetime('now'))
);

CREATE INDEX IF NOT EXISTS idx_usage_records_date ON usage_records(usage_date);
CREATE INDEX IF NOT EXISTS idx_usage_records_office ON usage_records(office_id);
CREATE INDEX IF NOT EXISTS idx_usage_records_consumable ON usage_records(consumable_id);
CREATE INDEX IF NOT EXISTS idx_tickets_status ON tickets(status);
CREATE INDEX IF NOT EXISTS idx_tickets_office ON tickets(office_id);
CREATE INDEX IF NOT EXISTS idx_notes_created ON notes(created_at);

INSERT OR IGNORE INTO consumable_models (name, unit, is_default) VALUES ('联想LJ2600D/2605D/2400/2600通用墨盒', '个', 1);

INSERT OR IGNORE INTO problem_types (device_type, name, sort_order, is_default) VALUES
('printer', '卡纸', 1, 0),
('printer', '无墨粉', 2, 0),
('printer', '不响应', 3, 0),
('printer', '其它', 99, 1),
('scanner', '不响应', 1, 0),
('scanner', '扫描模糊', 2, 0),
('scanner', '其它', 99, 1),
('computer', '不响应', 1, 0),
('computer', '运行缓慢', 2, 0),
('computer', '其它', 99, 1),
('copier', '卡纸', 1, 0),
('copier', '无墨粉', 2, 0),
('copier', '不响应', 3, 0),
('copier', '其它', 99, 1);
