package delivery

type CreateProjectIn struct {
	Name string `json:"name" validate:"required" example:"Работа"` // Название проекта.
}

type CreateProjectOut struct {
	ID int64 `json:"id" validate:"required" example:"1"` // Идентификатор проекта.
}

type ProjectOut struct {
	ID   int64  `json:"id" example:"1"`        // Идентификатор проекта.                         // Идентификатор проекта.
	Name string `json:"name" example:"Работа"` // Название проекта.
}

type ProjectsStatOut struct {
	TotalDurationInSec float64       `json:"total_duration_in_sec" example:"3600"` // Суммарное время (в сек.) потраченное на все проекты.
	Projects           []ProjectStat `json:"projects"`                             // Список проектов.
}

type ProjectStat struct {
	ID              int64   `json:"id" example:"1"`                // Идентификатор проекта.
	Name            string  `json:"name" example:"Работа"`         // Название проекта.
	DurationInSec   float64 `json:"duration_in_sec" example:"360"` // Суммарное время (в сек.) потраченное на проект.
	PercentDuration float64 `json:"percent_duration" example:"10"` // Доля (в процентах) длительности проекта от суммарной длительности.
}

type ProjectEntriesStatOut struct {
	TotalDurationInSec float64              `json:"total_duration_in_sec" example:"60"` // Суммарное время (в сек.) потраченное на проект.
	Entries            []ProjectEntriesStat `json:"entries"`                            // Записи времени.
}

type ProjectEntriesStat struct {
	Name            string  `json:"name" example:"task1"`          // Название записи.
	DurationInSec   float64 `json:"duration_in_sec" example:"360"` // Суммарное время (в сек.) потраченное на запись.
	PercentDuration float64 `json:"percent_duration"`              // Доля (в процентах) длительности записи от длительности проекта.
}
