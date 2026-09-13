package service

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"regexp"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/williamsebastianliman/WEB-WS-242/services/user/internal/model"
	"github.com/williamsebastianliman/WEB-WS-242/services/user/internal/repository"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type UserService interface {
    Register(
        ctx context.Context,
        name, username, email, password, gender,
        birthYear, birthMonth, birthDay string,
        subscribedNews bool,
        profilePictureId, bannerMediaId *uuid.UUID,
    ) (uuid.UUID, error)

    IsUsernameUnique(ctx context.Context, username string) (bool, error)
    GetUserByUsername(ctx context.Context, username string) (*model.User, error)
    GetUserByEmail(ctx context.Context, username string) (*model.User, error)
    GetMe(ctx context.Context) (*model.User, error)
    ActivateAccount(ctx context.Context, user_id string) (*model.User, error)
    GetUserByID(ctx context.Context, id string) (*model.User, error)
    Follow(ctx context.Context, userID, followerID string) error
    GetFollowers(ctx context.Context, userID string) ([]*model.User, error)
    GetFollowing(ctx context.Context, userID string) ([]*model.User, error)
    UnfollowUser(ctx context.Context, userID, followerID string) error
    CountFollowers(ctx context.Context, userID string) (int64, error)
    CountFollowing(ctx context.Context, userID string) (int64, error)
    IsFollowed(ctx context.Context, userID, followerID string) (bool, error)
    RequestPremium(ctx context.Context, user_id string, card_number string, reason string, face_image string) error
    IsPremium(ctx context.Context, userID string) (bool, error)
    GetAllPremiumRequests(ctx context.Context) ([]*model.PremiumRequest, error)
    AcceptPremium(ctx context.Context, userID string) (error)
    RejectPremium(ctx context.Context, userID string) (error)
    GetUserByName(ctx context.Context, name string) ([]*model.User, error)
    CreateUserReport(reportedUserID, reason string) (string, error)
    ChangePassword(ctx context.Context, userID string, newPassword string) error
    UpdateProfile(
        ctx context.Context,
        userId string,
        name string,
        username string,
        email string,
        password string,
        gender string,
        birthYear string,
        birthMonth string,
        birthDay string,
        profilePictureId *uuid.UUID,
        bannerMediaId *uuid.UUID,
    ) error
    BlockUser(ctx context.Context, userID, targetID string) error
    UnblockUser(ctx context.Context, userID, targetID string) error
    IsUserBlocked(ctx context.Context, userID, targetID string) (bool, error)
    GetAllUsers(ctx context.Context) ([]*model.User, error)
    GetAllUserReports() ([]*model.UserReport, error)
    ExploreUserByName(query string, threshold, page, size int) ([]model.User, int, error)
}
    
type userService struct {
    repo repository.UserRepository
}
var secretKey = []byte("myverystrongpasswordo32bitlength")

func Encrypt(plaintext string) (string, error) {
	block, err := aes.NewCipher(secretKey)
	if err != nil {
		return "", err
	}
	b := []byte(plaintext)
	ciphertext := make([]byte, aes.BlockSize+len(b))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}
	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], b)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func Decrypt(cryptoText string) (string, error) {
	ciphertext, err := base64.StdEncoding.DecodeString(cryptoText)
	if err != nil {
		return "", err
	}
	block, err := aes.NewCipher(secretKey)
	if err != nil {
		return "", err
	}
	if len(ciphertext) < aes.BlockSize {
		return "", errors.New("ciphertext too short")
	}
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]
	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(ciphertext, ciphertext)
	return string(ciphertext), nil
}

func NewUserService(repo repository.UserRepository) UserService {
    return &userService{repo: repo}
}

func (s *userService) ActivateAccount(ctx context.Context, user_id string) (*model.User, error){
    user, err := s.repo.ActivateUserById(user_id)
    if err != nil{
        return nil,err
    }
    return user, nil;
}

