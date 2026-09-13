package database

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/williamsebastianliman/WEB-WS-242/services/community/internal/model"
	"gorm.io/gorm"
)

func Seeder(db *gorm.DB) error {
	var communities []*model.Community
	now := time.Now()

	for i := 1; i <= 40; i++ {
		communities = append(communities, &model.Community{
			CommunityID:          uuid.New(),
			CommunityName:        fmt.Sprintf("Community %02d", i),
			CommunityRules:       fmt.Sprintf("Rule set for community %02d.", i),
			CommunityDescription: fmt.Sprintf("This is a description for community %02d.", i),
			Status:               "active",
			CreatedAt:            now.Add(-time.Duration(i) * time.Hour * 24),
			CommunityLogo:        "1c2d3bd4-89f1-417c-9647-af99b4ac831f",
			CommunityBanner:      "9073750f-e4ab-48d3-8b5f-af6b03fd2fd6",
		})
	}

	hardCodedCommunityID := uuid.MustParse("5ba990c8-4581-4851-b87b-c8bc8bc11c77")
	hardCodedCommunity := &model.Community{
		CommunityID:          hardCodedCommunityID,
		CommunityName:        "Hardcoded Community",
		CommunityRules:       "These are the rules for the hardcoded community.",
		CommunityDescription: "This is a special community with a fixed ID.",
		Status:               "active",
		CreatedAt:            now,
		CommunityLogo:        "1c2d3bd4-89f1-417c-9647-af99b4ac831f",
		CommunityBanner:      "9073750f-e4ab-48d3-8b5f-af6b03fd2fd6",
	}
	communities = append(communities, hardCodedCommunity)
	if err := db.Create(communities).Error; err != nil {
		return err
	}

	memberID := uuid.MustParse("ec724072-ddf8-4ebe-a6ec-ce041375a1e1")
	ownerMember := &model.CommunityMember{
		CommunityID: hardCodedCommunityID,
		MemberID:    memberID,
		Status:      "owner",
		CreatedAt:   now,
	}

	return db.Create(ownerMember).Error
}