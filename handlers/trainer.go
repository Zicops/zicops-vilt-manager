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

	qryStr := fmt.Sprintf(`SELECT * FROM vendorz.trainer WHERE id='%s'`, *input.ID)
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

func GetTrainerData(ctx context.Context, lspID *string, vendorID *string) ([]*model.Trainer, error) {
	claims, err := identity.GetClaimsFromContext(ctx)
	if err != nil {
		return nil, err
	}
	lsp := claims["lsp_id"].(string)
	if lspID != nil {
		lsp = *lspID
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
	qryStr += `ALLOW FILTERING`

	getTrainers := func() (trainersData []viltz.ViltTrainer, err error) {
		q := CassSession.Query(qryStr, nil)
		defer q.Release()
		iter := q.Iter()
		return trainersData, iter.Select(&trainersData)
	}

	trainers, err := getTrainers()
	if err != nil {
		return nil, err
	}
	if len(trainers) == 0 {
		return nil, nil
	}

	var res []*model.Trainer
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

	return res, nil
}