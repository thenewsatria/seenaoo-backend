package messages

var FLASHCARD_NOT_FOUND_ERROR_MESSAGE = "Flashcard with provided id is not found"
var FLASHCARD_BODY_PARSER_ERROR_MESSAGE = "Error parsing flashcard body request"
var FLASHCARD_FAIL_TO_INSERT_ERROR_MESSAGE = "Failed to insert flashcard to database"
var FLASHCARD_FAIL_TO_FETCH_ERROR_MESSAGE = "Failed to fetch flashcard to database"
var FLASHCARD_FAIL_TO_POPULATE_HINTS_ERROR_MESSAGE = "failed to populate flashcard's hints"

var FLASHCARD_HINT_BODY_PARSER_ERROR_MESSAGE = "Error parsing hint body request"
var FLASHCARD_HINT_FAIL_TO_INSERT_ERROR_MESSAGE = "Error parsing hint body request"

var USER_BODY_PARSER_ERROR_MESSAGE = "Error parsing user body request"
var USER_FAIL_TO_HASH_PASSWORD = "Error hashing user password"
var USER_FAIL_TO_INSERT_ERROR_MESSAGE = "Failed to insert user to database"
var USER_EMAIL_ALREADY_USED = "Email already used"
var USER_USERNAME_ALREADY_USED = "Username already used"

var AUTH_FAIL_TO_GENERATE_TOKEN = "Failed to generate token"
var AUTH_BODY_PARSER_ERROR_MESSAGE = "Error parsing user credential body request"
var AUTH_EMAIL_INCORRECT = "Provided email is incorrect"
var AUTH_USERNAME_INCORRECT = "Provided username is incorrect"
var AUTH_PASSWORD_INCORRECT = "Provided password is incorrect"
