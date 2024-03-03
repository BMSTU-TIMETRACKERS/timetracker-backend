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
