// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

type NewTodo struct {
	Text   string `json:"text"`
	UserID string `json:"userId"`
}

type Todo struct {
	ID   string `json:"id"`
	Text string `json:"text"`
	Done bool   `json:"done"`
	User *User  `json:"user"`
}

type User struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Vilt struct {
	LspID           *string   `json:"lsp_id"`
	CourseID        *string   `json:"course_id"`
	NoOfLearners    *int      `json:"no_of_learners"`
	Trainers        []*string `json:"trainers"`
	Moderators      []*string `json:"moderators"`
	CourseStartDate *string   `json:"course_start_date"`
	CourseEndDate   *string   `json:"course_end_date"`
	Curriculum      *string   `json:"curriculum"`
	CreatedAt       *string   `json:"created_at"`
	CreatedBy       *string   `json:"created_by"`
	UpdatedAt       *string   `json:"updated_at"`
	UpdatedBy       *string   `json:"updated_by"`
	Status          *string   `json:"status"`
}

type ViltInput struct {
	LspID           *string   `json:"lsp_id"`
	CourseID        *string   `json:"course_id"`
	NoOfLearners    *int      `json:"no_of_learners"`
	Trainers        []*string `json:"trainers"`
	Moderators      []*string `json:"moderators"`
	CourseStartDate *string   `json:"course_start_date"`
	CourseEndDate   *string   `json:"course_end_date"`
	Curriculum      *string   `json:"curriculum"`
	Status          *string   `json:"status"`
}
