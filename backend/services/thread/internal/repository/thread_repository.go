package repository

import (
	"errors"
	"fmt"
	"log"
	"strings"
	"unicode"

	"github.com/google/uuid"
	"github.com/williamsebastianliman/WEB-WS-242/services/thread/internal/model"
	"gorm.io/gorm"
)

type ThreadRepository interface {
	Create(thread *model.Thread) error
	GetByID(id string) (*model.Thread, error)
	Update(thread *model.Thread) error
	Delete(id string) error
	ListAll(page, size int) ([]model.Thread, error)
	ListByUser(userID uuid.UUID) ([]model.Thread, error)
	ListByCommunityID(communityID uuid.UUID, page, size int) ([]model.Thread, error)
	SearchByContent(keyword string, page, size int) ([]model.Thread, error)
	IncrementViewCount(id string) error
	IncrementLikeCount(id string) error
	IncrementCommentCount(id string) error
	IncrementShareCount(id string) error
	CountAll() (int64, error)
	CountLikes(threadID string) (int64, error)
    CountReplies(threadID string) (int64, error)
    CountBookmarks(threadID string) (int64, error)
	CreateThreadLike(like *model.ThreadLike) error
    CreateThreadBookmark(bm *model.ThreadBookmark) error
	ListThreadMedia(threadID string) ([]model.ThreadMedia, error)
	DeleteThreadLike(threadID string, userID uuid.UUID) error
	DeleteThreadBookmark(threadID string, userID uuid.UUID) error
	CreateThreadMedia(tm *model.ThreadMedia) error
	DeleteThreadMedia(threadID, mediaID uuid.UUID) error
	CountThreadMedia(threadID string) (int64, error)
	LikeExists(threadID string, userID uuid.UUID) (bool, error)
	BookmarkExists(threadID string, userID uuid.UUID) (bool, error)
	ListLikedThreadIDStrings(userID uuid.UUID) ([]string, error)
	ListBookmarkedThreadIDStrings(userID uuid.UUID, query string) ([]string, error)
	ListByParentID(parentID uuid.UUID, page, size int) ([]model.Thread, error)
	ListRepliesByUser(userID uuid.UUID) ([]model.Thread, error)
	GetMediaByUser(userID string) ([]model.ThreadMedia, error)
	ListByUserIDs(userIDs []string, page, size int) ([]model.Thread, error)

	ListAllCategories() ([]model.ThreadCategories, error)
	CreateCategory(name string) (*model.ThreadCategories, error)
	ListThreadsByCategory(categoryID uuid.UUID, page, size int) ([]model.Thread, error)

	ListThreadsByCommunityByLike(communityID uuid.UUID) ([]model.Thread, error)
	ListThreadsByCommunityByTime(communityID uuid.UUID) ([]model.Thread, error)
	GetMediaByCommunity(communityID string, page, size int) ([]model.ThreadMedia, error)
	CountReposts(threadID string) (int64, error)
	RepostExist(threadID string, userID uuid.UUID) (bool, error)
	PinThread(threadID string) (error)

	DeleteCategory(categoryID uuid.UUID) error
	UpdateCategory(categoryID uuid.UUID, newName string) error

	CreateThreadHashtag(threadID uuid.UUID, hashtag string) error
	GetThreadsByHashtag(hashtag string, page, size int) ([]model.Thread, error)
	GetTopHashtags(limit int) ([]HashtagCount, error)
	SearchThreadByContent(query string, threshold, page, size int) ([]model.Thread, int, error)
	SearchMediaByContent(query string, threshold, page, size int) ([]model.ThreadMedia, int, error)
}

type threadRepository struct {
	db *gorm.DB
}

func NewThreadRepository(db *gorm.DB) ThreadRepository {
	return &threadRepository{db: db}
}


func (r *threadRepository) Create(thread *model.Thread) error {
	if thread.ID == uuid.Nil {
		thread.ID = uuid.New()
	}
	return r.db.Create(thread).Error
}

func (r *threadRepository) GetByID(id string) (*model.Thread, error) {
	var t model.Thread
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}
	err = r.db.First(&t, "id = ?", uid).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &t, err
}

