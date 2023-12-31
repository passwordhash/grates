// Code generated by MockGen. DO NOT EDIT.
// Source: service.go

// Package mock_service is a generated GoMock package.
package mock_service

import (
	domain "grates/internal/domain"
	service "grates/internal/service"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockUser is a mock of User interface.
type MockUser struct {
	ctrl     *gomock.Controller
	recorder *MockUserMockRecorder
}

// MockUserMockRecorder is the mock recorder for MockUser.
type MockUserMockRecorder struct {
	mock *MockUser
}

// NewMockUser creates a new mock instance.
func NewMockUser(ctrl *gomock.Controller) *MockUser {
	mock := &MockUser{ctrl: ctrl}
	mock.recorder = &MockUserMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUser) EXPECT() *MockUserMockRecorder {
	return m.recorder
}

// AuthenticateUser mocks base method.
func (m *MockUser) AuthenticateUser(email, password string) (service.Tokens, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AuthenticateUser", email, password)
	ret0, _ := ret[0].(service.Tokens)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AuthenticateUser indicates an expected call of AuthenticateUser.
func (mr *MockUserMockRecorder) AuthenticateUser(email, password interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AuthenticateUser", reflect.TypeOf((*MockUser)(nil).AuthenticateUser), email, password)
}

// CreateUser mocks base method.
func (m *MockUser) CreateUser(user domain.UserSignUpInput) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateUser", user)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateUser indicates an expected call of CreateUser.
func (mr *MockUserMockRecorder) CreateUser(user interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateUser", reflect.TypeOf((*MockUser)(nil).CreateUser), user)
}

// GenerateTokens mocks base method.
func (m *MockUser) GenerateTokens(user domain.User) (service.Tokens, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GenerateTokens", user)
	ret0, _ := ret[0].(service.Tokens)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GenerateTokens indicates an expected call of GenerateTokens.
func (mr *MockUserMockRecorder) GenerateTokens(user interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GenerateTokens", reflect.TypeOf((*MockUser)(nil).GenerateTokens), user)
}

// GetAllUsers mocks base method.
func (m *MockUser) GetAllUsers() ([]domain.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllUsers")
	ret0, _ := ret[0].([]domain.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAllUsers indicates an expected call of GetAllUsers.
func (mr *MockUserMockRecorder) GetAllUsers() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllUsers", reflect.TypeOf((*MockUser)(nil).GetAllUsers))
}

// GetUserByEmail mocks base method.
func (m *MockUser) GetUserByEmail(email string) (domain.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserByEmail", email)
	ret0, _ := ret[0].(domain.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserByEmail indicates an expected call of GetUserByEmail.
func (mr *MockUserMockRecorder) GetUserByEmail(email interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserByEmail", reflect.TypeOf((*MockUser)(nil).GetUserByEmail), email)
}

// GetUserById mocks base method.
func (m *MockUser) GetUserById(id int) (domain.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserById", id)
	ret0, _ := ret[0].(domain.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserById indicates an expected call of GetUserById.
func (mr *MockUserMockRecorder) GetUserById(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserById", reflect.TypeOf((*MockUser)(nil).GetUserById), id)
}

// ParseToken mocks base method.
func (m *MockUser) ParseToken(token string) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ParseToken", token)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ParseToken indicates an expected call of ParseToken.
func (mr *MockUserMockRecorder) ParseToken(token interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ParseToken", reflect.TypeOf((*MockUser)(nil).ParseToken), token)
}

// RefreshTokens mocks base method.
func (m *MockUser) RefreshTokens(refreshToken string) (service.Tokens, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RefreshTokens", refreshToken)
	ret0, _ := ret[0].(service.Tokens)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// RefreshTokens indicates an expected call of RefreshTokens.
func (mr *MockUserMockRecorder) RefreshTokens(refreshToken interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RefreshTokens", reflect.TypeOf((*MockUser)(nil).RefreshTokens), refreshToken)
}

// UpdateProfile mocks base method.
func (m *MockUser) UpdateProfile(userId int, newProfile domain.ProfileUpdateInput) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateProfile", userId, newProfile)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateProfile indicates an expected call of UpdateProfile.
func (mr *MockUserMockRecorder) UpdateProfile(userId, newProfile interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateProfile", reflect.TypeOf((*MockUser)(nil).UpdateProfile), userId, newProfile)
}

// MockPost is a mock of Post interface.
type MockPost struct {
	ctrl     *gomock.Controller
	recorder *MockPostMockRecorder
}

// MockPostMockRecorder is the mock recorder for MockPost.
type MockPostMockRecorder struct {
	mock *MockPost
}

// NewMockPost creates a new mock instance.
func NewMockPost(ctrl *gomock.Controller) *MockPost {
	mock := &MockPost{ctrl: ctrl}
	mock.recorder = &MockPostMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockPost) EXPECT() *MockPostMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockPost) Create(post domain.Post) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", post)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *MockPostMockRecorder) Create(post interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockPost)(nil).Create), post)
}

// Delete mocks base method.
func (m *MockPost) Delete(id int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", id)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockPostMockRecorder) Delete(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockPost)(nil).Delete), id)
}

// Get mocks base method.
func (m *MockPost) Get(postId int) (domain.Post, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", postId)
	ret0, _ := ret[0].(domain.Post)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockPostMockRecorder) Get(postId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockPost)(nil).Get), postId)
}

// GetFriendsPosts mocks base method.
func (m *MockPost) GetFriendsPosts(userId, limit, offset int) ([]domain.Post, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetFriendsPosts", userId, limit, offset)
	ret0, _ := ret[0].([]domain.Post)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetFriendsPosts indicates an expected call of GetFriendsPosts.
func (mr *MockPostMockRecorder) GetFriendsPosts(userId, limit, offset interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetFriendsPosts", reflect.TypeOf((*MockPost)(nil).GetFriendsPosts), userId, limit, offset)
}

// GetUsersPosts mocks base method.
func (m *MockPost) GetUsersPosts(userId int) ([]domain.Post, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUsersPosts", userId)
	ret0, _ := ret[0].([]domain.Post)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUsersPosts indicates an expected call of GetUsersPosts.
func (mr *MockPostMockRecorder) GetUsersPosts(userId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUsersPosts", reflect.TypeOf((*MockPost)(nil).GetUsersPosts), userId)
}

// IsPostBelongsToUser mocks base method.
func (m *MockPost) IsPostBelongsToUser(userId, postId int) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsPostBelongsToUser", userId, postId)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// IsPostBelongsToUser indicates an expected call of IsPostBelongsToUser.
func (mr *MockPostMockRecorder) IsPostBelongsToUser(userId, postId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsPostBelongsToUser", reflect.TypeOf((*MockPost)(nil).IsPostBelongsToUser), userId, postId)
}

// Update mocks base method.
func (m *MockPost) Update(id int, newPost domain.PostUpdateInput) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", id, newPost)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update.
func (mr *MockPostMockRecorder) Update(id, newPost interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockPost)(nil).Update), id, newPost)
}

// MockComment is a mock of Comment interface.
type MockComment struct {
	ctrl     *gomock.Controller
	recorder *MockCommentMockRecorder
}

// MockCommentMockRecorder is the mock recorder for MockComment.
type MockCommentMockRecorder struct {
	mock *MockComment
}

// NewMockComment creates a new mock instance.
func NewMockComment(ctrl *gomock.Controller) *MockComment {
	mock := &MockComment{ctrl: ctrl}
	mock.recorder = &MockCommentMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockComment) EXPECT() *MockCommentMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockComment) Create(comment domain.CommentCreateInput) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", comment)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *MockCommentMockRecorder) Create(comment interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockComment)(nil).Create), comment)
}

// Delete mocks base method.
func (m *MockComment) Delete(userId, commentId int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", userId, commentId)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockCommentMockRecorder) Delete(userId, commentId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockComment)(nil).Delete), userId, commentId)
}

// GetPostComments mocks base method.
func (m *MockComment) GetPostComments(postId int) ([]domain.Comment, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPostComments", postId)
	ret0, _ := ret[0].([]domain.Comment)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetPostComments indicates an expected call of GetPostComments.
func (mr *MockCommentMockRecorder) GetPostComments(postId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPostComments", reflect.TypeOf((*MockComment)(nil).GetPostComments), postId)
}

// Update mocks base method.
func (m *MockComment) Update(userId, commentId int, newComment domain.CommentUpdateInput) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", userId, commentId, newComment)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update.
func (mr *MockCommentMockRecorder) Update(userId, commentId, newComment interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockComment)(nil).Update), userId, commentId, newComment)
}

