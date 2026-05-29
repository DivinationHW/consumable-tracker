package models

import "time"

type User struct {
	ID           int       `json:"id"`
	Username     string    `json:"username"`
	PasswordHash string    `json:"-"`
	Role         string    `json:"role"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type UserCreate struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

type UserUpdate struct {
	Username *string `json:"username"`
	Password *string `json:"password"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type ChangePasswordRequest struct {
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}

type Consumable struct {
	ID         int       `json:"id"`
	Name       string    `json:"name"`
	Unit       string    `json:"unit"`
	IsDefault  bool      `json:"is_default"`
	CreatedAt  time.Time `json:"created_at"`
}

type ConsumableCreate struct {
	Name string `json:"name"`
	Unit string `json:"unit"`
}

type Office struct {
	ID          int       `json:"id"`
	RoomNumber  string    `json:"room_number"`
	DeviceType  *string   `json:"device_type"`
	DeviceModel *string   `json:"device_model"`
	CreatedAt   time.Time `json:"created_at"`
}

type OfficeCreate struct {
	RoomNumber  string  `json:"room_number"`
	DeviceType  *string `json:"device_type"`
	DeviceModel *string `json:"device_model"`
}

type OfficeUpdate struct {
	RoomNumber  *string `json:"room_number"`
	DeviceType  *string `json:"device_type"`
	DeviceModel *string `json:"device_model"`
}

type UsageRecord struct {
	ID            int       `json:"id"`
	OfficeID      int       `json:"office_id"`
	ConsumableID  int       `json:"consumable_id"`
	Quantity      int       `json:"quantity"`
	UsageDate     string    `json:"usage_date"`
	Note          string    `json:"note"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	OfficeNumber  string    `json:"office_number,omitempty"`
	ConsumableName string   `json:"consumable_name,omitempty"`
}

type UsageRecordCreate struct {
	OfficeID     int    `json:"office_id"`
	ConsumableID int    `json:"consumable_id"`
	Quantity     int    `json:"quantity"`
	UsageDate    string `json:"usage_date"`
	Note         string `json:"note"`
}

type Note struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type NoteCreate struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

type ProblemType struct {
	ID         int       `json:"id"`
	DeviceType string    `json:"device_type"`
	Name       string    `json:"name"`
	SortOrder  int       `json:"sort_order"`
	IsDefault  bool      `json:"is_default"`
	CreatedAt  time.Time `json:"created_at"`
}

type ProblemTypeCreate struct {
	DeviceType string `json:"device_type"`
	Name       string `json:"name"`
	SortOrder  int    `json:"sort_order"`
}

type QRCode struct {
	ID          int       `json:"id"`
	Code        string    `json:"code"`
	OfficeID    *int      `json:"office_id"`
	DeviceType  *string   `json:"device_type"`
	DeviceModel *string   `json:"device_model"`
	CreatedAt   time.Time `json:"created_at"`
	RoomNumber  string    `json:"room_number,omitempty"`
	IsConfigured bool     `json:"is_configured"`
}

type QRCodeCreate struct {
	Count int `json:"count"`
}

type QRCodeUpdate struct {
	OfficeID    *int    `json:"office_id"`
	DeviceType  *string `json:"device_type"`
	DeviceModel *string `json:"device_model"`
}

type Ticket struct {
	ID              string     `json:"id"`
	OfficeID        int        `json:"office_id"`
	DeviceType      string     `json:"device_type"`
	DeviceModel     string     `json:"device_model"`
	ProblemType     string     `json:"problem_type"`
	Description     string     `json:"description"`
	Contact         string     `json:"contact"`
	Status          string     `json:"status"`
	ConsumableUsed  string     `json:"consumable_used"`
	ConsumableQuantity int     `json:"consumable_quantity"`
	HandledByUserID *int       `json:"handled_by_user_id"`
	HandleNote      string     `json:"handle_note"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
	RoomNumber      string     `json:"room_number,omitempty"`
}

type TicketCreate struct {
	OfficeID    int    `json:"office_id"`
	ProblemType string `json:"problem_type"`
	Description string `json:"description"`
	Contact     string `json:"contact"`
}

type TicketComplete struct {
	ConsumableUsed  string `json:"consumable_used"`
	ConsumableQuantity int `json:"consumable_quantity"`
	HandleNote      string `json:"handle_note"`
}

type TicketStatusUpdate struct {
	Status string `json:"status"`
}

type StatsSummary struct {
	TotalUsage       int    `json:"total_usage"`
	CurrentMonthUsage int   `json:"current_month_usage"`
	MostUsedOffice   string `json:"most_used_office"`
	MostUsedOfficeCount int `json:"most_used_office_count"`
	DateRangeStart   string `json:"date_range_start"`
	DateRangeEnd     string `json:"date_range_end"`
	OfficeStats      []OfficeStat `json:"office_stats"`
	MonthlyStats     []MonthlyStat `json:"monthly_stats"`
}

type OfficeStat struct {
	OfficeID   int    `json:"office_id"`
	RoomNumber string `json:"room_number"`
	Total      int    `json:"total"`
}

type MonthlyStat struct {
	Month string `json:"month"`
	Year  int    `json:"year"`
	Total int    `json:"total"`
}

type BackupFile struct {
	Name string `json:"name"`
	Size int64  `json:"size"`
	Date string `json:"date"`
}

type BackupConfig struct {
	Frequency string `json:"frequency"`
	KeepDays  int    `json:"keep_days"`
}

type TokenResponse struct {
	Token    string `json:"token"`
	UserID   int    `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}
