// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

type PaginatedTrainer struct {
	Trainers   []*Trainer `json:"trainers"`
	PageCursor *string    `json:"pageCursor"`
	Direction  *string    `json:"Direction"`
	PageSize   *int       `json:"pageSize"`
}

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
	ModuleID             *string   `json:"module_id"`
	CourseID             *string   `json:"course_id"`
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

type TrainerFilters struct {
	Name *string `json:"name"`
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
	ID                    *string   `json:"id"`
	LspID                 *string   `json:"lsp_id"`
	CourseID              *string   `json:"course_id"`
	NoOfLearners          *int      `json:"no_of_learners"`
	Trainers              []*string `json:"trainers"`
	Moderators            []*string `json:"moderators"`
	CourseStartDate       *string   `json:"course_start_date"`
	CourseEndDate         *string   `json:"course_end_date"`
	IsTrainerDecided      *bool     `json:"is_trainer_decided"`
	IsModeratorDecided    *bool     `json:"is_moderator_decided"`
	IsStartDateDecided    *bool     `json:"is_start_date_decided"`
	IsEndDateDecided      *bool     `json:"is_end_date_decided"`
	Curriculum            *string   `json:"curriculum"`
	PricingType           *string   `json:"pricing_type"`
	PricePerSeat          *int      `json:"price_per_seat"`
	Currency              *string   `json:"currency"`
	TaxPercentage         *float64  `json:"tax_percentage"`
	IsRegistrationOpen    *bool     `json:"is_registration_open"`
	IsBookingOpen         *bool     `json:"is_booking_open"`
	MaxRegistrations      *int      `json:"max_registrations"`
	RegistrationEndDate   *int      `json:"registration_end_date"`
	BookingStartDate      *int      `json:"booking_start_date"`
	BookingEndDate        *int      `json:"booking_end_date"`
	RegistrationPublishBy *string   `json:"registration_publish_by"`
	RegistrationPublishOn *int      `json:"registration_publish_on"`
	BookingPublishOn      *int      `json:"booking_publish_on"`
	BookingPublishBy      *string   `json:"booking_publish_by"`
	RegistrationStartDate *int      `json:"registration_start_date"`
	CreatedAt             *string   `json:"created_at"`
	CreatedBy             *string   `json:"created_by"`
	UpdatedAt             *string   `json:"updated_at"`
	UpdatedBy             *string   `json:"updated_by"`
	Status                *string   `json:"status"`
}

type ViltInput struct {
	ID                    *string   `json:"id"`
	LspID                 *string   `json:"lsp_id"`
	CourseID              *string   `json:"course_id"`
	NoOfLearners          *int      `json:"no_of_learners"`
	Trainers              []*string `json:"trainers"`
	Moderators            []*string `json:"moderators"`
	CourseStartDate       *string   `json:"course_start_date"`
	CourseEndDate         *string   `json:"course_end_date"`
	Curriculum            *string   `json:"curriculum"`
	IsTrainerDecided      *bool     `json:"is_trainer_decided"`
	IsModeratorDecided    *bool     `json:"is_moderator_decided"`
	IsStartDateDecided    *bool     `json:"is_start_date_decided"`
	IsEndDateDecided      *bool     `json:"is_end_date_decided"`
	PricingType           *string   `json:"pricing_type"`
	PricePerSeat          *int      `json:"price_per_seat"`
	Currency              *string   `json:"currency"`
	TaxPercentage         *float64  `json:"tax_percentage"`
	IsRegistrationOpen    *bool     `json:"is_registration_open"`
	IsBookingOpen         *bool     `json:"is_booking_open"`
	MaxRegistrations      *int      `json:"max_registrations"`
	RegistrationEndDate   *int      `json:"registration_end_date"`
	BookingStartDate      *int      `json:"booking_start_date"`
	BookingEndDate        *int      `json:"booking_end_date"`
	RegistrationPublishBy *string   `json:"registration_publish_by"`
	RegistrationPublishOn *int      `json:"registration_publish_on"`
	BookingPublishOn      *int      `json:"booking_publish_on"`
	BookingPublishBy      *string   `json:"booking_publish_by"`
	RegistrationStartDate *int      `json:"registration_start_date"`
	Status                *string   `json:"status"`
}
