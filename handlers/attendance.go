package handlers

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"sync"

	"github.com/zicops/contracts/userz"
	"github.com/zicops/contracts/viltz"
	"github.com/zicops/zicops-vilt-manager/global"
	"github.com/zicops/zicops-vilt-manager/graph/model"
	"github.com/zicops/zicops-vilt-manager/lib/identity"
)

func GetTopicAttendance(ctx context.Context, topicID string) ([]*model.TopicAttendance, error) {
	_, err := identity.GetClaimsFromContext(ctx)
	if err != nil {
		return nil, err
	}

	//for getting duration of topic
	session, err := global.CassPool.GetSession(ctx, "viltz")
	if err != nil {
		return nil, err
	}
	CassViltSession := session

	//for getting course views, and user information
	sessionU, err := global.CassPool.GetSession(ctx, "userz")
	if err != nil {
		return nil, err
	}
	CassUserSession := sessionU

	query := fmt.Sprintf(`SELECT * FROM userz.user_course_views WHERE topic_id='%s' ALLOW FILTERING`, topicID)
	getCourseViews := func() (courseData []userz.UserCourseViews, err error) {
		q := CassUserSession.Query(query, nil)
		defer q.Release()
		iter := q.Iter()
		return courseData, iter.Select(&courseData)
	}

	courseViews, err := getCourseViews()
	if err != nil {
		return nil, err
	}
	if len(courseViews) == 0 {
		return nil, err
	}

	/*
		dataChan := make(chan model.TopicAttendance, len(courseViews))
		flagChan := make(chan bool, len(courseViews))

		//for the given topic id, get all the data, for specific users
		UserTime := make(map[string]int, 0)
		res := make([]*model.TopicAttendance, 0)
		for _, vv := range courseViews {
			v := vv
			go getAttendanceData(v, dataChan, flagChan)
		}

		for i := 0; i < len(courseViews); i++ {
			select {
			case data := <- dataChan:
				//data sent from channels, append that to res

			}
		}*/

	res := make([]*model.TopicAttendance, len(courseViews))
	var wg sync.WaitGroup
	for kk, vvv := range courseViews {
		vv := vvv
		wg.Add(1)
		go func(k int, v userz.UserCourseViews) {
			defer wg.Done()

			ca := strconv.Itoa(int(v.CreatedAt))
			ua := strconv.Itoa(int(v.UpdatedAt))
			var subCat []*string
			for _, datas := range v.SubCategories {
				data := datas
				subCat = append(subCat, &data)
			}
			timeInt := int(v.Time)
			attendance := model.TopicAttendance{
				TopicID:       &topicID,
				CourseID:      &v.CourseId,
				UserID:        &v.Users,
				FirstJoinTime: &ca,
				LastLeaveTime: &ua,
				Duration:      &timeInt,
				Category:      &v.Category,
				SubCategories: subCat,
				DateValue:     &v.DateValue,
			}
			qryStr := fmt.Sprintf(`SELECT * FROM viltz.topic_classroom WHERE topic_id='%s' ALLOW FILTERING`, topicID)
			getTopicDetails := func() (details []viltz.TopicClassroom, err error) {
				q := CassViltSession.Query(qryStr, nil)
				defer q.Release()
				iter := q.Iter()
				return details, iter.Select(&details)
			}
			topics, err := getTopicDetails()
			if err != nil {
				log.Printf("Got error: %v", err)
				return
			}
			if len(topics) == 0 {
				return
			}
			topic := topics[0]
			ret := v.Time * 100 / topic.Duration
			retStr := strconv.Itoa(int(ret))
			attendance.Retention = &retStr
			res[k] = &attendance

		}(kk, vv)
	}
	wg.Wait()

	return res, nil
}

// func getAttendanceData(view userz.UserCourseViews, dataChan chan model.TopicAttendance, flagChan chan bool) {
// 	ca := strconv.Itoa(int(view.CreatedAt))
// 	ua := strconv.Itoa(int(view.UpdatedAt))
// 	timeInt := int(view.Time)
// 	var subCat []*string
// 	for _, vv := range view.SubCategories {
// 		v := vv
// 		subCat = append(subCat, &v)
// 	}

// 	tmp := model.TopicAttendance{
// 		CourseID:      &view.CourseId,
// 		TopicID:       &view.TopicId,
// 		UserID:        &view.Users,
// 		FirstJoinTime: &ca,
// 		LastLeaveTime: &ua,
// 		Duration:      &timeInt,
// 		Category:      &view.Category,
// 		SubCategories: subCat,
// 		DateValue:     &view.DateValue,
// 	}

// 	dataChan <- tmp
// 	flagChan <- true

// }
