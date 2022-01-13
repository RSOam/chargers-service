package chargers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-kit/log"
	consulapi "github.com/hashicorp/consul/api"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type database struct {
	db     *mongo.Database
	logger log.Logger
	consul consulapi.Client
}

func NewDatabase(db *mongo.Database, logger log.Logger, consul consulapi.Client) ChargerDB {
	return &database{
		db:     db,
		logger: log.With(logger, "database", "mongoDB"),
		consul: consul,
	}
}

func (dat *database) CreateCharger(ctx context.Context, charger Charger) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := dat.db.Collection("Chargers").InsertOne(ctx, charger)
	if err != nil {
		dat.logger.Log("Error inserting charger into DB: ", err.Error())
		return err
	}
	return nil
}

/*
func (dat *database) GetCharger(ctx context.Context, id string) (Charger, error) {
	tempCharger := Charger{}
	objectID, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		dat.logger.Log("Error getting charger from DB: ", err.Error())
		return tempCharger, err
	}
	val, _ := getConsulValue(dat.consul, dat.logger, "commratService")
	val2, _ := getConsulValue(dat.consul, dat.logger, "reservationsService")
	ratings, err := getChargerRatings(val, dat.logger, id)
	if err != nil {
		dat.logger.Log("Error getting charger from DB: ", err.Error())
		return tempCharger, err
	}
	comments, err := getChargerComments(val, dat.logger, id)
	if err != nil {
		dat.logger.Log("Error getting charger from DB: ", err.Error())
		return tempCharger, err
	}
	reservations, err := getChargerReservations(val2, dat.logger, id)
	if err != nil {
		dat.logger.Log("Error getting charger from DB: ", err.Error())
		return tempCharger, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = dat.db.Collection("Chargers").FindOne(ctx, bson.M{"_id": objectID}).Decode(&tempCharger)
	if err != nil {
		dat.logger.Log("Error getting charger from DB: ", err.Error())
		return tempCharger, err
	}
	tempCharger.Ratings = ratings
	tempCharger.Comments = comments
	tempCharger.Reservations = reservations

	return tempCharger, nil
}
*/

func (dat *database) GetCharger(ctx context.Context, id string) (Charger, ChargerExtra, error) {
	//ASYNC
	tempCharger := Charger{}
	tempExtra := ChargerExtra{}
	objectID, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		dat.logger.Log("Error getting charger from DB: ", err.Error())
		return tempCharger, tempExtra, err
	}
	ratChan := make(chan RatResult)
	comChan := make(chan ComResult)
	resChan := make(chan ResResult)
	val, _ := getConsulValue(dat.consul, dat.logger, "commratService")
	val2, _ := getConsulValue(dat.consul, dat.logger, "reservationsService")

	go getChargerRatingsAsync(val, dat.logger, id, ratChan)
	dat.logger.Log("Async getting chargers")
	go getChargerCommentsAsync(val, dat.logger, id, comChan)
	dat.logger.Log("Async getting comments")
	go getChargerReservationsAsync(val2, dat.logger, id, resChan)
	dat.logger.Log("Async getting reservations")

	ratRes := <-ratChan
	dat.logger.Log("Async got chargers")
	if ratRes.err != nil {
		dat.logger.Log("Error getting charger from DB: ", ratRes.err.Error())
		return tempCharger, tempExtra, ratRes.err
	}
	comRes := <-comChan
	dat.logger.Log("Async got comments")
	if err != nil {
		dat.logger.Log("Error getting charger from DB: ", comRes.err.Error())
		return tempCharger, tempExtra, comRes.err
	}
	resRes := <-resChan
	dat.logger.Log("Async got reservations")
	if resRes.err != nil {
		dat.logger.Log("Error getting charger from DB: ", resRes.err.Error())
		return tempCharger, tempExtra, resRes.err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = dat.db.Collection("Chargers").FindOne(ctx, bson.M{"_id": objectID}).Decode(&tempCharger)
	if err != nil {
		dat.logger.Log("Error getting charger from DB: ", err.Error())
		return tempCharger, tempExtra, err
	}
	city, temp, weather, err := callAPIs(tempCharger.Location, dat.logger, dat.consul)
	tempCharger.Ratings = ratRes.ratings
	tempCharger.Comments = comRes.comments
	tempCharger.Reservations = resRes.reservations
	tempExtra.City = city
	tempExtra.Weather = weather
	tempExtra.Temp = temp
	return tempCharger, tempExtra, nil
}

