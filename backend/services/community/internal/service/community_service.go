package service

import (
	"context"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/williamsebastianliman/WEB-WS-242/services/community/internal/model"
	"github.com/williamsebastianliman/WEB-WS-242/services/community/internal/repository"
)

type CommunityService interface {
	ListCommunities(ctx context.Context, filter string, categoryID *string, page, size int, userId string) ([]model.Community, error)
	SearchCommunities(ctx context.Context, keyword string, page, size int) ([]model.Community, error)
	GetCommunityByID(ctx context.Context, communityID string) (*model.Community, error)
	ListJoinedCommunities(ctx context.Context, userID string, page, size int) ([]model.Community, error)
	ListRequestedCommunities(ctx context.Context, userID string, page, size int) ([]model.Community, error)
	ListModeratedCommunities(ctx context.Context, userID string) ([]model.Community, error)
	CreateCommunity(ctx context.Context, name, desc, rules, logo, banner, status string, categoryIDs []string, ownerID string) (string, error)
	RequestJoinCommunity(ctx context.Context, communityID, userID string) error
	ListCommunityMembers(ctx context.Context, communityID string, role string, page, size int, memberIds []string) ([]model.CommunityMember, error)
	PromoteMember(ctx context.Context, communityID, userID string) error
	DemoteMember(ctx context.Context, communityID, userID string) error
	AcceptJoinRequest(ctx context.Context, communityID, userID string) error
	RejectJoinRequest(ctx context.Context, communityID, userID string) error
	GetCommunityAbout(ctx context.Context, communityID string) (*model.Community, error)
	ListCategories(ctx context.Context) ([]model.CommunityCategory, error)
	CreateCategory(ctx context.Context, name string) (string, error)
	UpdateCategory(ctx context.Context, categoryID, name string) error
	DeleteCategory(ctx context.Context, categoryID string) error
	AddCommunityCategory(ctx context.Context, communityID, categoryID string) error
	RemoveCommunityCategory(ctx context.Context, communityID, categoryID string) error
	ListCommunityCategories(ctx context.Context, communityID string) ([]model.CommunityCategory, error)
	ListPendingCommunities() ([]*model.Community, error)

	AcceptCommunityCreation(ctx context.Context, communityID string) error
    RejectCommunityCreation(ctx context.Context, communityID string) error
	GetTopCommunityMembers(ctx context.Context, communityID string) ([]string, error)
	GetUserRoleInCommunity(ctx context.Context, communityID, userID string) (string, error)
	GetCommunityMemberCount(ctx context.Context, communityID string) (int64, error)
	ExploreCommunityByName(query string, threshold, page, size int) ([]model.Community, int, error)
}

type communityService struct {
	repo repository.CommunityRepository
}

func NewCommunityService(repo repository.CommunityRepository) CommunityService {
	return &communityService{repo: repo}
}
func (s *communityService) GetUserRoleInCommunity(ctx context.Context, communityID, userID string) (string, error) {
	return s.repo.GetUserRoleInCommunity(communityID, userID)
}

func (s *communityService) GetTopCommunityMembers(ctx context.Context, communityID string) ([]string, error) {
	return s.repo.GetTopMemberCommunity(communityID)
}

func (s *communityService) ListPendingCommunities() ([]*model.Community, error) {
	return s.repo.ListPendingCommunities()
}
func (s *communityService) ListCommunities(ctx context.Context, filter string, categoryID *string, page, size int, userId string) ([]model.Community, error) {
	f := repository.CommunityListFilter{}
	if categoryID != nil {
		f.CategoryID = *categoryID
	}
	f.Query = filter

	comms, err := s.repo.ListCommunities(f, page, size, userId)
	if err != nil {
		return nil, err
	}

	var res []model.Community
	for _, c := range comms {
		res = append(res, *c)
	}
	return res, nil
}

func (s *communityService) SearchCommunities(ctx context.Context, keyword string, page, size int) ([]model.Community, error) {
	comms, err := s.repo.SearchCommunities(keyword, "")
	if err != nil {
		return nil, err
	}
	start := (page - 1) * size
	end := start + size
	if start > len(comms) {
		return []model.Community{}, nil
	}
	if end > len(comms) {
		end = len(comms)
	}
	var res []model.Community
	for _, c := range comms[start:end] {
		res = append(res, *c)
	}
	return res, nil
}

func (s *communityService) GetCommunityByID(ctx context.Context, communityID string) (*model.Community, error) {
	return s.repo.GetCommunityByID(communityID)
}

func (s *communityService) ListJoinedCommunities(ctx context.Context, userID string, page, size int) ([]model.Community, error) {
	log.Printf("User: %s, page:%d, size:%d\n", userID, page, size);
	
	comms, err := s.repo.ListUserJoinedCommunities(userID, page, size)
	if err != nil {
		return nil, err
	}
	var res []model.Community
	for _, c := range comms {
		res = append(res, *c)
	}
	return res, nil
}

func (s *communityService) ListRequestedCommunities(ctx context.Context, userID string, page, size int) ([]model.Community, error) {
	comms, err := s.repo.ListUserJoinRequests(userID, page, size)
	if err != nil {
		return nil, err
	}
	var res []model.Community
	for _, c := range comms {
		res = append(res, *c)
	}
	return res, nil
}

