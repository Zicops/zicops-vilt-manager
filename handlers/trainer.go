package handlers

import (
	"context"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/zicops/contracts/viltz"
	"github.com/zicops/zicops-vilt-manager/global"
	"github.com/zicops/zicops-vilt-manager/graph/model"
	"github.com/zicops/zicops-vilt-manager/lib/identity"
)

func CreateTrainerData(ctx context.Context, input *model.TrainerInput) (*model.Trainer, error) {
	if input.UserID == nil {
		return nil, fmt.Errorf("please pass user id")
	}
	claims, err := identity.GetClaimsFromContext(ctx)
	if err != nil {
		return nil, err
	}

	email := claims["email"].(string)
	lsp := claims["lsp_id"].(string)
	if input.LspID != nil {
		lsp = *input.LspID
	}

	session, err := global.CassPool.GetSession(ctx, "viltz")
	if err != nil {
		return nil, err
	}
	CassSession := session

	createdAt := time.Now().Unix()
	id := uuid.New().String()
	trainer := viltz.ViltTrainer{
		TrainerId: id,
		LspId:     lsp,
		UserId:    *input.UserID,
		CreatedAt: createdAt,
		CreatedBy: email,
		UpdatedAt: createdAt,
		UpdatedBy: email,
	}

	if input.VendorID != nil {
		trainer.VendorId = *input.VendorID
	}
	if input.Expertise != nil {
		var tmp []string
		for _, vv := range input.Expertise {
			v := vv
			tmp = append(tmp, *v)
		}
		trainer.Expertise = tmp
	}
	if input.Status != nil {
		trainer.Status = *input.Status
	}

	insertQuery := CassSession.Query(viltz.ViltTrainerTable.Insert()).BindStruct(trainer)
	if err = insertQuery.Exec(); err != nil {
		return nil, err
	}

	ca := strconv.Itoa(int(createdAt))
	res := model.Trainer{
		ID:        &id,
		LspID:     &lsp,
		UserID:    input.UserID,
		VendorID:  input.VendorID,
		Expertise: input.Expertise,
		Status:    input.Status,
		CreatedAt: &ca,
		CreatedBy: &email,
		UpdatedAt: &ca,
		UpdatedBy: &email,
	}

	return &res, nil
}

func UpdateTrainerData(ctx context.Context, input *model.TrainerInput) (*model.Trainer, error) {
	if input.ID == nil {
		return nil, fmt.Errorf("please enter id")
	}
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

	qryStr := fmt.Sprintf(`SELECT * FROM viltz.trainer WHERE id='%s'`, *input.ID)
	getTrainerData := func() (data []viltz.ViltTrainer, err error) {
		q := CassSession.Query(qryStr, nil)
		defer q.Release()
		iter := q.Iter()
		return data, iter.Select(&data)
	}

	trainers, err := getTrainerData()
	if err != nil {
		return nil, err
	}
	if len(trainers) == 0 {
		return nil, nil
	}

	trainer := trainers[0]
	var updatedCols []string
	if input.Expertise != nil {
		var tmp []string
		for _, vv := range input.Expertise {
			v := vv
			tmp = append(tmp, *v)
		}
		trainer.Expertise = tmp
		updatedCols = append(updatedCols, "expertise")
	}
	if input.Status != nil {
		trainer.Status = *input.Status
		updatedCols = append(updatedCols, "status")
	}

	ua := time.Now().Unix()
	if len(updatedCols) > 0 {
		trainer.UpdatedAt = ua
		trainer.UpdatedBy = email
		updatedCols = append(updatedCols, "updated_at", "updated_by")
		stmt, names := viltz.ViltTrainerTable.Update(updatedCols...)
		updatedQuery := CassSession.Query(stmt, names).BindStruct(&trainer)
		if err = updatedQuery.ExecRelease(); err != nil {
			return nil, err
		}
	}

	var exp []*string
	for _, vv := range trainer.Expertise {
		v := vv
		exp = append(exp, &v)
	}

	updatedAt := strconv.Itoa(int(ua))
	ca := strconv.Itoa(int(trainer.CreatedAt))
	res := model.Trainer{
		ID:        &trainer.TrainerId,
		LspID:     &trainer.LspId,
		UserID:    &trainer.UserId,
		VendorID:  &trainer.VendorId,
		Expertise: exp,
		Status:    &trainer.Status,
		UpdatedAt: &updatedAt,
		UpdatedBy: &email,
		CreatedAt: &ca,
		CreatedBy: &trainer.CreatedBy,
	}

	return &res, nil
}

