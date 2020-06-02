package routes

import (
	configuration "coupanda/configuration"
	"reflect"

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

	// userErr := json.NewDecoder(c.Request.Body).Decode(&userData)

	email := c.PostForm("email")
	firstName := c.PostForm("firstName")
	lastName := c.PostForm("lastName")
	password := c.PostForm("password")
	mobile := c.PostForm("mobile")

	userData.Email = email
	userData.FirstName = firstName
	userData.LastName = lastName
	userData.Password = password
	userData.Mobile = mobile

	// if userErr != nil {

	// 	helper.RespondWithError(c, http.StatusBadRequest, CONSTANTS.EmptyReqBody)
	// 	return
	// }
	validateInputErr := helper.ValidateSignupInput(userData)
	if validateInputErr != "" {
		helper.RespondWithError(c, http.StatusBadRequest, validateInputErr)
		return
	}
	hashedPassword, hashError := bcrypt.GenerateFromPassword([]byte(userData.Password), bcrypt.DefaultCost)
	if hashError != nil {
		helper.RespondWithError(c, http.StatusBadRequest, hashError)
		return
	}
	userData.Password = string(hashedPassword)
	userData.Date = time.Now().Round(time.Millisecond).UnixNano() / (int64(time.Millisecond) / int64(time.Nanosecond))

	mongoSession := configuration.ConnectDb(config.Database)
	defer mongoSession.Close()

	sessionCopy := mongoSession.Copy()
	defer sessionCopy.Close()

	getCollection := sessionCopy.DB(config.Database).C("advertisers")

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

	file, FileErr := c.FormFile("file")
	if file == nil {
		helper.RespondWithError(c, http.StatusBadRequest, CONSTANTS.FileMissing)
		return
	}
	// fmt.Print(file.Filename)

	if FileErr != nil {

		helper.RespondWithError(c, http.StatusBadRequest, FileErr)
		return

	}
	fileName := "static/" + time.Now().Format("2006-01-02 15:04:05.000000") + file.Filename
	userData.ProfilePhoto = "/" + fileName
	uploadErr := c.SaveUploadedFile(file, fileName)
	if uploadErr != nil {
		helper.RespondWithError(c, http.StatusBadRequest, uploadErr)
		return
	}
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

	getCollection := sessionCopy.DB(config.Database).C("advertisers")

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

// CreateAdvertisement for advertiser create
func CreateAdvertisement(c *gin.Context) {

	var advertismentModel models.Advertisment

	val := reflect.ValueOf(c.Keys["user_id"])

	advertismentModel.ID = bson.NewObjectId()
	advertismentModel.Advertiser = bson.ObjectIdHex(val.String())
	addErr := json.NewDecoder(c.Request.Body).Decode(&advertismentModel)
	if addErr != nil {

		helper.RespondWithError(c, http.StatusBadRequest, addErr)
		return
	}

	mongoSession := configuration.ConnectDb(config.Database)
	defer mongoSession.Close()

	sessionCopy := mongoSession.Copy()
	defer sessionCopy.Close()

	getCollection := sessionCopy.DB(config.Database).C("advertisment")
	err := getCollection.Insert(advertismentModel)

	if err != nil {
		helper.RespondWithError(c, http.StatusBadRequest, err)
		return
	}

	helper.RespondWithSuccess(c, http.StatusOK, CONSTANTS.AdvertiseCreatedSuccssfully, advertismentModel)

}

// GetAdvertisements to get all advertisments
func GetAdvertisements(c *gin.Context) {

	resp := []bson.M{}
	mongoSession := configuration.ConnectDb(config.Database)
	defer mongoSession.Close()

	sessionCopy := mongoSession.Copy()
	defer sessionCopy.Close()

	getCollection := sessionCopy.DB(config.Database).C("advertisment")

	email := c.Query("email")

	// var e string
	// if email != "" {
	// 	e = email
	// }

	query := []bson.M{{
		"$lookup": bson.M{
			"from":         "advertisers",
			"localField":   "advertiser",
			"foreignField": "_id",
			"as":           "advertisers",
		}},
		{
			"$unwind": "$advertisers",
		},
	}

	if email != "" {

		query = append(query, bson.M{
			"$match": bson.M{
				"advertisers.email": email,
			},
		})
	}

	pipe := getCollection.Pipe(query)
	err := pipe.All(&resp)
	if err != nil {
		helper.RespondWithError(c, http.StatusBadRequest, err)
		return
	}

	helper.RespondWithSuccess(c, http.StatusOK, CONSTANTS.ListFetchedSuccess, resp)

}
