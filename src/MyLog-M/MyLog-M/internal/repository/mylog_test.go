package repository

import (
	"MyLog-M/driver/sqlite"
	"MyLog-M/internal/domain"
	"reflect"
	"testing"
	"time"
)

func TestDatastore(t *testing.T) {
	db := sqlite.Open("../../log.db")
	defer sqlite.Close(db)

	r := NewLogRepo(db)
	testLogRepo_Get(t, *r)
	//testLogRepo_Insert(t, *r)
	testLogRepo_Tail(t, *r)

}

//"LogType" : "DEBUG","LogSeverity" : 1,"LogText": "TEST 1"
func testLogRepo_Insert(t *testing.T, db LogRepo) {
	testcases := []struct {
		req domain.Data
		id  int64
	}{
		{domain.Data{LogType: "DEBUG", LogSeverity: 1, LogText: "TEST 202"}, 202},
	}
	for i, v := range testcases {
		id, _ := db.Insert(&v.req)

		if !reflect.DeepEqual(id, v.id) {
			t.Errorf("[TEST%d]Failed. Got %v\tExpected %v\n", i+1, id, v.id)
		}
	}
}

//{"RID":1,"UnixTime":1628949212,"LocalTime":"2021-08-14T21:53:32+08:00","LogType":"DEBUG","LogSeverity":1,"LogText":"this is a","Status":{"Code":0,"Text":""}}]
func testLogRepo_Get(t *testing.T, db LogRepo) {
	testcases := []struct {
		id   int
		resp *domain.Data
	}{
		//{1, []domain.Data{{1, 1628949212, time.Date(2021, 8, 14, 21, 53, 32, 0, time.Local), "DEBUG", 1, "this is a", Status{0, ""}}}},
		{1, &domain.Data{1, 1628949212, time.Time{}, "DEBUG", 1, "this is a", domain.Status{0, ""}}},
	}
	for i, v := range testcases {
		resp, _ := db.Get(int64(v.id))

		if !reflect.DeepEqual(resp, v.resp) {
			t.Errorf("[TEST%d]Failed. Got %v\tExpected %v\n", i+1, resp, v.resp)
		}
	}
}

func testLogRepo_Tail(t *testing.T, db LogRepo) {
	testcases := []struct {
		limit int64
		resp  *[]domain.Data
	}{

		{1, &[]domain.Data{{202, 1630303817, time.Date(2021, 8, 30, 14, 10, 17, 0, time.Local), "DEBUG", 1, "TEST 202", domain.Status{0, ""}}}},
		{2, &[]domain.Data{

			{202, 1630303817, time.Date(2021, 8, 30, 14, 10, 17, 0, time.Local), "DEBUG", 1, "TEST 202", domain.Status{0, ""}},
			{201, 1630303785, time.Date(2021, 8, 30, 14, 9, 45, 0, time.Local), "DEBUG", 1, "TEST 2", domain.Status{0, ""}},
		}},
	}
	for i, v := range testcases {
		resp, _ := db.Tail(v.limit)

		if !reflect.DeepEqual(resp, v.resp) {
			t.Errorf("[TEST%d]Failed. Got %v\tExpected %v\n", i+1, resp, v.resp)
		}
	}
}
