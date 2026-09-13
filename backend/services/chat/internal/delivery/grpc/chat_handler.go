package grpc

import (
	"context"
	"time"

	"github.com/google/uuid"
	chatpb "github.com/williamsebastianliman/WEB-WS-242/proto/gen/go/chatpb/proto/chat"
	"github.com/williamsebastianliman/WEB-WS-242/services/chat/internal/service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ChatHandler struct {
    chatpb.UnimplementedChatServiceServer
    svc service.ChatService
}

func NewChatHandler(svc service.ChatService) *ChatHandler {
    return &ChatHandler{svc: svc}
}

func (h *ChatHandler) CreateChat(ctx context.Context, req *chatpb.CreateChatRequest) (*chatpb.CreateChatResponse, error) {
    groupID, err := uuid.Parse(req.GetGroupId())
    if err != nil {
        return nil, status.Errorf(codes.InvalidArgument, "invalid group_id")
    }
    senderID, err := uuid.Parse(req.GetSenderId())
    if err != nil {
        return nil, status.Errorf(codes.InvalidArgument, "invalid sender_id")
    }
    if err := h.svc.CreateChat(groupID, senderID, req.GetMessage()); err != nil {
        return nil, status.Errorf(codes.Internal, err.Error())
    }
    return &chatpb.CreateChatResponse{}, nil
}

func (h *ChatHandler) CreateGroup(ctx context.Context, req *chatpb.CreateGroupRequest) (*chatpb.CreateGroupResponse, error) {
    memberIDs := make([]uuid.UUID, len(req.MemberIds))
    for i, s := range req.MemberIds {
        id, err := uuid.Parse(s)
        if err != nil {
            return nil, status.Errorf(codes.InvalidArgument, "invalid member id: %s", s)
        }
        memberIDs[i] = id
    }
    groupID, err := h.svc.CreateGroup(memberIDs)
    if err != nil {
        return nil, status.Errorf(codes.Internal, err.Error())
    }
    return &chatpb.CreateGroupResponse{GroupId: groupID.String()}, nil
}

func (h *ChatHandler) AddMember(ctx context.Context, req *chatpb.AddMemberRequest) (*chatpb.AddMemberResponse, error) {
    groupID, err := uuid.Parse(req.GetGroupId())
    if err != nil {
        return nil, status.Errorf(codes.InvalidArgument, "invalid group_id")
    }
    userID, err := uuid.Parse(req.GetUserId())
    if err != nil {
        return nil, status.Errorf(codes.InvalidArgument, "invalid user_id")
    }
    if err := h.svc.AddMember(groupID, userID); err != nil {
        return nil, status.Errorf(codes.Internal, err.Error())
    }
    return &chatpb.AddMemberResponse{}, nil
}

func (h *ChatHandler) RemoveMember(ctx context.Context, req *chatpb.RemoveMemberRequest) (*chatpb.RemoveMemberResponse, error) {
    groupID, err := uuid.Parse(req.GetGroupId())
    if err != nil {
        return nil, status.Errorf(codes.InvalidArgument, "invalid group_id")
    }
    userID, err := uuid.Parse(req.GetUserId())
    if err != nil {
        return nil, status.Errorf(codes.InvalidArgument, "invalid user_id")
    }
    if err := h.svc.RemoveMember(groupID, userID); err != nil {
        return nil, status.Errorf(codes.Internal, err.Error())
    }
    return &chatpb.RemoveMemberResponse{}, nil
}

func (h *ChatHandler) CreateConversation(ctx context.Context, req *chatpb.CreateConversationRequest) (*chatpb.CreateConversationResponse, error) {
    userA, err := uuid.Parse(req.GetUserA())
    if err != nil {
        return nil, status.Errorf(codes.InvalidArgument, "invalid user_a")
    }
    userB, err := uuid.Parse(req.GetUserB())
    if err != nil {
        return nil, status.Errorf(codes.InvalidArgument, "invalid user_b")
    }
    groupID, err := h.svc.CreateConversation(userA, userB)
    if err != nil {
        return nil, status.Errorf(codes.Internal, err.Error())
    }
    return &chatpb.CreateConversationResponse{GroupId: groupID.String()}, nil
}

