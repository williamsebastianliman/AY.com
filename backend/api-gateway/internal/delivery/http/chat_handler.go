package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	chatpb "github.com/williamsebastianliman/WEB-WS-242/proto/gen/go/chatpb/proto/chat"
	"google.golang.org/grpc/status"
)

func CreateChatHandler(client chatpb.ChatServiceClient) gin.HandlerFunc {
    return func(c *gin.Context) {
        var req struct {
            GroupID  string `json:"group_id" binding:"required"`
            SenderID string `json:"sender_id" binding:"required"`
            Message  string `json:"message" binding:"required"`
        }
        if err := c.ShouldBindJSON(&req); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }
        grpcReq := &chatpb.CreateChatRequest{
            GroupId:  req.GroupID,
            SenderId: req.SenderID,
            Message:  req.Message,
        }
        _, err := client.CreateChat(c.Request.Context(), grpcReq)
        if err != nil {
            st := status.Convert(err)
            c.JSON(http.StatusInternalServerError, gin.H{"error": st.Message()})
            return
        }
        c.Status(http.StatusCreated)
    }
}

func CreateGroupHandler(client chatpb.ChatServiceClient) gin.HandlerFunc {
    return func(c *gin.Context) {
        var req struct {
            MemberIDs []string `json:"member_ids" binding:"required"`
            IsPrivate bool     `json:"is_private"`
        }
        if err := c.ShouldBindJSON(&req); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }
        grpcReq := &chatpb.CreateGroupRequest{
            MemberIds: req.MemberIDs,
        }
        resp, err := client.CreateGroup(c.Request.Context(), grpcReq)
        if err != nil {
            st := status.Convert(err)
            c.JSON(http.StatusInternalServerError, gin.H{"error": st.Message()})
            return
        }
        c.JSON(http.StatusOK, gin.H{"group_id": resp.GroupId})
    }
}

func AddMemberHandler(client chatpb.ChatServiceClient) gin.HandlerFunc {
    return func(c *gin.Context) {
        var req struct {
            GroupID string `json:"group_id" binding:"required"`
            UserID  string `json:"user_id" binding:"required"`
        }
        if err := c.ShouldBindJSON(&req); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }
        grpcReq := &chatpb.AddMemberRequest{
            GroupId: req.GroupID,
            UserId:  req.UserID,
        }
        _, err := client.AddMember(c.Request.Context(), grpcReq)
        if err != nil {
            st := status.Convert(err)
            c.JSON(http.StatusInternalServerError, gin.H{"error": st.Message()})
            return
        }
        c.Status(http.StatusOK)
    }
}

func RemoveMemberHandler(client chatpb.ChatServiceClient) gin.HandlerFunc {
    return func(c *gin.Context) {
        var req struct {
            GroupID string `json:"group_id" binding:"required"`
            UserID  string `json:"user_id" binding:"required"`
        }
        if err := c.ShouldBindJSON(&req); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }
        grpcReq := &chatpb.RemoveMemberRequest{
            GroupId: req.GroupID,
            UserId:  req.UserID,
        }
        _, err := client.RemoveMember(c.Request.Context(), grpcReq)
        if err != nil {
            st := status.Convert(err)
            c.JSON(http.StatusInternalServerError, gin.H{"error": st.Message()})
            return
        }
        c.Status(http.StatusOK)
    }
}

func CreateConversationHandler(client chatpb.ChatServiceClient) gin.HandlerFunc {
    return func(c *gin.Context) {
        var req struct {
            UserA string `json:"user_a" binding:"required"`
            UserB string `json:"user_b" binding:"required"`
        }
        if err := c.ShouldBindJSON(&req); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }
        grpcReq := &chatpb.CreateConversationRequest{
            UserA: req.UserA,
            UserB: req.UserB,
        }
        resp, err := client.CreateConversation(c.Request.Context(), grpcReq)
        if err != nil {
            st := status.Convert(err)
            c.JSON(http.StatusInternalServerError, gin.H{"error": st.Message()})
            return
        }
        c.JSON(http.StatusOK, gin.H{"group_id": resp.GroupId})
    }
}

func GetAllGroupsByUserHandler(client chatpb.ChatServiceClient) gin.HandlerFunc {
    return func(c *gin.Context) {
        userID := c.Param("user_id")
        if userID == "" {
            c.JSON(http.StatusBadRequest, gin.H{"error": "user_id required"})
            return
        }
        grpcReq := &chatpb.GetAllGroupsByUserRequest{UserId: userID}
        resp, err := client.GetAllGroupsByUser(c.Request.Context(), grpcReq)
        if err != nil {
            st := status.Convert(err)
            c.JSON(http.StatusInternalServerError, gin.H{"error": st.Message()})
            return
        }
        c.JSON(http.StatusOK, gin.H{"groups": resp.Groups})
    }
}

func GetAllChatsByGroupHandler(client chatpb.ChatServiceClient) gin.HandlerFunc {
    return func(c *gin.Context) {
        groupID := c.Param("group_id")
        if groupID == "" {
            c.JSON(http.StatusBadRequest, gin.H{"error": "group_id required"})
            return
        }
        grpcReq := &chatpb.GetAllChatsByGroupRequest{GroupId: groupID}
        resp, err := client.GetAllChatsByGroup(c.Request.Context(), grpcReq)
        if err != nil {
            st := status.Convert(err)
            c.JSON(http.StatusInternalServerError, gin.H{"error": st.Message()})
            return
        }
        c.JSON(http.StatusOK, gin.H{"chats": resp.Chats})
    }
}

func DeleteChatHandler(client chatpb.ChatServiceClient) gin.HandlerFunc {
    return func(c *gin.Context) {
        var req struct {
            ChatID string `json:"chat_id" binding:"required"`
        }
        if err := c.ShouldBindJSON(&req); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }
        grpcReq := &chatpb.DeleteChatRequest{ChatId: req.ChatID}
        _, err := client.DeleteChat(c.Request.Context(), grpcReq)
        if err != nil {
            st := status.Convert(err)
            c.JSON(http.StatusInternalServerError, gin.H{"error": st.Message()})
            return
        }
        c.Status(http.StatusOK)
    }
}

func MaskChatHandler(client chatpb.ChatServiceClient) gin.HandlerFunc {
    return func(c *gin.Context) {
        var req struct {
            ChatID string `json:"chat_id" binding:"required"`
        }
        if err := c.ShouldBindJSON(&req); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }
        grpcReq := &chatpb.MaskChatRequest{ChatId: req.ChatID}
        _, err := client.MaskChat(c.Request.Context(), grpcReq)
        if err != nil {
            st := status.Convert(err)
            c.JSON(http.StatusInternalServerError, gin.H{"error": st.Message()})
            return
        }
        c.Status(http.StatusOK)
    }
}

func GetAllMembersByGroupHandler(client chatpb.ChatServiceClient) gin.HandlerFunc {
    return func(c *gin.Context) {
        groupID := c.Param("group_id")
        if groupID == "" {
            c.JSON(http.StatusBadRequest, gin.H{"error": "group_id required"})
            return
        }
        grpcReq := &chatpb.GetAllMembersByGroupRequest{GroupId: groupID}
        resp, err := client.GetAllMembersByGroup(c.Request.Context(), grpcReq)
        if err != nil {
            st := status.Convert(err)
            c.JSON(http.StatusInternalServerError, gin.H{"error": st.Message()})
            return
        }
        c.JSON(http.StatusOK, gin.H{"members": resp.Members})
    }
}