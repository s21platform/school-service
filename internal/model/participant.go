package model

// v1/participants/{login}

type CampusParticipant struct {
	Uuid      string `json:"id"`
	ShortName string `json:"shortName"`
}

type Participant struct {
	Login          string   `json:"login"`
	ClassName      string   `json:"className"`
	ParallelName   string   `json:"parallelName"`
	ExpValue       int      `json:"expValue"`
	Level          int32    `json:"level"`
	ExpToNextLevel int64    `json:"expToNextLevel"`
	Campuses       []Campus `json:"campuses"`
	Status         string   `json:"status"`
}

type ErrorOfGettingParticipant struct {
	Status        int32  `json:"status"`
	ExceptionUuid string `json:"exception_uuid"`
	Code          string `json:"code"`
	Message       string `json:"message"`
}

///v1/participants/{login}/skills

type SkillsParticipant struct {
	Name   string `json:"name"`
	Points int32  `json:"points"`
}

type ErrorOfGettingSkills struct {
	Status        int32  `json:"status"`
	ExceptionUuid string `json:"exception_uuid"`
	Code          string `json:"code"`
	Message       string `json:"message"`
}

//v1/participants/{login}/points

type PointsParticipant struct {
	PeerReviewPoints int32 `json:"peerReviewPoints"`
	CodeReviewPoints int32 `json:"codeReviewPoints"`
	Coins            int32 `json:"coins"`
}
type ErrorOfGettingPoints struct {
	Status        int32  `json:"status"`
	ExceptionUuid string `json:"exception_uuid"`
	Code          string `json:"code"`
	Message       string `json:"message"`
}

//v1/participants/{login}/badges

type BadgesParticipant struct {
	Name            string `json:"name"`
	ReceiptDataTime string `json:"receiptDataTime"`
	IconUrl         string `json:"iconUrl"`
}

type ErrorOfGettingBadges struct {
	Status        int32  `json:"status"`
	ExceptionUuid string `json:"exception_uuid"`
	Code          string `json:"code"`
	Message       string `json:"message"`
}

type ParticipantDataResponse struct {
	ClassName            string   `json:"className"`
	ParallelName         string   `json:"parallelName"`
	ExpValue             int64    `json:"expValue"`
	Level                int32    `json:"level"`
	ExpToNextLevel       int64    `json:"expToNextLevel"`
	CampusUuid           string   `json:" campusUuid"`
	Status               string   `json:"status"`
	Skills               []string `json:"skills"`
	PeerReviewPoints     int64    `json:"peerReviewPoints"`
	PeerCodeReviewPoints int64    `json:"peerCodeReviewPoints"`
	Coins                int64    `json:"coins"`
	Badges               []string `json:"badges"`
}

type Skills struct {
	Name   string `json:"name"`
	Points int32  `json:"points"`
}

type Badges struct {
	Name            string `json:"name"`
	ReceiptDateTime string `json:"receiptDateTime"`
	IconURL         string `json:"iconUrl"`
}