func (h *ChatHandler) GetAllGroupsByUser(ctx context.Context, req *chatpb.GetAllGroupsByUserRequest) (*chatpb.GetAllGroupsByUserResponse, error) {
    userID, err := uuid.Parse(req.GetUserId())
    if err != nil {
        return nil, status.Errorf(codes.InvalidArgument, "invalid user_id")
    }
    groups, err := h.svc.GetAllGroupsByUser(userID)
    if err != nil {
        return nil, status.Errorf(codes.Internal, err.Error())
    }
    res := &chatpb.GetAllGroupsByUserResponse{}
    for _, g := range groups {
        res.Groups = append(res.Groups, &chatpb.Group{
            GroupId: g.GroupID.String(),
            IsPrivate: g.IsPrivate,
        })
    }
    return res, nil
}

func (h *ChatHandler) GetAllChatsByGroup(ctx context.Context, req *chatpb.GetAllChatsByGroupRequest) (*chatpb.GetAllChatsByGroupResponse, error) {
    groupID, err := uuid.Parse(req.GetGroupId())
    if err != nil {
        return nil, status.Errorf(codes.InvalidArgument, "invalid group_id")
    }
    chats, err := h.svc.GetAllChatsByGroup(groupID)
    if err != nil {
        return nil, status.Errorf(codes.Internal, err.Error())
    }
    res := &chatpb.GetAllChatsByGroupResponse{}
    for _, c := range chats {
        res.Chats = append(res.Chats, &chatpb.Chat{
            ChatId: c.ChatID.String(),
            GroupId: c.GroupID.String(),
            SenderId: c.SenderID.String(),
            Message: c.Message,
            IsMasked: c.IsMasked,
            CreatedAt: c.CreatedAt.Format(time.RFC3339),
        })
    }
    return res, nil
}

func (h *ChatHandler) DeleteChat(ctx context.Context, req *chatpb.DeleteChatRequest) (*chatpb.DeleteChatResponse, error) {
    chatID, err := uuid.Parse(req.GetChatId())
    if err != nil {
        return nil, status.Errorf(codes.InvalidArgument, "invalid chat_id")
    }
    if err := h.svc.DeleteChat(chatID); err != nil {
        return nil, status.Errorf(codes.Internal, err.Error())
    }
    return &chatpb.DeleteChatResponse{}, nil
}

func (h *ChatHandler) MaskChat(ctx context.Context, req *chatpb.MaskChatRequest) (*chatpb.MaskChatResponse, error) {
    chatID, err := uuid.Parse(req.GetChatId())
    if err != nil {
        return nil, status.Errorf(codes.InvalidArgument, "invalid chat_id")
    }
    if err := h.svc.MaskChat(chatID); err != nil {
        return nil, status.Errorf(codes.Internal, err.Error())
    }
    return &chatpb.MaskChatResponse{}, nil
}

func (h *ChatHandler) GetAllMembersByGroup(ctx context.Context, req *chatpb.GetAllMembersByGroupRequest) (*chatpb.GetAllMembersByGroupResponse, error) {
    groupID, err := uuid.Parse(req.GetGroupId())
    if err != nil {
        return nil, status.Errorf(codes.InvalidArgument, "invalid group_id")
    }
    members, err := h.svc.GetAllMembersByGroup(groupID)
    if err != nil {
        return nil, status.Errorf(codes.Internal, err.Error())
    }
    res := &chatpb.GetAllMembersByGroupResponse{}
    for _, m := range members {
        res.Members = append(res.Members, &chatpb.GroupMember{
            GroupId:  m.GroupID.String(),
            MemberId: m.MemberID.String(),
        })
    }
    return res, nil
}