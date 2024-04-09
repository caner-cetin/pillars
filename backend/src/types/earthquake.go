package types

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Metadata struct {
	Generated int64  `json:"generated"`
	URL       string `json:"url"`
	Title     string `json:"title"`
	Status    int    `json:"status"`
	API       string `json:"api"`
	Count     int    `json:"count"`
}

type USGSEarthquake struct {
	Type     string   `json:"type"`
	Metadata Metadata `json:"metadata"`
	Features []struct {
		Type       string `json:"type"`
		Properties struct {
			Mag     float64 `json:"mag"`
			Place   string  `json:"place"`
			Time    int64   `json:"time"`
			Updated int64   `json:"updated"`
			Tz      any     `json:"tz"`
			URL     string  `json:"url"`
			Detail  string  `json:"detail"`
			Felt    any     `json:"felt"`
			Cdi     any     `json:"cdi"`
			Mmi     any     `json:"mmi"`
			Alert   any     `json:"alert"`
			Status  string  `json:"status"`
			Tsunami int     `json:"tsunami"`
			Sig     int     `json:"sig"`
			Net     string  `json:"net"`
			Code    string  `json:"code"`
			Ids     string  `json:"ids"`
			Sources string  `json:"sources"`
			Types   string  `json:"types"`
			Nst     any     `json:"nst"`
			Dmin    any     `json:"dmin"`
			Rms     float64 `json:"rms"`
			Gap     any     `json:"gap"`
			MagType string  `json:"magType"`
			Type    string  `json:"type"`
			Title   string  `json:"title"`
		} `json:"properties"`
		Geometry struct {
			Type        string    `json:"type"`
			Coordinates []float64 `json:"coordinates"`
		} `json:"geometry"`
		ID string `json:"id"`
	} `json:"features"`
	Bbox []float64 `json:"bbox"`
}

type EarthquakeLocation struct {
	Type        string    `json:"type" bson:"type"`
	Coordinates []float64 `json:"coordinates" bson:"coordinates`
}

type EarthquakeDB struct {
	ID       primitive.ObjectID `json:"_id" bson:"_id"`
	Mag      float64            `json:"mag" bson:"mag"`
	Place    string             `json:"place" bson:"place"`
	Time     primitive.DateTime `json:"time" bson:"time"`
	Location EarthquakeLocation `json:"location" bson:"location"`
}

type EarthquakeCount struct {
	Count int64 `json:"count"`
}
