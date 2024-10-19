package edu_school

import (
	"encoding/json"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"io"
	"net/http"
)

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

func GetAllCampuses(token string) (*CampusesResponse, error) {
	requestUrl := "https://edu-api.21-school.ru/services/21-school/api/v1/campuses"

	client := &http.Client{}
	req, err := http.NewRequest("GET", requestUrl, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+token)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		var tribesError ErrorOfGettingTribes
		err = json.Unmarshal(body, &tribesError)
		if err != nil {
			return nil, err
		}
		return nil, status.Errorf(codes.Unknown, "Error getting list of tribes. Status code: %d: %s: %s", tribesError.Status, tribesError.Code, tribesError.Message)
	}

	var result CampusesResponse
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}
