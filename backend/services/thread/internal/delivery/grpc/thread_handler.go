package handler

import (
	"context"
	"log"
	"time"

	mediapb "github.com/williamsebastianliman/WEB-WS-242/proto/gen/go/mediapb/proto/media"
	threadpb "github.com/williamsebastianliman/WEB-WS-242/proto/gen/go/threadpb/proto/thread"
	userpb "github.com/williamsebastianliman/WEB-WS-242/proto/gen/go/userpb/proto/user"
	"github.com/williamsebastianliman/WEB-WS-242/services/thread/internal/model"
	"github.com/williamsebastianliman/WEB-WS-242/services/thread/internal/service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ThreadHandler struct {
	threadpb.UnimplementedThreadServiceServer
	svc service.ThreadService
	mediaClient mediapb.MediaServiceClient
	userClient userpb.UserServiceClient
}

func toPbCategory(cat *model.ThreadCategories) *threadpb.Category {
    return &threadpb.Category{
        CategoryId:   cat.CategoryID.String(),
        CategoryName: cat.CategoryName,
    }
}

func NewThreadHandler(svc service.ThreadService, mediaClient mediapb.MediaServiceClient, userClient userpb.UserServiceClient) *ThreadHandler {
	return &ThreadHandler{
        svc:         svc,
        mediaClient: mediaClient,
		userClient: userClient,
    }
}

func toProtoTimestamp(t time.Time) *threadpb.Timestamp {
	return &threadpb.Timestamp{
		Year:      int32(t.Year()),
		Month:     int32(t.Month()),
		Day:       int32(t.Day()),
		Hour:      int32(t.Hour()),
		Minute:    int32(t.Minute()),
		Second:    int32(t.Second()),
		Timezone:  t.Location().String(),
	}
}

func (h *ThreadHandler) CreateThread(ctx context.Context, req *threadpb.CreateThreadRequest) (*threadpb.CreateThreadResponse, error) {
	var comm *string
	if req.CommunityId != "" {
		comm = &req.CommunityId
	}

	var par *string
	if req.ParentId != "" {
		par = &req.ParentId
	}

	id, err := h.svc.CreateThread(ctx, req.UserId, req.Content, comm, par)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "%v", err)
	}
	return &threadpb.CreateThreadResponse{Id: id}, nil
}

func (h *ThreadHandler) CreateRepost(ctx context.Context, req *threadpb.RepostThreadRequest) (*threadpb.CreateThreadResponse, error) {
	var comm *string
	if req.CommunityId != "" {
		comm = &req.CommunityId
	}

	var par *string
	if req.ParentId != "" {
		par = &req.ParentId
	}

	var repost *string
	if req.RepostId != "" {
		repost = &req.RepostId
	}

	id, err := h.svc.CreateRepost(ctx, req.UserId, req.Content, comm, par, repost)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "%v", err)
	}
	return &threadpb.CreateThreadResponse{Id: id}, nil
}

func (h *ThreadHandler) GetThreadByID(ctx context.Context, req *threadpb.GetThreadByIDRequest) (*threadpb.GetThreadByIDResponse, error) {
	t, err := h.svc.GetThreadByID(ctx, req.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "%v", err)
	}
	if t == nil {
		return nil, status.Error(codes.NotFound, "thread not found")
	}
	var comm string
	if t.CommunityID != nil {
		comm = t.CommunityID.String()
	}
	var rpid string
	if t.RepostOfID != nil {
		rpid = t.RepostOfID.String()
	}
	return &threadpb.GetThreadByIDResponse{Thread: &threadpb.Thread{
		Id:          t.ID.String(),
		UserId:      t.UserID.String(),
		CommunityId: comm,
		Content:     t.Content,
		LikeCount:   t.LikeCount,
		CommentCount: t.CommentCount,
		ShareCount:  t.ShareCount,
		ViewCount:   t.ViewCount,
		RepostId: rpid,
		CreatedAt:   toProtoTimestamp(t.CreatedAt),
		UpdatedAt:   toProtoTimestamp(t.UpdatedAt),
	}}, nil
}

func (h *ThreadHandler) UpdateThread(ctx context.Context, req *threadpb.UpdateThreadRequest) (*threadpb.Empty, error) {
	var contentPtr *string
	if req.Content != "" {
		contentPtr = &req.Content
	}
	var commPtr *string
	if req.CommunityId != "" {
		commPtr = &req.CommunityId
	}
	err := h.svc.UpdateThread(ctx, req.Id, contentPtr, commPtr)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "%v", err)
	}
	return &threadpb.Empty{}, nil
}

func (h *ThreadHandler) DeleteThread(ctx context.Context, req *threadpb.DeleteThreadRequest) (*threadpb.Empty, error) {
	err := h.svc.DeleteThread(ctx, req.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "%v", err)
	}
	return &threadpb.Empty{}, nil
}

