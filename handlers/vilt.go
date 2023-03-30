package handlers

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/zicops/contracts/viltz"
	"github.com/zicops/zicops-vilt-manager/global"
	"github.com/zicops/zicops-vilt-manager/graph/model"
	"github.com/zicops/zicops-vilt-manager/lib/identity"
)

func CreateViltData(ctx context.Context, input *model.ViltInput) (*model.Vilt, error) {
	claims, err := identity.GetClaimsFromContext(ctx)
	if err != nil {
		return nil, err
	}

	lsp := claims["lsp_id"].(string)
	email := claims["email"].(string)

	session, err := global.CassPool.GetSession(ctx, "viltz")
	if err != nil {
		return nil, err
	}
	CassViltSession := session

	ca := time.Now().Unix()
	id := uuid.New().String()
	vilt := viltz.ViltMaster{
		Id:        id,
		LspId:     lsp,
		CourseId:  *input.CourseID,
		CreatedAt: ca,
		CreatedBy: email,
		UpdatedAt: ca,
		UpdatedBy: email,
	}

	if input.NoOfLearners != nil {
		vilt.NoOfLearners = int64(*input.NoOfLearners)
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
		vilt.Trainers = tmp
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
		vilt.Moderators = tmp
	}
	if input.CourseStartDate != nil && *input.CourseStartDate != "" {
		sd, err := strconv.Atoi(*input.CourseStartDate)
		if err != nil {
			return nil, err
		}
		sdInt := int64(sd)
		vilt.CourseStartDate = sdInt
	}
	if input.CourseEndDate != nil && *input.CourseEndDate != "" {
		sd, err := strconv.Atoi(*input.CourseEndDate)
		if err != nil {
			return nil, err
		}
		sdInt := int64(sd)
		vilt.CourseEndDate = sdInt
	}
	if input.Curriculum != nil {
		vilt.Curriculum = *input.Curriculum
	}
	if input.Status != nil {
		vilt.Status = *input.Status
	}
	if input.IsEndDateDecided != nil {
		vilt.IsEndDateDecided = *input.IsEndDateDecided
	}
	if input.IsModeratorDecided != nil {
		vilt.IsModeratorDecided = *input.IsModeratorDecided
	}
	if input.IsStartDateDecided != nil {
		vilt.IsStartDateDecided = *input.IsStartDateDecided
	}
	if input.IsTrainerDecided != nil {
		vilt.IsTrainerDecided = *input.IsTrainerDecided
	}

	insertQuery := CassViltSession.Query(viltz.ViltMasterTable.Insert()).BindStruct(vilt)
	if err = insertQuery.Exec(); err != nil {
		return nil, err
	}

	createdAt := strconv.Itoa(int(ca))
	res := model.Vilt{
		ID:                 &id,
		LspID:              &lsp,
		CourseID:           input.CourseID,
		NoOfLearners:       input.NoOfLearners,
		Trainers:           input.Trainers,
		Moderators:         input.Moderators,
		CourseStartDate:    input.CourseStartDate,
		CourseEndDate:      input.CourseEndDate,
		Curriculum:         input.Curriculum,
		IsTrainerDecided:   input.IsTrainerDecided,
		IsModeratorDecided: input.IsModeratorDecided,
		IsStartDateDecided: input.IsStartDateDecided,
		IsEndDateDecided:   input.IsEndDateDecided,
		CreatedAt:          &createdAt,
		CreatedBy:          &email,
		UpdatedAt:          &createdAt,
		UpdatedBy:          &email,
		Status:             input.Status,
	}

	return &res, nil
}

