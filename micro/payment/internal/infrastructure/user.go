package infrastructure

import (
    "fmt"
    "log"
    "payment/internal/config"

    "github.com/go-resty/resty/v2"
)

// User provides methods to interact with the user service from the payment
// service.  It allows creating activity logs on behalf of a user without
// requiring authentication by calling the user service's internal
// endpoints.
type User struct {
    restyClient *resty.Client
    serviceUser config.Service
}

// NewUser constructs a new User client.  It receives a Resty client and
// service configuration (host and port) for the user service.  The
// Resty client is shared across services to reuse underlying HTTP
// connections.
func NewUser(restyClient *resty.Client, serviceUser config.Service) *User {
    return &User{
        restyClient: restyClient,
        serviceUser: config.Service{
            Host: serviceUser.Host,
            Port: serviceUser.Port,
        },
    }
}

// CreateActivityLog sends a request to the user service to create a recent
// activity log for the given user.  The userID identifies the owner of
// the activity, while action and description describe the event.  If the
// user service returns an error status code, this function returns an
// error containing the status.  Otherwise, it returns nil to indicate
// success.
func (u *User) CreateActivityLog(userID uint, action, description string) error {
    url := fmt.Sprintf("%s:%s/api/v1/internal/activity", u.serviceUser.Host, u.serviceUser.Port)
    log.Printf("calling user activity endpoint %s", url)
    payload := map[string]interface{}{
        "user_id":    userID,
        "action":     action,
        "description": description,
    }
    resp, err := u.restyClient.R().
        SetHeader("Content-Type", "application/json").
        SetBody(payload).
        Post(url)
    if err != nil {
        return err
    }
    if resp.IsError() {
        return fmt.Errorf("failed to create activity log: %s", resp.Status())
    }
    return nil
}