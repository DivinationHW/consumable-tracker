package models

type User struct {
	ID           int    `json:"id"`
	Username     string `json:"username"`
	PasswordHash string `json:"-"`
	Role         string `json:"role"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
}

type ConsumableModel struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	Unit       string `json:"unit"`
	IsDefault  bool   `json:"is_default"`
	CreatedAt  string `json:"created_at"`
}

type Office struct {
	ID          int    `json:"id"`
	RoomNumber  string `json:"room_number"`
	DeviceType  string `json:"device_type"`
	DeviceModel string `json:"device_model"`
	CreatedAt   string `json:"created_at"`
}

type UsageRecord struct {
	ID            int    `json:"id"`
	OfficeID      int    `json:"office_id"`
	ConsumableID  int    `json:"consumable_id"`
	Quantity      int    `json:"quantity"`
	UsageDate     string `json:"usage_date"`
	Note          string `json:"note"`
	CreatedAt     string `json:"created_at"`
	UpdatedAt     string `json:"updated_at"`
	OfficeName    string `json:"office_name,omitempty"`
	ConsumableName string `json:"consumable_name,omitempty"`
	DeviceType    string `json:"device_type,omitempty"`
	DeviceModel   string `json:"device_model,omitempty"`
}

type Note struct {
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Content   string `json:"content"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type ProblemType struct {
	ID         int    `json:"id"`
	DeviceType string `json:"device_type"`
	Name       string `json:"name"`
	SortOrder  int    `json:"sort_order"`
	IsDefault  bool   `json:"is_default"`
	CreatedAt  string `json:"created_at"`
}

type QRCode struct {
	ID          int    `json:"id"`
	Code        string `json:"code"`
	OfficeID    *int   `json:"office_id"`
	DeviceType  string `json:"device_type"`
	DeviceModel string `json:"device_model"`
	CreatedAt   string `json:"created_at"`
	OfficeName  string `json:"office_name,omitempty"`
}

type Ticket struct {
	ID               string `json:"id"`
	OfficeID         int    `json:"office_id"`
	DeviceType       string `json:"device_type"`
	DeviceModel      string `json:"device_model"`
	ProblemType      string `json:"problem_type"`
	Description      string `json:"description"`
	Contact          string `json:"contact"`
	Status           string `json:"status"`
	ConsumableUsed   string `json:"consumable_used"`
	ConsumableQty    int    `json:"consumable_quantity"`
	HandledByUserID  *int   `json:"handled_by_user_id"`
	HandleNote       string `json:"handle_note"`
	CreatedAt        string `json:"created_at"`
	UpdatedAt        string `json:"updated_at"`
	OfficeName       string `json:"office_name,omitempty"`
	HandledByUser    string `json:"handled_by_user,omitempty"`
}

type StatsSummary struct {
	TotalUsage     int            `json:"total_usage"`
	CurrentMonth   int            `json:"current_month"`
	TopOffice      string         `json:"top_office"`
	TopConsumable  string         `json:"top_consumable"`
	ByOffice       []StatsItem    `json:"by_office"`
	MonthlyTrend   []StatsItem    `json:"monthly_trend"`
}

type StatsItem struct {
	Label string `json:"label"`
	Value int    `json:"value"`
}

type BackupInfo struct {
	Filename    string `json:"filename"`
	Size        int64  `json:"size"`
	CreatedAt   string `json:"created_at"`
}