func (s *communityService) ListModeratedCommunities(ctx context.Context, userID string) ([]model.Community, error) {
	comms, err := s.repo.ListUserModeratedCommunities(userID)
	if err != nil {
		return nil, err
	}
	var res []model.Community
	for _, c := range comms {
		res = append(res, *c)
	}
	return res, nil
}

func (s *communityService) CreateCommunity(ctx context.Context, name, desc, rules, logo, banner, status string, categoryIDs []string, ownerID string) (string, error) {
	id := uuid.New()
	community := &model.Community{
		CommunityID:          id,
		CommunityName:        name,
		CommunityDescription: desc,
		CommunityRules:       rules,
		CommunityLogo:        logo,
		CommunityBanner:      banner,
		Status:               status,
		CreatedAt:            time.Now(),
	}
	if err := s.repo.CreateCommunity(community); err != nil {
		return "", err
	}
	for _, catID := range categoryIDs {
		_ = s.repo.AddCommunityCategory(id.String(), catID)
	}
	s.repo.CreateOwner(id.String(), ownerID)
	return id.String(), nil
}

func (s *communityService) RequestJoinCommunity(ctx context.Context, communityID, userID string) error {
	log.Printf("joined id: %s",userID)
	return s.repo.RequestJoinCommunity(communityID, userID)
}

func (s *communityService) ListCommunityMembers(ctx context.Context, communityID string, role string, page, size int, memberIds []string) ([]model.CommunityMember, error) {
	members, err := s.repo.ListCommunityMembers(communityID, role, page, size, memberIds)
	if err != nil {
		return nil, err
	}
	var res []model.CommunityMember
	for _, m := range members {
		res = append(res, *m)
	}
	return res, nil
}

func (s *communityService) PromoteMember(ctx context.Context, communityID, userID string) error {
	return s.repo.PromoteToModerator(communityID, userID)
}

func (s *communityService) DemoteMember(ctx context.Context, communityID, userID string) error {
	return s.repo.DemoteFromModerator(communityID, userID)
}

func (s *communityService) AcceptJoinRequest(ctx context.Context, communityID, userID string) error {
	return s.repo.AcceptJoinRequest(communityID, userID)
}

func (s *communityService) RejectJoinRequest(ctx context.Context, communityID, userID string) error {
	return s.repo.RejectJoinRequest(communityID, userID)
}

func (s *communityService) GetCommunityAbout(ctx context.Context, communityID string) (*model.Community, error) {
	return s.repo.GetCommunityByID(communityID)
}

func (s *communityService) ListCategories(ctx context.Context) ([]model.CommunityCategory, error) {
	cats, err := s.repo.ListCategories()
	if err != nil {
		return nil, err
	}
	var res []model.CommunityCategory
	for _, c := range cats {
		res = append(res, *c)
	}
	return res, nil
}

func (s *communityService) CreateCategory(ctx context.Context, name string) (string, error) {
	id := uuid.New()
	cat := &model.CommunityCategory{
		CategoryID:   id,
		CategoryName: name,
		CreatedAt:    time.Now(),
	}
	if err := s.repo.CreateCategory(cat); err != nil {
		return "", err
	}
	return id.String(), nil
}

func (s *communityService) UpdateCategory(ctx context.Context, categoryID, name string) error {
	cid, _ := uuid.Parse(categoryID)
	cat := &model.CommunityCategory{
		CategoryID:   cid,
		CategoryName: name,
	}
	return s.repo.UpdateCategory(cat)
}

func (s *communityService) DeleteCategory(ctx context.Context, categoryID string) error {
	return s.repo.DeleteCategory(categoryID)
}

func (s *communityService) AddCommunityCategory(ctx context.Context, communityID, categoryID string) error {
	return s.repo.AddCommunityCategory(communityID, categoryID)
}

func (s *communityService) RemoveCommunityCategory(ctx context.Context, communityID, categoryID string) error {
	return s.repo.RemoveCommunityCategory(communityID, categoryID)
}

func (s *communityService) ListCommunityCategories(ctx context.Context, communityID string) ([]model.CommunityCategory, error) {
	cats, err := s.repo.ListCommunityCategories(communityID)
	if err != nil {
		return nil, err
	}
	var res []model.CommunityCategory
	for _, c := range cats {
		res = append(res, *c)
	}
	return res, nil
}

func (s *communityService) AcceptCommunityCreation(ctx context.Context, communityID string) error {
    return s.repo.AcceptCommunityCreation(communityID)
}

func (s *communityService) RejectCommunityCreation(ctx context.Context, communityID string) error {
    return s.repo.RejectCommunityCreation(communityID)
}

func (s *communityService) GetCommunityMemberCount(ctx context.Context, communityID string) (int64, error) {
	return s.repo.GetCommunityMemberCount(communityID)
}

func (s *communityService) ExploreCommunityByName(query string, threshold, page, size int) ([]model.Community, int, error){
	return s.repo.ExploreCommunityByName(query, threshold, page, size)
}