func (s *userService) Register(
    ctx context.Context,
    name, username, email, password, gender,
    birthYear, birthMonth, birthDay string,
    subscribedNews bool,
    profilePictureId, bannerMediaId *uuid.UUID,
) (uuid.UUID, error) {
    hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    if err != nil {
        return uuid.Nil, err
    }
    y, err := strconv.Atoi(birthYear)
    if err != nil {
        return uuid.Nil, err
    }
    m, err := strconv.Atoi(birthMonth)
    if err != nil {
        return uuid.Nil, err
    }
    d, err := strconv.Atoi(birthDay)
    if err != nil {
        return uuid.Nil, err
    }
    role := "User"

    dob := time.Date(y, time.Month(m), d, 0, 0, 0, 0, time.UTC)

    //Validation Section
    if len(name) < 5 || !regexp.MustCompile(`^[A-Za-z\s]+$`).MatchString(name) {
        return uuid.Nil, fmt.Errorf("name must be at least 5 letters and contain only alphabetic characters")
    }

    if existing, _ := s.repo.SearchUserByUsername(username); existing != nil {
        return uuid.Nil, fmt.Errorf("username %q is already taken", username)
    }

    emailRe := regexp.MustCompile(`^[A-Za-z0-9._%+\-]+@[A-Za-z0-9.\-]+\.com$`)
    if !emailRe.MatchString(email) {
        return uuid.Nil, fmt.Errorf("email must be of the form example@domain.com")
    }

    if existing, _ := s.repo.SearchUserByEmail(email); existing != nil {
        return uuid.Nil, fmt.Errorf("email %q is already registered", email)
    }

    passwordRules := []struct {
        pattern string
        message string
    }{
        {`.{8,}`, "password must be at least 8 characters long"},
        {`[A-Z]`, "password must contain at least one uppercase letter"},
        {`[a-z]`, "password must contain at least one lowercase letter"},
        {`\d`,    "password must contain at least one digit"},
        {`\W`,    "password must contain at least one special character"},
    }
    for _, rule := range passwordRules {
        if !regexp.MustCompile(rule.pattern).MatchString(password) {
            return uuid.Nil, fmt.Errorf(rule.message)
        }
    }

    if gender != "Male" && gender != "Female" {
        return uuid.Nil, fmt.Errorf("gender must be either \"male\" or \"female\"")
    }

    if time.Since(dob) < 13*365*24*time.Hour {
        return uuid.Nil, fmt.Errorf("you must be at least 13 years old")
    }

    u := &model.User{
        Name:             name,
        Username:         username,
        Email:            email,
        PasswordHash:     string(hash),
        Gender:           gender,
        DateOfBirth:      dob,
        SubscribedNews:   subscribedNews,
        ProfilePictureID: profilePictureId,
        BannerMediaID:    bannerMediaId,
        IsBanned:         false,
        IsActivated:      false,
        Role: role,
    }
    
    if err := s.repo.Create(u); err != nil {
        return uuid.Nil, err
    }
    return u.ID, nil
}

func (s *userService) IsUsernameUnique(ctx context.Context, username string) (bool, error) {
    user, err := s.repo.SearchUserByUsername(username)
    if err != nil {
        return false, err
    }
    return user == nil, nil
}

func (s *userService) GetUserByUsername(ctx context.Context, username string) (*model.User, error) {
    return s.repo.SearchUserByUsername(username)
}

func (s *userService) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
    return s.repo.SearchUserByEmail(email)
}

func (s *userService) GetMe(ctx context.Context) (*model.User, error) {
    md, ok := metadata.FromIncomingContext(ctx)
    if !ok || len(md["user_id"]) == 0 {
        return nil, status.Error(codes.Unauthenticated, "missing user_id in metadata")
    }
    userID := md["user_id"][0]
    user, err := s.repo.SearchUserById(userID)
    if err != nil {
        return nil, err
    }
    if user == nil {
        return nil, status.Error(codes.NotFound, "user not found")
    }
    return user, nil
}

func (s *userService) GetUserByID(ctx context.Context, id string) (*model.User, error) {
    user, err := s.repo.SearchUserById(id)
    if err != nil {
        return nil, err
    }
    if user == nil {
        return nil, status.Errorf(codes.NotFound, "user %s not found", id)
    }
    return user, nil
}

func (s *userService) Follow(ctx context.Context, userID, followerID string) error {
    if userID == followerID {
        return fmt.Errorf("you cannot follow yourself")
    }
    return s.repo.Follow(userID, followerID)
}

func (s *userService) GetFollowers(ctx context.Context, userID string) ([]*model.User, error) {
    return s.repo.GetFollowers(userID)
}

