package model

type ErrorOfGettingCampuses struct {
	Status        int    `json:"status"`
	ExceptionUuid string `json:"exceptionUUID"`
	Code          string `json:"code"`
	Message       string `json:"message"`
}

type Campus struct {
	Uuid      string `json:"id"`
	ShortName string `json:"shortName"`
	FullName  string `json:"fullName"`
}

type CampusesResponse struct {
	Campuses []Campus `json:"campuses"`
}
