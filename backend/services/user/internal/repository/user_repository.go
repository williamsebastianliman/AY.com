package repository

import (
	"errors"
	"fmt"
	"strings"
	"unicode"

	"github.com/google/uuid"
	"github.com/williamsebastianliman/WEB-WS-242/services/user/internal/model"
	"gorm.io/gorm"
)

    type UserRepository interface {
        Create(user *model.User) (error)
        SearchUserByUsername(username string) (*model.User, error)
        SearchUserById(id string) (*model.User, error)
        ActivateUserById(id string) (*model.User, error)
        SearchUserByEmail(email string) (*model.User, error)
        Follow(userID, followerID string) error
        GetFollowers(userID string) ([]*model.User, error)
        GetFollowing(userID string) ([]*model.User, error)
        Unfollow(userID, followerID string) error
        CountFollowers(userID string) (int64, error)
        CountFollowing(userID string) (int64, error)
        IsFollowed(userID, followerID string) (bool, error)
        RequestPremium(pq *model.PremiumRequest) (error)
        IsBlocked(userID string, blocked_id string) (bool, error)
        Unblock(userID, followerID string) error
        Block(userID, followerID string) error
        IsRequestingPremium(userID string) (bool, error)
        GetAllPremiumRequests() ([]*model.PremiumRequest, error)
        AcceptPremium(userID string) error
        RejectPremium (userID string) error
        GetUsersByName(name string) ([]*model.User, error)
        UpdateProfile(user *model.User) error
        GetAllUsers() ([]*model.User, error)
        CreateUserReport(user_report *model.UserReport) (error)
        GetAllUserReports() ([]*model.UserReport, error)
        ChangePassword(userID string, newPasswordHash string) error
        ExploreUserByName(query string, threshold, page, size int) ([]model.User, int, error)
    }

    type userRepository struct {db *gorm.DB}

    func NewUserRepository(db *gorm.DB) UserRepository{
        return &userRepository{db: db}
    }

    func (r *userRepository) Create(user *model.User) (error){
        return r.db.Create(user).Error
    }

    func (r *userRepository) CreateUserReport(user_report *model.UserReport) (error){
        return r.db.Create(user_report).Error
    }

    func (r *userRepository) SearchUserByUsername(username string) (*model.User, error){
        var user model.User
        err := r.db.First(&user, model.User{Username: username}).Error
        if errors.Is(err, gorm.ErrRecordNotFound){
            return nil, nil
        }
        if err != nil{
            return nil, err
        }
        return &user, nil
    }

    func (r *userRepository) SearchUserByEmail(email string) (*model.User, error){
        var user model.User
        err := r.db.First(&user, model.User{Email: email}).Error
        if errors.Is(err, gorm.ErrRecordNotFound){
            return nil, nil
        }
        if err != nil{
            return nil, err
        }
        return &user, nil
    }

    func (r *userRepository) SearchUserById(id string) (*model.User, error){
        var user model.User
        uid , err := uuid.Parse(id)
        if err != nil{
            return nil, err
        }
        err = r.db.First(&user, model.User{ID: uid}).Error
        if errors.Is(err, gorm.ErrRecordNotFound){
            return nil, nil
        }
        if err != nil{
            return nil, err
        }
        return &user, nil
    }

    func (r *userRepository) GetUserByID(id string) (*model.User, error) {
        uid, err := uuid.Parse(id)
        if err != nil {
            return nil, fmt.Errorf("invalid uuid: %w", err)
        }

        var user model.User
        err = r.db.
            Where("id = ?", uid).
            First(&user).
            Error

        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, fmt.Errorf("user %s not found", uid)
        }
        if err != nil {
            return nil, err
        }
        return &user, nil
    }

    func (r *userRepository) ActivateUserById(id string) (*model.User, error) {
        uid, err := uuid.Parse(id)
        if err != nil {
            return nil, fmt.Errorf("invalid uuid: %w", err)
        }

        var user model.User
        err = r.db.
            Where("id = ?", uid).
            First(&user).
            Error
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, fmt.Errorf("user %s not found", uid)
        }
        
        if err != nil {
            return nil, err
        }
        user.IsActivated = true
        fmt.Printf("User Is Activated Is %t", user.IsActivated)
        r.db.Save(&user)
        return &user, nil
    }
    func (r *userRepository) Follow(userID, followerID string) error {
        uid, err := uuid.Parse(userID)
        if err != nil {
            return fmt.Errorf("invalid userID: %w", err)
        }
        fid, err := uuid.Parse(followerID)
        if err != nil {
            return fmt.Errorf("invalid followerID: %w", err)
        }

        follow := model.UserFollower{
            UserID:     uid,
            FollowerID: fid,
        }

        err = r.db.
            Where(&follow).
            FirstOrCreate(&follow).
            Error
        return err
    }

    func (r *userRepository) Block(userID, followerID string) error {
        uid, err := uuid.Parse(userID)
        if err != nil {
            return fmt.Errorf("invalid userID: %w", err)
        }
        fid, err := uuid.Parse(followerID)
        if err != nil {
            return fmt.Errorf("invalid followerID: %w", err)
        }

        follow := model.UserBlock{
            UserID:     uid,
            BlockedUserID: fid,
        }

        err = r.db.
            Where(&follow).
            FirstOrCreate(&follow).
            Error
        return err
    }

    func (r *userRepository) GetFollowers(userID string) ([]*model.User, error) {
        uid, err := uuid.Parse(userID)
        if err != nil {
            return nil, fmt.Errorf("invalid userID: %w", err)
        }
        var followers []model.User
        err = r.db.
            Table("user_followers").
            Select("users.*").
            Joins("join users on user_followers.follower_id = users.id").
            Where("user_followers.user_id = ?", uid).
            Scan(&followers).Error
        if err != nil {
            return nil, err
        }
        result := make([]*model.User, len(followers))
        for i := range followers {
            result[i] = &followers[i]
        }
        return result, nil
    }

    func (r *userRepository) GetFollowing(userID string) ([]*model.User, error) {
        uid, err := uuid.Parse(userID)
        if err != nil {
            return nil, fmt.Errorf("invalid userID: %w", err)
        }
        var following []model.User
        err = r.db.
            Table("user_followers").
            Select("users.*").
            Joins("join users on user_followers.user_id = users.id").
            Where("user_followers.follower_id = ?", uid).
            Scan(&following).Error
        if err != nil {
            return nil, err
        }
        result := make([]*model.User, len(following))
        for i := range following {
            result[i] = &following[i]
        }
        return result, nil
    }

    func (r *userRepository) Unfollow(userID, followerID string) error {
        uid, err := uuid.Parse(userID)
        if err != nil {
            return fmt.Errorf("invalid user_id: %w", err)
        }
        fid, err := uuid.Parse(followerID)
        if err != nil {
            return fmt.Errorf("invalid follower_id: %w", err)
        }

        result := r.db.Delete(&model.UserFollower{}, "user_id = ? AND follower_id = ?", uid, fid)
        return result.Error
    }

    func (r *userRepository) Unblock(userID, followerID string) error {
        uid, err := uuid.Parse(userID)
        if err != nil {
            return fmt.Errorf("invalid user_id: %w", err)
        }
        fid, err := uuid.Parse(followerID)
        if err != nil {
            return fmt.Errorf("invalid follower_id: %w", err)
        }

        result := r.db.Delete(&model.UserBlock{}, "user_id = ? AND blocked_user_id = ?", uid, fid)
        return result.Error
    }

    func (r *userRepository) CountFollowers(userID string) (int64, error) {
        uid, err := uuid.Parse(userID)
        if err != nil {
            return 0, fmt.Errorf("invalid userID: %w", err)
        }
        var count int64
        err = r.db.Model(&model.UserFollower{}).Where("user_id = ?", uid).Count(&count).Error
        if err != nil {
            return 0, err
        }
        return count, nil
    }

    func (r *userRepository) CountFollowing(userID string) (int64, error) {
        uid, err := uuid.Parse(userID)
        if err != nil {
            return 0, fmt.Errorf("invalid userID: %w", err)
        }
        var count int64
        err = r.db.Model(&model.UserFollower{}).Where("follower_id = ?", uid).Count(&count).Error
        if err != nil {
            return 0, err
        }
        return count, nil
    }

    func (r *userRepository) IsFollowed(userID string, followerID string) (bool, error) {
        uid, err := uuid.Parse(userID)
        if err != nil {
            return false, fmt.Errorf("invalid userID: %w", err)
        }
        fid, err := uuid.Parse(followerID)
        if err != nil {
            return false, fmt.Errorf("invalid followerID: %w", err)
        }

        var count int64
        err = r.db.
            Model(&model.UserFollower{}).
            Where("user_id = ? AND follower_id = ?", uid, fid).
            Count(&count).Error
        if err != nil {
            return false, err
        }
        return count > 0, nil
    }

    func (r *userRepository) IsBlocked(userID string, blocked_id string) (bool, error) {
        uid, err := uuid.Parse(userID)
        if err != nil {
            return false, fmt.Errorf("invalid userID: %w", err)
        }
        fid, err := uuid.Parse(blocked_id)
        if err != nil {
            return false, fmt.Errorf("invalid blockedID: %w", err)
        }

        var count int64
        err = r.db.
            Model(&model.UserBlock{}).
            Where("user_id = ? AND blocked_user_id = ?", uid, fid).
            Count(&count).Error
        if err != nil {
            return false, err
        }
        return count > 0, nil
    }

    func (r *userRepository) IsRequestingPremium(userID string) (bool, error) {
        uid, err := uuid.Parse(userID)
        if err != nil {
            return false, fmt.Errorf("invalid userID: %w", err)
        }

        var count int64
        err = r.db.
            Model(&model.PremiumRequest{}).
            Where("user_id = ?", uid).
            Count(&count).Error
        if err != nil {
            return false, err
        }
        return count > 0, nil
    }

    func (r *userRepository) GetAllPremiumRequests() ([]*model.PremiumRequest, error) {
        var requests []*model.PremiumRequest
        err := r.db.Find(&requests).Error
        if err != nil {
            return nil, err
        }
        return requests, nil
    }
    func (r *userRepository) AcceptPremium(userID string) error {
        uid, err := uuid.Parse(userID)
        if err != nil {
            return fmt.Errorf("invalid userID: %w", err)
        }
        var req model.PremiumRequest
        err = r.db.Where("user_id = ?", uid).First(&req).Error
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return fmt.Errorf("premium request not found for user %s", userID)
        }
        if err != nil {
            return err
        }
        if err := r.db.Delete(&req).Error; err != nil {
            return err
        }
        if err := r.db.Model(&model.User{}).
            Where("id = ?", uid).
            Update("is_premium", true).Error; err != nil {
            return err
        }
        return nil
    }

    func (r *userRepository) RejectPremium(userID string) error {
        uid, err := uuid.Parse(userID)
        if err != nil {
            return fmt.Errorf("invalid userID: %w", err)
        }
        var req model.PremiumRequest
        err = r.db.Where("user_id = ?", uid).First(&req).Error
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return fmt.Errorf("premium request not found for user %s", userID)
        }
        if err != nil {
            return err
        }
        if err := r.db.Delete(&req).Error; err != nil {
            return err
        }
        return nil
    }

    func (r *userRepository) RequestPremium(pq *model.PremiumRequest) (error){
        return r.db.Create(pq).Error
    }

    func (r *userRepository) GetUsersByName(name string) ([]*model.User, error) {
        var users []*model.User
        err := r.db.
            Where("name ILIKE ?", "%"+name+"%").
            Find(&users).
            Error
        if err != nil {
            return nil, err
        }
        return users, nil
    }

    func (r *userRepository) UpdateProfile(user *model.User) error {
        updateFields := map[string]interface{}{
            "name":         user.Name,
            "username":     user.Username,
            "email":        user.Email,
            "password_hash":     user.PasswordHash,
            "gender":       user.Gender,
            "profile_picture_id": user.ProfilePictureID,
            "date_of_birth": user.DateOfBirth,
            "banner_media_id": user.BannerMediaID,
        }

        return r.db.Model(&model.User{}).
            Where("id = ?", user.ID).
            Updates(updateFields).Error
    }

    func (r *userRepository) GetAllUsers() ([]*model.User, error) {
        var users []*model.User
        err := r.db.Find(&users).Error
        if err != nil {
            return nil, err
        }
        return users, nil
    }

    func (r *userRepository) GetAllUserReports() ([]*model.UserReport, error) {
        var reports []*model.UserReport
        if err := r.db.Find(&reports).Error; err != nil {
            return nil, err
        }
        return reports, nil
    }

    func (r *userRepository) ChangePassword(userID string, newPasswordHash string) error {
        uid, err := uuid.Parse(userID)
        if err != nil {
            return fmt.Errorf("invalid userID: %w", err)
        }
        result := r.db.Model(&model.User{}).Where("id = ?", uid).Update("password_hash", newPasswordHash)
        return result.Error
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

    func (r *userRepository) ExploreUserByName(query string, threshold, page, size int) ([]model.User, int, error) {
        var users []model.User

        err := r.db.
            Table("users").
            Select("users.*, COUNT(user_followers.follower_id) as follower_count").
            Joins("LEFT JOIN user_followers ON users.id = user_followers.user_id").
            Group("users.id").
            Order("follower_count DESC, users.name ASC").
            Find(&users).Error
        if err != nil {
            return nil, 0, err
        }

        var matched []model.User
        for _, user := range users {
            if TokenMatching(query, user.Name, threshold) {
                matched = append(matched, user)
            }
        }

        total := len(matched)

        start := (page - 1) * size
        end := start + size
        if start >= total {
            return []model.User{}, total, nil
        }
        if end > total {
            end = total
        }

        return matched[start:end], total, nil
    }