func (r *threadRepository) Update(thread *model.Thread) error {
	if thread.ID == uuid.Nil {
		return fmt.Errorf("cannot update thread with empty ID")
	}
	return r.db.Save(thread).Error
}

func (r *threadRepository) Delete(id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return err
	}
	return r.db.Delete(&model.Thread{}, "id = ?", uid).Error
}

func (r *threadRepository) ListAll(page, size int) ([]model.Thread, error) {
	var list []model.Thread
	offset := (page - 1) * size
	err := r.db.
		Where("parent_id IS NULL").
		Order("created_at DESC").
		Offset(offset).
		Limit(size).
		Find(&list).
		Error
	return list, err
}

func (r *threadRepository) ListByUserIDs(userIDs []string, page, size int) ([]model.Thread, error) {
	var list []model.Thread
	offset := (page - 1) * size

	if len(userIDs) == 0 {
		return []model.Thread{}, nil
	}

	err := r.db.
		Where("parent_id IS NULL").
		Where("user_id IN ?", userIDs).
		Order("created_at DESC").
		Offset(offset).
		Limit(size).
		Find(&list).Error

	return list, err
}

func (r *threadRepository) ListByUser(userID uuid.UUID) ([]model.Thread, error) {
	var list []model.Thread
	err := r.db.
		Where("user_id = ? AND parent_id IS NULL", userID).
		Order("is_pinned DESC").              
		Order("created_at DESC").           
		Find(&list).
		Error
	log.Printf("ListByUser: userID=%s returned %d threads\n", userID, len(list))
	return list, err
}


func (r *threadRepository) ListRepliesByUser(userID uuid.UUID) ([]model.Thread, error) {
	var list []model.Thread
	err := r.db.
		Where("user_id = ? AND parent_id IS NOT NULL", userID).
		Order("is_pinned DESC").
		Order("created_at DESC").
		Find(&list).
		Error
	log.Printf("ListByUser: userID=%s returned %d threads\n", userID, len(list))
	return list, err
}

func (r *threadRepository) ListLikedThreadIDStrings(userID uuid.UUID) ([]string, error) {
    var ids []string

    err := r.db.
        Table("thread_likes").
        Where("user_id = ?", userID).
        Pluck("thread_id::text", &ids).
        Error

    if err != nil {
        return nil, err
    }

    return ids, nil
}

func (r *threadRepository) ListBookmarkedThreadIDStrings(userID uuid.UUID, query string) ([]string, error) {
    var ids []string

    err := r.db.
        Table("thread_bookmarks").
        Joins("JOIN threads ON threads.id = thread_bookmarks.thread_id").
        Where("thread_bookmarks.user_id = ?", userID).
        Where("threads.content ILIKE ?", "%"+query+"%").
        Pluck("thread_bookmarks.thread_id::text", &ids).
        Error

    if err != nil {
        return nil, err
    }

    return ids, nil
}


func (r *threadRepository) ListByCommunityID(communityID uuid.UUID, page, size int) ([]model.Thread, error) {
	var list []model.Thread
	offset := (page - 1) * size
	err := r.db.
		Where("community_id = ?", communityID).
		Order("created_at DESC").
		Offset(offset).
		Limit(size).
		Find(&list).
		Error
	return list, err
}

func (r *threadRepository) ListByParentID(parentID uuid.UUID, page, size int) ([]model.Thread, error) {
	var list []model.Thread
	offset := (page - 1) * size
	err := r.db.
		Where("parent_id = ?", parentID).
		Order("created_at DESC").
		Offset(offset).
		Limit(size).
		Find(&list).
		Error
	return list, err
}

func (r *threadRepository) SearchByContent(keyword string, page, size int) ([]model.Thread, error) {
	var list []model.Thread
	offset := (page - 1) * size
	pattern := fmt.Sprintf("%%%s%%", keyword)
	err := r.db.
		Where("content ILIKE ?", pattern).
		Order("created_at DESC").
		Offset(offset).
		Limit(size).
		Find(&list).
		Error
	return list, err
}

func (r *threadRepository) IncrementViewCount(id string) error {
	return r.incrementCounter(id, "view_count")
}

