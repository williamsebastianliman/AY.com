package event

import (
	"context"
	"encoding/json"
	"log"

	buspkg "github.com/williamsebastianliman/WEB-WS-242/services/user/internal/event"
	"github.com/williamsebastianliman/WEB-WS-242/services/user/internal/service"
)


type userActivatePayload struct {
    Id string `json:"id"`
    Email             string `json:"email"`
}

func StartEventConsumer(ctx context.Context, bus buspkg.EventBus, userSvc service.UserService) error{
    return bus.Subscribe(context.Background(), "user.activate", func(ctx context.Context, msg []byte) {
        var p userActivatePayload
        if err := json.Unmarshal(msg, &p); err != nil {
            log.Printf("[user.activate] invalid payload: %v", err)
            return
        }

        _, err := userSvc.ActivateAccount(
            ctx,
            p.Id,
        )
        if err != nil {
            log.Printf("[user.activate] failed to register: %v", err)
        } else {
            log.Printf("[user.activate] user %s created", p.Id)
        }
    })
}