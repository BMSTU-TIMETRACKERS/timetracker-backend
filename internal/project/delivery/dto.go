package delivery

type CreateProjectIn struct {
	Name string `json:"name" validate:"required"`
}

type CreateProjectOut struct {
	ID int64 `json:"id" validate:"required"`
}

type ProjectOut struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type ProjectsStatOut struct {
	TotalDurationInSec float64       `json:"total_duration_in_sec"`
	Projects           []ProjectStat `json:"projects"`
}

type ProjectStat struct {
	ID              int64   `json:"id"`
	Name            string  `json:"name"`
	DurationInSec   float64 `json:"duration_in_sec"`
	PercentDuration float64 `json:"percent_duration"`
}

type ProjectEntriesStatOut struct {
	TotalDurationInSec float64              `json:"total_duration_in_sec"`
	Entries            []ProjectEntriesStat `json:"entries"`
}

type ProjectEntriesStat struct {
	Name            string  `json:"name"`
	DurationInSec   float64 `json:"duration_in_sec"`
	PercentDuration float64 `json:"percent_duration"`
}
