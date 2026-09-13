package repository

import (
	"errors"
	"log"
	"strings"
	"unicode"

	"github.com/google/uuid"
	"github.com/williamsebastianliman/WEB-WS-242/services/community/internal/model"
	"gorm.io/gorm"
)

type communityRepository struct {
	db *gorm.DB
}
type CommunityListFilter struct {
	CategoryID string
	Query      string
}

type CommunityRepository interface {
	CreateCommunity(comm *model.Community) error

	GetCommunityByID(id string) (*model.Community, error)

	ListCommunities(filter CommunityListFilter, page, size int, user_id string) ([]*model.Community, error)

	SearchCommunities(query, category string) ([]*model.Community, error)

	ListCategories() ([]*model.CommunityCategory, error)

	CreateCategory(cat *model.CommunityCategory) error

	UpdateCategory(cat *model.CommunityCategory) error

	DeleteCategory(categoryID string) error

	ListUserJoinedCommunities(userID string, page, size int) ([]*model.Community, error)

	ListUserJoinRequests(userID string, page, size int) ([]*model.Community, error)

	ListUserModeratedCommunities(userID string) ([]*model.Community, error)

	RequestJoinCommunity(communityID, userID string) error

	ListCommunityMembers(communityID, role string, page, size int, memberIds []string) ([]*model.CommunityMember, error)

	PromoteToModerator(communityID, memberID string) error

	DemoteFromModerator(communityID, memberID string) error

	AcceptJoinRequest(communityID, memberID string) error

	RejectJoinRequest(communityID, memberID string) error

	AddCommunityCategory(communityID, categoryID string) error

	RemoveCommunityCategory(communityID, categoryID string) error

	ListCommunityCategories(communityID string) ([]*model.CommunityCategory, error)

	CreateOwner(communityID, ownerID string) error

	ListPendingCommunities() ([]*model.Community, error)

	AcceptCommunityCreation(communityID string) error

	RejectCommunityCreation(communityID string) error
	
	GetTopMemberCommunity(communityID string) ([]string, error)

	GetUserRoleInCommunity(communityID string, userID string) (string, error)

	GetCommunityMemberCount(communityID string) (int64, error)
	
	ExploreCommunityByName(query string, threshold, page, size int) ([]model.Community, int, error)
}
func NewCommunityRepository(db *gorm.DB) CommunityRepository {
	return &communityRepository{db: db}
}

func (r *communityRepository) ListPendingCommunities() ([]*model.Community, error) {
	var comms []*model.Community
	query := r.db.Model(&model.Community{}).Where("status = ?", "pending")

	err := query.Find(&comms).Error
	return comms, err
}

func (r *communityRepository) CreateCommunity(c *model.Community) error {
	return r.db.Create(c).Error
}

func (r *communityRepository) GetCommunityByID(id string) (*model.Community, error) {
	var comm model.Community
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}
	err = r.db.First(&comm, "community_id = ?", uid).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &comm, err
}

func (r *communityRepository) ListCommunities(filter CommunityListFilter, page, size int, user_id string) ([]*model.Community, error) {
	var comms []*model.Community
	query := r.db.Model(&model.Community{})

	if user_id != "" {
		query = query.Where(`
			NOT EXISTS (
				SELECT 1 FROM community_members 
				WHERE community_members.community_id = communities.community_id 
				AND community_members.member_id = ?
			)
		`, user_id)
	}

	if filter.CategoryID != "" {
		query = query.Joins("JOIN community_details ON community_details.community_id = communities.community_id").
			Where("community_details.category_id = ?", filter.CategoryID)
	}

	if filter.Query != "" {
		query = query.Where("community_name ILIKE ?", "%"+filter.Query+"%")
	}

	if page > 0 && size > 0 {
		offset := (page - 1) * size
		query = query.Offset(offset).Limit(size)
	}

	err := query.Find(&comms).Error
	return comms, err
}

func (r *communityRepository) SearchCommunities(query, category string) ([]*model.Community, error) {
	var comms []*model.Community
	db := r.db.Model(&model.Community{})
	if query != "" {
		db = db.Where("community_name ILIKE ?", "%"+query+"%")
	}
	if category != "" {
		db = db.Joins("JOIN community_details ON community_details.community_id = communities.community_id").
			Where("community_details.category_id = ?", category)
	}
	err := db.Find(&comms).Error
	return comms, err
}

func (r *communityRepository) ListUserJoinedCommunities(userID string, page, size int) ([]*model.Community, error) {
	var comms []*model.Community
	query := r.db.Joins("JOIN community_members ON community_members.community_id = communities.community_id").
		Where(
			"community_members.member_id = ? AND community_members.status != ? AND communities.status = ?",
			userID, "pending", "active",
		)

	if page > 0 && size > 0 {
		offset := (page - 1) * size
		query = query.Offset(offset).Limit(size)
	}

	err := query.Find(&comms).Error
	return comms, err
}

