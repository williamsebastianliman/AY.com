package service

import (
	"context"
	"fmt"
	"regexp"

	"github.com/google/uuid"
	"github.com/williamsebastianliman/WEB-WS-242/services/thread/internal/model"
	"github.com/williamsebastianliman/WEB-WS-242/services/thread/internal/repository"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ThreadService interface {
	CreateThread(ctx context.Context, userID, content string, communityID *string, parentID *string) (string, error)
	GetThreadByID(ctx context.Context, id string) (*model.Thread, error)
	UpdateThread(ctx context.Context, id string, content *string, communityID *string) error
	DeleteThread(ctx context.Context, id string) error
	ListThreads(ctx context.Context, page, size int) ([]model.Thread, error)
	ListThreadsByUser(ctx context.Context, userID string) ([]model.Thread, error)
	ListThreadsByUserLike(ctx context.Context, userID string) ([]model.Thread, error)
	ListThreadsByCommunity(ctx context.Context, communityID string, page, size int) ([]model.Thread, error)
	SearchThreads(ctx context.Context, keyword string, page, size int) ([]model.Thread, error)
	IncrementView(ctx context.Context, id string) error
	IncrementLike(ctx context.Context, id string) error
	IncrementComment(ctx context.Context, id string) error
	IncrementShare(ctx context.Context, id string) error
	CountThreads(ctx context.Context) (int64, error)
    LikeThread(ctx context.Context, threadID, userID string) error
    UnlikeThread(ctx context.Context, threadID, userID string) error
    CountThreadLikes(ctx context.Context, threadID string) (int64, error)

    BookmarkThread(ctx context.Context, threadID, userID string) error
    UnbookmarkThread(ctx context.Context, threadID, userID string) error
    CountThreadBookmarks(ctx context.Context, threadID string) (int64, error)

    CountThreadReplies(ctx context.Context, threadID string) (int64, error)

    ListThreadMedia(ctx context.Context, threadID string) ([]model.ThreadMedia, error)
    AddThreadMedia(ctx context.Context, threadID, mediaID string) error
    RemoveThreadMedia(ctx context.Context, threadID, mediaID string) error
    CountThreadMedia(ctx context.Context, threadID string) (int64, error)

	HasLikedThread(ctx context.Context, threadID, userID string) (bool, error)
	HasBookmarkedThread(ctx context.Context, threadID, userID string) (bool, error)
	ListThreadsByUserBookmark(ctx context.Context,userID string, query string) ([]model.Thread, error)
    ListThreadsReplies(ctx context.Context, parent_id string ,page int, size int) ([]model.Thread, error)
    ListRepliesByUser(ctx context.Context, user_id string) (map[*model.Thread][]*model.Thread, error)
    GetMediaByUserGrouped(userID string) (map[*model.Thread][]model.ThreadMedia, error)
    ListThreadsByUserIDs(ctx context.Context, userIDs []string, page, size int) ([]model.Thread, error)

    ListAllCategories(ctx context.Context) ([]model.ThreadCategories, error)
    CreateCategory(ctx context.Context, name string) (*model.ThreadCategories, error)
    ListThreadsByCategory(ctx context.Context, categoryID string, page, size int) ([]model.Thread, error)

    ListThreadsByCommunityByLike(ctx context.Context, communityID string) ([]model.Thread, error)
	ListThreadsByCommunityByTime(ctx context.Context, communityID string) ([]model.Thread, error)
	GetMediaByCommunity(communityID string, page, size int) (map[*model.Thread][]model.ThreadMedia, error)

    CreateRepost(ctx context.Context, userID, content string, communityID *string, parentID *string, repostID *string) (string, error)
    CountThreadRepost(ctx context.Context, threadID string) (int64, error)
    HasRepostedThread(ctx context.Context, threadID, userID string) (bool, error)
    PinThread(thread_id string) (error)
    DeleteCategory(ctx context.Context, categoryID string) error
    UpdateCategory(ctx context.Context, categoryID, newName string) (error)

    CreateThreadHashtag(ctx context.Context, threadID string, hashtag string) error
    GetThreadsByHashtag(ctx context.Context, hashtag string, page, size int) ([]model.Thread, error)
    GetTopHashtags(ctx context.Context, limit int) ([]repository.HashtagCount, error)
    SearchThreadByContent(query string, threshold, page, size int) ([]model.Thread, int, error)
    SearchMediaByContent(query string, threshold, page, size int) (map[*model.Thread][]model.ThreadMedia, int, error)
}

type threadService struct {
	repo repository.ThreadRepository
}

func NewThreadService(repo repository.ThreadRepository) ThreadService {
	return &threadService{repo: repo}
}

func (r *threadService) GetMediaByUserGrouped(userID string) (map[*model.Thread][]model.ThreadMedia, error) {
    tms, err := r.repo.GetMediaByUser(userID)
    if err != nil {
        return nil, err
    }

    threadIDSet := make(map[uuid.UUID]struct{})
    for _, tm := range tms {
        threadIDSet[tm.ThreadID] = struct{}{}
    }

    grouped := make(map[*model.Thread][]model.ThreadMedia)

    threadCache := make(map[uuid.UUID]*model.Thread)

    for _, tm := range tms {
        thread, ok := threadCache[tm.ThreadID]
        if !ok {
            threadObj, err := r.repo.GetByID(tm.ThreadID.String())
            if err != nil {
                return nil, err
            }
            if threadObj == nil {
                continue
            }
            thread = threadObj
            threadCache[tm.ThreadID] = threadObj
        }
        grouped[thread] = append(grouped[thread], tm)
    }

    return grouped, nil
} 

func (s *threadService) CreateThread(ctx context.Context, userID, content string, communityID *string, parentID *string) (string, error) {
    uid, err := uuid.Parse(userID)
    if err != nil {
        return "", status.Error(codes.InvalidArgument, "invalid user ID")
    }
    if content == "" {
        return "", status.Error(codes.InvalidArgument, "content cannot be empty")
    }

    var commUUID *uuid.UUID
    if communityID != nil {
        u, err := uuid.Parse(*communityID)
        if err != nil {
            return "", status.Error(codes.InvalidArgument, "invalid community ID")
        }
        commUUID = &u
    }

    var parUUID *uuid.UUID
    if parentID != nil {
        p, err := uuid.Parse(*parentID)
        if err != nil {
            return "", status.Error(codes.InvalidArgument, "Invalid Parent ID!")
        }
        parUUID = &p
    }

    thread := &model.Thread{
        UserID:      uid,
        CommunityID: commUUID,
        Content:     content,
        ParentID:    parUUID,
    }

    if err := s.repo.Create(thread); err != nil {
        return "", err
    }

    hashtagRegex := regexp.MustCompile(`#([a-zA-Z0-9_]+)`)
    matches := hashtagRegex.FindAllStringSubmatch(content, -1)

    uniqueHashtags := make(map[string]struct{})
    for _, match := range matches {
        hashtag := match[1]
        if _, exists := uniqueHashtags[hashtag]; !exists {
            uniqueHashtags[hashtag] = struct{}{}
        }
    }

    for hashtag := range uniqueHashtags {
        _ = s.CreateThreadHashtag(ctx, thread.ID.String(), hashtag)
    }
    
    return thread.ID.String(), nil
}

func (s *threadService) CreateRepost(ctx context.Context, userID, content string, communityID *string, parentID *string, repostID *string) (string, error) {
	uid, err := uuid.Parse(userID)
	if err != nil {
		return "", status.Error(codes.InvalidArgument, "invalid user ID")
	}
	if content == "" {
		return "", status.Error(codes.InvalidArgument, "content cannot be empty")
	}

	var commUUID *uuid.UUID
	if communityID != nil {
		u, err := uuid.Parse(*communityID)
		if err != nil {
			return "", status.Error(codes.InvalidArgument, "invalid community ID")
		}
		commUUID = &u
	}

    var parUUID *uuid.UUID
    if parentID != nil{
        p, err := uuid.Parse(*parentID)
        if err != nil{
            return "", status.Error(codes.InvalidArgument, "Invlaid Parent ID!")
        }
        parUUID = &p
    }

    var repostUUID *uuid.UUID
    if repostID != nil{
        p, err := uuid.Parse(*repostID)
        if err != nil{
            return "", status.Error(codes.InvalidArgument, "Invlaid Parent ID!")
        }
        repostUUID= &p
    }


	thread := &model.Thread{
		UserID:      uid,
		CommunityID: commUUID,
		Content:     content,
        ParentID: parUUID,
        RepostOfID: repostUUID,
	}

	if err := s.repo.Create(thread); err != nil {
		return "", err
	}
	return thread.ID.String(), nil
}

func (s *threadService) GetThreadByID(ctx context.Context, id string) (*model.Thread, error) {
	thread, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if thread == nil {
		return nil, status.Error(codes.NotFound, "thread not found")
	}
	return thread, nil
}

func (s *threadService) UpdateThread(ctx context.Context, id string, content *string, communityID *string) error {
	thread, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}
	if thread == nil {
		return status.Error(codes.NotFound, "thread not found")
	}

	if content != nil {
		if *content == "" {
			return status.Error(codes.InvalidArgument, "content cannot be empty")
		}
		thread.Content = *content
	}

	if communityID != nil {
		u, err := uuid.Parse(*communityID)
		if err != nil {
			return status.Error(codes.InvalidArgument, "invalid community ID")
		}
		thread.CommunityID = &u
	} else {
		thread.CommunityID = nil
	}

	return s.repo.Update(thread)
}