func (s *userService) GetFollowing(ctx context.Context, userID string) ([]*model.User, error) {
    return s.repo.GetFollowing(userID)
}

func (s *userService) UnfollowUser(ctx context.Context, userID, followerID string) error {
    err := s.repo.Unfollow(userID, followerID)
    if err != nil {
        return fmt.Errorf("could not unfollow: %w", err)
    }
    return nil
}

func (s *userService) CountFollowers(ctx context.Context, userID string) (int64, error) {
    return s.repo.CountFollowers(userID)
}

func (s *userService) CountFollowing(ctx context.Context, userID string) (int64, error) {
    return s.repo.CountFollowing(userID)
}

func (s *userService) IsFollowed(ctx context.Context, userID, followerID string) (bool, error) {
    return s.repo.IsFollowed(userID, followerID)
}

func (s *userService) BlockUser(ctx context.Context, userID, targetID string) error {
	if userID == targetID {
		return fmt.Errorf("you cannot block yourself")
	}
	err := s.repo.Block(userID, targetID)
	if err != nil {
		return fmt.Errorf("failed to block user: %w", err)
	}
	return nil
}

func (s *userService) UnblockUser(ctx context.Context, userID, targetID string) error {
	if userID == targetID {
		return fmt.Errorf("you cannot unblock yourself")
	}
	err := s.repo.Unblock(userID, targetID)
	if err != nil {
		return fmt.Errorf("failed to unblock user: %w", err)
	}
	return nil
}

func (s *userService) IsUserBlocked(ctx context.Context, userID, targetID string) (bool, error) {
	if userID == targetID {
		return false, nil
	}
	return s.repo.IsBlocked(userID, targetID)
}

func (s *userService) IsPremium(ctx context.Context, userID string) (bool, error) {
	return s.repo.IsRequestingPremium(userID)
}



func (s *userService) RequestPremium(ctx context.Context, user_id string, card_number string, reason string, face_image string) error {
    encryptedCard, err := Encrypt(card_number)
    if err != nil {
        return err
    }
    uid, err := uuid.Parse(user_id)
    if err != nil {
        return err
    }
    uid2, err := uuid.Parse(face_image)
    if err != nil {
        return err
    }
    pq := &model.PremiumRequest{
        UserID:     uid,
        CardNumber: encryptedCard,
        Reason:     reason,
        FaceImage:  &uid2,
    }
    return s.repo.RequestPremium(pq)
}

func (s *userService) GetAllPremiumRequests(ctx context.Context) ([]*model.PremiumRequest, error) {
    requests, err := s.repo.GetAllPremiumRequests()
    if err != nil {
        return nil, err
    }
    for _, req := range requests {
        decrypted, decErr := Decrypt(req.CardNumber)
        if decErr != nil {
            req.CardNumber = "[decrypt failed]"
        } else {
            req.CardNumber = decrypted
        }
    }
    return requests, nil
}

func (s *userService) AcceptPremium(ctx context.Context, userID string) (error) {
    return s.repo.AcceptPremium(userID)
}

func (s *userService) RejectPremium(ctx context.Context, userID string) (error) {
    return s.repo.RejectPremium(userID)
}

func (s* userService) GetUserByName(ctx context.Context, name string) ([]*model.User, error) {
    return s.repo.GetUsersByName(name)
}