func (h *ThreadHandler) ListThreads(ctx context.Context, req *threadpb.ListThreadsRequest) (*threadpb.ListThreadsResponse, error) {
	threads, err := h.svc.ListThreads(ctx, int(req.Page), int(req.Size))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "%v", err)
	}
	resp := make([]*threadpb.Thread, len(threads))
	for i, t := range threads {
		var comm string
		if t.CommunityID != nil {
			comm = t.CommunityID.String()
		}
		var rpid string
		if t.RepostOfID != nil {
			rpid = t.RepostOfID.String()
		}
		resp[i] = &threadpb.Thread{
			Id:          t.ID.String(),
			UserId:      t.UserID.String(),
			CommunityId: comm,
			Content:     t.Content,
			LikeCount:   t.LikeCount,
			CommentCount: t.CommentCount,
			ShareCount:  t.ShareCount,
			ViewCount:   t.ViewCount,
			RepostId: rpid,
			CreatedAt:   toProtoTimestamp(t.CreatedAt),
			UpdatedAt:   toProtoTimestamp(t.UpdatedAt),
		}
	}
	return &threadpb.ListThreadsResponse{Threads: resp}, nil
}

func (h *ThreadHandler) ListThreadsByUserFollower(ctx context.Context, req *threadpb.ListByUserRequest) (*threadpb.ListThreadsResponse, error) {
    usersResp, err := h.userClient.GetFollowing(ctx, &userpb.GetFollowingRequest{
        UserId: req.UserId,
    })
    if err != nil {
        return nil, status.Errorf(codes.Internal, "failed to get following: %v", err)
    }

    var followingIDs []string
    for _, u := range usersResp.Following {
        followingIDs = append(followingIDs, u.Id)
    }
    if len(followingIDs) == 0 {
        return &threadpb.ListThreadsResponse{Threads: []*threadpb.Thread{}}, nil
    }

    threads, err := h.svc.ListThreadsByUserIDs(ctx, followingIDs, int(req.Page), int(req.Size))
    if err != nil {
        return nil, status.Errorf(codes.Internal, "%v", err)
    }

    resp := make([]*threadpb.Thread, len(threads))
    for i, t := range threads {
        var comm string
        if t.CommunityID != nil {
            comm = t.CommunityID.String()
        }
		var rpid string
		if t.RepostOfID != nil {
			rpid = t.RepostOfID.String()
		}
        resp[i] = &threadpb.Thread{
            Id:          t.ID.String(),
            UserId:      t.UserID.String(),
            CommunityId: comm,
            Content:     t.Content,
            LikeCount:   t.LikeCount,
            CommentCount: t.CommentCount,
            ShareCount:  t.ShareCount,
            ViewCount:   t.ViewCount,
			RepostId: rpid,
            CreatedAt:   toProtoTimestamp(t.CreatedAt),
            UpdatedAt:   toProtoTimestamp(t.UpdatedAt),
        }
    }
    return &threadpb.ListThreadsResponse{Threads: resp}, nil
}

func (h *ThreadHandler) ListThreadsByUser(ctx context.Context, req *threadpb.ListByUserRequest) (*threadpb.ListThreadsResponse, error) {
	threads, err := h.svc.ListThreadsByUser(ctx, req.UserId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "%v", err)
	}
	resp := make([]*threadpb.Thread, len(threads))
	for i, t := range threads {
		var comm string
		if t.CommunityID != nil {
			comm = t.CommunityID.String()
		}
		var rpid string
		if t.RepostOfID != nil {
			rpid = t.RepostOfID.String()
		}
		// test:= toProtoTimestamp(t.CreatedAt)
		// log.Printf("Hour: %s", test.Hour)
		resp[i] = &threadpb.Thread{
			Id:          t.ID.String(),
			UserId:      t.UserID.String(),
			CommunityId: comm,
			Content:     t.Content,
			LikeCount:   t.LikeCount,
			CommentCount: t.CommentCount,
			ShareCount:  t.ShareCount,
			ViewCount:   t.ViewCount,
			RepostId: rpid,
			IsPinned: t.IsPinned,
			CreatedAt:   toProtoTimestamp(t.CreatedAt),
			UpdatedAt:   toProtoTimestamp(t.UpdatedAt),
		}
	}
	return &threadpb.ListThreadsResponse{Threads: resp},nil
}

func (h *ThreadHandler) ListThreadsByUserLike(ctx context.Context, req *threadpb.ListByUserRequest) (*threadpb.ListThreadsResponse, error) {
	threads, err := h.svc.ListThreadsByUserLike(ctx, req.UserId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "%v", err)
	}
	resp := make([]*threadpb.Thread, len(threads))
	for i, t := range threads {
		var comm string
		if t.CommunityID != nil {
			comm = t.CommunityID.String()
		}
		var rpid string
		if t.RepostOfID != nil {
			rpid = t.RepostOfID.String()
		}
		resp[i] = &threadpb.Thread{
			Id:          t.ID.String(),
			UserId:      t.UserID.String(),
			CommunityId: comm,
			Content:     t.Content,
			LikeCount:   t.LikeCount,
			CommentCount: t.CommentCount,
			ShareCount:  t.ShareCount,
			ViewCount:   t.ViewCount,
			RepostId: rpid,
			CreatedAt:   toProtoTimestamp(t.CreatedAt),
			UpdatedAt:   toProtoTimestamp(t.UpdatedAt),
		}
	}
	return &threadpb.ListThreadsResponse{Threads: resp},nil
}