func (r *threadRepository) IncrementLikeCount(id string) error {
	return r.incrementCounter(id, "like_count")
}

func (r *threadRepository) IncrementCommentCount(id string) error {
	return r.incrementCounter(id, "comment_count")
}

func (r *threadRepository) IncrementShareCount(id string) error {
	return r.incrementCounter(id, "share_count")
}

// incrementCounter is a helper to atomically add 1 to the given numeric column.
func (r *threadRepository) incrementCounter(id, column string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return err
	}
	return r.db.
		Model(&model.Thread{}).
		Where("id = ?", uid).
		UpdateColumn(column, gorm.Expr(fmt.Sprintf("%s + 1", column))).
		Error
}

func (r *threadRepository) CountAll() (int64, error) {
	var count int64
	err := r.db.Model(&model.Thread{}).Count(&count).Error
	return count, err
}

func (r *threadRepository) CountLikes(threadID string) (int64, error) {
    uid, err := uuid.Parse(threadID)
    if err != nil {
        return 0, err
    }
    var count int64
    err = r.db.
        Model(&model.ThreadLike{}).
        Where("thread_id = ?", uid).
        Count(&count).
        Error
    return count, err
}

func (r *threadRepository) CountReplies(threadID string) (int64, error) {
    uid, err := uuid.Parse(threadID)
    if err != nil {
        return 0, err
    }

    var count int64
    err = r.db.
        Model(&model.Thread{}).
        Where("parent_id = ?", uid).
        Count(&count).
        Error

    return count, err
}

func (r *threadRepository) CountBookmarks(threadID string) (int64, error) {
    uid, err := uuid.Parse(threadID)
    if err != nil {
        return 0, err
    }
    var count int64
    err = r.db.
        Model(&model.ThreadBookmark{}).
        Where("thread_id = ?", uid).
        Count(&count).
        Error
    return count, err
}

func (r *threadRepository) CreateThreadLike(like *model.ThreadLike) error {
    if like.ThreadID == uuid.Nil || like.UserID == uuid.Nil {
        return fmt.Errorf("threadID and userID must be set")
    }
    return r.db.Create(like).Error
}

func (r *threadRepository) CreateThreadBookmark(bm *model.ThreadBookmark) error {
    if bm.ThreadID == uuid.Nil || bm.UserID == uuid.Nil {
        return fmt.Errorf("threadID and userID must be set")
    }
    return r.db.Create(bm).Error
}

func (r *threadRepository) ListThreadMedia(threadID string) ([]model.ThreadMedia, error) {
    uid, err := uuid.Parse(threadID)
    if err != nil {
        return nil, err
    }

    var mediaList []model.ThreadMedia
    err = r.db.
        Where("thread_id = ?", uid).
        Find(&mediaList).
        Error

    return mediaList, err
}

func (r *threadRepository) DeleteThreadLike(threadID string, userID uuid.UUID) error {
    tid, err := uuid.Parse(threadID)
    if err != nil {
        return err
    }
    return r.db.
        Delete(&model.ThreadLike{}, "thread_id = ? AND user_id = ?", tid, userID).
        Error
}

func (r *threadRepository) DeleteThreadBookmark(threadID string, userID uuid.UUID) error {
    tid, err := uuid.Parse(threadID)
    if err != nil {
        return err
    }
    return r.db.
        Delete(&model.ThreadBookmark{}, "thread_id = ? AND user_id = ?", tid, userID).
        Error
}

func (r *threadRepository) CreateThreadMedia(tm *model.ThreadMedia) error {
    return r.db.Create(tm).Error
}

func (r *threadRepository) DeleteThreadMedia(threadID, mediaID uuid.UUID) error {
    return r.db.
        Delete(&model.ThreadMedia{}, "thread_id = ? AND media_id = ?", threadID, mediaID).
        Error
}

func (r *threadRepository) CountThreadMedia(threadID string) (int64, error) {
    uid, err := uuid.Parse(threadID)
    if err != nil {
        return 0, err
    }
    var count int64
    err = r.db.
        Model(&model.ThreadMedia{}).
        Where("thread_id = ?", uid).
        Count(&count).
        Error
    return count, err
}

