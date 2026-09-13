package routes

import (
	"github.com/gin-gonic/gin"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	authpb "github.com/williamsebastianliman/WEB-WS-242/proto/gen/go/authpb/proto/auth"
	chatpb "github.com/williamsebastianliman/WEB-WS-242/proto/gen/go/chatpb/proto/chat"
	communitypb "github.com/williamsebastianliman/WEB-WS-242/proto/gen/go/communitypb/proto/community"
	mediapb "github.com/williamsebastianliman/WEB-WS-242/proto/gen/go/mediapb/proto/media"
	threadpb "github.com/williamsebastianliman/WEB-WS-242/proto/gen/go/threadpb/proto/thread"
	userpb "github.com/williamsebastianliman/WEB-WS-242/proto/gen/go/userpb/proto/user"

	"github.com/williamsebastianliman/WEB-WS-242/api-gateway/internal/delivery/http"
	"github.com/williamsebastianliman/WEB-WS-242/api-gateway/internal/middleware"
)

func InitRoutes(
	router *gin.Engine,
	userClient userpb.UserServiceClient,
	authClient authpb.AuthServiceClient,
	securityAnswerClient authpb.SecurityAnswerServiceClient,
	mediaClient mediapb.MediaServiceClient,
	threadClient threadpb.ThreadServiceClient,
	communityClient communitypb.CommunityServiceClient,
	chatClient chatpb.ChatServiceClient,
) {
	api := router.Group("/api/v1")
	{
		api.POST("/auth/send-otp", http.SendOTPHandler(authClient))
		api.POST("/auth/verify-otp", http.VerifyOTPHandler(authClient))
		api.POST("/auth/login", http.LoginHandler(authClient))
		api.POST("/auth/refresh", http.RefreshHandler(authClient))

		api.POST("/user/register", http.RegisterUser(userClient))
		api.GET("/user/validate-username", http.IsUsernameUnique(userClient))
		api.POST("/user/get-by-username", http.GetUserByUsernameHandler(userClient))
		api.POST("/user/get-by-email", http.GetUserByEmail(userClient))
		api.GET("/user/get-me", middleware.AuthMiddleware("p9sD#7vZ!qE3rC@X1tL$zM4nA&bUoJ8w"), http.GetMeHandler(userClient))
		api.GET("/user/:id", http.GetUserByIDHandler(userClient))
		api.POST("/user/block", http.BlockUserHandler(userClient))
		api.POST("/user/unblock", http.UnblockUserHandler(userClient))
		api.POST("/user/is-blocked", http.IsUserBlockedHandler(userClient))
		api.POST("/user/request-premium", http.RequestPremiumHandler(userClient))

		api.GET("/users/communities/joined/:user_id",    http.ListJoinedCommunities(communityClient))
		api.GET("/users/communities/requests/:user_id",  http.ListRequestedCommunities(communityClient))
		api.GET("/users/communities/moderated/:user_id", http.ListModeratedCommunities(communityClient))


		api.POST("/user/follow", http.FollowUserHandler(userClient))
		api.POST("/user/unfollow", http.UnfollowUserHandler(userClient))
		api.GET("/user/followers/:user_id", http.GetFollowersHandler(userClient))
		api.GET("/user/following/:user_id", http.GetFollowingHandler(userClient))
		api.GET("/user/follower_count/:user_id", http.FollowerCountHandler(userClient))
		api.GET("/user/following_count/:user_id", http.FollowingCountHandler(userClient))
		api.POST("/user/is-followed", http.IsFollowedHandler(userClient))

		api.POST("/media/upload", http.UploadHandler(mediaClient))
		api.POST("/media/get-media", http.GetMediaByIDHandler(mediaClient))

		api.POST("/security-answer/create", http.CreateSecurityAnswerHandler(securityAnswerClient))
		api.POST("/security-answer", http.GetSecurityAnswersHandler(securityAnswerClient))

		api.POST("/threads", http.CreateThread(threadClient))
		api.GET("/threads/following", http.ListThreadsByUserFollower(threadClient))
		api.GET("/threads/:id", http.GetThreadByID(threadClient))
		api.PUT("/threads/:id", http.UpdateThread(threadClient))
		api.DELETE("/threads/:id", middleware.ThreadOwnerOnly(threadClient,"p9sD#7vZ!qE3rC@X1tL$zM4nA&bUoJ8w"),http.DeleteThread(threadClient))

		api.GET("/threads", http.ListThreads(threadClient))
		api.GET("/users/:user_id/threads", http.ListThreadsByUser(threadClient))
		api.GET("/users/:user_id/threads-like", http.ListThreadsByUserLike(threadClient))
		api.GET("/users/:user_id/threads-bookmark", http.ListThreadsByUserBookmark(threadClient))
		api.GET("/communities/threads/:community_id", http.ListThreadsByCommunity(threadClient))
		api.GET("/threads/search", http.SearchThreads(threadClient))
		api.GET("/users/:user_id/media", http.GetMediaByUser(threadClient))
		api.GET("/replies", http.ListReplies(threadClient))
		api.GET("/replies/user/:user_id", http.ListRepliesByUserId(threadClient))

		api.GET("/threads/count", http.CountThreads(threadClient))
		api.POST("/threads/:id/like", http.LikeThread(threadClient))
		api.DELETE("/threads/:id/like", http.UnlikeThread(threadClient))
		api.GET("/threads/:id/likes", http.CountThreadLikes(threadClient))
		api.POST("/threads/:id/bookmark", http.BookmarkThread(threadClient))
		api.DELETE("/threads/:id/bookmark", http.UnbookmarkThread(threadClient))
		api.GET("/threads/:id/bookmarks", http.CountThreadBookmarks(threadClient))
		api.GET("/threads/:id/replies/count", http.CountThreadReplies(threadClient))

		api.GET("/threads/:id/repost/count", http.CountThreadRepost(threadClient))


		api.GET("/threads/:id/media", http.ListThreadMedia(threadClient))
		
		api.POST("/threads/:id/media", http.AddThreadMedia(threadClient))
		api.DELETE("/threads/:id/media", http.RemoveThreadMedia(threadClient))
		api.GET("/threads/:id/media/count", http.CountThreadMedia(threadClient))

		api.POST("/threads/:id/liked",       http.HasLikedThread(threadClient))
		api.POST("/threads/:id/bookmarked",  http.HasBookmarkedThread(threadClient))
		api.POST("/threads/:id/reposted",  http.HasRepostedThread(threadClient))
		
		api.POST("/ml/predict", http.PredictHandler())

		
		api.GET("/communities/pending", http.ListPendingCommunities(communityClient))
		api.GET("/communities",             http.ListCommunities(communityClient))
		api.GET("/communities/search",      http.SearchCommunities(communityClient))

		api.GET("/communities/:id",         http.GetCommunityByID(communityClient))

		api.POST("/communities",            http.CreateCommunity(communityClient))
		api.POST("/communities/join/:community_id",   http.RequestJoinCommunity(communityClient))
		api.GET("/communities/members/:community_id",    http.ListCommunityMembers(communityClient))
		api.POST("/communities/members/promote/:community_id/:user_id", http.PromoteMember(communityClient))
		api.POST("/communities/members/demote/:community_id/:user_id", http.DemoteMember(communityClient))

		api.POST("/communities/requests/accept/:community_id/:user_id", http.AcceptJoinRequest(communityClient))
		api.POST("/communities/requests/reject/:community_id/:user_id", http.RejectJoinRequest(communityClient))

		api.GET("/communities/about/:community_id",      http.CommunityAbout(communityClient))
		api.GET("/communities/categories/:community_id", http.ListCommunityCategories(communityClient))
		api.POST("/communities/categories/:community_id", http.AddCommunityCategory(communityClient))
		api.DELETE("/communities/categories/:community_id", http.RemoveCommunityCategory(communityClient))
		api.GET("/community-categories",     http.ListCategories(communityClient))
		api.POST("/community-categories",    http.CreateCategory(communityClient))
		api.PUT("/community-categories/:category_id", http.UpdateCategory(communityClient))
		api.DELETE("/community-categories/:category_id", http.DeleteCategory(communityClient))
		api.GET("/communities/top-members/:community_id")

		api.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

		api.POST("/thread-categories", http.CreateThreadCategory(threadClient))
		api.GET("/thread-categories", http.ListThreadCategories(threadClient))
		api.GET("/thread-categories/threads/:category_id", http.ListThreadsByCategory(threadClient))

		api.GET("/community/pending", http.ListPendingCommunities(communityClient))
		api.POST("/community/accept/:community_id", http.AcceptCommunityRequest(communityClient))
		api.POST("/community/reject/:community_id", http.RejectCommunityRequest(communityClient))

		api.GET("/user/get-premium", http.GetAllPremiumRequestsHandler(userClient))
		api.POST("/user/premium-request/accept", http.AcceptPremiumHandler(userClient))
		api.POST("/user/premium-request/reject", http.RejectPremiumHandler(userClient))
		api.GET("/user/is-premium", http.IsPremium(userClient))
		api.GET("/user/community/role/:community_id/:user_id", http.GetUserRoleInCommunity(communityClient))
		
		api.GET("/threads/community/:community_id/by-like", http.ListThreadsByCommunityByLike(threadClient))
		api.GET("/threads/community/:community_id/by-time", http.ListThreadsByCommunityByTime(threadClient))
		api.GET("/threads/community/:community_id/media", http.GetMediaByCommunity(threadClient))
		api.GET("/community/top-user/:community_id", http.GetTopCommunityMembers(communityClient))
		api.POST("/threads/repost", http.CreateRepost(threadClient))
		
		api.PATCH("/threads/pin/:id", middleware.ThreadOwnerOnly(threadClient,"p9sD#7vZ!qE3rC@X1tL$zM4nA&bUoJ8w"), http.PinThread(threadClient))

		api.GET("/community/count/:community_id", http.GetCommunityMemberCountHandler(communityClient))
		api.PATCH("/user/update-profile", http.UpdateProfileHandler(userClient))

		api.PUT("/thread-categories/:category_id", http.UpdateThreadCategory(threadClient))
		api.DELETE("/thread-categories/:category_id", http.DeleteThreadCategory(threadClient))
		api.GET("/users", http.GetAllUsers(userClient))

		api.POST("/user-reports", http.CreateUserReportHandler(userClient))
		api.GET("/user-reports", http.GetAllUserReportsHandler(userClient))

		api.POST("/auth/check-security-answer", http.CheckSecurityAnswerHandler(securityAnswerClient))

		api.GET("/threads/hashtag", http.GetThreadsByHashtag(threadClient))
		api.GET("/top-hashtag", http.GetTopHashtags(threadClient))

		//Explore Page
		api.GET("/threads/explore", http.SearchThreadByContent(threadClient))
		api.GET("threads-media/explore", http.SearchMediaByContent(threadClient))
		api.GET("/user/explore", http.ExploreUserByNameHandler(userClient))
		api.GET("/community/explore", http.ExploreCommunityByNameHandler(communityClient))

		//Chat Page
		api.POST("/chat", http.CreateChatHandler(chatClient))
		api.POST("/group", http.CreateGroupHandler(chatClient))
		api.POST("/group/add-member", http.AddMemberHandler(chatClient))
		api.POST("/group/remove-member", http.RemoveMemberHandler(chatClient))
		api.POST("/conversation", http.CreateConversationHandler(chatClient))
		api.GET("/user/groups/:user_id", http.GetAllGroupsByUserHandler(chatClient))
		api.GET("/group/chats/:group_id", http.GetAllChatsByGroupHandler(chatClient))
		api.POST("/chat/delete", http.DeleteChatHandler(chatClient))
		api.POST("/chat/mask", http.MaskChatHandler(chatClient))
		api.GET("/group/members/:group_id", http.GetAllMembersByGroupHandler(chatClient))

		//Logout
		api.POST("/logout", middleware.LogoutHandler)
	}
}