func (dat *database) DeleteCharger(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		dat.logger.Log("Error deleting charger from DB: ", err.Error())
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	filter := bson.M{"_id": objectID}
	res := dat.db.Collection("Chargers").FindOneAndDelete(ctx, filter)
	if res.Err() == mongo.ErrNoDocuments {
		dat.logger.Log("Error deleting charger from DB: ", err.Error())
		return err
	}
	return nil
}
func (dat *database) GetChargers(ctx context.Context) ([]Charger, error) {
	tempCharger := Charger{}
	tempChargers := []Charger{}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	cursor, err := dat.db.Collection("Chargers").Find(ctx, bson.D{})
	if err != nil {
		dat.logger.Log("Error getting chargers from DB: ", err.Error())
		return tempChargers, err
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		err := cursor.Decode(&tempCharger)
		if err != nil {
			dat.logger.Log("Error getting chargers from DB: ", err.Error())
			return tempChargers, err
		}
		tempChargers = append(tempChargers, tempCharger)
	}
	return tempChargers, nil
}
func (dat *database) UpdateCharger(ctx context.Context, id string, charger Charger) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		dat.logger.Log("Error updating charger: ", err.Error())
		return err
	}
	update := bson.M{
		"$set": bson.M{
			"averageRating": charger.AverageRating,
			"modified":      time.Now().Format(time.RFC3339),
		},
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err = dat.db.Collection("Chargers").UpdateOne(ctx, bson.M{"_id": objectID}, update)
	if err != nil {
		dat.logger.Log("Error updating charger: ", err.Error())
		return err
	}

	return nil
}
func getChargerRatings(commratAddr string, logger log.Logger, chargerID string) ([]Rating, error) {
	requestBody, err := json.Marshal(GetChargerRatingsRequest{})
	tempResponse := GetChargerRatingsResponse{}
	tempRatings := []Rating{}
	if err != nil {
		return tempRatings, err
	}
	client := &http.Client{}
	commratAddr = "http://20.124.136.184:8080"
	commratUri := commratAddr + "/ratings/"
	req, err := http.NewRequest(http.MethodGet, commratUri+"?charger="+chargerID, bytes.NewBuffer(requestBody))
	if err != nil {
		return tempRatings, err
	}
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	resp, err := client.Do(req)
	if err != nil {
		return tempRatings, err
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&tempResponse)

	tempRatings = tempResponse.Ratings
	if err != nil {
		return tempRatings, err
	}
	client.CloseIdleConnections()
	return tempRatings, nil
}

type RatResult struct {
	ratings []Rating
	err     error
}

func getChargerRatingsAsync(commratAddr string, logger log.Logger, chargerID string, response chan RatResult) {
	requestBody, err := json.Marshal(GetChargerRatingsRequest{})
	tempResponse := GetChargerRatingsResponse{}
	tempRatings := []Rating{}
	tempRR := RatResult{}
	if err != nil {
		tempRR.ratings = tempRatings
		tempRR.err = err
		response <- tempRR
	}
	client := &http.Client{}

	commratUri := commratAddr + "/ratings/"
	req, err := http.NewRequest(http.MethodGet, commratUri+"?charger="+chargerID, bytes.NewBuffer(requestBody))
	if err != nil {
		tempRR.ratings = tempRatings
		tempRR.err = err
		response <- tempRR
	}
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	resp, err := client.Do(req)
	if err != nil {
		tempRR.ratings = tempRatings
		tempRR.err = err
		response <- tempRR
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&tempResponse)

	tempRatings = tempResponse.Ratings
	if err != nil {
		tempRR.ratings = tempRatings
		tempRR.err = err
		response <- tempRR
	}
	client.CloseIdleConnections()
	tempRR.ratings = tempRatings
	tempRR.err = nil
	response <- tempRR
}
func getChargerComments(commratAddr string, logger log.Logger, chargerID string) ([]Comment, error) {
	requestBody, err := json.Marshal(GetChargerCommentsRequest{})
	tempResponse := GetChargerCommentsResponse{}
	tempComments := []Comment{}
	if err != nil {
		return tempComments, err
	}
	client := &http.Client{}
	commratUri := commratAddr + "/comments/"
	req, err := http.NewRequest(http.MethodGet, commratUri+"?charger="+chargerID, bytes.NewBuffer(requestBody))
	if err != nil {
		return tempComments, err
	}
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	resp, err := client.Do(req)
	if err != nil {
		return tempComments, err
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&tempResponse)

	tempComments = tempResponse.Comments
	if err != nil {
		return tempComments, err
	}
	client.CloseIdleConnections()
	return tempComments, nil
}

type ComResult struct {
	comments []Comment
	err      error
}

func getChargerCommentsAsync(commratAddr string, logger log.Logger, chargerID string, response chan ComResult) {
	requestBody, err := json.Marshal(GetChargerCommentsRequest{})
	tempResponse := GetChargerCommentsResponse{}
	tempComments := []Comment{}
	tempCR := ComResult{}
	if err != nil {
		tempCR.err = err
		tempCR.comments = tempComments
		response <- tempCR
	}
	client := &http.Client{}
	commratUri := commratAddr + "/comments/"
	req, err := http.NewRequest(http.MethodGet, commratUri+"?charger="+chargerID, bytes.NewBuffer(requestBody))
	if err != nil {
		tempCR.err = err
		tempCR.comments = tempComments
		response <- tempCR
	}
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	resp, err := client.Do(req)
	if err != nil {
		tempCR.err = err
		tempCR.comments = tempComments
		response <- tempCR
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&tempResponse)

	tempComments = tempResponse.Comments
	if err != nil {
		tempCR.err = err
		tempCR.comments = tempComments
		response <- tempCR
	}
	client.CloseIdleConnections()
	tempCR.err = nil
	tempCR.comments = tempComments
	response <- tempCR
}
func getChargerReservations(reservationsAddr string, logger log.Logger, chargerID string) ([]Reservation, error) {
	requestBody, err := json.Marshal(GetChargerReservationsRequest{})
	tempResponse := GetChargerReservationsResponse{}
	tempReservations := []Reservation{}
	if err != nil {
		return tempReservations, err
	}
	client := &http.Client{}
	reservationsUri := reservationsAddr + "/reservations/"
	req, err := http.NewRequest(http.MethodGet, reservationsUri+"?charger="+chargerID, bytes.NewBuffer(requestBody))
	if err != nil {
		return tempReservations, err
	}
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	resp, err := client.Do(req)
	if err != nil {
		return tempReservations, err
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&tempResponse)

	tempReservations = tempResponse.Reservations
	if err != nil {
		return tempReservations, err
	}
	client.CloseIdleConnections()
	return tempReservations, nil
}