// MockEmail is a mock of Email interface.
type MockEmail struct {
	ctrl     *gomock.Controller
	recorder *MockEmailMockRecorder
}

// MockEmailMockRecorder is the mock recorder for MockEmail.
type MockEmailMockRecorder struct {
	mock *MockEmail
}

// NewMockEmail creates a new mock instance.
func NewMockEmail(ctrl *gomock.Controller) *MockEmail {
	mock := &MockEmail{ctrl: ctrl}
	mock.recorder = &MockEmailMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockEmail) EXPECT() *MockEmailMockRecorder {
	return m.recorder
}

// ConfirmEmail mocks base method.
func (m *MockEmail) ConfirmEmail(hash string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ConfirmEmail", hash)
	ret0, _ := ret[0].(error)
	return ret0
}

// ConfirmEmail indicates an expected call of ConfirmEmail.
func (mr *MockEmailMockRecorder) ConfirmEmail(hash interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ConfirmEmail", reflect.TypeOf((*MockEmail)(nil).ConfirmEmail), hash)
}

// ReplaceConfirmationEmail mocks base method.
func (m *MockEmail) ReplaceConfirmationEmail(userId int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ReplaceConfirmationEmail", userId)
	ret0, _ := ret[0].(error)
	return ret0
}

// ReplaceConfirmationEmail indicates an expected call of ReplaceConfirmationEmail.
func (mr *MockEmailMockRecorder) ReplaceConfirmationEmail(userId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReplaceConfirmationEmail", reflect.TypeOf((*MockEmail)(nil).ReplaceConfirmationEmail), userId)
}

