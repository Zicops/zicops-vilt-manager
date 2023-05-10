package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"log"

	"github.com/zicops/zicops-vilt-manager/graph/model"
	"github.com/zicops/zicops-vilt-manager/handlers"
)

// CreateViltData is the resolver for the createViltData field.
func (r *mutationResolver) CreateViltData(ctx context.Context, input *model.ViltInput) (*model.Vilt, error) {
	resp, err := handlers.CreateViltData(ctx, input)
	if err != nil {
		log.Printf("Got error while creating vilt data: %v", err)
		return nil, err
	}
	return resp, nil
}

// UpdateViltData is the resolver for the updateViltData field.
func (r *mutationResolver) UpdateViltData(ctx context.Context, input *model.ViltInput) (*model.Vilt, error) {
	resp, err := handlers.UpdateViltData(ctx, input)
	if err != nil {
		log.Printf("Got error while updating vilt data: %v", err)
		return nil, err
	}
	return resp, nil
}

// CreateTopicClassroom is the resolver for the createTopicClassroom field.
func (r *mutationResolver) CreateTopicClassroom(ctx context.Context, input *model.TopicClassroomInput) (*model.TopicClassroom, error) {
	resp, err := handlers.CreateTopicClassroom(ctx, input)
	if err != nil {
		log.Printf("Got error while creating topic classroom: %v", err)
		return nil, err
	}
	return resp, nil
}

// UpdateTopicClassroom is the resolver for the updateTopicClassroom field.
func (r *mutationResolver) UpdateTopicClassroom(ctx context.Context, input *model.TopicClassroomInput) (*model.TopicClassroom, error) {
	resp, err := handlers.UpdateTopicClassroom(ctx, input)
	if err != nil {
		log.Printf("Got error while updating topic classroom: %v", err)
		return nil, err
	}
	return resp, nil
}

// CreateTrainerData is the resolver for the createTrainerData field.
func (r *mutationResolver) CreateTrainerData(ctx context.Context, input *model.TrainerInput) (*model.Trainer, error) {
	resp, err := handlers.CreateTrainerData(ctx, input)
	if err != nil {
		log.Printf("Got error while creating trainer data: %v", err)
		return nil, err
	}
	return resp, nil
}

// UpdateTrainerData is the resolver for the updateTrainerData field.
func (r *mutationResolver) UpdateTrainerData(ctx context.Context, input *model.TrainerInput) (*model.Trainer, error) {
	resp, err := handlers.UpdateTrainerData(ctx, input)
	if err != nil {
		log.Printf("Got error while updating trainer data: %v", err)
		return nil, err
	}
	return resp, nil
}

// RegisterUserForCourse is the resolver for the registerUserForCourse field.
func (r *mutationResolver) RegisterUserForCourse(ctx context.Context, input *model.UserCourseRegisterInput) (*model.UserCourseRegister, error) {
	resp, err := handlers.RegisterUserForCourse(ctx, input)
	if err != nil {
		log.Printf("Got error while registering user to course: %v", err)
		return nil, err
	}
	return resp, nil
}

// UpdateRegistrationForCourse is the resolver for the updateRegistrationForCourse field.
func (r *mutationResolver) UpdateRegistrationForCourse(ctx context.Context, input *model.UserCourseRegisterInput) (*model.UserCourseRegister, error) {
	resp, err := handlers.UpdateRegistrationForCourse(ctx, input)
	if err != nil {
		log.Printf("Got error while updating registeration for course: %v", err)
		return nil, err
	}
	return resp, nil
}

// GetViltData is the resolver for the getViltData field.
func (r *queryResolver) GetViltData(ctx context.Context, courseID *string) ([]*model.Vilt, error) {
	resp, err := handlers.GetViltData(ctx, courseID)
	if err != nil {
		log.Printf("Got error while getting vilt data: %v", err)
		return nil, err
	}
	return resp, nil
}

// GetViltDataByID is the resolver for the getViltDataById field.
func (r *queryResolver) GetViltDataByID(ctx context.Context, id *string) (*model.Vilt, error) {
	resp, err := handlers.GetViltDataByID(ctx, id)
	if err != nil {
		log.Printf("Got error while getting vilt data: %v", err)
		return nil, err
	}
	return resp, nil
}

// GetTopicClassroom is the resolver for the getTopicClassroom field.
func (r *queryResolver) GetTopicClassroom(ctx context.Context, topicID *string) (*model.TopicClassroom, error) {
	resp, err := handlers.GetTopicClassroom(ctx, topicID)
	if err != nil {
		log.Printf("Got error while getting classroom data: %v", err)
		return nil, err
	}
	return resp, nil
}

// GetTopicClassroomsByTopicIds is the resolver for the getTopicClassroomsByTopicIds field.
func (r *queryResolver) GetTopicClassroomsByTopicIds(ctx context.Context, topicIds []*string) ([]*model.TopicClassroom, error) {
	resp, err := handlers.GetTopicClassroomsByTopicIds(ctx, topicIds)
	if err != nil {
		log.Printf("Got error while getting topics: %v", err)
		return nil, err
	}
	return resp, nil
}

// GetTrainerData is the resolver for the getTrainerData field.
func (r *queryResolver) GetTrainerData(ctx context.Context, lspID *string, vendorID *string, pageCursor *string, direction *string, pageSize *int, filters *model.TrainerFilters) (*model.PaginatedTrainer, error) {
	resp, err := handlers.GetTrainerData(ctx, lspID, vendorID, pageCursor, direction, pageSize, filters)
	if err != nil {
		log.Printf("Got error while getting trainer: %v", err)
		return nil, err
	}
	return resp, nil
}

// GetTrainerByID is the resolver for the getTrainerById field.
func (r *queryResolver) GetTrainerByID(ctx context.Context, id *string) (*model.Trainer, error) {
	resp, err := handlers.GetTrainerByID(ctx, id)
	if err != nil {
		log.Printf("Got error while getting trainer: %v", err)
		return nil, err
	}
	return resp, nil
}

// GetAllRegistrations is the resolver for the getAllRegistrations field.
func (r *queryResolver) GetAllRegistrations(ctx context.Context, courseID *string, pageCursor *string, direction *string, pageSize *int) (*model.PaginatedUserCourseRegister, error) {
	resp, err := handlers.GetAllRegistrations(ctx, courseID, pageCursor, direction, pageSize)
	if err != nil {
		log.Printf("Got error while getting all registrations: %v", err)
		return nil, err
	}
	return resp, nil
}

// GetRegistrationDetails is the resolver for the getRegistrationDetails field.
func (r *queryResolver) GetRegistrationDetails(ctx context.Context, id *string) (*model.UserCourseRegister, error) {
	resp, err := handlers.GetRegistrationDetails(ctx, id)
	if err != nil {
		log.Printf("Got error while getting registration details: %v", err)
		return nil, err
	}
	return resp, nil
}

// GetTrainerCourses is the resolver for the getTrainerCourses field.
func (r *queryResolver) GetTrainerCourses(ctx context.Context, userID *string) ([]*model.Course, error) {
	resp, err := handlers.GetTrainerCourses(ctx, userID)
	if err != nil {
		log.Printf("Got error while getting trainer courses: %v", err)
		return nil, err
	}
	return resp, nil
}

// GetTopicAttendance is the resolver for the getTopicAttendance field.
func (r *queryResolver) GetTopicAttendance(ctx context.Context, topicID string) ([]*model.TopicAttendance, error) {
	resp, err := handlers.GetTopicAttendance(ctx, topicID)
	if err != nil {
		log.Printf("Got error while getting attendance: %v", err)
		return nil, err
	}
	return resp, nil
}

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
