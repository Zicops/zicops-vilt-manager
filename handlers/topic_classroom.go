package handlers

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/google/uuid"

	"github.com/zicops/contracts/viltz"
	"github.com/zicops/zicops-vilt-manager/global"
	"github.com/zicops/zicops-vilt-manager/graph/model"
	"github.com/zicops/zicops-vilt-manager/lib/identity"
)

func CreateTopicClassroom(ctx context.Context, input *model.TopicClassroomInput) (*model.TopicClassroom, error) {
	claims, err := identity.GetClaimsFromContext(ctx)
	if err != nil {
		return nil, err
	}
	email := claims["email"].(string)

	session, err := global.CassPool.GetSession(ctx, "viltz")
	if err != nil {
		return nil, err
	}
	CassSession := session

	id := uuid.New().String()
	createdAt := time.Now().Unix()
	topic := viltz.TopicClassroom{
		Id:        id,
		TopicId:   *input.TopicID,
		CreatedAt: createdAt,
		CreatedBy: email,
		UpdatedAt: createdAt,
		UpdatedBy: email,
	}
	if input.Trainers != nil {
		var tmp []string
		for _, vv := range input.Trainers {
			v := vv
			if v == nil {
				continue
			}
			tmp = append(tmp, *v)
		}
		topic.Trainer = tmp
	}

	if input.Moderators != nil {
		var tmp []string
		for _, vv := range input.Moderators {
			v := vv
			if v == nil {
				continue
			}
			tmp = append(tmp, *v)
		}
		topic.Moderator = tmp
	}

	if input.TrainingStartTime != nil {
		tst, err := strconv.Atoi(*input.TrainingStartTime)
		if err != nil {
			return nil, err
		}
		topic.TrainingStartTime = int64(tst)
	}
	if input.TrainingEndTime != nil {
		tet, err := strconv.Atoi(*input.TrainingEndTime)
		if err != nil {
			return nil, err
		}
		topic.TrainingEndTime = int64(tet)
	}
	if input.Duration != nil {
		d, err := strconv.Atoi(*input.Duration)
		if err != nil {
			return nil, err
		}
		topic.Duration = int64(d)
	}
	if input.Breaktime != nil {
		topic.Breaktime = *input.Breaktime
	}
	if input.Language != nil {
		topic.Language = *input.Language
	}
	if input.IsCameraEnabled != nil {
		topic.IsCameraEnabled = *input.IsCameraEnabled
	}
	if input.IsChatEnabled != nil {
		topic.IsChatEnabled = *input.IsChatEnabled
	}
	if input.IsMicrophoneEnabled != nil {
		topic.IsMicrophoneEnabled = *input.IsMicrophoneEnabled
	}
	if input.IsOverrideConfig != nil {
		topic.IsOverrideConfig = *input.IsOverrideConfig
	}
	if input.IsQaEnabled != nil {
		topic.IsQAEnabled = *input.IsQaEnabled
	}
	if input.IsScreenShareEnabled != nil {
		topic.IsScreenShareEnabled = *input.IsScreenShareEnabled
	}
	if input.Status != nil {
		topic.Status = *input.Status
	}

	insertQuery := CassSession.Query(viltz.TopicClassroomTable.Insert()).BindStruct(topic)
	if err = insertQuery.Exec(); err != nil {
		return nil, err
	}

	ca := strconv.Itoa(int(createdAt))
	res := model.TopicClassroom{
		ID:                   &id,
		TopicID:              input.TopicID,
		Trainers:             input.Trainers,
		Moderators:           input.Moderators,
		TrainingStartTime:    input.TrainingStartTime,
		TrainingEndTime:      input.TrainingEndTime,
		Duration:             input.Duration,
		Breaktime:            input.Breaktime,
		Language:             input.Language,
		IsScreenShareEnabled: input.IsScreenShareEnabled,
		IsChatEnabled:        input.IsChatEnabled,
		IsMicrophoneEnabled:  input.IsMicrophoneEnabled,
		IsQaEnabled:          input.IsQaEnabled,
		IsCameraEnabled:      input.IsCameraEnabled,
		IsOverrideConfig:     input.IsOverrideConfig,
		CreatedAt:            &ca,
		CreatedBy:            &email,
		UpdatedAt:            &ca,
		UpdatedBy:            &email,
		Status:               input.Status,
	}
	return &res, nil
}

func GetTopicClassroom(ctx context.Context, topicID *string) (*model.TopicClassroom, error) {
	_, err := identity.GetClaimsFromContext(ctx)
	if err != nil {
		return nil, err
	}

	session, err := global.CassPool.GetSession(ctx, "viltz")
	if err != nil {
		return nil, err
	}
	CassSession := session
	qryStr := fmt.Sprintf(`SELECT * FROM viltz.topic_classroom WHERE topic_id='%s' ALLOW FILTERING`, *topicID)
	getTopicDetails := func() (topics []viltz.TopicClassroom, err error) {
		q := CassSession.Query(qryStr, nil)
		defer q.Release()
		iter := q.Iter()
		return topics, iter.Select(&topics)
	}
	topics, err := getTopicDetails()
	if err != nil {
		return nil, err
	}
	if len(topics) == 0 {
		return nil, nil
	}

	topic := topics[0]
	var trainers []*string
	for _, vv := range topic.Trainer {
		v := vv
		trainers = append(trainers, &v)
	}
	var moderators []*string
	for _, vv := range topic.Moderator {
		v := vv
		moderators = append(moderators, &v)
	}
	tst := strconv.Itoa(int(topic.TrainingStartTime))
	tet := strconv.Itoa(int(topic.TrainingEndTime))
	duration := strconv.Itoa(int(topic.Duration))
	ca := strconv.Itoa(int(topic.CreatedAt))
	ua := strconv.Itoa(int(topic.UpdatedAt))
	res := model.TopicClassroom{
		ID:                   &topic.Id,
		TopicID:              &topic.TopicId,
		Trainers:             trainers,
		Moderators:           moderators,
		TrainingStartTime:    &tst,
		TrainingEndTime:      &tet,
		Duration:             &duration,
		Breaktime:            &topic.Breaktime,
		Language:             &topic.Language,
		IsScreenShareEnabled: &topic.IsScreenShareEnabled,
		IsChatEnabled:        &topic.IsChatEnabled,
		IsMicrophoneEnabled:  &topic.IsMicrophoneEnabled,
		IsQaEnabled:          &topic.IsQAEnabled,
		IsCameraEnabled:      &topic.IsCameraEnabled,
		IsOverrideConfig:     &topic.IsOverrideConfig,
		CreatedAt:            &ca,
		CreatedBy:            &topic.CreatedBy,
		UpdatedAt:            &ua,
		UpdatedBy:            &topic.UpdatedBy,
		Status:               &topic.Status,
	}
	return &res, nil
}