func (s *threadService) DeleteThread(ctx context.Context, id string) error {
	thread, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}
	if thread == nil {
		return status.Error(codes.NotFound, "thread not found")
	}
	return s.repo.Delete(id)
}

func (s *threadService) ListThreads(ctx context.Context, page, size int) ([]model.Thread, error) {
	if page < 1 || size < 1 {
		return nil, status.Error(codes.InvalidArgument, "page and size must be >= 1")
	}
	return s.repo.ListAll(page, size)
}

func (s *threadService) ListThreadsReplies(ctx context.Context, parent_id string ,page int, size int) ([]model.Thread, error) {
    uid, err := uuid.Parse(parent_id)
    if err != nil{
        return nil, err
    }
	if page < 1 || size < 1 {
		return nil, status.Error(codes.InvalidArgument, "page and size must be >= 1")
	}
	return s.repo.ListByParentID(uid, page, size)
}

func(s *threadService) ListRepliesByUser(ctx context.Context, user_id string) (map[*model.Thread][]*model.Thread, error) {
    grouped := make(map[*model.Thread][]*model.Thread)
    parent_cache := make(map[string]*model.Thread)
    uid, err := uuid.Parse(user_id)
    if err != nil{
        return nil, err
    }
    threads, err := s.repo.ListRepliesByUser(uid)
    if err != nil{
        return nil, err
    }

    for idx := range threads{
        thread := &threads[idx]
        if (thread.ParentID!=nil && (*thread.ParentID).String()!=""){
            parent_id := (*thread.ParentID).String()
            par,ok := parent_cache[parent_id]
            if !ok{
                par, err = s.repo.GetByID(parent_id)
                if err != nil{
                    return nil, err
                }
                parent_cache[parent_id] = par
            }
            grouped[par] = append(grouped[par], thread)   
        }
        
    }
    return grouped, nil
}