func GetTrainerData(ctx context.Context, lspID *string, vendorID *string, pageCursor *string, direction *string, pageSize *int, filters *model.TrainerFilters) (*model.PaginatedTrainer, error) {
	claims, err := identity.GetClaimsFromContext(ctx)
	if err != nil {
		return nil, err
	}
	lsp := claims["lsp_id"].(string)
	if lspID != nil {
		lsp = *lspID
	}

	var newPage []byte
	if pageCursor != nil && *pageCursor != "" {
		page, err := global.CryptSession.DecryptString(*pageCursor, nil)
		if err != nil {
			return nil, err
		}
		newPage = page
	}
	var pageSizeInt int
	if pageSize != nil {
		pageSizeInt = *pageSize
	} else {
		pageSizeInt = 10
	}

	session, err := global.CassPool.GetSession(ctx, "viltz")
	if err != nil {
		return nil, err
	}
	CassSession := session

	qryStr := fmt.Sprintf(`SELECT * FROM viltz.trainer WHERE lsp_id='%s' `, lsp)
	if vendorID != nil {
		qryStr += fmt.Sprintf(`AND vendor_id='%s' `, *vendorID)
	}
	// //list of users ids from trainer
	// //among those, search on basis of name
	// if filters != nil && filters.Name != nil {
	// 	name := strings.ToLower(*filters.Name)
	// 	namesArray := strings.Fields(name)
	// 	for _, vv := range namesArray {
	// 		v := vv
	// 		qryStr += fmt.Sprintf(` AND name CONTAINS '%s' `, v)
	// 	}
	// }
	qryStr += `ALLOW FILTERING`

	getTrainers := func(page []byte) (trainersData []viltz.ViltTrainer, newPage []byte, err error) {
		q := CassSession.Query(qryStr, nil)
		defer q.Release()
		q.PageState(page)
		q.PageSize(pageSizeInt)
		iter := q.Iter()
		return trainersData, iter.PageState(), iter.Select(&trainersData)
	}

	trainers, newPage, err := getTrainers(newPage)
	if err != nil {
		return nil, err
	}

	var newCursor string
	if len(newPage) != 0 {
		newCursor, err = global.CryptSession.EncryptAsString(newPage, nil)
		if err != nil {
			return nil, err
		}
	}

	if len(trainers) == 0 {
		return nil, nil
	}

	res := make([]*model.Trainer, len(trainers))
	var wg sync.WaitGroup
	for kk, vv := range trainers {
		wg.Add(1)
		go func(k int, trainer viltz.ViltTrainer) {
			defer wg.Done()

			var exp []*string
			for _, vv := range trainer.Expertise {
				v := vv
				exp = append(exp, &v)
			}
			ua := strconv.Itoa(int(trainer.UpdatedAt))
			ca := strconv.Itoa(int(trainer.CreatedAt))
			tmp := model.Trainer{
				ID:        &trainer.TrainerId,
				LspID:     &trainer.LspId,
				UserID:    &trainer.UserId,
				VendorID:  &trainer.VendorId,
				Expertise: exp,
				Status:    &trainer.Status,
				UpdatedAt: &ua,
				UpdatedBy: &trainer.UpdatedBy,
				CreatedAt: &ca,
				CreatedBy: &trainer.CreatedBy,
			}
			res[k] = &tmp
		}(kk, vv)
	}
	wg.Wait()

	result := model.PaginatedTrainer{
		Trainers:   res,
		PageCursor: &newCursor,
		PageSize:   &pageSizeInt,
		Direction:  direction,
	}
	return &result, nil
}

func GetTrainerByID(ctx context.Context, id *string) (*model.Trainer, error) {
	if id == nil {
		return nil, fmt.Errorf("please enter id")
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

	qryStr := fmt.Sprintf(`SELECT * FROM viltz.trainer WHERE id='%s' ALLOW FILTERING`, *id)
	getTrainer := func() (trainerData []viltz.ViltTrainer, err error) {
		q := CassSession.Query(qryStr, nil)
		defer q.Release()
		iter := q.Iter()
		return trainerData, iter.Select(&trainerData)
	}
	trainers, err := getTrainer()
	if err != nil {
		return nil, err
	}
	if len(trainers) == 0 {
		return nil, nil
	}
	trainer := trainers[0]

	var exp []*string
	for _, vv := range trainer.Expertise {
		v := vv
		exp = append(exp, &v)
	}
	ua := strconv.Itoa(int(trainer.UpdatedAt))
	ca := strconv.Itoa(int(trainer.CreatedAt))
	res := model.Trainer{
		ID:        &trainer.TrainerId,
		LspID:     &trainer.LspId,
		UserID:    &trainer.UserId,
		VendorID:  &trainer.VendorId,
		Expertise: exp,
		Status:    &trainer.Status,
		UpdatedAt: &ua,
		UpdatedBy: &trainer.UpdatedBy,
		CreatedAt: &ca,
		CreatedBy: &trainer.CreatedBy,
	}
	return &res, nil
}

//low standard high hd
//tile view, grid view

//media server- gcs bucket redirection, link
//local recording, to gcs