type ResResult struct {
	reservations []Reservation
	err          error
}

func getChargerReservationsAsync(reservationsAddr string, logger log.Logger, chargerID string, response chan ResResult) {
	requestBody, err := json.Marshal(GetChargerReservationsRequest{})
	tempResponse := GetChargerReservationsResponse{}
	tempReservations := []Reservation{}
	tempRR := ResResult{}
	if err != nil {
		tempRR.err = err
		tempRR.reservations = tempReservations
		response <- tempRR
	}
	client := &http.Client{}
	reservationsUri := reservationsAddr + "/reservations/"
	req, err := http.NewRequest(http.MethodGet, reservationsUri+"?charger="+chargerID, bytes.NewBuffer(requestBody))
	if err != nil {
		tempRR.err = err
		tempRR.reservations = tempReservations
		response <- tempRR
	}
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	resp, err := client.Do(req)
	if err != nil {
		tempRR.err = err
		tempRR.reservations = tempReservations
		response <- tempRR
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&tempResponse)

	tempReservations = tempResponse.Reservations
	if err != nil {
		tempRR.err = err
		tempRR.reservations = tempReservations
		response <- tempRR
	}
	client.CloseIdleConnections()
	tempRR.err = nil
	tempRR.reservations = tempReservations
	response <- tempRR
}
func getConsulValue(consul consulapi.Client, logger log.Logger, key string) (string, error) {
	kv := consul.KV()
	keyPair, _, err := kv.Get(key, nil)
	if err != nil {
		logger.Log("msg", "Failed getting consul key")
		return "", err
	}
	return string(keyPair.Value), nil
}

type geoAPI struct {
	Data []geoAPIdata `json:"data"`
}
type geoAPIdata struct {
	Name string `json:"name"`
}

type weatherAPI struct {
	Current weatherAPIcurrent
}

type weatherAPIcurrent struct {
	Temp      float64             `json:"temp_c"`
	Condition weatherAPIcondition `json:"condition"`
}
type weatherAPIcondition struct {
	Weather string `json:"text"`
	Icon    string `json:"icon"`
}

func callAPIs(location Location, logger log.Logger, consul consulapi.Client) (string, float64, weatherAPIcondition, error) {

	geoURL, _ := getConsulValue(consul, logger, "geoAPIurl")
	weaURL, _ := getConsulValue(consul, logger, "weatherAPIurl")
	rapidKey, _ := getConsulValue(consul, logger, "rapidAPIkey")

	lat := fmt.Sprintf("%f", location.Latitude)
	long := fmt.Sprintf("%f", location.Longitude)
	url := "https://" + geoURL + "/v1/geo/locations/" + lat + "+" + long + "/nearbyCities?radius=10"
	tmpCond := weatherAPIcondition{}
	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("x-rapidapi-host", geoURL)
	req.Header.Add("x-rapidapi-key", rapidKey)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		logger.Log("msg", "Problem getting Geo API")
		return "", 0, tmpCond, err
	}
	defer res.Body.Close()
	tmpGeo := geoAPI{}
	err = json.NewDecoder(res.Body).Decode(&tmpGeo)
	if err != nil {
		logger.Log("msg", "Problem getting Geo API")
		return "", 0, tmpCond, err
	}

	url = "https://" + weaURL + "/current.json?q=" + lat + "%2C" + long

	req, _ = http.NewRequest("GET", url, nil)

	req.Header.Add("x-rapidapi-host", weaURL)
	req.Header.Add("x-rapidapi-key", rapidKey)

	res, err = http.DefaultClient.Do(req)
	if err != nil {
		logger.Log("msg", "Problem getting Weather API")
		return "", 0, tmpCond, err
	}
	defer res.Body.Close()
	tmpWea := weatherAPI{}
	err = json.NewDecoder(res.Body).Decode(&tmpWea)
	if err != nil {
		logger.Log("msg", "Problem getting Weather API")
		return "", 0, tmpCond, err
	}
	return tmpGeo.Data[0].Name, tmpWea.Current.Temp, tmpWea.Current.Condition, nil
}