func (h *ThreadHandler) ListThreadsByUserBookmark(ctx context.Context, req *threadpb.ListByUserRequest) (*threadpb.ListThreadsResponse, error) {
	threads, err := h.svc.ListThreadsByUserBookmark(ctx, req.UserId, req.Query)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "%v", err)
	}
	resp := make([]*threadpb.Thread, len(threads))
	for i, t := range threads {
		var comm string
		if t.CommunityID != nil {
			comm = t.CommunityID.String()
		}
		var rpid string
		if t.RepostOfID != nil {
			rpid = t.RepostOfID.String()
		}
		resp[i] = &threadpb.Thread{
			Id:          t.ID.String(),
			UserId:      t.UserID.String(),
			CommunityId: comm,
			Content:     t.Content,
			LikeCount:   t.LikeCount,
			CommentCount: t.CommentCount,
			ShareCount:  t.ShareCount,
			ViewCount:   t.ViewCount,
			RepostId: rpid,
			CreatedAt:   toProtoTimestamp(t.CreatedAt),
			UpdatedAt:   toProtoTimestamp(t.UpdatedAt),
		}
	}
	return &threadpb.ListThreadsResponse{Threads: resp},nil
}

func (h *ThreadHandler) ListThreadsByCommunity(ctx context.Context, req *threadpb.ListByCommunityRequest) (*threadpb.ListThreadsResponse, error) {
	_, err := h.svc.ListThreadsByCommunity(ctx, req.CommunityId, int(req.Page), int(req.Size))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "%v", err)
	}
	return h.ListThreads(ctx, &threadpb.ListThreadsRequest{Page: req.Page, Size: req.Size})
}

func (h *ThreadHandler) SearchThreads(ctx context.Context, req *threadpb.SearchThreadsRequest) (*threadpb.SearchThreadsResponse, error) {
	threads, err := h.svc.SearchThreads(ctx, req.Keyword, int(req.Page), int(req.Size))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "%v", err)
	}
	reps := make([]*threadpb.Thread, len(threads))
	for i, t := range threads {
		var comm string
		if t.CommunityID != nil {
			comm = t.CommunityID.String()
		}
		var rpid string
		if t.RepostOfID != nil {
			rpid = t.RepostOfID.String()
		}
		reps[i] = &threadpb.Thread{
			Id:          t.ID.String(),
			UserId:      t.UserID.String(),
			CommunityId: comm,
			Content:     t.Content,
			LikeCount:   t.LikeCount,
			CommentCount: t.CommentCount,
			ShareCount:  t.ShareCount,
			ViewCount:   t.ViewCount,
			RepostId: rpid,
			CreatedAt:   toProtoTimestamp(t.CreatedAt),
			UpdatedAt:   toProtoTimestamp(t.UpdatedAt),
		}
	}
	return &threadpb.SearchThreadsResponse{Threads: reps}, nil
}

func (h *ThreadHandler) IncrementView(ctx context.Context, req *threadpb.IncrementRequest) (*threadpb.Empty, error) {
	err := h.svc.IncrementView(ctx, req.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "%v", err)
	}
	return &threadpb.Empty{}, nil
}

func (h *ThreadHandler) IncrementLike(ctx context.Context, req *threadpb.IncrementRequest) (*threadpb.Empty, error) {
	err := h.svc.IncrementLike(ctx, req.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "%v", err)
	}
	return &threadpb.Empty{}, nil
}

func (h *ThreadHandler) IncrementComment(ctx context.Context, req *threadpb.IncrementRequest) (*threadpb.Empty, error) {
	err := h.svc.IncrementComment(ctx, req.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "%v", err)
	}
	return &threadpb.Empty{}, nil
}

func (h *ThreadHandler) IncrementShare(ctx context.Context, req *threadpb.IncrementRequest) (*threadpb.Empty, error) {
	err := h.svc.IncrementShare(ctx, req.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "%v", err)
	}
	return &threadpb.Empty{}, nil
}

func (h *ThreadHandler) CountThreads(ctx context.Context, _ *threadpb.Empty) (*threadpb.CountThreadsResponse, error) {
	count, err := h.svc.CountThreads(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "%v", err)
	}
	return &threadpb.CountThreadsResponse{Count: count}, nil
}

func (h *ThreadHandler) LikeThread(ctx context.Context, req *threadpb.ThreadUserRequest) (*threadpb.Empty, error) {
    if err := h.svc.LikeThread(ctx, req.ThreadId, req.UserId); err != nil {
        return nil, status.Errorf(codes.Internal, "%v", err)
    }
    return &threadpb.Empty{}, nil
}

