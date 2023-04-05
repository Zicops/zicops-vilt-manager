package handlers

import (
	"fmt"

	"github.com/scylladb/gocqlx/v2"
	"github.com/zicops/contracts/coursez"
)

func GetModule(CassCoursezSession *gocqlx.Session, moduleId string, lsp string) (coursez.Module, error) {
	getQueryStr := fmt.Sprintf("SELECT * FROM coursez.module WHERE id='%s' and lsp_id='%s' and is_active=true", moduleId, lsp)
	getModule := func() (moduleData []coursez.Module, err error) {
		q := CassCoursezSession.Query(getQueryStr, nil)
		defer q.Release()
		iter := q.Iter()
		return moduleData, iter.Select(&moduleData)
	}
	n := coursez.Module{}
	modules, err := getModule()
	if err != nil {
		return n, err
	}
	if len(modules) == 0 {
		return n, fmt.Errorf("no module found")
	}

	return modules[0], nil
}

func GetCourse(CassCoursezSession *gocqlx.Session, courseId string, lsp string) (coursez.Course, error) {
	getQueryStr := fmt.Sprintf("SELECT * FROM coursez.course WHERE id='%s' and lsp_id='%s' and is_active=true", courseId, lsp)
	getCourse := func() (courseData []coursez.Course, err error) {
		q := CassCoursezSession.Query(getQueryStr, nil)
		defer q.Release()
		iter := q.Iter()
		return courseData, iter.Select(&courseData)
	}
	n := coursez.Course{}
	courses, err := getCourse()
	if err != nil {
		return n, err
	}
	if len(courses) == 0 {
		return n, fmt.Errorf("no course found")
	}

	return courses[0], nil
}
