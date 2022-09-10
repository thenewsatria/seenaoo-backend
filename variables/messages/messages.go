package messages

var FLASHCARD_NOT_FOUND_ERROR_MESSAGE = "Flashcard with provided id is not found"
var FLASHCARD_BODY_PARSER_ERROR_MESSAGE = "Error parsing flashcard body request"
var FLASHCARD_FAIL_TO_INSERT_ERROR_MESSAGE = "Failed to insert flashcard to database"
var FLASHCARD_FAIL_TO_FETCH_ERROR_MESSAGE = "Failed to fetch flashcard from database"
var FLASHCARD_FAIL_TO_POPULATE_HINTS_ERROR_MESSAGE = "Failed to populate flashcard's hints"
var FLASHCARD_FAIL_TO_UPDATE_ERROR_MESSAGE = "Failed to update flashcard"
var FLASHCARD_FAIL_TO_DELETE_ERROR_MESSAGE = "Failed to delete flashcard"
var FLASHCARD_HAS_EMPTY_HINTS_ERROR_MESSAGE = "Flashcard has empty hints"

var FLASHCARD_HINT_NOT_FOUND_ERROR_MESSAGE = "Flashcard hint with provided id is not found"
var FLASHCARD_HINT_BODY_PARSER_ERROR_MESSAGE = "Error parsing hint body request"
var FLASHCARD_HINT_FAIL_TO_INSERT_ERROR_MESSAGE = "Failed to insert flashcard hint to database"
var FLASHCARD_HINT_FAIL_TO_FETCH_ERROR_MESSAGE = "Failed to fetch flashcard hint from database"
var FLASHCARD_HINT_FAIL_TO_UPDATE_ERROR_MESSAGE = "Failed to update flashcard hint"
var FLASHCARD_HINT_FAIL_TO_DELETE_ERROR_MESSAGE = "Failed to delete flashcard hint"

var USER_BODY_PARSER_ERROR_MESSAGE = "Error parsing user body request"
var USER_FAIL_TO_HASH_PASSWORD_ERROR_MESSAGE = "Error hashing user password"
var USER_FAIL_TO_INSERT_ERROR_MESSAGE = "Failed to insert user to database"
var USER_EMAIL_ALREADY_USED_ERROR_MESSAGE = "Email already used"
var USER_USERNAME_ALREADY_USED_ERROR_MESSAGE = "Username already used"
var USER_USERNAME_NOT_FOUND_ERROR_MESSAGE = "user with speficied username is not found"
var USER_FAIL_TO_FETCH_ERROR_MESSAGE = "Failed to fetch user hint from database"

var AUTH_FAIL_TO_GENERATE_ACCESS_TOKEN_ERROR_MESSAGE = "Failed to generate access token"
var AUTH_FAIL_TO_GENERATE_REFRESH_TOKEN_ERROR_MESSAGE = "Failed to generate refresh token"
var AUTH_BODY_PARSER_ERROR_MESSAGE = "Error parsing user credential body request"
var AUTH_EMAIL_INCORRECT_ERROR_MESSAGE = "Provided email is incorrect"
var AUTH_USERNAME_INCORRECT_ERROR_MESSAGE = "Provided username is incorrect"
var AUTH_PASSWORD_INCORRECT_ERROR_MESSAGE = "Provided password is incorrect"
var AUTH_TOKEN_NOT_EXIST_ERROR_MESSAGE = "Token is not exist"
var AUTH_TOKEN_INVALID_ERROR_MESSAGE = "Token Invalid, try login again"
var AUTH_ACCESS_TOKEN_EXPIRED_ERROR_MESSAGE = "Access token expired please refresh access token"
var AUTH_REFRESH_TOKEN_EXPIRED_ERROR_MESSAGE = "Refresh token expired please login again"
var AUTH_REFRESH_TOKEN_BLOCKED_ERROR_MESSAGE = "Refresh token blocked please refresh access token"
var AUTH_STORED_REFRESH_TOKEN_IS_EMPTY_ERROR_MESSAGE = "Stored refresh token is empty"
var AUTH_REFRESH_TOKEN_DIFFERENT_FROM_STORED_ERROR_MESSAGE = "Refresh token given is different from the storedc refresh token"
var AUTH_COLLABORATION_UNAUTHORIZED_ERROR_MESSAGE = "Unauthorized to access this collaboration"
var AUTH_FLASHCARD_COVER_UNATHORIZED_ERROR_MESSAGE = "Unauthorized to access this flashcard cover"
var AUTH_FLASHCARD_HINT_UNATHORIZED_ERROR_MESSAGE = "Unauthorized to access this flashcard hint"
var AUTH_FLASHCARD_UNAUTHORIZED_ERROR_MESSAGE = "Unauthorized to access this flashcard"

