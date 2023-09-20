package errorconstants

import "errors"

var (
	RecordNotFoundError           = errors.New("record not found")
	EditConflictError             = errors.New("edit conflict")
	DBConnectionError             = errors.New("Db connection error")
	InvalidJSONFormatError        = errors.New("Invalid JSON format")
	InvalidBase64ImagePrefixError = errors.New("Invalid base64 image format prefix")
	InternalServerError           = errors.New("Internal server error")
)

// Env errors
var (
	LoadingEnvFileError     = errors.New("Error loading .env file")
	FirebaseURLError        = errors.New("FIREBASE_URL environment variable is not set")
	FirebaseBucketNameError = errors.New("FIREBASE_BUCKET_NAME environment variable is not set")
	AWSAccessKeyError       = errors.New("AWS_ACCESS_KEY_ID environment variable is not set")
	AWSSecretKeyError       = errors.New("AWS_SECRET_KEY environment variable is not set")
	JWTPrivateKeyError      = errors.New("JWT_PRIVATE_KEY environment variable is not set")
	SMTPHostError           = errors.New("SMTP_HOST environment variable is not set")
	SMTPPortError           = errors.New("SMTP_PORT environment variable is not set")
	SMTPUsernameError       = errors.New("SMTP_USERNAME environment variable is not set")
	SMTPPasswordError       = errors.New("SMTP_PASSWORD environment variable is not set")
	SMTPSenderError         = errors.New("SMTP_SENDER environment variable is not set")
	WebAppBaseUrlError      = errors.New("WEB_APP_BASE_URL environment variable is not set")
)

// Authentication errors
var (
	MissingAuthorizationHeaderError       = errors.New("Missing authorization header")
	InvalidAuthorizationHeaderFormatError = errors.New("Invalid authorization header format")
	InvalidTokenError                     = errors.New("Invalid token")
	InvalidTokenClaimsError               = errors.New("Invalid token claims")
)

// User Firebase errors
var (
	FirebaseClientError  = errors.New("Failed to initialize Firebase Storage client")
	FileFolderEmptyError = errors.New("FileFolder is empty")
	FileNameEmptyError   = errors.New("FileName is empty")
)

// User errors
var (
	RequiredFieldError        = errors.New("Field is required")
	EmailFormatError          = errors.New("Email must be in the correct email format")
	PasswordMinLengthError    = errors.New("Password must be at least 8 symbols")
	PasswordMaxLengthError    = errors.New("Password must be less than 72 symbols")
	UserNameMinLengthError    = errors.New("Name must be at least 5 symbols")
	UserNameMaxLengthError    = errors.New("Name must be less than 50 symbols")
	UserNameNoWhitespaceError = errors.New("Must contain two names seperated by whitespace")
	RoleNotPermittedError     = errors.New("Role can only be Maintainer or Owner")
	UserIsNotAuthorizedError  = errors.New("User is not authorized")
	DuplicateEmailError       = errors.New("Duplicate email")
	UserNotFoundError         = errors.New("User not found")
	FailedLoginError          = errors.New("Invalid email or password")
)

// Facility errors
var (
	NameMinLengthError             = errors.New("Name must be at least 2 symbols")
	NameMaxLengthError             = errors.New("Name must be less than 50 symbols")
	AddressMinLengthError          = errors.New("Address must be at least 6 symbols")
	AddressMaxLengthError          = errors.New("Address must be less than 100 symbols")
	CityMinLengthError             = errors.New("City must be at least 3 symbols")
	CityMaxLengthError             = errors.New("City must be less than 100 symbols")
	UserAlreadyInFacilityError     = errors.New("User already exists in the facility")
	AssetAlreadyInFacilityError    = errors.New("Facility already contains asset")
	AssetNotInFacilityError        = errors.New("Facility doesn't contain asset")
	UserFacilityRelashionshipError = errors.New("User - Facility relationship does not exist")
	FailedToInsertFacilityError    = errors.New("Failed to insert facility")
)

// Space errors
var (
	SpaceNameMinLengthError     = errors.New("Name must be at least 2 symbols")
	SpaceNameMaxLengthError     = errors.New("Name must be less than 50 symbols")
	SpaceLocationMinLengthError = errors.New("Location must be at least 6 symbols")
	SpaceLocationMaxLengthError = errors.New("Location must be less than 100 symbols")
	FailedToInsertSpaceError    = errors.New("Failed to insert space")
)

// Punch errors
var (
	PunchTitleMinLengthError       = errors.New("Title must be at least 5 symbols")
	PunchTitleMaxLengthError       = errors.New("Title must be less than 100 symbols")
	PunchDescriptionMaxLengthError = errors.New("Description must be less than 500 symbols")
	PunchCoordXMinValueError       = errors.New("CoordX must be equal to or greater than 0")
	PunchCoordXMaxValueError       = errors.New("CoordX must be equal to or less than 100")
	PunchCoordYMinValueError       = errors.New("CoordY must be equal to or greater than 0")
	PunchCoordYMaxValueError       = errors.New("CoordY must be equal to or less than 100")
	InvalidDateTimeFormatError     = errors.New("Invalid datetime format")
	InvalidDateTimeRangeError      = errors.New("Invalid datetime range")
	InvalidPunchStatusError        = errors.New("Invalid punch status value")
	AssigneeIsNotMaintainerError   = errors.New("Assignee must be a maintainer in the facility")
	PunchNotExistError             = errors.New("Punch doesn't exist")
	FailedToInsertPunchError       = errors.New("Failed to insert punch")
)

// Comment errors
var (
	CommentTextMinLengthError  = errors.New("Text must be longer than 5 symbols")
	CommentTextMaxLengthError  = errors.New("Text must be shorter than 500 symbols")
	FailedToInsertCommentError = errors.New("Failed to insert comment")
)
