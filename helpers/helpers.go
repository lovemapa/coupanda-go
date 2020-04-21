package helper

import (
	"coupanda/models"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

var jwtKey = []byte("my_secret_key")

type claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

// CreateToken fore creating token
func CreateToken(_id string) (string, error) {

	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["user_id"] = _id
	atClaims["exp"] = time.Now().Add(time.Minute * 15).Unix()
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	token, err := at.SignedString([]byte("my_secret_key"))
	if err != nil {
		return "", err
	}
	return token, nil

}

//RespondWithError for sending errors
func RespondWithError(c *gin.Context, code int, message interface{}) {
	c.AbortWithStatusJSON(code, gin.H{
		"error":  message,
		"status": false,
	})
}

//RespondWithSuccess for success response
func RespondWithSuccess(c *gin.Context, code int, message string, data interface{}) {
	c.AbortWithStatusJSON(code, gin.H{
		"data":    data,
		"message": message,
		"status":  true,
	})
}

//ValidateSignupInput for validting user Input
func ValidateSignupInput(user models.UserSignup) string {
	var errMsg string

	// Check first_name if valid.
	if user.FirstName == "" {
		errMsg = "Please enter a valid first name."
		return errMsg
	}

	// Check last_name if valid.
	if user.LastName == "" {
		errMsg = "Please enter a valid last name."
		return errMsg

	}

	// Check email_id if valid.
	if user.Email == "" {
		errMsg = "Please enter a valid email id."
		return errMsg
	}

	//check password if valid
	if user.Password == "" {
		errMsg = "Please enter a valid password"
		return errMsg
	}

	//check password if valid
	if user.Mobile == "" {
		errMsg = "Please enter a valid mobile"
		return errMsg
	}
	return ""
}

// ValidateLoginInput for login Input
func ValidateLoginInput(user models.UserLogin) string {
	var errMsg string

	if user.Email == "" {
		errMsg = "Please provide email"
		return errMsg
	}

	if user.Password == "" {
		errMsg = "Please provide password"
		return errMsg
	}

	return errMsg
}
