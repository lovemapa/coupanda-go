package middleware

import (
	"coupanda/configuration"
	"coupanda/models"
	"fmt"
	"instituteNew/config"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2/bson"
)

func hello(auths ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Println("Print from  middlware")
		c.Next()
	}
}

func respondWithError(c *gin.Context, code int, message interface{}) {
	c.AbortWithStatusJSON(code, gin.H{
		"error":  message,
		"status": false})
}

// TokenAuthMiddleware for auth
func TokenAuthMiddleware() gin.HandlerFunc {
	// requiredToken := "pawan"

	// if requiredToken == "" {
	// 	log.Fatal("Please set API_TOKEN environment variable")
	// }
	fmt.Print("YES")
	return func(c *gin.Context) {

		if c.Request.Header["Token"] == nil {
			respondWithError(c, 401, "Authorization Missing")
			return
		}
		t := c.Request.Header["Token"][0]

		if t == "" {
			respondWithError(c, 401, "Auth token required")
			return
		}

		tkn, _ := jwt.Parse(t, func(token *jwt.Token) (interface{}, error) {

			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				respondWithError(c, 401, "unexpected signing method")
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte("my_secret_key"), nil
		})

		var userData models.UserSignup
		mongoSession := configuration.ConnectDb(config.Database)
		defer mongoSession.Close()

		sessionCopy := mongoSession.Copy()
		defer sessionCopy.Close()

		getCollection := sessionCopy.DB(config.Database).C("advertisement")

		queryErr := getCollection.Find(bson.M{"token": t}).One(&userData)

		if queryErr != nil {

			respondWithError(c, 401, "Token is not correct")
			return
		}

		claims, ok := tkn.Claims.(jwt.MapClaims)

		if !ok || !tkn.Valid {
			respondWithError(c, 401, "Unauthorized access")
			return
		}
		c.Set("user_id", claims)
		c.Next()
	}
}

//VerifyToken token
func VerifyToken(t string) (jwt.MapClaims, bool) {
	tkn, _ := jwt.Parse(t, func(token *jwt.Token) (interface{}, error) {

		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte("my_secret_key"), nil
	})
	claims, ok := tkn.Claims.(jwt.MapClaims)
	if ok && tkn.Valid {

		return claims, ok
	}
	return claims, ok
}