func (h *ThreadHandler) UnlikeThread(ctx context.Context, req *threadpb.ThreadUserRequest) (*threadpb.Empty, error) {
    if err := h.svc.UnlikeThread(ctx, req.ThreadId, req.UserId); err != nil {
        return nil, status.Errorf(codes.Internal, "%v", err)
    }
    return &threadpb.Empty{}, nil
}

func (h *ThreadHandler) CountThreadLikes(ctx context.Context, req *threadpb.ThreadRequest) (*threadpb.CountThreadsResponse, error) {
    count, err := h.svc.CountThreadLikes(ctx, req.ThreadId)
    if err != nil {
        return nil, status.Errorf(codes.Internal, "%v", err)
    }
    return &threadpb.CountThreadsResponse{Count: count}, nil
}

func (h *ThreadHandler) BookmarkThread(ctx context.Context, req *threadpb.ThreadUserRequest) (*threadpb.Empty, error) {
    if err := h.svc.BookmarkThread(ctx, req.ThreadId, req.UserId); err != nil {
        return nil, status.Errorf(codes.Internal, "%v", err)
    }
    return &threadpb.Empty{}, nil
}

func (h *ThreadHandler) UnbookmarkThread(ctx context.Context, req *threadpb.ThreadUserRequest) (*threadpb.Empty, error) {
    if err := h.svc.UnbookmarkThread(ctx, req.ThreadId, req.UserId); err != nil {
        return nil, status.Errorf(codes.Internal, "%v", err)
    }
    return &threadpb.Empty{}, nil
}

func (h *ThreadHandler) CountThreadBookmarks(ctx context.Context, req *threadpb.ThreadRequest) (*threadpb.CountThreadsResponse, error) {
    count, err := h.svc.CountThreadBookmarks(ctx, req.ThreadId)
    if err != nil {
        return nil, status.Errorf(codes.Internal, "%v", err)
    }
    return &threadpb.CountThreadsResponse{Count: count}, nil
}

func (h *ThreadHandler) CountThreadReplies(ctx context.Context, req *threadpb.ThreadRequest) (*threadpb.CountThreadsResponse, error) {
    count, err := h.svc.CountThreadReplies(ctx, req.ThreadId)
    if err != nil {
        return nil, status.Errorf(codes.Internal, "%v", err)
    }
    return &threadpb.CountThreadsResponse{Count: count}, nil
}

func (h *ThreadHandler) ListThreadMedia(ctx context.Context, req *threadpb.ThreadRequest) (*threadpb.ListThreadMediaResponse, error) {
    records, err := h.svc.ListThreadMedia(ctx, req.ThreadId)
    if err != nil {
        return nil, status.Errorf(codes.Internal, "%v", err)
    }

    resp := make([]*threadpb.ThreadMedia, len(records))
    for i, rec := range records {
        mediaResp, err := h.mediaClient.GetMediaById(ctx, &mediapb.GetMediaByIdRequest{Id: rec.MediaID.String()})
        if err != nil {
            return nil, status.Errorf(codes.Internal, "media lookup failed: %v", err)
        }
        resp[i] = &threadpb.ThreadMedia{
            ThreadId:   req.ThreadId,
            ImageUrl:  mediaResp.PublicUrl,
			Extension: mediaResp.Extension,
            CreatedAt: toProtoTimestamp(rec.CreatedAt),
        }
    }
    return &threadpb.ListThreadMediaResponse{Media: resp}, nil
}
func (h *ThreadHandler) ListReplies(ctx context.Context, req *threadpb.ListRepliesRequest) (*threadpb.ListThreadsResponse, error){
	threads, err := h.svc.ListThreadsReplies(ctx, req.ParentId, int(req.Page), int(req.Size))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "%v", err)
	}
	resp := make([]*threadpb.Thread, len(threads))
	for i, t := range threads {
		var comm string
		if t.CommunityID != nil {
			comm = t.CommunityID.String()
		}
		var rpid string
		if t.RepostOfID != nil {
			rpid = t.RepostOfID.String()
		}
		resp[i] = &threadpb.Thread{
			Id:          t.ID.String(),
			UserId:      t.UserID.String(),
			CommunityId: comm,
			Content:     t.Content,
			LikeCount:   t.LikeCount,
			CommentCount: t.CommentCount,
			ShareCount:  t.ShareCount,
			ViewCount:   t.ViewCount,
			RepostId: rpid,
			CreatedAt:   toProtoTimestamp(t.CreatedAt),
			UpdatedAt:   toProtoTimestamp(t.UpdatedAt),
		}
	}
	return &threadpb.ListThreadsResponse{Threads: resp}, nil
}


func (h *ThreadHandler) AddThreadMedia(ctx context.Context, req *threadpb.ThreadMediaRequest) (*threadpb.Empty, error) {
    if err := h.svc.AddThreadMedia(ctx, req.ThreadId, req.MediaId); err != nil {
        return nil, status.Errorf(codes.Internal, "%v", err)
    }
    return &threadpb.Empty{}, nil
}