var REFRESH_ACCESS_TOKEN_BODY_PARSER_ERROR_MESSAGE = "Error parsing access token body request"
var REFRESH_TOKEN_FAIL_TO_INSERT_ERROR_MESSAGE = "Failed to insert refresh token to database"
var REFRESH_TOKEN_USERNAME_NOT_FOUND_ERROR_MESSAGE = "Refresh token with specified username is not found"
var REFRESH_TOKEN_FAIL_TO_FETCH_ERROR_MESSAGE = "Failed to fetch refresh token from database"
var REFRESH_TOKEN_FAIL_TO_UPDATE_STORED_TOKEN_ERROR_MESSAGE = "Failed to update stored refresh token"

var COLLABORATION_BODY_PARSER_ERROR_MESSAGE = "Error parsing collaboration body request"
var COLLABORATION_FAIL_TO_INSERT_ERROR_MESSAGE = "Failed to insert collaboration to database"
var COLLABORATION_NOT_FOUND_ERROR_MESSAGE = "Collaboration with provided id is not found"
var COLLABORATION_INVITER_DOESNT_EXIST_ERROR_MESSAGE = "Can't send collaboration invite, inviter doesn't exist"
var COLLABORATION_COLLABORATOR_DOESNT_EXIST_ERROR_MESSAGE = "Can't send collaboration invite, collaborator doesn't exist"
var COLLABORATION_ITEM_TYPE_IS_UNKNOWN = "Collaboration item type is unknown"
var COLLABORATION_FAIL_TO_UPDATE_ERROR_MESSAGE = "Failed to update collaboration"
var COLLABORATION_FAIL_TO_FETCH_ERROR_MESSAGE = "Failed to fetch collaboration"
var COLLABORATION_FAIL_TO_DELETE_ERROR_MESSAGE = "Failed to fetch collaboration"

var FLASHCARD_COVER_NOT_FOUND_ERROR_MESSAGE = "Flashcard cover with provided id is not found"
var FLASHCARD_COVER_FAIL_TO_FETCH_ERROR_MESSAGE = "Failed to fetch flashcard from database"
var FLASHCARD_COVER_BODY_PARSER_ERROR_MESSAGE = "Error parsing flashcard cover body request"
var FLASHCARD_COVER_FAIL_TO_INSERT_ERROR_MESSAGE = "Failed to insert flashcard from database"
var FLASHCARD_COVER_FAIL_TO_UPDATE_ERROR_MESSAGE = "Failed to update flashcard from database"
var FLASHCARD_COVER_FAIL_TO_DELETE_ERROR_MESSAGE = "Failed to delete flashcard from database"
var FLASHCARD_COVER_FAIL_TO_POPULATE_FLASHCARDS_ERROR_MESSAGE = "Failed to populate flashcard covers's hints"

var TAG_FAIL_TO_FETCH_ERROR_MESSAGE = "Failed to fetch tag from database"
var TAG_FAIL_TO_INSERT_ERROR_MESSAGE = "Failed to insert tag from database"
var TAG_NOT_FOUND_ERROR_MESSAGE = "Tag with provided id is not found"

var MIDDLEWARE_ISAUTHOR_UNKNOWN_SERVICE_TYPE_ERROR_MESSAGE = "[isAuthor Middleware]: Unknown service type"
