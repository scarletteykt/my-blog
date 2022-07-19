package tag

type Tag struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Slug string `json:"slug"`
}

type CreateTag struct {
	Name string
	Slug string
}

type UpdateTag struct {
	ID   int
	Name string
	Slug string
}

type DeleteTag struct {
	ID int
}
