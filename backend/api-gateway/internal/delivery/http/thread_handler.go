package http

import (
	"net/http"
	"sort"
	"strconv"

	"github.com/gin-gonic/gin"
	threadpb "github.com/williamsebastianliman/WEB-WS-242/proto/gen/go/threadpb/proto/thread"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ThreadCreateRequest struct {
	UserID      string `json:"user_id" binding:"required"`
	Content     string `json:"content" binding:"required"`
	ParentID    string `json:"parent_id"`
	CommunityID string `json:"community_id"`
}

type ThreadJSON struct {
    ID          string    `json:"id"`
    UserID      string    `json:"user_id"`
    CommunityID string    `json:"community_id"`
    Content     string    `json:"content"`
    LikeCount   int64     `json:"like_count"`
    CommentCount int64    `json:"comment_count"`
    ShareCount  int64     `json:"share_count"`
    ViewCount   int64     `json:"view_count"`
    CreatedAt   TimeJSON  `json:"created_at"`
    UpdatedAt   TimeJSON  `json:"updated_at"`
	RepostID string `json:"repost_id"`
	IsPinned bool `json:"is_pinned"`
}

type CategoryJSON struct {
    CategoryId   string `json:"category_id"`
    CategoryName string `json:"category_name"`
}

type TimeJSON struct {
    Year     int32  `json:"year"`
    Month    int32  `json:"month"`
    Day      int32  `json:"day"`
    Hour     int32  `json:"hour"`
    Minute   int32  `json:"minute"`
    Second   int32  `json:"second"`
    Timezone string `json:"timezone"`
}

func toThreadCategoryJSON(pb *threadpb.Category) CategoryJSON {
    return CategoryJSON{
        CategoryId:   pb.GetCategoryId(),
        CategoryName: pb.GetCategoryName(),
    }
}

func toThreadJSON(proto *threadpb.Thread) ThreadJSON {
    created := proto.GetCreatedAt()
    updated := proto.GetUpdatedAt()
    return ThreadJSON{
        ID:          proto.GetId(),
        UserID:      proto.GetUserId(),
        CommunityID: proto.GetCommunityId(),
        Content:     proto.GetContent(),
        LikeCount:   proto.GetLikeCount(),
        CommentCount: proto.GetCommentCount(),
        ShareCount:  proto.GetShareCount(),
        ViewCount:   proto.GetViewCount(),
		RepostID: proto.RepostId,
		IsPinned: proto.IsPinned,
        CreatedAt: TimeJSON{
            Year:     created.GetYear(),
            Month:    created.GetMonth(),
            Day:      created.GetDay(),
            Hour:     created.GetHour(),
            Minute:   created.GetMinute(),
            Second:   created.GetSecond(),
            Timezone: created.GetTimezone(),
        },
        UpdatedAt: TimeJSON{
            Year:     updated.GetYear(),
            Month:    updated.GetMonth(),
            Day:      updated.GetDay(),
            Hour:     updated.GetHour(),
            Minute:   updated.GetMinute(),
            Second:   updated.GetSecond(),
            Timezone: updated.GetTimezone(),
        },
    }
}

// CreateThread godoc
// @Summary      Create a new thread
// @Description  Create a thread with user ID, content, and optional parent/community
// @Tags         thread
// @Accept       json
// @Produce      json
// @Param        body  body      ThreadCreateRequest  true  "Thread content"
// @Success      200   {object}  map[string]string
// @Failure      400   {object}  map[string]string
// @Router       /api/v1/threads [post]
func CreateThread(client threadpb.ThreadServiceClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			UserID      string `json:"user_id" binding:"required"`
			Content     string `json:"content" binding:"required"`
			ParentID string `json:"parent_id"`
			CommunityID string `json:"community_id"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		resp, err := client.CreateThread(c.Request.Context(), &threadpb.CreateThreadRequest{
			UserId:      req.UserID,
			Content:     req.Content,
			CommunityId: req.CommunityID,
			ParentId: req.ParentID,
		})
		if err != nil {
			st := status.Convert(err)
			c.JSON(codesToHTTP(st.Code()), gin.H{"error": st.Message()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"id": resp.Id})
	}
}

func CreateRepost(client threadpb.ThreadServiceClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			UserID      string `json:"user_id" binding:"required"`
			Content     string `json:"content"`
			ParentID string `json:"parent_id"`
			CommunityID string `json:"community_id"`
			RepostID string `json:"repost_id"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		resp, err := client.CreateRepost(c.Request.Context(), &threadpb.RepostThreadRequest{
			UserId:      req.UserID,
			Content:     req.Content,
			CommunityId: req.CommunityID,
			ParentId: req.ParentID,
			RepostId: req.RepostID,
		})
		if err != nil {
			st := status.Convert(err)
			c.JSON(codesToHTTP(st.Code()), gin.H{"error": st.Message()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"id": resp.Id})
	}
}

// GetThreadByID godoc
// @Summary      Get a thread by ID
// @Description  Retrieve a single thread using its unique ID
// @Tags         thread
// @Param        id   path      string  true  "Thread ID"
// @Success      200  {object}  ThreadJSON
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Router       /api/v1/threads/{id} [get]
func GetThreadByID(client threadpb.ThreadServiceClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if id == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "id required"})
			return
		}

		resp, err := client.GetThreadByID(c.Request.Context(), &threadpb.GetThreadByIDRequest{Id: id})
		if err != nil {
			st := status.Convert(err)
			c.JSON(codesToHTTP(st.Code()), gin.H{"error": st.Message()})
			return
		}
		c.JSON(http.StatusOK, toThreadJSON(resp.Thread))
	}
}


