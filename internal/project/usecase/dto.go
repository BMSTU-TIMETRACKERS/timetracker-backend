package usecase

type Project struct {
	ID     int64
	Name   string
	UserID int64
}

type ProjectStatInfo struct {
	ProjectID              int64
	ProjectName            string
	ProjectDurationInHours float64
	ProjectDurationPercent float64
}

type AllProjectsStat struct {
	TotalDurationInHours float64
	ProjectsStat         []ProjectStatInfo
}

type ProjectEntrieInfo struct {
	EntryName            string
	EntryDurationInHours float64
	EntryDurationPercent float64
}

type AllProjectEntriesStat struct {
	TotalDurationInHours float64
	EntriesStat          []ProjectEntrieInfo
}