func (r *communityRepository) ListUserJoinRequests(userID string, page, size int) ([]*model.Community, error) {
	var comms []*model.Community
	query := r.db.Joins("JOIN community_members ON community_members.community_id = communities.community_id").
		Where("community_members.member_id = ? AND community_members.status = ?", userID, "pending")

	if page > 0 && size > 0 {
		offset := (page - 1) * size
		query = query.Offset(offset).Limit(size)
	}

	err := query.Find(&comms).Error
	return comms, err
}

func (r *communityRepository) RequestJoinCommunity(communityID, userID string) error {
	cid, err := uuid.Parse(communityID)
	if err != nil {
		return err
	}
	log.Printf("repo uid: %s", userID)
	uid, err := uuid.Parse(userID)
	if err != nil {
		return err
	}
	log.Printf("repo uid after: %s", uid)
	member := model.CommunityMember{
		CommunityID: cid,
		MemberID:    uid,
		Status:      "pending",
	}
	return r.db.Create(&member).Error
}

func (r* communityRepository) CreateOwner(communityID, ownerID string) error {
	cid, err := uuid.Parse(communityID)
	if err !=nil{
		return err
	}
	oid, err := uuid.Parse(ownerID)
	if err != nil{
		return err
	}
	owner := model.CommunityMember{
		CommunityID: cid,
		MemberID:    oid,
		Status:      "owner",
	}
	return r.db.Create(&owner).Error
}

func (r *communityRepository) ListCommunityMembers(communityID, role string, page, size int, memberIds []string) ([]*model.CommunityMember, error) {
	var members []*model.CommunityMember
	offset := (page - 1) * size

	q := r.db.Where("community_id = ?", communityID)
	if role != "" {
		q = q.Where("status = ?", role)
	}
	q = q.Where("member_id IN (?)", memberIds)
	err := q.Offset(offset).Limit(size).Find(&members).Error
	return members, err
}

func (r *communityRepository) PromoteToModerator(communityID, memberID string) error {
	return r.db.Model(&model.CommunityMember{}).
		Where("community_id = ? AND member_id = ?", communityID, memberID).
		Update("status", "moderator").Error
}

func (r *communityRepository) DemoteFromModerator(communityID, memberID string) error {
	return r.db.Model(&model.CommunityMember{}).
		Where("community_id = ? AND member_id = ?", communityID, memberID).
		Update("status", "accepted").Error
}

func (r *communityRepository) AcceptJoinRequest(communityID, memberID string) error {
	return r.db.Model(&model.CommunityMember{}).
		Where("community_id = ? AND member_id = ? AND status = ?", communityID, memberID, "pending").
		Update("status", "accepted").Error
}

func (r *communityRepository) RejectJoinRequest(communityID, memberID string) error {
	return r.db.Where("community_id = ? AND member_id = ? AND status = ?", communityID, memberID, "pending").
		Delete(&model.CommunityMember{}).Error
}

func (r *communityRepository) ListUserModeratedCommunities(userID string) ([]*model.Community, error) {
	var comms []*model.Community
	err := r.db.Joins("JOIN community_members ON community_members.community_id = communities.community_id").
		Where("community_members.member_id = ? AND community_members.status = ?", userID, "moderator").
		Find(&comms).Error
	return comms, err
}

func (r *communityRepository) ListCategories() ([]*model.CommunityCategory, error) {
	var cats []*model.CommunityCategory
	err := r.db.Find(&cats).Error
	return cats, err
}

func (r *communityRepository) CreateCategory(cat *model.CommunityCategory) error {
	return r.db.Create(cat).Error
}

func (r *communityRepository) UpdateCategory(cat *model.CommunityCategory) error {
	return r.db.Model(cat).Where("category_id = ?", cat.CategoryID).Updates(cat).Error
}

func (r *communityRepository) DeleteCategory(categoryID string) error {
	return r.db.Delete(&model.CommunityCategory{}, "category_id = ?", categoryID).Error
}

func (r *communityRepository) AddCommunityCategory(communityID, categoryID string) error {
	cid, err := uuid.Parse(communityID)
	if err != nil {
		return err
	}
	catid, err := uuid.Parse(categoryID)
	if err != nil {
		return err
	}
	return r.db.Create(&model.CommunityDetail{
		CommunityID: cid,
		CategoryID:  catid,
	}).Error
}

func (r *communityRepository) ListCommunityCategories(communityID string) ([]*model.CommunityCategory, error) {
	cid, err := uuid.Parse(communityID)
	if err != nil {
		return nil, err
	}
	var categories []*model.CommunityCategory
	err = r.db.Table("community_details").
		Select("community_categories.*").
		Joins("join community_categories on community_details.category_id = community_categories.category_id").
		Where("community_details.community_id = ?", cid).
		Scan(&categories).Error
	if err != nil {
		return nil, err
	}
	return categories, nil
}