type UpdateThreadRequest struct {
    Content     string `json:"content"`
    CommunityID string `json:"community_id"`
}
// UpdateThread godoc
// @Summary      Update a thread by ID
// @Description  Update the content and/or community of an existing thread by its unique ID
// @Tags         thread
// @Accept       json
// @Produce      json
// @Param        id    path      string               true  "Thread ID"
// @Param        body  body      UpdateThreadRequest  true  "Updated content and/or community ID"
// @Success      204   "No Content"
// @Failure      400   {object}  map[string]string
// @Failure      404   {object}  map[string]string
// @Failure      500   {object}  map[string]string
// @Router       /api/v1/threads/{id} [put]
func UpdateThread(client threadpb.ThreadServiceClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if id == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "id required"})
			return
		}

		var req struct {
			Content     string `json:"content"`
			CommunityID string `json:"community_id"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		_, err := client.UpdateThread(c.Request.Context(), &threadpb.UpdateThreadRequest{
			Id:          id,
			Content:     req.Content,
			CommunityId: req.CommunityID,
		})
		if err != nil {
			st := status.Convert(err)
			c.JSON(codesToHTTP(st.Code()), gin.H{"error": st.Message()})
			return
		}
		c.Status(http.StatusNoContent)
	}
}



// DeleteThread godoc
// @Summary      Delete a thread by ID
// @Description  Delete a thread using its unique ID
// @Tags         thread
// @Param        id   path      string  true  "Thread ID"
// @Success      204  {string}  string  "No Content"
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Router       /api/v1/threads/{id} [delete]
func DeleteThread(client threadpb.ThreadServiceClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if id == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "id required"})
			return
		}

		_, err := client.DeleteThread(c.Request.Context(), &threadpb.DeleteThreadRequest{Id: id})
		if err != nil {
			st := status.Convert(err)
			c.JSON(codesToHTTP(st.Code()), gin.H{"error": st.Message()})
			return
		}
		c.Status(http.StatusNoContent)
	}
}

func ListThreads(client threadpb.ThreadServiceClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			Page int32 `form:"page" binding:"required"`
			Size int32 `form:"size" binding:"required"`
		}
		if err := c.ShouldBindQuery(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		resp, err := client.ListThreads(c.Request.Context(), &threadpb.ListThreadsRequest{
			Page: req.Page,
			Size: req.Size,
		})
		if err != nil {
			st := status.Convert(err)
			c.JSON(codesToHTTP(st.Code()), gin.H{"error": st.Message()})
			return
		}

		var threads []ThreadJSON
		for _, t := range resp.Threads {
			threads = append(threads, toThreadJSON(t))
		}
		c.JSON(http.StatusOK, threads)
	}
}

func ListThreadsByUser(client threadpb.ThreadServiceClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.Param("user_id")
		if userID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "user_id required"})
			return
		}

		resp, err := client.ListThreadsByUser(c.Request.Context(), &threadpb.ListByUserRequest{
			UserId: userID,
		})
		if err != nil {
			st := status.Convert(err)
			c.JSON(codesToHTTP(st.Code()), gin.H{"error": st.Message()})
			return
		}
		var threads []ThreadJSON
		for _, t := range resp.Threads {
			threads = append(threads, toThreadJSON(t))
		}
		c.JSON(http.StatusOK, threads)
	}
}

func ListThreadsByUserLike(client threadpb.ThreadServiceClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.Param("user_id")
		if userID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "user_id required"})
			return
		}

		resp, err := client.ListThreadsByUserLike(c.Request.Context(), &threadpb.ListByUserRequest{
			UserId: userID,
		})
		if err != nil {
			st := status.Convert(err)
			c.JSON(codesToHTTP(st.Code()), gin.H{"error": st.Message()})
			return
		}
		var threads []ThreadJSON
		for _, t := range resp.Threads {
			threads = append(threads, toThreadJSON(t))
		}
		c.JSON(http.StatusOK, threads)
	}
}

func ListThreadsByUserBookmark(client threadpb.ThreadServiceClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.Param("user_id")
		query := c.Query("query");
		if userID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "user_id required"})
			return
		}

		resp, err := client.ListThreadsByUserBookmark(c.Request.Context(), &threadpb.ListByUserRequest{
			UserId: userID,
			Query: query,
		})
		if err != nil {
			st := status.Convert(err)
			c.JSON(codesToHTTP(st.Code()), gin.H{"error": st.Message()})
			return
		}
		var threads []ThreadJSON
		for _, t := range resp.Threads {
			threads = append(threads, toThreadJSON(t))
		}
		c.JSON(http.StatusOK, threads)
	}
}

func ListThreadsByCommunity(client threadpb.ThreadServiceClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		communityID := c.Param("community_id")
		if communityID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "community_id required"})
			return
		}
		var req struct {
			Page int32 `form:"page" binding:"required"`
			Size int32 `form:"size" binding:"required"`
		}
		if err := c.ShouldBindQuery(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		resp, err := client.ListThreadsByCommunity(c.Request.Context(), &threadpb.ListByCommunityRequest{
			CommunityId: communityID,
			Page:        req.Page,
			Size:        req.Size,
		})
		if err != nil {
			st := status.Convert(err)
			c.JSON(codesToHTTP(st.Code()), gin.H{"error": st.Message()})
			return
		}
		var threads []ThreadJSON
		for _, t := range resp.Threads {
			threads = append(threads, toThreadJSON(t))
		}
		c.JSON(http.StatusOK, threads)
	}
}

func SearchThreads(client threadpb.ThreadServiceClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			Keyword string `form:"keyword" binding:"required"`
			Page    int32  `form:"page" binding:"required"`
			Size    int32  `form:"size" binding:"required"`
		}
		if err := c.ShouldBindQuery(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		resp, err := client.SearchThreads(c.Request.Context(), &threadpb.SearchThreadsRequest{
			Keyword: req.Keyword,
			Page:    req.Page,
			Size:    req.Size,
		})
		if err != nil {
			st := status.Convert(err)
			c.JSON(codesToHTTP(st.Code()), gin.H{"error": st.Message()})
			return
		}
		var threads []ThreadJSON
		for _, t := range resp.Threads {
			threads = append(threads, toThreadJSON(t))
		}
		c.JSON(http.StatusOK, threads)
	}
}

func IncrementView(client threadpb.ThreadServiceClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if id == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "id required"})
			return
		}
		_, err := client.IncrementView(c.Request.Context(), &threadpb.IncrementRequest{Id: id})
		if err != nil {
			st := status.Convert(err)
			c.JSON(codesToHTTP(st.Code()), gin.H{"error": st.Message()})
			return
		}
		c.Status(http.StatusNoContent)
	}
}

func IncrementLike(client threadpb.ThreadServiceClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if id == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "id required"})
			return
		}
		_, err := client.IncrementLike(c.Request.Context(), &threadpb.IncrementRequest{Id: id})
		if err != nil {
			st := status.Convert(err)
			c.JSON(codesToHTTP(st.Code()), gin.H{"error": st.Message()})
			return
		}
		c.Status(http.StatusNoContent)
	}
}

func IncrementComment(client threadpb.ThreadServiceClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if id == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "id required"})
			return
		}
		_, err := client.IncrementComment(c.Request.Context(), &threadpb.IncrementRequest{Id: id})
		if err != nil {
			st := status.Convert(err)
			c.JSON(codesToHTTP(st.Code()), gin.H{"error": st.Message()})
			return
		}
		c.Status(http.StatusNoContent)
	}
}

func IncrementShare(client threadpb.ThreadServiceClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if id == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "id required"})
			return
		}
		_, err := client.IncrementShare(c.Request.Context(), &threadpb.IncrementRequest{Id: id})
		if err != nil {
			st := status.Convert(err)
			c.JSON(codesToHTTP(st.Code()), gin.H{"error": st.Message()})
			return
		}
		c.Status(http.StatusNoContent)
	}
}

func CountThreads(client threadpb.ThreadServiceClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		resp, err := client.CountThreads(c.Request.Context(), &threadpb.Empty{})
		if err != nil {
			st := status.Convert(err)
			c.JSON(codesToHTTP(st.Code()), gin.H{"error": st.Message()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"count": resp.Count})
	}
}

func LikeThread(client threadpb.ThreadServiceClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		threadID := c.Param("id")
		if threadID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "thread id required"})
			return
		}
		var body struct {
			UserID string `json:"user_id" binding:"required"`
		}
		if err := c.ShouldBindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		_, err := client.LikeThread(c.Request.Context(), &threadpb.ThreadUserRequest{
			ThreadId: threadID,
			UserId:   body.UserID,
		})
		if err != nil {
			st := status.Convert(err)
			c.JSON(codesToHTTP(st.Code()), gin.H{"error": st.Message()})
			return
		}
		c.Status(http.StatusNoContent)
	}
}

func UnlikeThread(client threadpb.ThreadServiceClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		threadID := c.Param("id")
		if threadID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "thread id required"})
			return
		}
		var body struct {
			UserID string `json:"user_id" binding:"required"`
		}
		if err := c.ShouldBindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		_, err := client.UnlikeThread(c.Request.Context(), &threadpb.ThreadUserRequest{
			ThreadId: threadID,
			UserId:   body.UserID,
		})
		if err != nil {
			st := status.Convert(err)
			c.JSON(codesToHTTP(st.Code()), gin.H{"error": st.Message()})
			return
		}
		c.Status(http.StatusNoContent)
	}
}

func CountThreadLikes(client threadpb.ThreadServiceClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		threadID := c.Param("id")
		if threadID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "thread id required"})
			return
		}
		resp, err := client.CountThreadLikes(c.Request.Context(), &threadpb.ThreadRequest{ThreadId: threadID})
		if err != nil {
			st := status.Convert(err)
			c.JSON(codesToHTTP(st.Code()), gin.H{"error": st.Message()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"count": resp.Count})
	}
}

func BookmarkThread(client threadpb.ThreadServiceClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		threadID := c.Param("id")
		if threadID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "thread id required"})
			return
		}
		var body struct {
			UserID string `json:"user_id" binding:"required"`
		}
		if err := c.ShouldBindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		_, err := client.BookmarkThread(c.Request.Context(), &threadpb.ThreadUserRequest{
			ThreadId: threadID,
			UserId:   body.UserID,
		})
		if err != nil {
			st := status.Convert(err)
			c.JSON(codesToHTTP(st.Code()), gin.H{"error": st.Message()})
			return
		}
		c.Status(http.StatusNoContent)
	}
}

func UnbookmarkThread(client threadpb.ThreadServiceClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		threadID := c.Param("id")
		if threadID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "thread id required"})
			return
		}
		var body struct {
			UserID string `json:"user_id" binding:"required"`
		}
		if err := c.ShouldBindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		_, err := client.UnbookmarkThread(c.Request.Context(), &threadpb.ThreadUserRequest{
			ThreadId: threadID,
			UserId:   body.UserID,
		})
		if err != nil {
			st := status.Convert(err)
			c.JSON(codesToHTTP(st.Code()), gin.H{"error": st.Message()})
			return
		}
		c.Status(http.StatusNoContent)
	}
}

func CountThreadBookmarks(client threadpb.ThreadServiceClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		threadID := c.Param("id")
		if threadID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "thread id required"})
			return
		}
		resp, err := client.CountThreadBookmarks(c.Request.Context(), &threadpb.ThreadRequest{ThreadId: threadID})
		if err != nil {
			st := status.Convert(err)
			c.JSON(codesToHTTP(st.Code()), gin.H{"error": st.Message()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"count": resp.Count})
	}
}

func CountThreadReplies(client threadpb.ThreadServiceClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		threadID := c.Param("id")
		if threadID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "thread id required"})
			return
		}
		resp, err := client.CountThreadReplies(c.Request.Context(), &threadpb.ThreadRequest{ThreadId: threadID})
		if err != nil {
			st := status.Convert(err)
			c.JSON(codesToHTTP(st.Code()), gin.H{"error": st.Message()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"count": resp.Count})
	}
}

func ListThreadMedia(client threadpb.ThreadServiceClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		threadID := c.Param("id")
		if threadID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "thread id required"})
			return
		}
		resp, err := client.ListThreadMedia(c.Request.Context(), &threadpb.ThreadRequest{ThreadId: threadID})
		if err != nil {
			st := status.Convert(err)
			c.JSON(codesToHTTP(st.Code()), gin.H{"error": st.Message()})
			return
		}
		c.JSON(http.StatusOK, resp.Media)
	}
}

func AddThreadMedia(client threadpb.ThreadServiceClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		threadID := c.Param("id")
		if threadID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "thread id required"})
			return
		}
		var body struct {
			MediaID string `json:"media_id" binding:"required"`
		}
		if err := c.ShouldBindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		_, err := client.AddThreadMedia(c.Request.Context(), &threadpb.ThreadMediaRequest{
			ThreadId: threadID,
			MediaId:  body.MediaID,
		})
		if err != nil {
			st := status.Convert(err)
			c.JSON(codesToHTTP(st.Code()), gin.H{"error": st.Message()})
			return
		}
		c.Status(http.StatusNoContent)
	}
}

func RemoveThreadMedia(client threadpb.ThreadServiceClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		threadID := c.Param("id")
		if threadID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "thread id required"})
			return
		}
		var body struct {
			MediaID string `json:"media_id" binding:"required"`
		}
		if err := c.ShouldBindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		_, err := client.RemoveThreadMedia(c.Request.Context(), &threadpb.ThreadMediaRequest{
			ThreadId: threadID,
			MediaId:  body.MediaID,
		})
		if err != nil {
			st := status.Convert(err)
			c.JSON(codesToHTTP(st.Code()), gin.H{"error": st.Message()})
			return
		}
		c.Status(http.StatusNoContent)
	}
}

func CountThreadMedia(client threadpb.ThreadServiceClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		threadID := c.Param("id")
		if threadID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "thread id required"})
			return
		}
		resp, err := client.CountThreadMedia(c.Request.Context(), &threadpb.ThreadRequest{ThreadId: threadID})
		if err != nil {
			st := status.Convert(err)
			c.JSON(codesToHTTP(st.Code()), gin.H{"error": st.Message()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"count": resp.Count})
	}
}

func HasLikedThread(client threadpb.ThreadServiceClient) gin.HandlerFunc {
  return func(c *gin.Context) {
    threadID := c.Param("id")
    if threadID == "" {
      c.JSON(http.StatusBadRequest, gin.H{"error": "thread id required"})
      return
    }
    var body struct {
      UserID string `json:"user_id" binding:"required"`
    }
    if err := c.ShouldBindJSON(&body); err != nil {
      c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
      return
    }

    resp, err := client.HasLikedThread(c.Request.Context(), &threadpb.ThreadUserRequest{
      ThreadId: threadID,
      UserId:   body.UserID,
    })
    if err != nil {
      st := status.Convert(err)
      c.JSON(codesToHTTP(st.Code()), gin.H{"error": st.Message()})
      return
    }
    c.JSON(http.StatusOK, gin.H{"liked": resp.Exists})
  }
}

func HasBookmarkedThread(client threadpb.ThreadServiceClient) gin.HandlerFunc {
  return func(c *gin.Context) {
    threadID := c.Param("id")
    if threadID == "" {
      c.JSON(http.StatusBadRequest, gin.H{"error": "thread id required"})
      return
    }
    var body struct {
      UserID string `json:"user_id" binding:"required"`
    }
    if err := c.ShouldBindJSON(&body); err != nil {
      c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
      return
    }

    resp, err := client.HasBookmarkedThread(c.Request.Context(), &threadpb.ThreadUserRequest{
      ThreadId: threadID,
      UserId:   body.UserID,
    })
    if err != nil {
      st := status.Convert(err)
      c.JSON(codesToHTTP(st.Code()), gin.H{"error": st.Message()})
      return
    }
    c.JSON(http.StatusOK, gin.H{"bookmarked": resp.Exists})
  }
}
func ListReplies(client threadpb.ThreadServiceClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			ParentID string `form:"parent_id" binding:"required"`
			Page int32 `form:"page" binding:"required"`
			Size int32 `form:"size" binding:"required"`
		}
		if err := c.ShouldBindQuery(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		
		resp, err := client.ListReplies(c.Request.Context(), &threadpb.ListRepliesRequest{
			Page: req.Page,
			Size: req.Size,
			ParentId: req.ParentID,
		})
		if err != nil {
			st := status.Convert(err)
			c.JSON(codesToHTTP(st.Code()), gin.H{"error": st.Message()})
			return
		}

		var threads []ThreadJSON
		for _, t := range resp.Threads {
			threads = append(threads, toThreadJSON(t))
		}
		c.JSON(http.StatusOK, threads)
	}
}

func ListRepliesByUserId(client threadpb.ThreadServiceClient) gin.HandlerFunc {
    return func(c *gin.Context) {
        userID := c.Param("user_id")
        if userID == "" {
            c.JSON(http.StatusBadRequest, gin.H{"error": "user_id required"})
            return
        }

        resp, err := client.ListRepliesByUser(c.Request.Context(), &threadpb.GetThreadByIDRequest{Id: userID})
        if err != nil {
            st := status.Convert(err)
            c.JSON(codesToHTTP(st.Code()), gin.H{"error": st.Message()})
            return
        }

        type RepliesGroupJSON struct {
            Parent  ThreadJSON   `json:"parent_thread"`
            Replies []ThreadJSON `json:"replies"`
        }
        var result []RepliesGroupJSON

        for _, grp := range resp.UserReplies {
            parent := toThreadJSON(grp.ParentThread)
            var replies []ThreadJSON
            for _, rep := range grp.Replies {
                replies = append(replies, toThreadJSON(rep))
            }
            result = append(result, RepliesGroupJSON{
                Parent:  parent,
                Replies: replies,
            })
        }

        sort.Slice(result, func(i, j int) bool {
			if result[i].Parent.IsPinned != result[j].Parent.IsPinned {
				return result[i].Parent.IsPinned
			}
			hasPinnedReplyI := false
			for _, reply := range result[i].Replies {
				if reply.IsPinned {
					hasPinnedReplyI = true
					break
				}
			}
			hasPinnedReplyJ := false
			for _, reply := range result[j].Replies {
				if reply.IsPinned {
					hasPinnedReplyJ = true
					break
				}
			}
			if hasPinnedReplyI != hasPinnedReplyJ {
				return hasPinnedReplyI
			}
			return false
		})
		
        c.JSON(http.StatusOK, result)
    }
}

type ThreadMediaJSON struct {
    ThreadID  string    `json:"thread_id"`
    ImageUrl  string    `json:"image_url"`
    Extension string    `json:"extension"`
    CreatedAt TimeJSON  `json:"created_at"`
}


func GetMediaByUser(client threadpb.ThreadServiceClient) gin.HandlerFunc {
    return func(c *gin.Context) {
        userID := c.Param("user_id")
        if userID == "" {
            c.JSON(http.StatusBadRequest, gin.H{"error": "user_id required"})
            return
        }

        resp, err := client.GetMediaByUser(
            c.Request.Context(),
            &threadpb.GetThreadByIDRequest{Id: userID},
        )
        if err != nil {
            st := status.Convert(err)
            c.JSON(codesToHTTP(st.Code()), gin.H{"error": st.Message()})
            return
        }

        var result []GroupedMediaJSON
        for _, group := range resp.Groups {
            threadJson := toThreadJSON(group.Thread)
            var mediaList []ThreadMediaJSON
            for _, m := range group.MediaList {
                created := m.GetCreatedAt()
                mediaList = append(mediaList, ThreadMediaJSON{
                    ThreadID:  m.GetThreadId(),
                    ImageUrl:  m.GetImageUrl(),
                    Extension: m.GetExtension(),
                    CreatedAt: TimeJSON{
                        Year:     created.GetYear(),
                        Month:    created.GetMonth(),
                        Day:      created.GetDay(),
                        Hour:     created.GetHour(),
                        Minute:   created.GetMinute(),
                        Second:   created.GetSecond(),
                        Timezone: created.GetTimezone(),
                    },
                })
            }
            result = append(result, GroupedMediaJSON{
                Thread:    threadJson,
                MediaList: mediaList,
            })
        }

        c.JSON(http.StatusOK, result)
    }
}

func ListThreadsByUserFollower(client threadpb.ThreadServiceClient) gin.HandlerFunc {
    return func(c *gin.Context) {
        userID := c.Query("user_id")
        pageStr := c.DefaultQuery("page", "1")
        sizeStr := c.DefaultQuery("size", "10")

        if userID == "" {
            c.JSON(http.StatusBadRequest, gin.H{"error": "user_id required"})
            return
        }

        page, err := strconv.Atoi(pageStr)
        if err != nil || page < 1 {
            c.JSON(http.StatusBadRequest, gin.H{"error": "invalid page"})
            return
        }
        size, err := strconv.Atoi(sizeStr)
        if err != nil || size < 1 {
            c.JSON(http.StatusBadRequest, gin.H{"error": "invalid size"})
            return
        }

        resp, err := client.ListThreadsByUserFollower(c.Request.Context(), &threadpb.ListByUserRequest{
            UserId: userID,
            Page:   int32(page),
            Size:   int32(size),
        })
        if err != nil {
            st := status.Convert(err)
            c.JSON(codesToHTTP(st.Code()), gin.H{"error": st.Message()})
            return
        }

        threads := make([]ThreadJSON, 0, len(resp.Threads))
        for _, t := range resp.Threads {
            threads = append(threads, toThreadJSON(t))
        }
        c.JSON(http.StatusOK, threads)
    }
}

func CreateThreadCategory(client threadpb.ThreadServiceClient) gin.HandlerFunc {
    return func(c *gin.Context) {
        var req struct {
            CategoryName string `json:"category_name" binding:"required"`
        }
        if err := c.ShouldBindJSON(&req); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }
        resp, err := client.CreateCategory(c.Request.Context(), &threadpb.CreateCategoryRequest{
            CategoryName: req.CategoryName,
        })
        if err != nil {
            st := status.Convert(err)
            c.JSON(codesToHTTP(st.Code()), gin.H{"error": st.Message()})
            return
        }
        c.JSON(http.StatusOK, gin.H{"category_id": resp.CategoryId})
    }
}

func ListThreadCategories(client threadpb.ThreadServiceClient) gin.HandlerFunc {
    return func(c *gin.Context) {
        resp, err := client.ListCategories(c.Request.Context(), &threadpb.Empty{})
        if err != nil {
            st := status.Convert(err)
            c.JSON(codesToHTTP(st.Code()), gin.H{"error": st.Message()})
            return
        }
        cats := make([]CategoryJSON, 0, len(resp.Categories))
        for _, cat := range resp.Categories {
            cats = append(cats, toThreadCategoryJSON(cat))
        }
        c.JSON(http.StatusOK, cats)
    }
}

func ListThreadsByCategory(client threadpb.ThreadServiceClient) gin.HandlerFunc {
    return func(c *gin.Context) {
        categoryId := c.Param("category_id")
        if categoryId == "" {
            c.JSON(http.StatusBadRequest, gin.H{"error": "category_id required"})
            return
        }
        var req struct {
            Page int32 `form:"page" binding:"required"`
            Size int32 `form:"size" binding:"required"`
        }
        if err := c.ShouldBindQuery(&req); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }
        resp, err := client.ListThreadsByCategory(c.Request.Context(), &threadpb.ListThreadsByCategoryRequest{
            CategoryId: categoryId,
            Page: req.Page,
            Size: req.Size,
        })
        if err != nil {
            st := status.Convert(err)
            c.JSON(codesToHTTP(st.Code()), gin.H{"error": st.Message()})
            return
        }
        var threads []ThreadJSON
        for _, t := range resp.Threads {
            threads = append(threads, toThreadJSON(t))
        }
        c.JSON(http.StatusOK, threads)
    }
}

func ListThreadsByCommunityByLike(client threadpb.ThreadServiceClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		communityID := c.Param("community_id")
		if communityID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "community_id required"})
			return
		}
		var req struct {
			Page int32 `form:"page" binding:"required"`
			Size int32 `form:"size" binding:"required"`
		}
		if err := c.ShouldBindQuery(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		resp, err := client.ListThreadsByCommunityByLike(c.Request.Context(), &threadpb.ListThreadsByCommunityRequest{
			CommunityId: communityID,
			Page:        req.Page,
			Size:        req.Size,
		})
		if err != nil {
			st := status.Convert(err)
			c.JSON(codesToHTTP(st.Code()), gin.H{"error": st.Message()})
			return
		}
		var threads []ThreadJSON
		for _, t := range resp.Threads {
			threads = append(threads, toThreadJSON(t))
		}
		c.JSON(http.StatusOK, threads)
	}
}
func ListThreadsByCommunityByTime(client threadpb.ThreadServiceClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		communityID := c.Param("community_id")
		if communityID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "community_id required"})
			return
		}
		var req struct {
			Page int32 `form:"page" binding:"required"`
			Size int32 `form:"size" binding:"required"`
		}
		if err := c.ShouldBindQuery(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		resp, err := client.ListThreadsByCommunityByTime(c.Request.Context(), &threadpb.ListThreadsByCommunityRequest{
			CommunityId: communityID,
			Page:        req.Page,
			Size:        req.Size,
		})
		if err != nil {
			st := status.Convert(err)
			c.JSON(codesToHTTP(st.Code()), gin.H{"error": st.Message()})
			return
		}
		var threads []ThreadJSON
		for _, t := range resp.Threads {
			threads = append(threads, toThreadJSON(t))
		}
		c.JSON(http.StatusOK, threads)
	}
}
type GroupedMediaJSON struct {
	Thread    ThreadJSON        `json:"thread"`
	MediaList []ThreadMediaJSON `json:"media_list"`
}

func GetMediaByCommunity(client threadpb.ThreadServiceClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		communityID := c.Param("community_id")
		if communityID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "community_id required"})
			return
		}
		var req struct {
			Page int32 `form:"page" binding:"required"`
			Size int32 `form:"size" binding:"required"`
		}
		if err := c.ShouldBindQuery(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		resp, err := client.GetMediaByCommunity(
			c.Request.Context(),
			&threadpb.GetMediaByCommunityRequest{
				CommunityId: communityID,
				Page:        req.Page,
				Size:        req.Size,
			},
		)
		if err != nil {
			st := status.Convert(err)
			c.JSON(codesToHTTP(st.Code()), gin.H{"error": st.Message()})
			return
		}

		var result []GroupedMediaJSON
		for _, group := range resp.Groups {
			threadJson := toThreadJSON(group.Thread)
			var mediaList []ThreadMediaJSON
			for _, m := range group.MediaList {
				created := m.GetCreatedAt()
				mediaList = append(mediaList, ThreadMediaJSON{
					ThreadID:  m.GetThreadId(),
					ImageUrl:  m.GetImageUrl(),
					Extension: m.GetExtension(),
					CreatedAt: TimeJSON{
						Year:     created.GetYear(),
						Month:    created.GetMonth(),
						Day:      created.GetDay(),
						Hour:     created.GetHour(),
						Minute:   created.GetMinute(),
						Second:   created.GetSecond(),
						Timezone: created.GetTimezone(),
					},
				})
			}
			result = append(result, GroupedMediaJSON{
				Thread:    threadJson,
				MediaList: mediaList,
			})
		}

		c.JSON(http.StatusOK, result)
	}
}
func CountThreadRepost(client threadpb.ThreadServiceClient) gin.HandlerFunc {
    return func(c *gin.Context) {
        threadID := c.Param("id")
        if threadID == "" {
            c.JSON(http.StatusBadRequest, gin.H{"error": "thread id required"})
            return
        }
        resp, err := client.CountRepost(c.Request.Context(), &threadpb.ThreadRequest{ThreadId: threadID})
        if err != nil {
            st := status.Convert(err)
            c.JSON(codesToHTTP(st.Code()), gin.H{"error": st.Message()})
            return
        }
        c.JSON(http.StatusOK, gin.H{"count": resp.Count})
    }
}

func HasRepostedThread(client threadpb.ThreadServiceClient) gin.HandlerFunc {
    return func(c *gin.Context) {
        threadID := c.Param("id")
        if threadID == "" {
            c.JSON(http.StatusBadRequest, gin.H{"error": "thread id required"})
            return
        }
        var body struct {
            UserID string `json:"user_id" binding:"required"`
        }
        if err := c.ShouldBindJSON(&body); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }
        resp, err := client.HasReposted(c.Request.Context(), &threadpb.ThreadUserRequest{
            ThreadId: threadID,
            UserId:   body.UserID,
        })
        if err != nil {
            st := status.Convert(err)
            c.JSON(codesToHTTP(st.Code()), gin.H{"error": st.Message()})
            return
        }
        c.JSON(http.StatusOK, gin.H{"reposted": resp.Exists})
    }
}

func PinThread(client threadpb.ThreadServiceClient) gin.HandlerFunc{
	return func( c *gin.Context){
		threadID := c.Param("id")
		if threadID == "" {
            c.JSON(http.StatusBadRequest, gin.H{"error": "thread id required"})
            return
        }

		_,err := client.PinThread(c.Request.Context(), &threadpb.ThreadRequest{
			ThreadId: threadID,
		})
		if err != nil{
			st := status.Convert(err)
            c.JSON(codesToHTTP(st.Code()), gin.H{"error": st.Message()})
            return
		}
		c.JSON(http.StatusOK, gin.H{"ok" : "ok"})
	}
}

func UpdateThreadCategory(client threadpb.ThreadServiceClient) gin.HandlerFunc {
    return func(c *gin.Context) {
        categoryId := c.Param("category_id")
        if categoryId == "" {
            c.JSON(http.StatusBadRequest, gin.H{"error": "category_id required"})
            return
        }
        var req struct {
            CategoryName string `json:"category_name" binding:"required"`
        }
        if err := c.ShouldBindJSON(&req); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }

        _, err := client.UpdateCategory(c.Request.Context(), &threadpb.UpdateCategoryRequest{
            CategoryId: categoryId,
            CategoryName: req.CategoryName,
        })
        if err != nil {
            st := status.Convert(err)
            c.JSON(codesToHTTP(st.Code()), gin.H{"error": st.Message()})
            return
        }
        c.Status(http.StatusNoContent)
    }
}

func DeleteThreadCategory(client threadpb.ThreadServiceClient) gin.HandlerFunc {
    return func(c *gin.Context) {
        categoryId := c.Param("category_id")
        if categoryId == "" {
            c.JSON(http.StatusBadRequest, gin.H{"error": "category_id required"})
            return
        }

        _, err := client.DeleteCategory(c.Request.Context(), &threadpb.DeleteCategoryRequest{
            CategoryId: categoryId,
        })
        if err != nil {
            st := status.Convert(err)
            c.JSON(codesToHTTP(st.Code()), gin.H{"error": st.Message()})
            return
        }
        c.Status(http.StatusNoContent)
    }
}

func GetThreadsByHashtag(client threadpb.ThreadServiceClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		hashtag := c.Query("hashtag")
		if hashtag == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "hashtag is required"})
			return
		}
		pageStr := c.DefaultQuery("page", "1")
		sizeStr := c.DefaultQuery("size", "20")
		page, err := strconv.Atoi(pageStr)
		if err != nil || page < 1 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid page"})
			return
		}
		size, err := strconv.Atoi(sizeStr)
		if err != nil || size < 1 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid size"})
			return
		}

		resp, err := client.GetThreadsByHashtag(c.Request.Context(), &threadpb.GetThreadsByHashtagRequest{
			HashtagName: hashtag,
			Page:    int32(page),
			Size:    int32(size),
		})
		if err != nil {
			st := status.Convert(err)
			c.JSON(codesToHTTP(st.Code()), gin.H{"error": st.Message()})
			return
		}
		threads := make([]ThreadJSON, 0, len(resp.Threads))
		for _, t := range resp.Threads {
			threads = append(threads, toThreadJSON(t))
		}
		c.JSON(http.StatusOK, threads)
	}
}

func GetTopHashtags(client threadpb.ThreadServiceClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		limitStr := c.DefaultQuery("limit", "10")
		limit, err := strconv.Atoi(limitStr)
		if err != nil || limit < 1 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid limit"})
			return
		}
		resp, err := client.GetTopHashtags(c.Request.Context(), &threadpb.GetTopHashtagsRequest{
			Limit: int32(limit),
		})
		if err != nil {
			st := status.Convert(err)
			c.JSON(codesToHTTP(st.Code()), gin.H{"error": st.Message()})
			return
		}
		type HashtagJSON struct {
			Hashtag string `json:"hashtag"`
			Count   int64  `json:"count"`
		}
		var result []HashtagJSON
		for _, h := range resp.Hashtags {
			result = append(result, HashtagJSON{
				Hashtag: h.HashtagName,
				Count:   h.Count,
			})
		}
		c.JSON(http.StatusOK, result)
	}
}

func SearchThreadByContent(client threadpb.ThreadServiceClient) gin.HandlerFunc {
    return func(c *gin.Context) {
        var req struct {
            Keyword   string `form:"keyword" binding:"required"`
            Threshold int    `form:"threshold" binding:"required"`
            Page      int    `form:"page" binding:"required"`
            Size      int    `form:"size" binding:"required"`
        }
        if err := c.ShouldBindQuery(&req); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }

        resp, err := client.SearchThreadsByContent(c.Request.Context(), &threadpb.SearchThreadsByContentRequest{
            Keyword:   req.Keyword,
            Threshold: int32(req.Threshold),
            Page:      int32(req.Page),
            Size:      int32(req.Size),
        })
        if err != nil {
            st := status.Convert(err)
            c.JSON(codesToHTTP(st.Code()), gin.H{"error": st.Message()})
            return
        }

        threads := make([]ThreadJSON, 0, len(resp.Threads))
        for _, t := range resp.Threads {
            threads = append(threads, toThreadJSON(t))
        }

        c.JSON(http.StatusOK, gin.H{
            "threads": threads,
            "total":   resp.Total,
        })
    }
}
type SearchMediaGroupedMediaJSON struct {
	Thread    ThreadJSON        `json:"thread"`
	MediaList []ThreadMediaJSON `json:"media_list"`
}

func SearchMediaByContent(client threadpb.ThreadServiceClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			Keyword   string `form:"keyword" binding:"required"`
			Threshold int    `form:"threshold" binding:"required"`
			Page      int    `form:"page" binding:"required"`
			Size      int    `form:"size" binding:"required"`
		}
		if err := c.ShouldBindQuery(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		resp, err := client.SearchMediaByContent(c.Request.Context(), &threadpb.SearchMediaByContentRequest{
			Keyword:   req.Keyword,
			Threshold: int32(req.Threshold),
			Page:      int32(req.Page),
			Size:      int32(req.Size),
		})
		if err != nil {
			st := status.Convert(err)
			c.JSON(codesToHTTP(st.Code()), gin.H{"error": st.Message()})
			return
		}

		var groups []SearchMediaGroupedMediaJSON
		for _, group := range resp.Groups {
			thread := toThreadJSON(group.Thread)
			var mediaList []ThreadMediaJSON
			for _, m := range group.MediaList {
				created := m.GetCreatedAt()
				mediaList = append(mediaList, ThreadMediaJSON{
					ThreadID:  m.GetThreadId(),
					ImageUrl:  m.GetImageUrl(),
					Extension: m.GetExtension(),
					CreatedAt: TimeJSON{
						Year:     created.GetYear(),
						Month:    created.GetMonth(),
						Day:      created.GetDay(),
						Hour:     created.GetHour(),
						Minute:   created.GetMinute(),
						Second:   created.GetSecond(),
						Timezone: created.GetTimezone(),
					},
				})
			}
			groups = append(groups, SearchMediaGroupedMediaJSON{
				Thread:    thread,
				MediaList: mediaList,
			})
		}

		c.JSON(http.StatusOK, gin.H{
			"groups": groups,
			"total":  resp.Total,
		})
	}
}

func codesToHTTP(c codes.Code) int {
	switch c {
	case codes.OK:
		return http.StatusOK
	case codes.InvalidArgument:
		return http.StatusBadRequest
	case codes.NotFound:
		return http.StatusNotFound
	case codes.AlreadyExists:
		return http.StatusConflict
	case codes.Unauthenticated:
		return http.StatusUnauthorized
	case codes.PermissionDenied:
		return http.StatusForbidden
	case codes.ResourceExhausted:
		return http.StatusTooManyRequests
	case codes.Unimplemented:
		return http.StatusNotImplemented
	case codes.Unavailable:
		return http.StatusServiceUnavailable
	default:
		return http.StatusInternalServerError
	}
}