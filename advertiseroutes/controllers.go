package routes

import (
	configuration "coupanda/configuration"

	CONSTANTS "coupanda/constant"
	helper "coupanda/helpers"
	"coupanda/models"
	"encoding/json"
	"instituteNew/config"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

//Register advertisment
func Register(c *gin.Context) {
	var userData models.UserSignup

	userErr := json.NewDecoder(c.Request.Body).Decode(&userData)
	if userErr != nil {

		helper.RespondWithError(c, http.StatusBadRequest, userErr)
	}
	validateInputErr := helper.ValidateSignupInput(userData)
	if validateInputErr != "" {
		helper.RespondWithError(c, http.StatusBadRequest, validateInputErr)
	}
	hashedPassword, hashError := bcrypt.GenerateFromPassword([]byte(userData.Password), bcrypt.DefaultCost)
	if hashError != nil {
		helper.RespondWithError(c, http.StatusBadRequest, hashError)
	}
	userData.Password = string(hashedPassword)
	userData.Date = time.Now().Round(time.Millisecond).UnixNano() / (int64(time.Millisecond) / int64(time.Nanosecond))

	mongoSession := configuration.ConnectDb(config.Database)
	defer mongoSession.Close()

	sessionCopy := mongoSession.Copy()
	defer sessionCopy.Close()

	getCollection := sessionCopy.DB(config.Database).C("advertisement")

	index := mgo.Index{
		Key:        []string{"email"},
		Unique:     true,
		DropDups:   true,
		Background: true,
		Sparse:     true,
	}

	indexerr := getCollection.EnsureIndex(index)

	if indexerr != nil {
		helper.RespondWithError(c, http.StatusBadRequest, indexerr)
		return
	}

	userData.ID = bson.NewObjectId()

	token, _ := helper.CreateToken(userData.ID.Hex())
	userData.Token = token

	err := getCollection.Insert(userData)

	if err != nil {
		if mgo.IsDup(err) == true {

			helper.RespondWithError(c, http.StatusBadRequest, CONSTANTS.AccountAlreadyExists)
			return
		}
		helper.RespondWithError(c, http.StatusBadRequest, err)
		return
	}

	helper.RespondWithSuccess(c, http.StatusOK, CONSTANTS.CreatedSuccssfully, userData)

}

// Login for advertiser login
func Login(c *gin.Context) {
	var Login models.UserLogin
	var userData models.UserSignup
	userErr := json.NewDecoder(c.Request.Body).Decode(&Login)
	if userErr != nil {

		helper.RespondWithError(c, http.StatusBadRequest, userErr)
	}
	LoginErr := helper.ValidateLoginInput(Login)
	if LoginErr != "" {
		helper.RespondWithError(c, http.StatusBadRequest, LoginErr)
	}

	mongoSession := configuration.ConnectDb(config.Database)
	defer mongoSession.Close()

	sessionCopy := mongoSession.Copy()
	defer sessionCopy.Close()

	getCollection := sessionCopy.DB(config.Database).C("advertisement")

	err := getCollection.Find(bson.M{"email": Login.Email}).One(&userData)
	if err != nil {
		helper.RespondWithError(c, http.StatusBadRequest, CONSTANTS.AccountNotExists)
		return
	}

	PassErr := bcrypt.CompareHashAndPassword([]byte(userData.Password), []byte(Login.Password))
	if PassErr != nil {
		helper.RespondWithError(c, http.StatusBadRequest, CONSTANTS.IncorrectPassword)
		return
	}
	token, _ := helper.CreateToken(userData.ID.Hex())
	userData.Token = token
	err = getCollection.UpdateId(bson.ObjectIdHex(userData.ID.Hex()), userData)
	helper.RespondWithSuccess(c, http.StatusOK, CONSTANTS.LoggedInSuccess, userData)

}

// CreateTokenAdvertisement for advertiser create
func CreateTokenAdvertisement(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"claims": c.Keys["user_id"],
	})

}