func (h *ThreadHandler) RemoveThreadMedia(ctx context.Context, req *threadpb.ThreadMediaRequest) (*threadpb.Empty, error) {
    if err := h.svc.RemoveThreadMedia(ctx, req.ThreadId, req.MediaId); err != nil {
        return nil, status.Errorf(codes.Internal, "%v", err)
    }
    return &threadpb.Empty{}, nil
}

func (h *ThreadHandler) CountThreadMedia(ctx context.Context, req *threadpb.ThreadRequest) (*threadpb.CountThreadsResponse, error) {
    count, err := h.svc.CountThreadMedia(ctx, req.ThreadId)
    if err != nil {
        return nil, status.Errorf(codes.Internal, "%v", err)
    }
    return &threadpb.CountThreadsResponse{Count: count}, nil
}

func (h *ThreadHandler) HasLikedThread(ctx context.Context, req *threadpb.ThreadUserRequest) (*threadpb.ExistsResponse, error) {
    exists, err := h.svc.HasLikedThread(ctx, req.ThreadId, req.UserId)
    if err != nil {
        return nil, status.Errorf(codes.Internal, "%v", err)
    }
    return &threadpb.ExistsResponse{Exists: exists}, nil
}

func (h *ThreadHandler) HasBookmarkedThread(ctx context.Context, req *threadpb.ThreadUserRequest) (*threadpb.ExistsResponse, error) {
    exists, err := h.svc.HasBookmarkedThread(ctx, req.ThreadId, req.UserId)
    if err != nil {
        return nil, status.Errorf(codes.Internal, "%v", err)
    }
    return &threadpb.ExistsResponse{Exists: exists}, nil
}

func (h *ThreadHandler) ListRepliesByUser(ctx context.Context, req *threadpb.GetThreadByIDRequest) (*threadpb.GetRepliesByIdResponse, error){
	grouped_thread, err := h.svc.ListRepliesByUser(ctx, req.Id)
	if err!=nil{
		return nil,err
	}
	i:=0
	resp := make([]*threadpb.RepliesGroup, len(grouped_thread))
	for key, val := range grouped_thread{

		replies_list := make([]*threadpb.Thread, len(val))
		for item_idx := range replies_list{
			reply := val[item_idx]
			var comm string
			if reply.CommunityID != nil {
				comm = reply.CommunityID.String()
			}
			replies_list[item_idx] = &threadpb.Thread{
				Id:          reply.ID.String(),
				UserId:      reply.UserID.String(),
				CommunityId: comm,
				Content:     reply.Content,
				LikeCount:   reply.LikeCount,
				CommentCount: reply.CommentCount,
				ShareCount:  reply.ShareCount,
				ViewCount:   reply.ViewCount,
				IsPinned: reply.IsPinned,
				CreatedAt:   toProtoTimestamp(reply.CreatedAt),
				UpdatedAt:   toProtoTimestamp(reply.UpdatedAt),
			}
		}
		var par_comm string
		if key.CommunityID != nil {
			par_comm = key.CommunityID.String()
		}
		par_thread := &threadpb.Thread{
			Id:          key.ID.String(),
			UserId:      key.UserID.String(),
			CommunityId: par_comm,
			Content:     key.Content,
			LikeCount:   key.LikeCount,
			CommentCount: key.CommentCount,
			ShareCount:  key.ShareCount,
			ViewCount:   key.ViewCount,
			IsPinned: key.IsPinned,
			CreatedAt:   toProtoTimestamp(key.CreatedAt),
			UpdatedAt:   toProtoTimestamp(key.UpdatedAt),
		}
		resp[i] = &threadpb.RepliesGroup{
			ParentThread: par_thread,
			Replies: replies_list,
		}
		i+=1
	}
	return &threadpb.GetRepliesByIdResponse{UserReplies: resp}, nil
}

func (h *ThreadHandler) GetMediaByUser(ctx context.Context, req *threadpb.GetThreadByIDRequest) (*threadpb.GetMediaByUserResponse, error) {
    grouped, err := h.svc.GetMediaByUserGrouped(req.Id)
    if err != nil {
        return nil, status.Errorf(codes.Internal, "get media by user: %v", err)
    }

    var out []*threadpb.GroupedMediaResponse
	

    for thread, medias := range grouped {
		var comm string
		if thread.CommunityID != nil {
			comm = thread.CommunityID.String()
		}
        threadpbObj := &threadpb.Thread{
            Id:           thread.ID.String(),
            UserId:       thread.UserID.String(),
            CommunityId:  comm,
            Content:      thread.Content,
            LikeCount:    thread.LikeCount,
            CommentCount: thread.CommentCount,
            ShareCount:   thread.ShareCount,
            ViewCount:    thread.ViewCount,
            CreatedAt:    toProtoTimestamp(thread.CreatedAt),
            UpdatedAt:    toProtoTimestamp(thread.UpdatedAt),
        }

        var mediaList []*threadpb.ThreadMedia
        for _, media := range medias {
            mediaResp, err := h.mediaClient.GetMediaById(ctx, &mediapb.GetMediaByIdRequest{Id: media.MediaID.String()})
            if err != nil {
                return nil, status.Errorf(codes.Internal, "media lookup failed: %v", err)
            }
            mediaList = append(mediaList, &threadpb.ThreadMedia{
                ThreadId:  media.ThreadID.String(),
                ImageUrl:  mediaResp.PublicUrl,
                Extension: mediaResp.Extension,
                CreatedAt: toProtoTimestamp(media.CreatedAt),
            })
        }

        out = append(out, &threadpb.GroupedMediaResponse{
            Thread:    threadpbObj,
            MediaList: mediaList,
        })
    }

    return &threadpb.GetMediaByUserResponse{
        Groups: out,
    }, nil
}

