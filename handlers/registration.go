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

func RegisterUserForCourse(ctx context.Context, input *model.UserCourseRegisterInput) (*model.UserCourseRegister, error) {
	if input == nil || input.CourseID == nil || input.UserID == nil || input.RegistrationDate == nil {
		return nil, fmt.Errorf("please pass all the parameters to the query")
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
	CassViltSession := session

	id := uuid.New().String()
	reg := viltz.UserCourseRegister{
		Id:               id,
		CourseId:         *input.CourseID,
		UserId:           *input.UserID,
		RegistrationDate: int64(*input.RegistrationDate),
	}
	if input.Invoice != nil {
		reg.Invoice = *input.Invoice
	}
	if input.Status != nil {
		reg.Status = *input.Status
	}
	ca := time.Now().Unix()
	reg.CreatedAt = ca
	reg.UpdatedAt = ca
	reg.CreatedBy = email
	reg.UpdatedBy = email

	insertQuery := CassViltSession.Query(viltz.UserCourseRegisterTable.Insert()).BindStruct(reg)
	if err = insertQuery.Exec(); err != nil {
		return nil, err
	}

	createdAt := strconv.Itoa(int(ca))
	res := model.UserCourseRegister{
		ID:               &id,
		CourseID:         input.CourseID,
		UserID:           input.UserID,
		RegistrationDate: input.RegistrationDate,
		Invoice:          input.Invoice,
		Status:           input.Status,
		CreatedAt:        &createdAt,
		CreatedBy:        &email,
		UpdatedAt:        &createdAt,
		UpdatedBy:        &email,
	}
	return &res, nil
}

func UpdateRegistrationForCourse(ctx context.Context, input *model.UserCourseRegisterInput) (*model.UserCourseRegister, error) {
	if input.ID == nil {
		return nil, fmt.Errorf("please enter ID")
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
	CassViltSession := session

	qryStr := fmt.Sprintf(`SELECT * FROM viltz.user_course_register WHERE id='%s' ALLOW FILTERING`, *input.ID)
	getRegistration := func() (user_registrations []viltz.UserCourseRegister, err error) {
		q := CassViltSession.Query(qryStr, nil)
		defer q.Release()
		iter := q.Iter()
		return user_registrations, iter.Select(&user_registrations)
	}

	registrations, err := getRegistration()
	if err != nil {
		return nil, err
	}

	if len(registrations) == 0 {
		return nil, nil
	}

	registration := registrations[0]

	var updatedCols []string
	if input.CourseID != nil {
		registration.CourseId = *input.CourseID
		updatedCols = append(updatedCols, "course_id")
	}
	if input.Invoice != nil {
		registration.Invoice = *input.Invoice
		updatedCols = append(updatedCols, "invoice")
	}
	if input.RegistrationDate != nil {
		registration.RegistrationDate = int64(*input.RegistrationDate)
		updatedCols = append(updatedCols, "registration_date")
	}
	if input.Status != nil {
		registration.Status = *input.Status
		updatedCols = append(updatedCols, "status")
	}
	if input.UserID != nil {
		registration.UserId = *input.UserID
		updatedCols = append(updatedCols, "user_id")
	}
	var ua int64
	if len(updatedCols) > 0 {
		ua = time.Now().Unix()
		registration.UpdatedAt = ua
		registration.UpdatedBy = email
		updatedCols = append(updatedCols, "updated_at", "updated_by")
		stmt, names := viltz.UserCourseRegisterTable.Update(updatedCols...)
		updateQuery := CassViltSession.Query(stmt, names).BindStruct(&registration)
		if err = updateQuery.ExecRelease(); err != nil {
			return nil, err
		}
	}

	rd := int(registration.RegistrationDate)
	ca := strconv.Itoa(int(registration.CreatedAt))
	uaStr := strconv.Itoa(int(registration.UpdatedAt))
	res := model.UserCourseRegister{
		ID:               input.ID,
		CourseID:         &registration.CourseId,
		UserID:           &registration.UserId,
		RegistrationDate: &rd,
		Invoice:          &registration.Invoice,
		Status:           &registration.Status,
		CreatedAt:        &ca,
		CreatedBy:        &registration.CreatedBy,
		UpdatedAt:        &uaStr,
		UpdatedBy:        &registration.UpdatedBy,
	}
	return &res, nil
}

func GetAllRegistrations(ctx context.Context, courseID *string, pageCursor *string, direction *string, pageSize *int) (*model.PaginatedUserCourseRegister, error) {
	if courseID == nil {
		return nil, fmt.Errorf("please enter course id")
	}
	_, err := identity.GetClaimsFromContext(ctx)
	if err != nil {
		return nil, err
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
	CassViltSession := session

	qryStr := fmt.Sprintf(`SELECT * FROM viltz.user_course_register WHERE course_id='%s' ALLOW FILTERING`, *courseID)
	getRegistrations := func(page []byte) (registrations []viltz.UserCourseRegister, newPage []byte, err error) {
		q := CassViltSession.Query(qryStr, nil)
		defer q.Release()
		q.PageState(page)
		q.PageSize(pageSizeInt)
		iter := q.Iter()
		return registrations, iter.PageState(), iter.Select(&registrations)
	}

	reg, newPage, err := getRegistrations(newPage)
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

	if len(reg) == 0 {
		return nil, nil
	}

	res := make([]*model.UserCourseRegister, len(reg))
	var wg sync.WaitGroup
	for kk, vvv := range reg {
		vv := vvv
		wg.Add(1)
		go func(k int, v viltz.UserCourseRegister) {
			defer wg.Done()
			rd := int(v.RegistrationDate)
			ca := strconv.Itoa(int(v.CreatedAt))
			ua := strconv.Itoa(int(v.UpdatedAt))
			tmp := model.UserCourseRegister{
				ID:               &v.Id,
				CourseID:         &v.CourseId,
				UserID:           &v.UserId,
				RegistrationDate: &rd,
				Invoice:          &v.Invoice,
				Status:           &v.Status,
				CreatedAt:        &ca,
				CreatedBy:        &v.CreatedBy,
				UpdatedAt:        &ua,
				UpdatedBy:        &v.UpdatedBy,
			}
			res[k] = &tmp
		}(kk, vv)
	}
	wg.Wait()

	ucr := model.PaginatedUserCourseRegister{
		Data:       res,
		PageCursor: &newCursor,
		PageSize:   &pageSizeInt,
		Direction:  direction,
	}
	return &ucr, nil
}

func GetRegistrationDetails(ctx context.Context, id *string) (*model.UserCourseRegister, error) {
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
	CassViltSession := session
	qryStr := fmt.Sprintf(`SELECT * FROM viltz.user_course_register WHERE id='%s' ALLOW FILTERING`, *id)
	getRegistration := func() (user_registrations []viltz.UserCourseRegister, err error) {
		q := CassViltSession.Query(qryStr, nil)
		defer q.Release()
		iter := q.Iter()
		return user_registrations, iter.Select(&user_registrations)
	}

	registrations, err := getRegistration()
	if err != nil {
		return nil, err
	}

	if len(registrations) == 0 {
		return nil, nil
	}

	registration := registrations[0]

	rd := int(registration.RegistrationDate)
	ca := strconv.Itoa(int(registration.CreatedAt))
	ua := strconv.Itoa(int(registration.UpdatedAt))
	res := model.UserCourseRegister{
		ID:               &registration.Id,
		CourseID:         &registration.CourseId,
		UserID:           &registration.UserId,
		RegistrationDate: &rd,
		Invoice:          &registration.Invoice,
		Status:           &registration.Status,
		CreatedAt:        &ca,
		CreatedBy:        &registration.CreatedBy,
		UpdatedAt:        &ua,
		UpdatedBy:        &registration.UpdatedBy,
	}
	return &res, nil
}
