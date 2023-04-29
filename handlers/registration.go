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