func (r *communityRepository) RemoveCommunityCategory(communityID, categoryID string) error {
	cid, err := uuid.Parse(communityID)
	if err != nil {
		return err
	}
	catid, err := uuid.Parse(categoryID)
	if err != nil {
		return err
	}
	return r.db.Delete(&model.CommunityDetail{}, "community_id = ? AND category_id = ?", cid, catid).Error
}

func (r *communityRepository) AcceptCommunityCreation(communityID string) error {
    uid, err := uuid.Parse(communityID)
    if err != nil {
        return err
    }
    result := r.db.Model(&model.Community{}).
        Where("community_id = ? AND status = ?", uid, "pending").
        Update("status", "active")
    if result.Error != nil {
        return result.Error
    }
    if result.RowsAffected == 0 {
        return gorm.ErrRecordNotFound
    }
    return nil
}

func (r *communityRepository) RejectCommunityCreation(communityID string) error {
    uid, err := uuid.Parse(communityID)
    if err != nil {
        return err
    }
    result := r.db.Where("community_id = ? AND status = ?", uid, "pending").
        Delete(&model.Community{})
    if result.Error != nil {
        return result.Error
    }
    if result.RowsAffected == 0 {
        return gorm.ErrRecordNotFound
    }
    return nil
}

func (r *communityRepository) GetTopMemberCommunity(communityID string) ([]string, error) {
	uid, err := uuid.Parse(communityID)
	if err != nil {
		return nil, err
	}
	var members []model.CommunityMember
	err = r.db.
		Where("community_id = ? AND status IN ?", uid, []string{"accepted", "moderator", "owner"}).
		Order("created_at ASC").
		Find(&members).Error
	if err != nil {
		return nil, err
	}
	topIDs := make([]string, len(members))
	for i, m := range members {
		topIDs[i] = m.MemberID.String()
	}
	return topIDs, nil
}

func (r *communityRepository) GetUserRoleInCommunity(communityID string, userID string) (string, error) {
	var member model.CommunityMember
	err := r.db.
		Where("community_id = ? AND member_id = ?", communityID, userID).
		First(&member).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return "", nil
		}
		return "", err
	}
	return member.Status, nil
}

func (r *communityRepository) GetCommunityMemberCount(communityID string) (int64, error) {
	var count int64
	err := r.db.Model(&model.CommunityMember{}).
		Where("community_id = ? AND status != ?", communityID, "pending").
		Count(&count).Error
	return count, err
}

func DamerauLevenshtein(a, b string) int {
	n := len(a)
	m := len(b)

	dist := make([][]int, n+1)
	for i := range dist {
		dist[i] = make([]int, m+1)
	}

	for i := 0; i <= n; i++ {
		dist[i][0] = i
	}
	for j := 0; j <= m; j++ {
		dist[0][j] = j
	}

	for i := 1; i <= n; i++ {
		for j := 1; j <= m; j++ {
			cost := 0
			if a[i-1] != b[j-1] {
				cost = 1
			}
			dist[i][j] = min(
				dist[i-1][j]+1,
				dist[i][j-1]+1,
				dist[i-1][j-1]+cost,
			)
			if i > 1 && j > 1 && a[i-1] == b[j-2] && a[i-2] == b[j-1] {
				dist[i][j] = min2(dist[i][j], dist[i-2][j-2]+1)
			}
		}
	}
	return dist[n][m]
}

func min2(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func min(a, b, c int) int {
	if a < b {
		if a < c {
			return a
		}
		return c
	}
	if b < c {
		return b
	}
	return c
}

func TokenMatching(query, content string, threshold int) bool {
	queryTokens := tokenize(query)
	contentTokens := tokenize(content)

	for _, q := range queryTokens {
		found := false
		for _, c := range contentTokens {
			if DamerauLevenshtein(q, c) < threshold {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}

func tokenize(s string) []string {
	cleaned := strings.Map(func(r rune) rune {
		if unicode.IsLetter(r) || unicode.IsDigit(r) || unicode.IsSpace(r) {
			return r
		}
		return ' '
	}, strings.ToLower(s))
	tokens := strings.Fields(cleaned)
	return tokens
}

func (r *communityRepository) ExploreCommunityByName(query string, threshold, page, size int) ([]model.Community, int, error) {
	var communities []model.Community
	err := r.db.Order("created_at DESC").Find(&communities).Error
	if err != nil {
		return nil, 0, err
	}

	var matched []model.Community
	for _, community := range communities {
		if TokenMatching(query, community.CommunityName, threshold) {
			matched = append(matched, community)
		}
	}

	total := len(matched)

	start := (page - 1) * size
	end := start + size
	if start >= total {
		return []model.Community{}, total, nil
	}
	if end > total {
		end = total
	}

	return matched[start:end], total, nil
}