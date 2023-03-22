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

// GetViltData is the resolver for the getViltData field.
func (r *queryResolver) GetViltData(ctx context.Context, courseID *string) (*model.Vilt, error) {
	resp, err := handlers.GetViltData(ctx, courseID)
	if err != nil {
		log.Printf("Got error while getting vilt data: %v", err)
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

// !!! WARNING !!!
// The code below was going to be deleted when updating resolvers. It has been copied here so you have
// one last chance to move it out of harms way if you want. There are two reasons this happens:
//   - When renaming or deleting a resolver the old code will be put in here. You can safely delete
//     it when you're done.
//   - You have helper methods in this file. Move them out to keep these resolver files clean.