func (s *threadService) ListThreadsByUser(ctx context.Context, userID string) ([]model.Thread, error) {
	uid, err := uuid.Parse(userID)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid user ID")
	}
	return s.repo.ListByUser(uid)
}

func (s *threadService) ListThreadsByCommunity(ctx context.Context, communityID string, page, size int) ([]model.Thread, error) {
	uid, err := uuid.Parse(communityID)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid community ID")
	}
	return s.repo.ListByCommunityID(uid, page, size)
}

func (s *threadService) SearchThreads(ctx context.Context, keyword string, page, size int) ([]model.Thread, error) {
	if page < 1 || size < 1 {
		return nil, status.Error(codes.InvalidArgument, "page and size must be >= 1")
	}
	if keyword == "" {
		return nil, status.Error(codes.InvalidArgument, "keyword cannot be empty")
	}
	return s.repo.SearchByContent(keyword, page, size)
}

func (s *threadService) IncrementView(ctx context.Context, id string) error {
	return s.repo.IncrementViewCount(id)
}

func (s *threadService) IncrementLike(ctx context.Context, id string) error {
	return s.repo.IncrementLikeCount(id)
}

func (s *threadService) IncrementComment(ctx context.Context, id string) error {
	return s.repo.IncrementCommentCount(id)
}

