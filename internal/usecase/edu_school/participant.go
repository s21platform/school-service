package edu_school

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/s21platform/school-service/internal/model"
	"golang.org/x/sync/errgroup"
	_ "golang.org/x/sync/errgroup"
	"io"
	"log"
	"net/http"
)

func fetchFromAPI(ctx context.Context, url, token string, result interface{}, errResult interface{}) error {
	client := &http.Client{}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+token)

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		if err := json.Unmarshal(body, errResult); err == nil {
			return fmt.Errorf("API error: %s", string(body))
		}
		return fmt.Errorf("unknown API error: %s", string(body))
	}

	return json.Unmarshal(body, result)
}

func GetParticipantData(ctx context.Context, login, token string) (*model.ParticipantDataResponse, error) {
	var (
		participant model.Participant
		skillsResp  model.SkillsParticipantResponse
		points      model.PointsParticipant
		badgesResp  model.BadgesParticipantResponse

		errParticipant model.ErrorOfGettingParticipant
		errSkills      model.ErrorOfGettingSkills
		errPoints      model.ErrorOfGettingPoints
		errBadges      model.ErrorOfGettingBadges
	)

	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		url := "https://edu-api.21-school.ru/services/21-school/api/v1/participants/" + login
		return fetchFromAPI(ctx, url, token, &participant, &errParticipant)
	})

	g.Go(func() error {
		url := "https://edu-api.21-school.ru/services/21-school/api/v1/participants/" + login + "/skills"
		return fetchFromAPI(ctx, url, token, &skillsResp, &errSkills)
	})

	g.Go(func() error {
		url := "https://edu-api.21-school.ru/services/21-school/api/v1/participants/" + login + "/points"
		return fetchFromAPI(ctx, url, token, &points, &errPoints)
	})

	g.Go(func() error {
		url := "https://edu-api.21-school.ru/services/21-school/api/v1/participants/" + login + "/badges"
		return fetchFromAPI(ctx, url, token, &badgesResp, &errBadges)
	})

	if err := g.Wait(); err != nil {
		log.Printf("Ошибка получения данных: %v", err)
		return nil, err
	}
	result := &model.ParticipantDataResponse{
		ClassName:            participant.ClassName,
		ParallelName:         participant.ParallelName,
		ExpValue:             int64(participant.ExpValue),
		Level:                participant.Level,
		ExpToNextLevel:       participant.ExpToNextLevel,
		CampusUuid:           participant.Campus.Uuid,
		Status:               participant.Status,
		Skills:               skillsResp,
		PeerReviewPoints:     int64(points.PeerReviewPoints),
		PeerCodeReviewPoints: int64(points.CodeReviewPoints),
		Coins:                int64(points.Coins),
		Badges:               badgesResp,
	}

	return result, nil
}
