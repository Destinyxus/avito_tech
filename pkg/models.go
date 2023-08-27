package pkg

type Segment struct {
	Slug string `json:"slug"`
}

type User struct {
	Name string `json:"name"`
}

type SegmentToAdd struct {
	Id   int      `json:"id"`
	Slug []string `json:"slug"`
}

type RequestActive struct {
	Id int `json:"id"`
}

type CSV struct {
	UserID int `json:"user_id"`
	Year   int `json:"year"`
	Month  int `json:"month"`
}
