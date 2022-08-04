package messages

var FLASHCARD_NOT_FOUND_ERROR_MESSAGE = "Flashcard with provided id is not found"
var FLASHCARD_BODY_PARSER_ERROR_MESSAGE = "Error parsing flashcard body request"
var FLASHCARD_FAIL_TO_INSERT_ERROR_MESSAGE = "Failed to insert flashcard to database"
var FLASHCARD_FAIL_TO_FETCH_ERROR_MESSAGE = "Failed to fetch flashcard from database"
var FLASHCARD_FAIL_TO_POPULATE_HINTS_ERROR_MESSAGE = "failed to populate flashcard's hints"

var FLASHCARD_HINT_BODY_PARSER_ERROR_MESSAGE = "Error parsing hint body request"
var FLASHCARD_HINT_FAIL_TO_INSERT_ERROR_MESSAGE = "Error parsing hint body request"

var USER_BODY_PARSER_ERROR_MESSAGE = "Error parsing user body request"
var USER_FAIL_TO_HASH_PASSWORD_ERROR_MESSAGE = "Error hashing user password"
var USER_FAIL_TO_INSERT_ERROR_MESSAGE = "Failed to insert user to database"
var USER_EMAIL_ALREADY_USED_ERROR_MESSAGE = "Email already used"
var USER_USERNAME_ALREADY_USED_ERROR_MESSAGE = "Username already used"
var USER_USERNAME_NOT_FOUND_ERROR_MESSAGE = "user with speficied username is not found"

var AUTH_FAIL_TO_GENERATE_ACCESS_TOKEN_ERROR_MESSAGE = "Failed to generate access token"
var AUTH_FAIL_TO_GENERATE_REFRESH_TOKEN_ERROR_MESSAGE = "Failed to generate refresh token"
var AUTH_BODY_PARSER_ERROR_MESSAGE = "Error parsing user credential body request"
var AUTH_EMAIL_INCORRECT_ERROR_MESSAGE = "Provided email is incorrect"
var AUTH_USERNAME_INCORRECT_ERROR_MESSAGE = "Provided username is incorrect"
var AUTH_PASSWORD_INCORRECT_ERROR_MESSAGE = "Provided password is incorrect"
var AUTH_TOKEN_NOT_EXIST_ERROR_MESSAGE = "Token is not exist"
var AUTH_TOKEN_INVALID_ERROR_MESSAGE = "Token Invalid"
var AUTH_TOKEN_EXPIRED_ERROR_MESSAGE = "Token expired please login again"

var REFRESH_TOKEN_FAIL_TO_INSERT_ERROR_MESSAGE = "Failed to insert refresh token to database"
var REFRESH_TOKEN_USERNAME_NOT_FOUND_ERROR_MESSAGE = "Refresh token with specified username is not found"
var REFRESH_TOKEN_FAIL_TO_FETCH_ERROR_MESSAGE = "Failed to fetch refresh token from database"