func (s *userService) UpdateProfile(
    ctx context.Context,
    userId string,
    name string,
    username string,
    email string,
    password string,
    gender string,
    birthYear string,
    birthMonth string,
    birthDay string,
    profilePictureId *uuid.UUID,
    bannerMediaId *uuid.UUID,
) error {
    if len(name) < 5 || !regexp.MustCompile(`^[A-Za-z\s]+$`).MatchString(name) {
        return fmt.Errorf("name must be at least 5 letters and contain only alphabetic characters")
    }

    user, err := s.repo.SearchUserById(userId)
    if err != nil || user == nil {
        return fmt.Errorf("user not found")
    }
    if user.Username != username {
        if existing, _ := s.repo.SearchUserByUsername(username); existing != nil {
            return fmt.Errorf("username %q is already taken", username)
        }
    }

    emailRe := regexp.MustCompile(`^[A-Za-z0-9._%+\-]+@[A-Za-z0-9.\-]+\.com$`)
    if !emailRe.MatchString(email) {
        return fmt.Errorf("email must be of the form example@domain.com")
    }
    if user.Email != email {
        if existing, _ := s.repo.SearchUserByEmail(email); existing != nil {
            return fmt.Errorf("email %q is already registered", email)
        }
    }

    var hash string
    if password != "" {
        passwordRules := []struct {
            pattern string
            message string
        }{
            {`.{8,}`, "password must be at least 8 characters long"},
            {`[A-Z]`, "password must contain at least one uppercase letter"},
            {`[a-z]`, "password must contain at least one lowercase letter"},
            {`\d`,    "password must contain at least one digit"},
            {`\W`,    "password must contain at least one special character"},
        }
        for _, rule := range passwordRules {
            if !regexp.MustCompile(rule.pattern).MatchString(password) {
                return fmt.Errorf(rule.message)
            }
        }
        hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
        if err != nil {
            return fmt.Errorf("failed to hash password: %w", err)
        }
        hash = string(hashed)
    }

    if gender != "Male" && gender != "Female" {
        return fmt.Errorf("gender must be either \"male\" or \"female\"")
    }

    y, err := strconv.Atoi(birthYear)
    if err != nil {
        return fmt.Errorf("invalid birth year")
    }
    m, err := strconv.Atoi(birthMonth)
    if err != nil {
        return fmt.Errorf("invalid birth month")
    }
    d, err := strconv.Atoi(birthDay)
    if err != nil {
        return fmt.Errorf("invalid birth day")
    }
    dob := time.Date(y, time.Month(m), d, 0, 0, 0, 0, time.UTC)
    if time.Since(dob) < 13*365*24*time.Hour {
        return fmt.Errorf("you must be at least 13 years old")
    }

    user.Name = name
    user.Username = username
    user.Email = email
    user.Gender = gender
    user.DateOfBirth = dob
    if profilePictureId != nil {
        user.ProfilePictureID = profilePictureId
    }
    if bannerMediaId != nil {
        user.BannerMediaID = bannerMediaId
    }
    if hash != "" {
        user.PasswordHash = hash
    }

    if err := s.repo.UpdateProfile(user); err != nil {
        return err
    }
    return nil
}

func (s *userService) GetAllUsers(ctx context.Context) ([]*model.User, error) {
    return s.repo.GetAllUsers()
}

func (s *userService) CreateUserReport(reportedUserID, reason string) (string, error) {
	report := &model.UserReport{
		ReportID:     uuid.New(),
		ReportedUser: uuid.MustParse(reportedUserID),
		Reason:       reason,
	}
	err := s.repo.CreateUserReport(report)
	if err != nil {
		return "", err
	}
	return report.ReportID.String(), nil
}

func (s *userService) GetAllUserReports() ([]*model.UserReport, error) {
    return s.repo.GetAllUserReports()
}

func (s *userService) ChangePassword(ctx context.Context, userID string, newPassword string) error {
    rules := []struct {
        pattern string
        message string
    }{
        {`.{8,}`, "password must be at least 8 characters long"},
        {`[A-Z]`, "password must contain at least one uppercase letter"},
        {`[a-z]`, "password must contain at least one lowercase letter"},
        {`\d`,    "password must contain at least one digit"},
        {`\W`,    "password must contain at least one special character"},
    }
    for _, rule := range rules {
        if !regexp.MustCompile(rule.pattern).MatchString(newPassword) {
            return fmt.Errorf(rule.message)
        }
    }

    user, err := s.repo.SearchUserById(userID)
    if err != nil || user == nil {
        return fmt.Errorf("user not found")
    }

    if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(newPassword)); err == nil {
        return fmt.Errorf("new password cannot be the same as the old password")
    }

    hash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
    if err != nil {
        return fmt.Errorf("failed to hash password: %w", err)
    }

    if err := s.repo.ChangePassword(userID, string(hash)); err != nil {
        return fmt.Errorf("failed to update password: %w", err)
    }
    return nil
}

func (s *userService) ExploreUserByName(query string, threshold, page, size int) ([]model.User, int, error){
    return s.repo.ExploreUserByName(query, threshold, page, size)
}