func (h *ThreadHandler) CreateCategory(
    ctx context.Context,
    req *threadpb.CreateCategoryRequest,
) (*threadpb.CreateCategoryResponse, error) {
    mod, err := h.svc.CreateCategory(ctx, req.CategoryName)
    if err != nil {
        return nil, status.Errorf(codes.Internal, "%v", err)
    }
    return &threadpb.CreateCategoryResponse{CategoryId: mod.CategoryID.String()}, nil
}

func (h *ThreadHandler) ListCategories(
    ctx context.Context,
    req *threadpb.Empty,
) (*threadpb.ListCategoriesResponse, error) {
    cats, err := h.svc.ListAllCategories(ctx)
    if err != nil {
        return nil, status.Errorf(codes.Internal, "%v", err)
    }
    resp := make([]*threadpb.Category, len(cats))
    for i, cat := range cats {
        resp[i] = toPbCategory(&cat)
    }
    return &threadpb.ListCategoriesResponse{Categories: resp}, nil
}

func (h *ThreadHandler) ListThreadsByCategory(
    ctx context.Context,
    req *threadpb.ListThreadsByCategoryRequest,
) (*threadpb.ListThreadsResponse, error) {
    threads, err := h.svc.ListThreadsByCategory(ctx, req.CategoryId, int(req.Page), int(req.Size))
    if err != nil {
        return nil, status.Errorf(codes.Internal, "%v", err)
    }
    resp := make([]*threadpb.Thread, len(threads))
    for i, t := range threads {
        var comm string
        if t.CommunityID != nil {
            comm = t.CommunityID.String()
        }
		var rpid string
		if t.RepostOfID != nil {
			rpid = t.RepostOfID.String()
		}
        resp[i] = &threadpb.Thread{
            Id:          t.ID.String(),
            UserId:      t.UserID.String(),
            CommunityId: comm,
            Content:     t.Content,
            LikeCount:   t.LikeCount,
            CommentCount: t.CommentCount,
            ShareCount:  t.ShareCount,
            ViewCount:   t.ViewCount,
			RepostId: rpid,
            CreatedAt:   toProtoTimestamp(t.CreatedAt),
            UpdatedAt:   toProtoTimestamp(t.UpdatedAt),
        }
    }
    return &threadpb.ListThreadsResponse{Threads: resp}, nil
}

func (h *ThreadHandler) ListThreadsByCommunityByLike(ctx context.Context, req *threadpb.ListThreadsByCommunityRequest) (*threadpb.ListThreadsByCommunityResponse, error) {
	threads, err := h.svc.ListThreadsByCommunityByLike(ctx, req.CommunityId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list threads by like: %v", err)
	}
	var resp threadpb.ListThreadsByCommunityResponse
	for _, t := range threads {
		var comm string
		if t.CommunityID != nil {
			comm = t.CommunityID.String()
		}
		var rpid string
		if t.RepostOfID != nil {
			rpid = t.RepostOfID.String()
		}
		resp.Threads = append(resp.Threads, &threadpb.Thread{
			Id:          t.ID.String(),
			UserId:      t.UserID.String(),
			CommunityId: comm,
			Content:     t.Content,
			LikeCount:   t.LikeCount,
			CommentCount: t.CommentCount,
			ShareCount:  t.ShareCount,
			ViewCount:   t.ViewCount,
			RepostId: rpid,
			CreatedAt:   toProtoTimestamp(t.CreatedAt),
			UpdatedAt:   toProtoTimestamp(t.UpdatedAt),
		})
	}
	return &resp, nil
}

func (h *ThreadHandler) ListThreadsByCommunityByTime(ctx context.Context, req *threadpb.ListThreadsByCommunityRequest) (*threadpb.ListThreadsByCommunityResponse, error) {
	threads, err := h.svc.ListThreadsByCommunityByTime(ctx, req.CommunityId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list threads by time: %v", err)
	}
	var resp threadpb.ListThreadsByCommunityResponse
	for _, t := range threads {
		var comm string
		if t.CommunityID != nil {
			comm = t.CommunityID.String()
		}
		var rpid string
		if t.RepostOfID != nil {
			rpid = t.RepostOfID.String()
		}
		resp.Threads = append(resp.Threads, &threadpb.Thread{
			Id:          t.ID.String(),
			UserId:      t.UserID.String(),
			CommunityId: comm,
			Content:     t.Content,
			LikeCount:   t.LikeCount,
			CommentCount: t.CommentCount,
			ShareCount:  t.ShareCount,
			ViewCount:   t.ViewCount,
			RepostId: rpid,
			CreatedAt:   toProtoTimestamp(t.CreatedAt),
			UpdatedAt:   toProtoTimestamp(t.UpdatedAt),
		})
	}
	return &resp, nil
}

