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

// GetViltData is the resolver for the getViltData field.
func (r *queryResolver) GetViltData(ctx context.Context, courseID *string) (*model.Vilt, error) {
	resp, err := handlers.GetViltData(ctx, courseID)
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

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