func (s *threadService) IncrementShare(ctx context.Context, id string) error {
	return s.repo.IncrementShareCount(id)
}

func (s *threadService) CountThreads(ctx context.Context) (int64, error) {
	return s.repo.CountAll()
}

func (s *threadService) UnlikeThread(ctx context.Context, threadID, userID string) error {
    uid, err := uuid.Parse(userID)
    if err != nil {
        return status.Error(codes.InvalidArgument, "invalid user ID")
    }
    return s.repo.DeleteThreadLike(threadID, uid)
}

func (s *threadService) UnbookmarkThread(ctx context.Context, threadID, userID string) error {
    uid, err := uuid.Parse(userID)
    if err != nil {
        return status.Error(codes.InvalidArgument, "invalid user ID")
    }
    return s.repo.DeleteThreadBookmark(threadID, uid)
}

func (s *threadService) LikeThread(ctx context.Context, threadID, userID string) error {
    tid, err := uuid.Parse(threadID)
    if err != nil {
        return status.Error(codes.InvalidArgument, "invalid thread ID")
    }
    uid, err := uuid.Parse(userID)
    if err != nil {
        return status.Error(codes.InvalidArgument, "invalid user ID")
    }
    like := &model.ThreadLike{
        ThreadID: tid,
        UserID:   uid,
    }
    return s.repo.CreateThreadLike(like)
}

func (s *threadService) CountThreadLikes(ctx context.Context, threadID string) (int64, error) {
    return s.repo.CountLikes(threadID)
}

func (s *threadService) BookmarkThread(ctx context.Context, threadID, userID string) error {
    tid, err := uuid.Parse(threadID)
    if err != nil {
        return status.Error(codes.InvalidArgument, "invalid thread ID")
    }
    uid, err := uuid.Parse(userID)
    if err != nil {
        return status.Error(codes.InvalidArgument, "invalid user ID")
    }
    bm := &model.ThreadBookmark{
        ThreadID: tid,
        UserID:   uid,
    }
    return s.repo.CreateThreadBookmark(bm)
}

func (s *threadService) CountThreadBookmarks(ctx context.Context, threadID string) (int64, error) {
    return s.repo.CountBookmarks(threadID)
}

func (s *threadService) CountThreadRepost(ctx context.Context, threadID string) (int64, error) {
    return s.repo.CountReposts(threadID)
}

func (s *threadService) CountThreadReplies(ctx context.Context, threadID string) (int64, error) {
    return s.repo.CountReplies(threadID)
}

func (s *threadService) ListThreadMedia(ctx context.Context, threadID string) ([]model.ThreadMedia, error) {
    return s.repo.ListThreadMedia(threadID)
}

func (s *threadService) AddThreadMedia(ctx context.Context, threadID, mediaID string) error {
    tid, err := uuid.Parse(threadID)
    if err != nil {
        return status.Error(codes.InvalidArgument, "invalid thread ID")
    }
    mid, err := uuid.Parse(mediaID)
    if err != nil {
        return status.Error(codes.InvalidArgument, "invalid media ID")
    }
    tm := &model.ThreadMedia{
        ThreadID: tid,
        MediaID:  mid,
    }
    return s.repo.CreateThreadMedia(tm)
}

func (s *threadService) RemoveThreadMedia(ctx context.Context, threadID, mediaID string) error {
    tid, err := uuid.Parse(threadID)
    if err != nil {
        return status.Error(codes.InvalidArgument, "invalid thread ID")
    }
    mid, err := uuid.Parse(mediaID)
    if err != nil {
        return status.Error(codes.InvalidArgument, "invalid media ID")
    }
    return s.repo.DeleteThreadMedia(tid, mid)
}

func (s *threadService) CountThreadMedia(ctx context.Context, threadID string) (int64, error) {
    return s.repo.CountThreadMedia(threadID)
}