func (r *threadRepository) LikeExists(threadID string, userID uuid.UUID) (bool, error) {
	tid, err := uuid.Parse(threadID)
	if err != nil {
		return false, err
	}
	var count int64
	err = r.db.Model(&model.ThreadLike{}).
		Where("thread_id = ? AND user_id = ?", tid, userID).
		Count(&count).Error
	return count > 0, err
}

func (r *threadRepository) BookmarkExists(threadID string, userID uuid.UUID) (bool, error) {
	tid, err := uuid.Parse(threadID)
	if err != nil {
		return false, err
	}
	var count int64
	err = r.db.Model(&model.ThreadBookmark{}).
		Where("thread_id = ? AND user_id = ?", tid, userID).
		Count(&count).Error
	return count > 0, err
}

func (r *threadRepository) GetMediaByUser(userID string) ([]model.ThreadMedia, error) {
    uid, err := uuid.Parse(userID)
    if err != nil {
        return nil, err
    }
    var medias []model.ThreadMedia

    err = r.db.
        Table("thread_media").
        Joins("JOIN threads ON threads.id = thread_media.thread_id").
        Where("threads.user_id = ?", uid).
        Find(&medias).Error

    return medias, err
}

func (r *threadRepository) ListAllCategories() ([]model.ThreadCategories, error) {
    var cats []model.ThreadCategories
    err := r.db.Order("created_at DESC").Find(&cats).Error
    return cats, err
}

func (r *threadRepository) CreateCategory(name string) (*model.ThreadCategories, error) {
    cat := &model.ThreadCategories{
        CategoryID:   uuid.New(),
        CategoryName: name,
    }
    if err := r.db.Create(cat).Error; err != nil {
        return nil, err
    }
    return cat, nil
}

func (r *threadRepository) ListThreadsByCategory(categoryID uuid.UUID, page, size int) ([]model.Thread, error) {
    var threads []model.Thread
    offset := (page - 1) * size

    err := r.db.Table("threads").
        Joins("JOIN thread_details ON threads.id = thread_details.thread_id").
        Where("thread_details.category_id = ?", categoryID).
        Order("threads.created_at DESC").
        Offset(offset).
        Limit(size).
        Find(&threads).Error

    return threads, err
}

func (r *threadRepository) ListThreadsByCommunityByLike(communityID uuid.UUID) ([]model.Thread, error) {
	var threads []model.Thread
	err := r.db.
		Model(&model.Thread{}).
		Select("threads.*, COUNT(thread_likes.thread_id) as like_count").
		Joins("LEFT JOIN thread_likes ON threads.id = thread_likes.thread_id").
		Where("threads.community_id = ?", communityID).
		Group("threads.id").
		Order("COUNT(thread_likes.thread_id) DESC, threads.created_at DESC").
		Find(&threads).Error
	return threads, err
}


func (r *threadRepository) ListThreadsByCommunityByTime(communityID uuid.UUID) ([]model.Thread, error) {
	var threads []model.Thread
	err := r.db.
		Model(&model.Thread{}).
		Where("community_id = ?", communityID).
		Order("created_at DESC").
		Find(&threads).Error
	return threads, err
}

func (r *threadRepository) GetMediaByCommunity(communityID string, page, size int) ([]model.ThreadMedia, error) {
    cid, err := uuid.Parse(communityID)
    if err != nil {
        return nil, err
    }

    var medias []model.ThreadMedia
    offset := (page - 1) * size
    log.Printf("offset: %d", offset)
    log.Printf("page: %d", page)
    log.Printf("size: %d", size)
    err = r.db.
        Table("thread_media").
        Joins("JOIN threads ON threads.id = thread_media.thread_id").
        Where("threads.community_id = ?", cid).
        Order("threads.created_at DESC").
        Offset(offset).
        Limit(size).
        Find(&medias).Error

    return medias, err
}

func (r *threadRepository) CountReposts(threadID string) (int64, error) {
    var count int64
    err:= r.db.
        Model(&model.Thread{}).
        Where("repost_of_id = ?", threadID).
        Count(&count).Error
    return count, err
}

