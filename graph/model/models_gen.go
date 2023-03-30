// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

type TopicClassroom struct {
	ID                   *string   `json:"id"`
	TopicID              *string   `json:"topic_id"`
	Trainers             []*string `json:"trainers"`
	Moderators           []*string `json:"moderators"`
	TrainingStartTime    *string   `json:"training_start_time"`
	TrainingEndTime      *string   `json:"training_end_time"`
	Duration             *string   `json:"duration"`
	Breaktime            *string   `json:"breaktime"`
	Language             *string   `json:"language"`
	IsScreenShareEnabled *bool     `json:"is_screen_share_enabled"`
	IsChatEnabled        *bool     `json:"is_chat_enabled"`
	IsMicrophoneEnabled  *bool     `json:"is_microphone_enabled"`
	IsQaEnabled          *bool     `json:"is_qa_enabled"`
	IsCameraEnabled      *bool     `json:"is_camera_enabled"`
	IsOverrideConfig     *bool     `json:"is_override_config"`
	CreatedAt            *string   `json:"created_at"`
	CreatedBy            *string   `json:"created_by"`
	UpdatedAt            *string   `json:"updated_at"`
	UpdatedBy            *string   `json:"updated_by"`
	Status               *string   `json:"status"`
}

type TopicClassroomInput struct {
	ID                   *string   `json:"id"`
	TopicID              *string   `json:"topic_id"`
	Trainers             []*string `json:"trainers"`
	Moderators           []*string `json:"moderators"`
	TrainingStartTime    *string   `json:"training_start_time"`
	TrainingEndTime      *string   `json:"training_end_time"`
	Duration             *string   `json:"duration"`
	Breaktime            *string   `json:"breaktime"`
	Language             *string   `json:"language"`
	IsScreenShareEnabled *bool     `json:"is_screen_share_enabled"`
	IsChatEnabled        *bool     `json:"is_chat_enabled"`
	IsMicrophoneEnabled  *bool     `json:"is_microphone_enabled"`
	IsQaEnabled          *bool     `json:"is_qa_enabled"`
	IsCameraEnabled      *bool     `json:"is_camera_enabled"`
	IsOverrideConfig     *bool     `json:"is_override_config"`
	Status               *string   `json:"status"`
}

type Trainer struct {
	ID        *string   `json:"id"`
	LspID     *string   `json:"lsp_id"`
	UserID    *string   `json:"user_id"`
	VendorID  *string   `json:"vendor_id"`
	Expertise []*string `json:"expertise"`
	Status    *string   `json:"status"`
	CreatedAt *string   `json:"created_at"`
	CreatedBy *string   `json:"created_by"`
	UpdatedAt *string   `json:"updated_at"`
	UpdatedBy *string   `json:"updated_by"`
}

type TrainerInput struct {
	ID        *string   `json:"id"`
	LspID     *string   `json:"lsp_id"`
	UserID    *string   `json:"user_id"`
	VendorID  *string   `json:"vendor_id"`
	Expertise []*string `json:"expertise"`
	Status    *string   `json:"status"`
}

type Vilt struct {
	ID                 *string   `json:"id"`
	LspID              *string   `json:"lsp_id"`
	CourseID           *string   `json:"course_id"`
	NoOfLearners       *int      `json:"no_of_learners"`
	Trainers           []*string `json:"trainers"`
	Moderators         []*string `json:"moderators"`
	CourseStartDate    *string   `json:"course_start_date"`
	CourseEndDate      *string   `json:"course_end_date"`
	IsTrainerDecided   *bool     `json:"is_trainer_decided"`
	IsModeratorDecided *bool     `json:"is_moderator_decided"`
	IsStartDateDecided *bool     `json:"is_start_date_decided"`
	IsEndDateDecided   *bool     `json:"is_end_date_decided"`
	Curriculum         *string   `json:"curriculum"`
	CreatedAt          *string   `json:"created_at"`
	CreatedBy          *string   `json:"created_by"`
	UpdatedAt          *string   `json:"updated_at"`
	UpdatedBy          *string   `json:"updated_by"`
	Status             *string   `json:"status"`
}

type ViltInput struct {
	ID                 *string   `json:"id"`
	LspID              *string   `json:"lsp_id"`
	CourseID           *string   `json:"course_id"`
	NoOfLearners       *int      `json:"no_of_learners"`
	Trainers           []*string `json:"trainers"`
	Moderators         []*string `json:"moderators"`
	CourseStartDate    *string   `json:"course_start_date"`
	CourseEndDate      *string   `json:"course_end_date"`
	Curriculum         *string   `json:"curriculum"`
	IsTrainerDecided   *bool     `json:"is_trainer_decided"`
	IsModeratorDecided *bool     `json:"is_moderator_decided"`
	IsStartDateDecided *bool     `json:"is_start_date_decided"`
	IsEndDateDecided   *bool     `json:"is_end_date_decided"`
	Status             *string   `json:"status"`
}