func UpdateViltData(ctx context.Context, input *model.ViltInput) (*model.Vilt, error) {
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
	qryStr := fmt.Sprintf(`SELECT * FROM viltz.vilt_master WHERE id='%s' ALLOW FILTERING`, *input.ID)
	getViltDetails := func() (viltMap []viltz.ViltMaster, err error) {
		q := CassSession.Query(qryStr, nil)
		defer q.Release()
		iter := q.Iter()
		return viltMap, iter.Select(&viltMap)
	}
	vilts, err := getViltDetails()
	if err != nil {
		return nil, err
	}
	if len(vilts) == 0 {
		return nil, nil
	}

	vilt := vilts[0]
	var updatedCols []string
	if input.NoOfLearners != nil {
		vilt.NoOfLearners = int64(*input.NoOfLearners)
		updatedCols = append(updatedCols, "no_of_learners")
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
		vilt.Trainers = tmp
		updatedCols = append(updatedCols, "trainers")
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
		vilt.Moderators = tmp
		updatedCols = append(updatedCols, "moderators")
	}
	if input.CourseStartDate != nil {
		cs, err := strconv.Atoi(*input.CourseStartDate)
		if err != nil {
			return nil, err
		}
		csInt := int64(cs)
		vilt.CourseStartDate = csInt
		updatedCols = append(updatedCols, "course_start_date")
	}
	if input.CourseEndDate != nil {
		ce, err := strconv.Atoi(*input.CourseEndDate)
		if err != nil {
			return nil, err
		}
		ceInt := int64(ce)
		vilt.CourseEndDate = ceInt
		updatedCols = append(updatedCols, "course_end_date")
	}
	if input.IsEndDateDecided != nil {
		vilt.IsEndDateDecided = *input.IsEndDateDecided
		updatedCols = append(updatedCols, "is_end_date_decided")
	}
	if input.IsModeratorDecided != nil {
		vilt.IsModeratorDecided = *input.IsModeratorDecided
		updatedCols = append(updatedCols, "is_moderator_decided")
	}
	if input.IsStartDateDecided != nil {
		vilt.IsStartDateDecided = *input.IsStartDateDecided
		updatedCols = append(updatedCols, "is_start_date_decided")
	}
	if input.IsTrainerDecided != nil {
		vilt.IsTrainerDecided = *input.IsTrainerDecided
		updatedCols = append(updatedCols, "is_trainer_decided")
	}
	if input.Curriculum != nil {
		vilt.Curriculum = *input.Curriculum
		updatedCols = append(updatedCols, "curriculum")
	}
	if input.Status != nil {
		vilt.Status = *input.Status
		updatedCols = append(updatedCols, "status")
	}
	ua := time.Now().Unix()
	if len(updatedCols) > 0 {
		updatedCols = append(updatedCols, "updated_at", "updated_by")
		vilt.UpdatedAt = ua
		vilt.UpdatedBy = email

		stmt, names := viltz.ViltMasterTable.Update(updatedCols...)
		updatedQuery := CassSession.Query(stmt, names).BindStruct(&vilt)
		if err = updatedQuery.ExecRelease(); err != nil {
			log.Printf("Error while updating vilt: %v", err)
			return nil, err
		}
	}

	learners := int(vilt.NoOfLearners)
	var trainers []*string
	for _, vv := range vilt.Trainers {
		v := vv
		trainers = append(trainers, &v)
	}
	var moderators []*string
	for _, vv := range vilt.Moderators {
		v := vv
		moderators = append(moderators, &v)
	}
	cs := strconv.Itoa(int(vilt.CourseStartDate))
	ce := strconv.Itoa(int(vilt.CourseEndDate))
	ca := strconv.Itoa(int(vilt.CreatedAt))
	uaStr := strconv.Itoa(int(ua))

	res := model.Vilt{
		ID:                 &vilt.Id,
		CourseID:           &vilt.CourseId,
		LspID:              &vilt.LspId,
		NoOfLearners:       &learners,
		Trainers:           trainers,
		Moderators:         moderators,
		CourseStartDate:    &cs,
		CourseEndDate:      &ce,
		Curriculum:         &vilt.Curriculum,
		IsTrainerDecided:   &vilt.IsTrainerDecided,
		IsModeratorDecided: &vilt.IsModeratorDecided,
		IsStartDateDecided: &vilt.IsStartDateDecided,
		IsEndDateDecided:   &vilt.IsEndDateDecided,
		CreatedAt:          &ca,
		CreatedBy:          &vilt.CreatedBy,
		UpdatedAt:          &uaStr,
		UpdatedBy:          &email,
		Status:             &vilt.Status,
	}
	return &res, nil
}