func (r *threadRepository) RepostExist(threadID string, userID uuid.UUID) (bool, error) {
	tid, err := uuid.Parse(threadID)
	if err != nil {
		return false, err
	}
	var count int64
	err = r.db.Model(&model.Thread{}).
		Where("repost_of_id = ? AND user_id = ?", tid, userID).
		Count(&count).Error
	return count > 0, err
}

func (r *threadRepository) PinThread(threadID string) (error){
	tid, err := uuid.Parse(threadID)
	if err != nil {
		return err
	}
	return r.db.Model(&model.Thread{}).
		Where("id = ?", tid).
		Update("is_pinned", true).Error
}

func (r *threadRepository) DeleteCategory(categoryID uuid.UUID) error {
	return r.db.Delete(&model.ThreadCategories{}, "category_id = ?", categoryID).Error
}

func (r *threadRepository) UpdateCategory(categoryID uuid.UUID, newName string) error {
	return r.db.Model(&model.ThreadCategories{}).
		Where("category_id = ?", categoryID).
		Update("category_name", newName).
		Error
}

func (r *threadRepository) CreateThreadHashtag(threadID uuid.UUID, hashtag string) error {
    tag := &model.ThreadHashtag{
        ThreadID:    threadID,
        HashtagName: hashtag,
    }
    return r.db.Create(tag).Error
}

func (r *threadRepository) GetThreadsByHashtag(hashtag string, page, size int) ([]model.Thread, error) {
    var threads []model.Thread
    offset := (page - 1) * size

    err := r.db.
        Table("threads").
        Select("threads.*").
        Joins("JOIN thread_hashtags ON threads.id = thread_hashtags.thread_id").
        Joins("LEFT JOIN thread_likes ON threads.id = thread_likes.thread_id").
        Where("thread_hashtags.hashtag_name = ?", hashtag).
        Group("threads.id").
        Order("COUNT(thread_likes.thread_id) DESC, threads.created_at DESC").
        Offset(offset).
        Limit(size).
        Find(&threads).Error

    return threads, err
}

type HashtagCount struct {
    HashtagName string
    Count       int64
}

func (r *threadRepository) GetTopHashtags(limit int) ([]HashtagCount, error) {
    var tags []HashtagCount
    err := r.db.
        Table("thread_hashtags").
        Select("hashtag_name, COUNT(*) as count").
        Group("hashtag_name").
        Order("count DESC").
        Limit(limit).
        Scan(&tags).Error
    return tags, err
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

func (r *threadRepository) SearchThreadByContent(query string, threshold, page, size int) ([]model.Thread, int, error) {
    var threads []model.Thread

    err := r.db.Order("created_at DESC").Find(&threads).Error
    if err != nil {
        return nil, 0, err
    }

    var matched []model.Thread
    for _, thread := range threads {
        if TokenMatching(query, thread.Content, threshold) {
            matched = append(matched, thread)
        }
    }

    total := len(matched)

    start := (page - 1) * size
    end := start + size
    if start >= total {
        return []model.Thread{}, total, nil
    }
    if end > total {
        end = total
    }

    return matched[start:end], total, nil
}

func (r *threadRepository) SearchMediaByContent(query string, threshold, page, size int) ([]model.ThreadMedia, int, error) {
    type MediaWithThread struct {
        model.ThreadMedia
        ThreadContent string
    }
    var allMedia []MediaWithThread

    err := r.db.
        Table("thread_media").
        Select("thread_media.*, threads.content as thread_content").
        Joins("JOIN threads ON threads.id = thread_media.thread_id").
        Order("thread_media.created_at DESC").
        Find(&allMedia).Error
    if err != nil {
        return nil, 0, err
    }

    var matched []model.ThreadMedia
    for _, mw := range allMedia {
        if TokenMatching(query, mw.ThreadContent, threshold) {
            matched = append(matched, mw.ThreadMedia)
        }
    }

    total := len(matched)
    start := (page - 1) * size
    end := start + size
    if start >= total {
        return []model.ThreadMedia{}, total, nil
    }
    if end > total {
        end = total
    }

    return matched[start:end], total, nil
}