func (h *ThreadHandler) GetMediaByCommunity(
    ctx context.Context,
    req *threadpb.GetMediaByCommunityRequest,
) (*threadpb.GetMediaByCommunityResponse, error) {
    grouped, err := h.svc.GetMediaByCommunity(req.CommunityId, int(req.Page), int(req.Size))
    if err != nil {
        return nil, status.Errorf(codes.Internal, "get media by community: %v", err)
    }

    var out []*threadpb.GroupedMediaResponse

    for thread, medias := range grouped {
        var comm string
        if thread.CommunityID != nil {
            comm = thread.CommunityID.String()
        }
        threadpbObj := &threadpb.Thread{
            Id:           thread.ID.String(),
            UserId:       thread.UserID.String(),
            CommunityId:  comm,
            Content:      thread.Content,
            LikeCount:    thread.LikeCount,
            CommentCount: thread.CommentCount,
            ShareCount:   thread.ShareCount,
            ViewCount:    thread.ViewCount,
            CreatedAt:    toProtoTimestamp(thread.CreatedAt),
            UpdatedAt:    toProtoTimestamp(thread.UpdatedAt),
        }

        var mediaList []*threadpb.ThreadMedia
        for _, media := range medias {
            mediaResp, err := h.mediaClient.GetMediaById(ctx, &mediapb.GetMediaByIdRequest{Id: media.MediaID.String()})
            if err != nil {
                return nil, status.Errorf(codes.Internal, "media lookup failed: %v", err)
            }
            mediaList = append(mediaList, &threadpb.ThreadMedia{
                ThreadId:  media.ThreadID.String(),
                ImageUrl:  mediaResp.PublicUrl,
                Extension: mediaResp.Extension,
                CreatedAt: toProtoTimestamp(media.CreatedAt),
            })
        }

        out = append(out, &threadpb.GroupedMediaResponse{
            Thread:    threadpbObj,
            MediaList: mediaList,
        })
    }

    return &threadpb.GetMediaByCommunityResponse{
        Groups: out,
    }, nil
}

func (h *ThreadHandler) CountRepost(ctx context.Context, req *threadpb.ThreadRequest) (*threadpb.CountThreadsResponse, error) {
    count, err := h.svc.CountThreadRepost(ctx, req.ThreadId)
    if err != nil {
        return nil, status.Errorf(codes.Internal, "%v", err)
    }
    return &threadpb.CountThreadsResponse{Count: count}, nil
}

func (h *ThreadHandler) HasReposted(ctx context.Context, req *threadpb.ThreadUserRequest) (*threadpb.ExistsResponse, error) {
    exists, err := h.svc.HasRepostedThread(ctx, req.ThreadId, req.UserId)
    if err != nil {
        return nil, status.Errorf(codes.Internal, "%v", err)
    }
    return &threadpb.ExistsResponse{Exists: exists}, nil
}

func (h *ThreadHandler) PinThread(ctx context.Context, req *threadpb.ThreadRequest) (*threadpb.Empty, error) {
	err := h.svc.PinThread(req.ThreadId)
	if err != nil{
		log.Printf("error: %s",err)
		return nil, err
	}
	return &threadpb.Empty{},nil
}

func (h *ThreadHandler) UpdateCategory(
    ctx context.Context,
    req *threadpb.UpdateCategoryRequest,
) (*threadpb.Empty, error) {
    err := h.svc.UpdateCategory(ctx, req.CategoryId, req.CategoryName)
    if err != nil {
        return nil, status.Errorf(codes.Internal, "%v", err)
    }
    return &threadpb.Empty{}, nil
}

func (h *ThreadHandler) DeleteCategory(
    ctx context.Context,
    req *threadpb.DeleteCategoryRequest,
) (*threadpb.Empty, error) {
    err := h.svc.DeleteCategory(ctx, req.CategoryId)
    if err != nil {
        return nil, status.Errorf(codes.Internal, "%v", err)
    }
    return &threadpb.Empty{}, nil
}

func (h *ThreadHandler) GetThreadsByHashtag(ctx context.Context, req *threadpb.GetThreadsByHashtagRequest) (*threadpb.GetThreadsByHashtagResponse, error) {
    threads, err := h.svc.GetThreadsByHashtag(ctx, req.HashtagName, int(req.Page), int(req.Size))
    if err != nil {
        return nil, status.Errorf(codes.Internal, "failed to get threads by hashtag: %v", err)
    }
    resp := make([]*threadpb.Thread, len(threads))
    for i, t := range threads {
        var comm string
        if t.CommunityID != nil {
            comm = t.CommunityID.String()
        }
        var rpid string
        if t.RepostOfID != nil {
            rpid = t.RepostOfID.String()
        }
        resp[i] = &threadpb.Thread{
            Id:          t.ID.String(),
            UserId:      t.UserID.String(),
            CommunityId: comm,
            Content:     t.Content,
            LikeCount:   t.LikeCount,
            CommentCount: t.CommentCount,
            ShareCount:  t.ShareCount,
            ViewCount:   t.ViewCount,
            RepostId:    rpid,
            IsPinned:    t.IsPinned,
            CreatedAt:   toProtoTimestamp(t.CreatedAt),
            UpdatedAt:   toProtoTimestamp(t.UpdatedAt),
        }
    }
    return &threadpb.GetThreadsByHashtagResponse{Threads: resp}, nil
}

