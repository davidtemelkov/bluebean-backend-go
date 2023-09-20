package generalconstants

// DB constants
const (
	TableName      = "Bluebean"
	PK             = "PK"
	SK             = "SK"
	GSI1PK         = "GSI1PK"
	GSI1SK         = "GSI1SK"
	GSI1           = "GSI1"
	UserPrefix     = "USER#"
	FacilityPrefix = "FACILITY#"
	SpacePrefix    = "SPACE#"
	PunchPrefix    = "PUNCH#"
	PunchSKPrefix  = "PUNCH##"
	CommentPrefix  = "COMMENT#"
)

// Email regex expressions
const (
	EmailRX = "^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$"
)

// DateTime regex expressions
const (
	ISO8601 = `^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}Z$`
)

// Punch constants
const (
	StatusUnassigned = "Unassigned"
	StatusInProgress = "In progress"
	StatusCompleted  = "Completed"

	AssetNone = "None"
)
