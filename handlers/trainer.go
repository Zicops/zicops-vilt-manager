package handlers

import (
	"context"
	"fmt"
	"strconv"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/google/uuid"
	"github.com/zicops/contracts/coursez"
	"github.com/zicops/contracts/viltz"
	"github.com/zicops/zicops-vilt-manager/global"
	"github.com/zicops/zicops-vilt-manager/graph/model"
	"github.com/zicops/zicops-vilt-manager/lib/db/bucket"
	"github.com/zicops/zicops-vilt-manager/lib/googleprojectlib"
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
	if input.YearsOfExperience != nil {
		trainer.YearsOfExperience = *input.YearsOfExperience
	}
	if input.Description != nil {
		trainer.Description = *input.Description
	}
	if input.Github != nil {
		trainer.Github = *input.Github
	}
	if input.Linkedin != nil {
		trainer.LinkedIn = *input.Linkedin
	}
	if input.Website != nil {
		trainer.Website = *input.Website
	}

	insertQuery := CassSession.Query(viltz.ViltTrainerTable.Insert()).BindStruct(trainer)
	if err = insertQuery.Exec(); err != nil {
		return nil, err
	}

	ca := strconv.Itoa(int(createdAt))
	res := model.Trainer{
		ID:                &id,
		LspID:             &lsp,
		UserID:            input.UserID,
		VendorID:          input.VendorID,
		Expertise:         input.Expertise,
		Status:            input.Status,
		YearsOfExperience: input.YearsOfExperience,
		Website:           input.Website,
		Linkedin:          input.Linkedin,
		Github:            input.Github,
		Description:       input.Description,
		CreatedAt:         &ca,
		CreatedBy:         &email,
		UpdatedAt:         &ca,
		UpdatedBy:         &email,
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
	if input.Website != nil {
		trainer.Website = *input.Website
		updatedCols = append(updatedCols, "website")
	}
	if input.Description != nil {
		trainer.Description = *input.Description
		updatedCols = append(updatedCols, "description")
	}
	if input.Github != nil {
		trainer.Github = *input.Github
		updatedCols = append(updatedCols, "github")
	}
	if input.Linkedin != nil {
		trainer.LinkedIn = *input.Linkedin
		updatedCols = append(updatedCols, "linkedin")
	}
	if input.YearsOfExperience != nil {
		trainer.YearsOfExperience = *input.YearsOfExperience
		updatedCols = append(updatedCols, "year_of_experience")
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
		ID:                &trainer.TrainerId,
		LspID:             &trainer.LspId,
		UserID:            &trainer.UserId,
		VendorID:          &trainer.VendorId,
		Expertise:         exp,
		Status:            &trainer.Status,
		Website:           &trainer.Website,
		Github:            &trainer.Github,
		Linkedin:          &trainer.LinkedIn,
		YearsOfExperience: &trainer.YearsOfExperience,
		Description:       &trainer.Description,
		UpdatedAt:         &updatedAt,
		UpdatedBy:         &email,
		CreatedAt:         &ca,
		CreatedBy:         &trainer.CreatedBy,
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
				ID:                &trainer.TrainerId,
				LspID:             &trainer.LspId,
				UserID:            &trainer.UserId,
				VendorID:          &trainer.VendorId,
				Expertise:         exp,
				Status:            &trainer.Status,
				YearsOfExperience: &trainer.YearsOfExperience,
				Website:           &trainer.Website,
				Linkedin:          &trainer.LinkedIn,
				Github:            &trainer.Github,
				Description:       &trainer.Description,
				UpdatedAt:         &ua,
				UpdatedBy:         &trainer.UpdatedBy,
				CreatedAt:         &ca,
				CreatedBy:         &trainer.CreatedBy,
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

func GetTrainerCourses(ctx context.Context, userID *string) ([]*model.Course, error) {
	if userID == nil || *userID == "" {
		return nil, fmt.Errorf("please enter userId")
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

	qryStr := fmt.Sprintf(`SELECT * FROM viltz.trainer WHERE user_id='%s' ALLOW FILTERING`, *userID)
	getTrainer := func() (trainersData []viltz.ViltTrainer, err error) {
		q := CassViltSession.Query(qryStr, nil)
		defer q.Release()
		iter := q.Iter()
		return trainersData, iter.Select(&trainersData)
	}

	trainers, err := getTrainer()
	if err != nil {
		return nil, err
	}
	if len(trainers) == 0 {
		return nil, nil
	}
	trainer := trainers[0]

	query := fmt.Sprintf(`SELECT * FROM viltz.topic_classroom WHERE trainer CONTAINS '%s' ALLOW FILTERING`, trainer.TrainerId)
	getTopicClassroomData := func() (topicsData []viltz.TopicClassroom, err error) {
		q := CassViltSession.Query(query, nil)
		defer q.Release()
		iter := q.Iter()
		return topicsData, iter.Select(&topicsData)
	}

	topicClassrooms, err := getTopicClassroomData()
	if err != nil {
		log.Printf("Got error while getting topic classroom data: %v", err)
		return nil, err
	}
	if len(topicClassrooms) == 0 {
		return nil, err
	}

	res := make([]*model.Course, len(topicClassrooms))
	var wg sync.WaitGroup
	for kk, vvv := range topicClassrooms {
		vv := vvv
		wg.Add(1)
		go func(k int, v viltz.TopicClassroom, ctx context.Context) {
			defer wg.Done()
			course := getCourse(ctx, v.TopicId)
			if course == nil {
				return
			}
			res[k] = course
		}(kk, vv, ctx)
	}
	wg.Wait()

	//remove redundant courses
	var response []*model.Course
	flags := make(map[string]bool)
	for _, vv := range res {
		v := vv
		if flags[*v.ID] {
			continue
		}
		response = append(response, v)
		flags[*v.ID] = true
	}

	return response, nil
}

func getCourse(ctx context.Context, topicId string) *model.Course {
	session, err := global.CassPool.GetSession(ctx, "coursez")
	if err != nil {
		log.Printf("Got error while getting session: %v", err)
		return nil
	}
	CassCoursezSession := session

	qryStr := fmt.Sprintf(`SELECT * FROM coursez.topic WHERE id='%s' ALLOW FILTERING`, topicId)
	getTopicDetails := func() (topics []coursez.Topic, err error) {
		q := CassCoursezSession.Query(qryStr, nil)
		defer q.Release()
		iter := q.Iter()
		return topics, iter.Select(&topics)
	}

	topics, err := getTopicDetails()
	if err != nil {
		log.Printf("Got error while getting topic details: %v", err)
		return nil
	}
	if len(topics) == 0 {
		return nil
	}
	query := fmt.Sprintf(`SELECT * FROM coursez.course WHERE id='%s' ALLOW FILTERING`, topics[0].CourseID)
	getCourseDetails := func() (courseDetails []coursez.Course, err error) {
		q := CassCoursezSession.Query(query, nil)
		defer q.Release()
		iter := q.Iter()
		return courseDetails, iter.Select(&courseDetails)
	}
	courses, err := getCourseDetails()
	if err != nil {
		log.Printf("Got error while getting course details: %v", err)
		return nil
	}
	if len(courses) == 0 {
		return nil
	}
	course := courses[0]
	createdAt := strconv.FormatInt(course.CreatedAt, 10)
	updatedAt := strconv.FormatInt(course.UpdatedAt, 10)
	language := make([]*string, 0)
	takeaways := make([]*string, 0)
	outcomes := make([]*string, 0)
	prequisites := make([]*string, 0)
	goodFor := make([]*string, 0)
	mustFor := make([]*string, 0)
	relatedSkills := make([]*string, 0)
	approvers := make([]*string, 0)
	subCatsRes := make([]*model.SubCategories, 0)

	for _, lang := range course.Language {
		langCopied := lang
		language = append(language, &langCopied)
	}
	for _, take := range course.Benefits {
		takeCopied := take
		takeaways = append(takeaways, &takeCopied)
	}
	for _, out := range course.Outcomes {
		outCopied := out
		outcomes = append(outcomes, &outCopied)
	}
	for _, preq := range course.Prequisites {
		preCopied := preq
		prequisites = append(prequisites, &preCopied)
	}
	for _, good := range course.GoodFor {
		goodCopied := good
		goodFor = append(goodFor, &goodCopied)
	}
	for _, must := range course.MustFor {
		mustCopied := must
		mustFor = append(mustFor, &mustCopied)
	}
	for _, relSkill := range course.RelatedSkills {
		relCopied := relSkill
		relatedSkills = append(relatedSkills, &relCopied)
	}
	for _, approver := range course.Approvers {
		appoverCopied := approver
		approvers = append(approvers, &appoverCopied)
	}
	for _, subCat := range course.SubCategories {
		subCopied := subCat
		var subCR model.SubCategories
		subCR.Name = &subCopied.Name
		subCR.Rank = &subCopied.Rank
		subCatsRes = append(subCatsRes, &subCR)
	}

	storageC := bucket.NewStorageHandler()
	gproject := googleprojectlib.GetGoogleProjectID()
	err = storageC.InitializeStorageClient(ctx, gproject)
	if err != nil {
		log.Errorf("Failed to initialize storage: %v", err.Error())
		return nil
	}
	tileUrl := course.TileImage
	if course.TileImageBucket != "" {
		tileUrl = storageC.GetSignedURLForObject(ctx, course.TileImageBucket)
	}
	imageUrl := course.Image
	if course.ImageBucket != "" {
		imageUrl = storageC.GetSignedURLForObject(ctx, course.ImageBucket)
	}
	previewUrl := course.PreviewVideo
	if course.PreviewVideoBucket != "" {
		previewUrl = storageC.GetSignedURLForObject(ctx, course.PreviewVideoBucket)
	}
	currentCourse := model.Course{
		ID:                 &course.ID,
		Name:               &course.Name,
		LspID:              &course.LspId,
		Publisher:          &course.Publisher,
		Description:        &course.Description,
		Summary:            &course.Summary,
		Instructor:         &course.Instructor,
		Owner:              &course.Owner,
		Duration:           &course.Duration,
		ExpertiseLevel:     &course.ExpertiseLevel,
		Language:           language,
		Benefits:           takeaways,
		Outcomes:           outcomes,
		CreatedAt:          &createdAt,
		UpdatedAt:          &updatedAt,
		Type:               &course.Type,
		Prequisites:        prequisites,
		GoodFor:            goodFor,
		MustFor:            mustFor,
		RelatedSkills:      relatedSkills,
		PublishDate:        &course.PublishDate,
		ExpiryDate:         &course.ExpiryDate,
		ExpectedCompletion: &course.ExpectedCompletion,
		QaRequired:         &course.QARequired,
		Approvers:          approvers,
		CreatedBy:          &course.CreatedBy,
		UpdatedBy:          &course.UpdatedBy,
		Status:             &course.Status,
		IsDisplay:          &course.IsDisplay,
		Category:           &course.Category,
		SubCategory:        &course.SubCategory,
		SubCategories:      subCatsRes,
		IsActive:           &course.IsActive,
	}
	if course.TileImageBucket != "" {
		currentCourse.TileImage = &tileUrl
	}
	if course.ImageBucket != "" {
		currentCourse.Image = &imageUrl
	}
	if course.PreviewVideoBucket != "" {
		currentCourse.PreviewVideo = &previewUrl
	}

	return &currentCourse
}
