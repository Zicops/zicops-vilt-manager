package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"

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
	lsp := claims["lsp_id"].(string)

	session, err := global.CassPool.GetSession(ctx, "viltz")
	if err != nil {
		return nil, err
	}
	CassSession := session
	session, err = global.CassPool.GetSession(ctx, "coursez")
	if err != nil {
		return nil, err
	}
	CassCoursezSession := session

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

	if input.TrainingStartTime != nil && *input.TrainingStartTime != "" {
		tst, err := strconv.Atoi(*input.TrainingStartTime)
		if err != nil {
			return nil, err
		}
		topic.TrainingStartTime = int64(tst)
	}
	if input.TrainingEndTime != nil && *input.TrainingEndTime != "" {
		tet, err := strconv.Atoi(*input.TrainingEndTime)
		if err != nil {
			return nil, err
		}
		topic.TrainingEndTime = int64(tet)
	}
	if input.Duration != nil && *input.Duration != "" {
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
	err = setFlagsInFirestore(ctx, id, input)
	if err != nil {
		return nil, err
	}
	if input.Status != nil {
		topic.Status = *input.Status
	}

	if input.Duration != nil && *input.Duration != "" && input.ModuleID != nil {
		module, err := GetModule(CassCoursezSession, *input.ModuleID, lsp)
		if err != nil {
			return nil, err
		}
		duration, _ := strconv.Atoi(*input.Duration)
		newDuration := module.Duration + duration
		queryStr := fmt.Sprintf("UPDATE coursez.module SET duration=%d WHERE id='%s' and lsp_id='%s' and is_active=true AND created_at=%d", newDuration, *input.ModuleID, lsp, module.CreatedAt)
		updateQ := CassSession.Query(queryStr, nil)
		if err := updateQ.ExecRelease(); err != nil {
			return nil, err
		}
	}

	if input.Duration != nil && *input.Duration != "" && input.CourseID != nil {
		course, err := GetCourse(CassCoursezSession, *input.CourseID, lsp)
		if err != nil {
			return nil, err
		}
		duration, _ := strconv.Atoi(*input.Duration)
		newDuration := course.Duration + duration
		queryStr := fmt.Sprintf("UPDATE coursez.course SET duration=%d WHERE id='%s' and lsp_id='%s' and is_active=true and created_at=%d", newDuration, *input.CourseID, lsp, course.CreatedAt)
		updateQ := CassSession.Query(queryStr, nil)
		if err := updateQ.ExecRelease(); err != nil {
			return nil, err
		}
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

func UpdateTopicClassroom(ctx context.Context, input *model.TopicClassroomInput) (*model.TopicClassroom, error) {
	if input.TopicID == nil {
		return nil, errors.New("please mention topic id")
	}
	claims, err := identity.GetClaimsFromContext(ctx)
	if err != nil {
		return nil, err
	}
	email := claims["email"].(string)
	lsp := claims["lsp_id"].(string)

	session, err := global.CassPool.GetSession(ctx, "viltz")
	if err != nil {
		return nil, err
	}
	CassSession := session

	session, err = global.CassPool.GetSession(ctx, "coursez")
	if err != nil {
		return nil, err
	}
	CassCoursezSession := session
	qryStr := fmt.Sprintf(`SELECT * FROM viltz.topic_classroom WHERE topic_id='%s' ALLOW FILTERING`, *input.TopicID)

	getTopicsData := func() (topicsData []viltz.TopicClassroom, err error) {
		q := CassSession.Query(qryStr, nil)
		defer q.Release()
		iter := q.Iter()
		return topicsData, iter.Select(&topicsData)
	}

	topics, err := getTopicsData()
	if err != nil {
		return nil, err
	}
	if len(topics) == 0 {
		return nil, nil
	}
	topic := topics[0]
	var updatedCols []string
	if input.Breaktime != nil {
		topic.Breaktime = *input.Breaktime
		updatedCols = append(updatedCols, "breaktime")
	}

	if input.Duration != nil {
		duration, err := strconv.Atoi(*input.Duration)
		if err != nil {
			return nil, err
		}
		topic.Duration = int64(duration)
		updatedCols = append(updatedCols, "duration")
	}

	if input.Duration != nil && *input.Duration != "" && input.ModuleID != nil {
		module, err := GetModule(CassCoursezSession, *input.ModuleID, lsp)
		if err != nil {
			return nil, err
		}
		duration, _ := strconv.Atoi(*input.Duration)
		newDuration := module.Duration - int(topic.Duration) + duration
		queryStr := fmt.Sprintf("UPDATE coursez.module SET duration=%d WHERE id='%s' and lsp_id='%s' and is_active=true AND created_at=%d", newDuration, *input.ModuleID, lsp, module.CreatedAt)
		updateQ := CassSession.Query(queryStr, nil)
		if err := updateQ.ExecRelease(); err != nil {
			return nil, err
		}
	}

	if input.Duration != nil && *input.Duration != "" && input.CourseID != nil {
		course, err := GetCourse(CassCoursezSession, *input.CourseID, lsp)
		if err != nil {
			return nil, err
		}
		duration, _ := strconv.Atoi(*input.Duration)
		newDuration := course.Duration - int(topic.Duration) + duration
		queryStr := fmt.Sprintf("UPDATE coursez.course SET duration=%d WHERE id='%s' and lsp_id='%s' and is_active=true and created_at=%d", newDuration, *input.CourseID, lsp, course.CreatedAt)
		updateQ := CassSession.Query(queryStr, nil)
		if err := updateQ.ExecRelease(); err != nil {
			return nil, err
		}
	}

	if input.IsCameraEnabled != nil {
		topic.IsCameraEnabled = *input.IsCameraEnabled
		updatedCols = append(updatedCols, "is_camera_enabled")
	}
	if input.IsChatEnabled != nil {
		topic.IsChatEnabled = *input.IsChatEnabled
		updatedCols = append(updatedCols, "is_chat_enabled")
	}
	if input.IsMicrophoneEnabled != nil {
		topic.IsMicrophoneEnabled = *input.IsMicrophoneEnabled
		updatedCols = append(updatedCols, "is_microphone_enabled")
	}
	if input.IsOverrideConfig != nil {
		topic.IsOverrideConfig = *input.IsOverrideConfig
		updatedCols = append(updatedCols, "is_override_config")
	}
	if input.IsQaEnabled != nil {
		topic.IsQAEnabled = *input.IsQaEnabled
		updatedCols = append(updatedCols, "is_qa_enabled")
	}
	if input.IsScreenShareEnabled != nil {
		topic.IsScreenShareEnabled = *input.IsScreenShareEnabled
		updatedCols = append(updatedCols, "is_screen_share_enabled")
	}
	if input.Language != nil {
		topic.Language = *input.Language
		updatedCols = append(updatedCols, "language")
	}
	err = setFlagsInFirestore(ctx, *input.ID, input)
	if err != nil {
		return nil, err
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
		updatedCols = append(updatedCols, "moderator")
	}
	if input.Status != nil {
		topic.Status = *input.Status
		updatedCols = append(updatedCols, "status")
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
		updatedCols = append(updatedCols, "trainer")
	}
	if input.TrainingEndTime != nil {
		tet, err := strconv.Atoi(*input.TrainingEndTime)
		if err != nil {
			return nil, err
		}
		topic.TrainingEndTime = int64(tet)
		updatedCols = append(updatedCols, "training_end_time")
	}
	if input.TrainingStartTime != nil {
		tst, err := strconv.Atoi(*input.TrainingStartTime)
		if err != nil {
			return nil, err
		}
		topic.TrainingStartTime = int64(tst)
		updatedCols = append(updatedCols, "training_start_time")
	}

	ua := time.Now().Unix()
	if len(updatedCols) > 0 {
		topic.UpdatedAt = ua
		topic.UpdatedBy = email
		updatedCols = append(updatedCols, "updated_at", "updated_by")
		stmt, names := viltz.TopicClassroomTable.Update(updatedCols...)
		updateQuery := CassSession.Query(stmt, names).BindStruct(&topic)
		if err = updateQuery.ExecRelease(); err != nil {
			log.Printf("Got error while updating topic classroom data: %v", err)
			return nil, err
		}
	}

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
	start := strconv.Itoa(int(topic.TrainingStartTime))
	end := strconv.Itoa(int(topic.TrainingEndTime))
	duration := strconv.Itoa(int(topic.Duration))
	createdAt := strconv.Itoa(int(topic.CreatedAt))
	updatedAt := strconv.Itoa(int(ua))
	res := model.TopicClassroom{
		ID:                   &topic.Id,
		TopicID:              &topic.TopicId,
		Trainers:             trainers,
		Moderators:           moderators,
		TrainingStartTime:    &start,
		TrainingEndTime:      &end,
		Duration:             &duration,
		Breaktime:            &topic.Breaktime,
		Language:             &topic.Language,
		IsScreenShareEnabled: &topic.IsScreenShareEnabled,
		IsChatEnabled:        &topic.IsChatEnabled,
		IsMicrophoneEnabled:  &topic.IsMicrophoneEnabled,
		IsQaEnabled:          &topic.IsQAEnabled,
		IsCameraEnabled:      &topic.IsCameraEnabled,
		IsOverrideConfig:     &topic.IsOverrideConfig,
		CreatedAt:            &createdAt,
		CreatedBy:            &topic.CreatedBy,
		UpdatedAt:            &updatedAt,
		UpdatedBy:            &email,
		Status:               &topic.Status,
	}
	return &res, nil
}

func GetTopicClassroomsByTopicIds(ctx context.Context, topicIds []*string) ([]*model.TopicClassroom, error) {
	if len(topicIds) == 0 {
		return nil, nil
	}
	_, err := identity.GetClaimsFromContext(ctx)
	if err != nil {
		return nil, err
	}

	session, err := global.CassPool.GetSession(ctx, "viltz")
	if err != nil {
		return nil, err
	}
	CassSession := session

	var wg sync.WaitGroup
	res := make([]*model.TopicClassroom, len(topicIds))
	for kk, vv := range topicIds {
		v := vv
		k := kk
		if v == nil {
			continue
		}
		wg.Add(1)
		go func(topicId string, val int) {
			defer wg.Done()
			qryStr := fmt.Sprintf(`SELECT * FROM viltz.topic_classroom WHERE topic_id='%s' ALLOW FILTERING`, topicId)
			getTopicDetails := func() (topicDetails []viltz.TopicClassroom, err error) {
				q := CassSession.Query(qryStr, nil)
				defer q.Release()
				iter := q.Iter()
				return topicDetails, iter.Select(&topicDetails)
			}
			topics, err := getTopicDetails()
			if err != nil {
				log.Printf("Got error while getting topic details: %v", err)
				return
			}
			if len(topics) == 0 {
				return
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
			tmp := model.TopicClassroom{
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

			res[val] = &tmp
		}(*v, k)
	}
	wg.Wait()
	return res, nil
}

func setFlagsInFirestore(ctx context.Context, id string, input *model.TopicClassroomInput) error {

	claims, err := identity.GetClaimsFromContext(ctx)
	if err != nil {
		return err
	}
	lsp := claims["lsp_id"].(string)

	query := fmt.Sprintf(`
	{
	   addClassroomFlags(input: {
		 is_microphone_enabled: %b
		 is_screen_sharing_enabled: %b
		 is_chat_enabled: %b
		 is_qa_enabled: %b
		 is_video_sharing_enabled: %b
		 
	   }) {
		 id
		 is_classroom_started
		 is_participants_present
		 is_ad_displayed
		 is_break
		 is_moderator_joined
		 is_trainer_joined
		 ad_video_url
		 is_microphone_enabled
		 is_video_sharing_enabled
		 is_screen_sharing_enabled
		 is_chat_enabled
		 is_qa_enabled
		 quiz
	   }
	 }
	 
   `, input.IsMicrophoneEnabled, input.IsScreenShareEnabled, input.IsChatEnabled, input.IsQaEnabled, input.IsCameraEnabled)
	jsonData := map[string]string{
		"mutation": query,
	}

	jsonValue, err := json.Marshal(jsonData)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", "https://demo.zicops.com/ns/query", bytes.NewBuffer(jsonValue))
	if err != nil {
		return err
	}
	ctxStr := fmt.Sprintf("%v", ctx)
	req.Header.Set("Authorization", ctxStr)
	req.Header.Set("fcm-token", "fcm-token")
	req.Header.Set("tenant", lsp)

	client := &http.Client{Timeout: time.Second * 30}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	log.Printf("Status code for %s is %v", *input.ID, resp.StatusCode)
	return nil
}