func (s *threadService) HasLikedThread(ctx context.Context, threadID, userID string) (bool, error) {
    uid, err := uuid.Parse(userID)
    if err != nil {
        return false, status.Error(codes.InvalidArgument, "invalid user ID")
    }
    return s.repo.LikeExists(threadID, uid)
}

func (s *threadService) HasBookmarkedThread(ctx context.Context, threadID, userID string) (bool, error) {
    uid, err := uuid.Parse(userID)
    if err != nil {
        return false, status.Error(codes.InvalidArgument, "invalid user ID")
    }
    return s.repo.BookmarkExists(threadID, uid)
}

func (s *threadService) HasRepostedThread(ctx context.Context, threadID, userID string) (bool, error) {
    uid, err := uuid.Parse(userID)
    if err != nil {
        return false, status.Error(codes.InvalidArgument, "invalid user ID")
    }
    return s.repo.RepostExist(threadID, uid)
}

func (s *threadService) ListThreadsByUserLike(
    ctx context.Context,
    userID string,
) ([]model.Thread, error) {
    uid, err := uuid.Parse(userID)
    if err != nil {
        return nil, fmt.Errorf("invalid userID %q: %w", userID, err)
    }

    idStrings, err := s.repo.ListLikedThreadIDStrings(uid)
    if err != nil {
        return nil, fmt.Errorf("failed to load liked IDs: %w", err)
    }

    var result []model.Thread
    for _, tid := range idStrings {
        t, err := s.repo.GetByID(tid)
        if err != nil {
            return nil, fmt.Errorf("failed to load thread %q: %w", tid, err)
        }
        if t != nil {
            result = append(result, *t)
        }
    }

    return result, nil
}

func (s *threadService) ListThreadsByUserBookmark(
    ctx context.Context,
    userID string,
    query string,
) ([]model.Thread, error) {
    uid, err := uuid.Parse(userID)
    if err != nil {
        return nil, fmt.Errorf("invalid userID %q: %w", userID, err)
    }

    idStrings, err := s.repo.ListBookmarkedThreadIDStrings(uid, query)
    if err != nil {
        return nil, fmt.Errorf("failed to load liked IDs: %w", err)
    }

    var result []model.Thread
    for _, tid := range idStrings {
        t, err := s.repo.GetByID(tid)
        if err != nil {
            return nil, fmt.Errorf("failed to load thread %q: %w", tid, err)
        }
        if t != nil {
            result = append(result, *t)
        }
    }

    return result, nil
}

func (r *threadService) GetMediaByCommunity(communityID string, page, size int) (map[*model.Thread][]model.ThreadMedia, error) {
    tms, err := r.repo.GetMediaByCommunity(communityID, page, size)
    if err != nil {
        return nil, err
    }

    grouped := make(map[*model.Thread][]model.ThreadMedia)
    threadCache := make(map[uuid.UUID]*model.Thread)

    for _, tm := range tms {
        thread, ok := threadCache[tm.ThreadID]
        if !ok {
            threadObj, err := r.repo.GetByID(tm.ThreadID.String())
            if err != nil {
                return nil, err
            }
            if threadObj == nil {
                continue
            }
            thread = threadObj
            threadCache[tm.ThreadID] = threadObj
        }
        grouped[thread] = append(grouped[thread], tm)
    }

    return grouped, nil
}

func (s *threadService) ListThreadsByUserIDs(ctx context.Context, userIDs []string, page, size int) ([]model.Thread, error) {
	if len(userIDs) == 0 {
		return []model.Thread{}, nil
	}
	if page < 1 || size < 1 {
		return nil, status.Error(codes.InvalidArgument, "page and size must be >= 1")
	}
	return s.repo.ListByUserIDs(userIDs, page, size)
}

func (s *threadService) ListAllCategories(ctx context.Context) ([]model.ThreadCategories, error) {
	return s.repo.ListAllCategories()
}

func (s *threadService) CreateCategory(ctx context.Context, name string) (*model.ThreadCategories, error) {
	if name == "" {
		return nil, status.Error(codes.InvalidArgument, "category name cannot be empty")
	}
	return s.repo.CreateCategory(name)
}

