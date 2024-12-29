package model

type Tribe struct {
	Id   int32  `json:"coalitionId"`
	Name string `json:"name"`
}

type TribesResponse struct {
	Tribes []Tribe `json:"coalitions"`
}

type ErrorOfGettingTribes struct {
	Status        int    `json:"status"`
	ExceptionUuid string `json:"exceptionUUID"`
	Code          string `json:"code"`
	Message       string `json:"message"`
}
