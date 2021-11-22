package service

import (
	"MyLog-M/internal/domain"
	"MyLog-M/internal/service/mock"
	"errors"
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"
)

func TestService_Tail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	m := mock.NewMockrepository(ctrl)

	tests := []struct {
		name    string
		args    int64
		mock    func(*mock.Mockrepository)
		want    *[]domain.Data
		wantErr bool
	}{
		{"success case", 1, func(m *mock.Mockrepository) {
			m.EXPECT().Tail(int64(1)).Return(&[]domain.Data{
				{RID: 1},
			}, nil)
		}, &[]domain.Data{{RID: 1}}, false},
		{"failure case", 1, func(m *mock.Mockrepository) {
			m.EXPECT().Tail(int64(1)).Return(nil, errors.New("anything"))
		}, &[]domain.Data{}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock(m)
			s := New(m)
			got, err := s.Tail(tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.Tail() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Service.Tail() = %v, want %v", got, tt.want)
			}
		})
	}
}