func (h *ThreadHandler) GetTopHashtags(ctx context.Context, req *threadpb.GetTopHashtagsRequest) (*threadpb.GetTopHashtagsResponse, error) {
    hashtags, err := h.svc.GetTopHashtags(ctx, int(req.Limit))
    if err != nil {
        return nil, status.Errorf(codes.Internal, "failed to get top hashtags: %v", err)
    }
    resp := make([]*threadpb.TopHashtag, len(hashtags))
    for i, tag := range hashtags {
        resp[i] = &threadpb.TopHashtag{
            HashtagName: tag.HashtagName,
            Count:       tag.Count,
        }
    }
    return &threadpb.GetTopHashtagsResponse{Hashtags: resp}, nil
}

func ModelToPBThread(t *model.Thread) *threadpb.Thread {

    pbThread := &threadpb.Thread{
        Id:           t.ID.String(),
        UserId:       t.UserID.String(),
        Content:      t.Content,
        LikeCount:    t.LikeCount,
        CommentCount: t.CommentCount,
        ShareCount:   t.ShareCount,
        ViewCount:    t.ViewCount,
        IsPinned:     t.IsPinned,
        CreatedAt:    toProtoTimestamp(t.CreatedAt),
        UpdatedAt:    toProtoTimestamp(t.UpdatedAt),
    }
    if t.RepostOfID != nil {
        pbThread.RepostId = t.RepostOfID.String()
    }
	if t.CommunityID != nil{
		pbThread.CommunityId = t.CommunityID.String()
	}
    return pbThread
}

func (s *ThreadHandler) SearchThreadsByContent(ctx context.Context, req *threadpb.SearchThreadsByContentRequest) (*threadpb.SearchThreadsByContentResponse, error) {
    threads, total, err := s.svc.SearchThreadByContent(req.Keyword, int(req.Threshold), int(req.Page), int(req.Size))
    if err != nil {
        return nil, status.Errorf(codes.Internal, "failed to search threads: %v", err)
    }

    var pbThreads []*threadpb.Thread
    for _, t := range threads {
        pbThreads = append(pbThreads, ModelToPBThread(&t))
    }

    return &threadpb.SearchThreadsByContentResponse{
        Threads: pbThreads,
        Total: int32(total),
    }, nil
}

func (h *ThreadHandler) SearchMediaByContent(
    ctx context.Context,
    req *threadpb.SearchMediaByContentRequest,
) (*threadpb.SearchMediaByContentResponse, error) {
    grouped, total, err := h.svc.SearchMediaByContent(
        req.Keyword,
        int(req.Threshold),
        int(req.Page),
        int(req.Size),
    )
    if err != nil {
        return nil, status.Errorf(codes.Internal, "failed to search media: %v", err)
    }

    var out []*threadpb.GroupedMediaResponse

    for thread, medias := range grouped {
        var comm string
        if thread.CommunityID != nil {
            comm = thread.CommunityID.String()
        }
        threadpbObj := &threadpb.Thread{
            Id:           thread.ID.String(),
            UserId:       thread.UserID.String(),
            CommunityId:  comm,
            Content:      thread.Content,
            LikeCount:    thread.LikeCount,
            CommentCount: thread.CommentCount,
            ShareCount:   thread.ShareCount,
            ViewCount:    thread.ViewCount,
            CreatedAt:    toProtoTimestamp(thread.CreatedAt),
            UpdatedAt:    toProtoTimestamp(thread.UpdatedAt),
        }

        var mediaList []*threadpb.ThreadMedia
        for _, media := range medias {
            mediaResp, err := h.mediaClient.GetMediaById(ctx, &mediapb.GetMediaByIdRequest{Id: media.MediaID.String()})
            if err != nil {
                return nil, status.Errorf(codes.Internal, "media lookup failed: %v", err)
            }
            mediaList = append(mediaList, &threadpb.ThreadMedia{
                ThreadId:  media.ThreadID.String(),
                ImageUrl:  mediaResp.PublicUrl,
                Extension: mediaResp.Extension,
                CreatedAt: toProtoTimestamp(media.CreatedAt),
            })
        }

        out = append(out, &threadpb.GroupedMediaResponse{
            Thread:    threadpbObj,
            MediaList: mediaList,
        })
    }

    return &threadpb.SearchMediaByContentResponse{
        Groups: out,
        Total: int32(total),
    }, nil
}
