package gist

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/magiconair/properties/assert"
	"github.com/nanmenkaimak/github-gist/internal/gist/entity"
	mock_repository "github.com/nanmenkaimak/github-gist/internal/gist/repository/mocks"
	"github.com/nanmenkaimak/github-gist/internal/gist/transport"
	"go.uber.org/mock/gomock"
)

func TestCreateGist(t *testing.T) {
	testUUID, _ := uuid.Parse("2d96eee6-c261-45c1-b9f5-2a10b088e68c")

	type mockBehavior func(s *mock_repository.MockRepository, newGist entity.GistRequest)

	testCreate := []struct {
		name               string
		inputUser          entity.GistRequest
		mockBehavior       mockBehavior
		expectedStatusBody *CreateGistResponse
		expectedError      error
	}{
		{
			name: "OK",
			inputUser: entity.GistRequest{
				Gist: entity.Gist{
					Name:        "first gist",
					Description: "che tam",
					Visible:     true,
				},
				Commit: entity.Commit{
					Comment: "men commit",
				},
				Files: []entity.File{
					{
						Name: "main.py",
						Code: "package main\n import \"fmt\"\nfunc main() {\n fmt.Println(\"ri e\")\n}  ",
					},
					{
						Name: "main2.py",
						Code: "package main\n import \"fmt\"\nfunc main() {\n fmt.Println(\"che \")\n}  ",
					},
				},
			},
			mockBehavior: func(s *mock_repository.MockRepository, newGist entity.GistRequest) {
				s.EXPECT().CreateGist(newGist).Return(testUUID, nil)
			},
			expectedStatusBody: &CreateGistResponse{
				GistID: testUUID,
			},
			expectedError: nil,
		},
		{
			name:      "empty request",
			inputUser: entity.GistRequest{},
			mockBehavior: func(s *mock_repository.MockRepository, newGist entity.GistRequest) {
				s.EXPECT().CreateGist(newGist).Return(uuid.Nil, errors.New("n"))
			},
			expectedStatusBody: nil,
			expectedError:      fmt.Errorf("creating gist err: %v", "n"),
		},
		{
			name: "missing values",
			inputUser: entity.GistRequest{
				Commit: entity.Commit{
					Comment: "men commit",
				},
				Files: []entity.File{
					{
						Name: "main.py",
						Code: "package main\n import \"fmt\"\nfunc main() {\n fmt.Println(\"ri e\")\n}  ",
					},
					{
						Name: "main2.py",
						Code: "package main\n import \"fmt\"\nfunc main() {\n fmt.Println(\"che \")\n}  ",
					},
				},
			},
			mockBehavior: func(s *mock_repository.MockRepository, newGist entity.GistRequest) {
				s.EXPECT().CreateGist(newGist).Return(uuid.Nil, errors.New("n"))
			},
			expectedStatusBody: nil,
			expectedError:      fmt.Errorf("creating gist err: %v", "n"),
		},
	}

	for _, tt := range testCreate {
		t.Run(tt.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			newRepo := mock_repository.NewMockRepository(c)
			tt.mockBehavior(newRepo, tt.inputUser)

			var userTransport *transport.UserGrpcTransport
			usecase := NewGistService(newRepo, userTransport)

			response, err := usecase.CreateGist(context.Background(), tt.inputUser)

			assert.Equal(t, response, tt.expectedStatusBody)
			assert.Equal(t, err, tt.expectedError)
		})
	}
}

func TestGetAllGists(t *testing.T) {
	type mockBehavior func(s *mock_repository.MockRepository, request GetAllGistsRequest)

	test := []struct {
		name               string
		inputUser          GetAllGistsRequest
		mockBehavior       mockBehavior
		expectedStatusBody *[]entity.GistRequest
		expectedError      error
	}{
		{
			name: "OK",
			inputUser: GetAllGistsRequest{
				Direction: "asc",
				Sort:      "created_at",
			},
			mockBehavior: func(s *mock_repository.MockRepository, request GetAllGistsRequest) {
				s.EXPECT().GetOtherAllGists(request.Sort, request.Direction).Return([]entity.GistRequest{}, nil)
			},
			expectedStatusBody: &[]entity.GistRequest{},
			expectedError:      nil,
		},
		{
			name: "missing sort",
			inputUser: GetAllGistsRequest{
				Direction: "asc",
				Sort:      "",
			},
			mockBehavior: func(s *mock_repository.MockRepository, request GetAllGistsRequest) {
				s.EXPECT().GetOtherAllGists("created_at", request.Direction).Return([]entity.GistRequest{}, nil)
			},
			expectedStatusBody: &[]entity.GistRequest{},
			expectedError:      nil,
		},
		{
			name: "missing direction",
			inputUser: GetAllGistsRequest{
				Direction: "",
				Sort:      "created_at",
			},
			mockBehavior: func(s *mock_repository.MockRepository, request GetAllGistsRequest) {
				s.EXPECT().GetOtherAllGists(request.Sort, "desc").Return([]entity.GistRequest{}, nil)
			},
			expectedStatusBody: &[]entity.GistRequest{},
			expectedError:      nil,
		},
		{
			name: "wrong input direction",
			inputUser: GetAllGistsRequest{
				Direction: "nanmenkaimak",
				Sort:      "created_at",
			},
			mockBehavior: func(s *mock_repository.MockRepository, request GetAllGistsRequest) {
				s.EXPECT().GetOtherAllGists(request.Sort, request.Direction).Return(nil, errors.New("katego"))
			},
			expectedStatusBody: nil,
			expectedError:      fmt.Errorf("getting all gists err: %v", "katego"),
		},
		{
			name: "wrong input sort",
			inputUser: GetAllGistsRequest{
				Direction: "asc",
				Sort:      "eine",
			},
			mockBehavior: func(s *mock_repository.MockRepository, request GetAllGistsRequest) {
				s.EXPECT().GetOtherAllGists(request.Sort, request.Direction).Return(nil, errors.New("katego"))
			},
			expectedStatusBody: nil,
			expectedError:      fmt.Errorf("getting all gists err: %v", "katego"),
		},
	}

	for _, tt := range test {
		t.Run(tt.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			newRepo := mock_repository.NewMockRepository(c)
			tt.mockBehavior(newRepo, tt.inputUser)

			var userTransport *transport.UserGrpcTransport
			usecase := NewGistService(newRepo, userTransport)

			response, err := usecase.GetAllGists(context.Background(), tt.inputUser)

			assert.Equal(t, response, tt.expectedStatusBody)
			assert.Equal(t, err, tt.expectedError)
		})
	}
}