func GetViltData(ctx context.Context, courseID *string) ([]*model.Vilt, error) {
	_, err := identity.GetClaimsFromContext(ctx)
	if err != nil {
		return nil, err
	}
	session, err := global.CassPool.GetSession(ctx, "viltz")
	if err != nil {
		return nil, err
	}
	CassSession := session
	qryStr := fmt.Sprintf(`SELECT * FROM viltz.vilt_master WHERE course_id='%s' ALLOW FILTERING`, *courseID)
	getViltDetails := func() (viltDetails []viltz.ViltMaster, err error) {
		q := CassSession.Query(qryStr, nil)
		defer q.Release()
		iter := q.Iter()
		return viltDetails, iter.Select(&viltDetails)
	}
	vilts, err := getViltDetails()
	if err != nil {
		return nil, err
	}
	if len(vilts) == 0 {
		return nil, nil
	}

	var wg sync.WaitGroup
	res := make([]*model.Vilt, len(vilts))
	for kk, vv := range vilts {
		v := vv
		wg.Add(1)
		go func(vilt viltz.ViltMaster, k int) {
			defer wg.Done()
			learners := int(vilt.NoOfLearners)
			var trainers []*string
			for _, vv := range vilt.Trainers {
				v := vv
				trainers = append(trainers, &v)
			}

			var moderators []*string
			for _, vv := range vilt.Moderators {
				v := vv
				moderators = append(moderators, &v)
			}
			cs := strconv.Itoa(int(vilt.CourseStartDate))
			ce := strconv.Itoa(int(vilt.CourseEndDate))
			ca := strconv.Itoa(int(vilt.CreatedAt))
			ua := strconv.Itoa(int(vilt.UpdatedAt))
			tmp := model.Vilt{
				ID:                 &vilt.Id,
				LspID:              &vilt.LspId,
				CourseID:           &vilt.CourseId,
				NoOfLearners:       &learners,
				Trainers:           trainers,
				Moderators:         moderators,
				CourseStartDate:    &cs,
				CourseEndDate:      &ce,
				Curriculum:         &vilt.Curriculum,
				IsTrainerDecided:   &vilt.IsTrainerDecided,
				IsModeratorDecided: &vilt.IsModeratorDecided,
				IsStartDateDecided: &vilt.IsStartDateDecided,
				IsEndDateDecided:   &vilt.IsEndDateDecided,
				CreatedAt:          &ca,
				CreatedBy:          &vilt.CreatedBy,
				UpdatedAt:          &ua,
				UpdatedBy:          &vilt.UpdatedBy,
				Status:             &vilt.Status,
			}
			res[k] = &tmp

		}(v, kk)
	}
	wg.Wait()

	return res, nil
}

func GetViltDataByID(ctx context.Context, id *string) (*model.Vilt, error) {
	_, err := identity.GetClaimsFromContext(ctx)
	if err != nil {
		return nil, err
	}
	session, err := global.CassPool.GetSession(ctx, "viltz")
	if err != nil {
		return nil, err
	}
	CassSession := session
	qryStr := fmt.Sprintf(`SELECT * FROM viltz.vilt_master WHERE id='%s' ALLOW FILTERING`, *id)
	getViltDetails := func() (viltDetails []viltz.ViltMaster, err error) {
		q := CassSession.Query(qryStr, nil)
		defer q.Release()
		iter := q.Iter()
		return viltDetails, iter.Select(&viltDetails)
	}
	vilts, err := getViltDetails()
	if err != nil {
		return nil, err
	}
	if len(vilts) == 0 {
		return nil, nil
	}

	vilt := vilts[0]
	learners := int(vilt.NoOfLearners)
	var trainers []*string
	for _, vv := range vilt.Trainers {
		v := vv
		trainers = append(trainers, &v)
	}

	var moderators []*string
	for _, vv := range vilt.Moderators {
		v := vv
		moderators = append(moderators, &v)
	}
	cs := strconv.Itoa(int(vilt.CourseStartDate))
	ce := strconv.Itoa(int(vilt.CourseEndDate))
	ca := strconv.Itoa(int(vilt.CreatedAt))
	ua := strconv.Itoa(int(vilt.UpdatedAt))
	res := model.Vilt{
		ID:                 &vilt.Id,
		LspID:              &vilt.LspId,
		CourseID:           &vilt.CourseId,
		NoOfLearners:       &learners,
		Trainers:           trainers,
		Moderators:         moderators,
		CourseStartDate:    &cs,
		CourseEndDate:      &ce,
		Curriculum:         &vilt.Curriculum,
		IsTrainerDecided:   &vilt.IsTrainerDecided,
		IsModeratorDecided: &vilt.IsModeratorDecided,
		IsStartDateDecided: &vilt.IsStartDateDecided,
		IsEndDateDecided:   &vilt.IsEndDateDecided,
		CreatedAt:          &ca,
		CreatedBy:          &vilt.CreatedBy,
		UpdatedAt:          &ua,
		UpdatedBy:          &vilt.UpdatedBy,
		Status:             &vilt.Status,
	}
	return &res, nil
}