// SendAuthEmail mocks base method.
func (m *MockEmail) SendAuthEmail(to, name, hash string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendAuthEmail", to, name, hash)
	ret0, _ := ret[0].(error)
	return ret0
}

// SendAuthEmail indicates an expected call of SendAuthEmail.
func (mr *MockEmailMockRecorder) SendAuthEmail(to, name, hash interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendAuthEmail", reflect.TypeOf((*MockEmail)(nil).SendAuthEmail), to, name, hash)
}

// MockLike is a mock of Like interface.
type MockLike struct {
	ctrl     *gomock.Controller
	recorder *MockLikeMockRecorder
}

// MockLikeMockRecorder is the mock recorder for MockLike.
type MockLikeMockRecorder struct {
	mock *MockLike
}

// NewMockLike creates a new mock instance.
func NewMockLike(ctrl *gomock.Controller) *MockLike {
	mock := &MockLike{ctrl: ctrl}
	mock.recorder = &MockLikeMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockLike) EXPECT() *MockLikeMockRecorder {
	return m.recorder
}

// LikePost mocks base method.
func (m *MockLike) LikePost(userId, postId int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LikePost", userId, postId)
	ret0, _ := ret[0].(error)
	return ret0
}

// LikePost indicates an expected call of LikePost.
func (mr *MockLikeMockRecorder) LikePost(userId, postId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LikePost", reflect.TypeOf((*MockLike)(nil).LikePost), userId, postId)
}

// UnlikePost mocks base method.
func (m *MockLike) UnlikePost(userId, postId int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UnlikePost", userId, postId)
	ret0, _ := ret[0].(error)
	return ret0
}

// UnlikePost indicates an expected call of UnlikePost.
func (mr *MockLikeMockRecorder) UnlikePost(userId, postId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UnlikePost", reflect.TypeOf((*MockLike)(nil).UnlikePost), userId, postId)
}

// MockFriend is a mock of Friend interface.
type MockFriend struct {
	ctrl     *gomock.Controller
	recorder *MockFriendMockRecorder
}

// MockFriendMockRecorder is the mock recorder for MockFriend.
type MockFriendMockRecorder struct {
	mock *MockFriend
}

// NewMockFriend creates a new mock instance.
func NewMockFriend(ctrl *gomock.Controller) *MockFriend {
	mock := &MockFriend{ctrl: ctrl}
	mock.recorder = &MockFriendMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockFriend) EXPECT() *MockFriendMockRecorder {
	return m.recorder
}

// AcceptFriendRequest mocks base method.
func (m *MockFriend) AcceptFriendRequest(fromId, toId int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AcceptFriendRequest", fromId, toId)
	ret0, _ := ret[0].(error)
	return ret0
}

// AcceptFriendRequest indicates an expected call of AcceptFriendRequest.
func (mr *MockFriendMockRecorder) AcceptFriendRequest(fromId, toId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AcceptFriendRequest", reflect.TypeOf((*MockFriend)(nil).AcceptFriendRequest), fromId, toId)
}

// FriendRequests mocks base method.
func (m *MockFriend) FriendRequests(userId int) ([]domain.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FriendRequests", userId)
	ret0, _ := ret[0].([]domain.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FriendRequests indicates an expected call of FriendRequests.
func (mr *MockFriendMockRecorder) FriendRequests(userId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FriendRequests", reflect.TypeOf((*MockFriend)(nil).FriendRequests), userId)
}

// GetFriends mocks base method.
func (m *MockFriend) GetFriends(userId int) ([]domain.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetFriends", userId)
	ret0, _ := ret[0].([]domain.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetFriends indicates an expected call of GetFriends.
func (mr *MockFriendMockRecorder) GetFriends(userId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetFriends", reflect.TypeOf((*MockFriend)(nil).GetFriends), userId)
}

// SendFriendRequest mocks base method.
func (m *MockFriend) SendFriendRequest(fromId, toId int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendFriendRequest", fromId, toId)
	ret0, _ := ret[0].(error)
	return ret0
}

// SendFriendRequest indicates an expected call of SendFriendRequest.
func (mr *MockFriendMockRecorder) SendFriendRequest(fromId, toId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendFriendRequest", reflect.TypeOf((*MockFriend)(nil).SendFriendRequest), fromId, toId)
}

// Unfriend mocks base method.
func (m *MockFriend) Unfriend(userId, friendId int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Unfriend", userId, friendId)
	ret0, _ := ret[0].(error)
	return ret0
}

// Unfriend indicates an expected call of Unfriend.
func (mr *MockFriendMockRecorder) Unfriend(userId, friendId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Unfriend", reflect.TypeOf((*MockFriend)(nil).Unfriend), userId, friendId)
}
