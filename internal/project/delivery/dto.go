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
	TotalDurationInHours float64       `json:"total_duration_in_hours"`
	Projects             []ProjectStat `json:"projects"`
}

type ProjectStat struct {
	ID              int64   `json:"id"`
	Name            string  `json:"name"`
	DurationInHours float64 `json:"duration_in_hours"`
	PercentDuration float64 `json:"percent_duration"`
}
