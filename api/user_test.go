package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	mock_db "simplebank/db/mock"
	db "simplebank/db/sqlc"
	db_utils "simplebank/db/utils"
	"simplebank/utils"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

type eqCreateUserParamsMatcher struct {
	arg      db.CreateUserParams
	password string
}

func (e eqCreateUserParamsMatcher) Matches(x interface{}) bool {
	arg, ok := x.(db.CreateUserParams)
	if !ok {
		return false
	}

	err := utils.CheckPasword(e.password, arg.HashedPassword)
	if err != nil {
		return false
	}

	e.arg.HashedPassword = arg.HashedPassword
	return reflect.DeepEqual(e.arg, x)
}

func (e eqCreateUserParamsMatcher) String() string {
	return fmt.Sprintf("matches arg %v and password %v", e.arg, e.password)
}

func eqCreateUserParams(arg db.CreateUserParams, password string) gomock.Matcher {
	return eqCreateUserParamsMatcher{arg, password}
}

func TestServer_createUser(t *testing.T) {
	user, password := randomUser(t)

	type args struct {
		ctx *gin.Context
	}
	tests := []struct {
		name          string
		body          gin.H
		buildStubs    func(store *mock_db.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: gin.H{
				"username":  user.Username,
				"password":  password,
				"full_name": user.FullName,
				"email":     user.Email,
			},
			buildStubs: func(store *mock_db.MockStore) {
				arg := db.CreateUserParams{
					Username: user.Username,
					FullName: user.FullName,
					Email:    user.Email,
				}

				store.EXPECT().
					CreateUser(gomock.Any(), eqCreateUserParams(arg, password)).
					Times(1).
					Return(user, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mock_db.NewMockStore(ctrl)
			tt.buildStubs(store)

			server := NewServer(store)
			recorder := httptest.NewRecorder()

			data, err := json.Marshal(tt.body)
			require.NoError(t, err)

			url := fmt.Sprintf("/users")
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			tt.checkResponse(t, recorder)
		})
	}
}

func randomUser(t *testing.T) (user db.User, password string) {
	password = db_utils.RandomString(6)
	hashedPassword, err := utils.HashPassword(password)
	require.NoError(t, err)

	user = db.User{
		Username:       db_utils.RandomOwner(),
		HashedPassword: hashedPassword,
		FullName:       db_utils.RandomOwner(),
		Email:          db_utils.RandomEmail(),
	}

	return
}