func (s *threadService) ListThreadsByCategory(ctx context.Context, categoryID string, page, size int) ([]model.Thread, error) {
	if categoryID == "" {
		return nil, status.Error(codes.InvalidArgument, "category ID cannot be empty")
	}
	cid, err := uuid.Parse(categoryID)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid category ID")
	}
	if page < 1 || size < 1 {
		return nil, status.Error(codes.InvalidArgument, "page and size must be >= 1")
	}
	return s.repo.ListThreadsByCategory(cid, page, size)
}

func (s *threadService) ListThreadsByCommunityByLike(ctx context.Context, communityID string) ([]model.Thread, error) {
	id, err := uuid.Parse(communityID)
	if err != nil {
		return nil, err
	}
	return s.repo.ListThreadsByCommunityByLike(id)
}

func (s *threadService) ListThreadsByCommunityByTime(ctx context.Context, communityID string) ([]model.Thread, error) {
	id, err := uuid.Parse(communityID)
	if err != nil {
		return nil, err
	}
	return s.repo.ListThreadsByCommunityByTime(id)
}

func (s *threadService) PinThread(thread_id string) (error){
    return s.repo.PinThread(thread_id)
}

func (s *threadService) DeleteCategory(ctx context.Context, categoryID string) error {
    if categoryID == "" {
        return status.Error(codes.InvalidArgument, "category ID cannot be empty")
    }
    cid, err := uuid.Parse(categoryID)
    if err != nil {
        return status.Error(codes.InvalidArgument, "invalid category ID")
    }
    return s.repo.DeleteCategory(cid)
}
func (s *threadService) UpdateCategory(ctx context.Context, categoryID, newName string) (error) {
    if categoryID == "" {
        return status.Error(codes.InvalidArgument, "category ID cannot be empty")
    }
    if newName == "" {
        return status.Error(codes.InvalidArgument, "category name cannot be empty")
    }
    cid, err := uuid.Parse(categoryID)
    if err != nil {
        return status.Error(codes.InvalidArgument, "invalid category ID")
    }
    return s.repo.UpdateCategory(cid, newName)
}

func (s *threadService) CreateThreadHashtag(ctx context.Context, threadID string, hashtag string) error {
    tid, err := uuid.Parse(threadID)
    if err != nil {
        return status.Error(codes.InvalidArgument, "invalid thread ID")
    }
    if hashtag == "" {
        return status.Error(codes.InvalidArgument, "hashtag cannot be empty")
    }
    return s.repo.CreateThreadHashtag(tid, hashtag)
}

func (s *threadService) GetThreadsByHashtag(ctx context.Context, hashtag string, page, size int) ([]model.Thread, error) {
    if hashtag == "" {
        return nil, status.Error(codes.InvalidArgument, "hashtag cannot be empty")
    }
    if page < 1 || size < 1 {
        return nil, status.Error(codes.InvalidArgument, "page and size must be >= 1")
    }
    return s.repo.GetThreadsByHashtag(hashtag, page, size)
}

func (s *threadService) GetTopHashtags(ctx context.Context, limit int) ([]repository.HashtagCount, error) {
    if limit < 1 {
        return nil, status.Error(codes.InvalidArgument, "limit must be >= 1")
    }
    return s.repo.GetTopHashtags(limit)
}

func (s *threadService) SearchThreadByContent(query string, threshold, page, size int) ([]model.Thread, int, error){
    return s.repo.SearchThreadByContent(query, threshold, page, size);
}

func (s *threadService) SearchMediaByContent(query string, threshold, page, size int) (map[*model.Thread][]model.ThreadMedia, int, error) {
    medias, total, err := s.repo.SearchMediaByContent(query, threshold, page, size)
    if err != nil {
        return nil, 0, err
    }

    grouped := make(map[*model.Thread][]model.ThreadMedia)
    threadCache := make(map[string]*model.Thread)

    for _, media := range medias {
        tid := media.ThreadID.String()
        thread, ok := threadCache[tid]
        if !ok {
            threadObj, err := s.repo.GetByID(tid)
            if err != nil {
                return nil, 0, err
            }
            if threadObj == nil {
                continue
            }
            thread = threadObj
            threadCache[tid] = threadObj
        }
        grouped[thread] = append(grouped[thread], media)
    }
    return grouped, total, nil
}
