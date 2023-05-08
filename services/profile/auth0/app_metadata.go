package auth0

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"net/http"
	"time"
)

func (c *Client) UpdateAppMetadata(userId string, profileId uuid.UUID) error {
	updateAppMetadataRequest := &UpdateAppMetadataRequest{
		ProfileAppMetadata{ProfileId: profileId},
	}
	body, err := json.Marshal(updateAppMetadataRequest)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("PATCH",
		fmt.Sprintf("%s/api/v2/users/%s", c.Auth0IssuerUrl, userId),
		bytes.NewBuffer(body),
	)
	if err != nil {
		return err
	}

	resp := &UpdateAppMetadataResponse{}

	return c.Execute(req, resp)
}

type UpdateAppMetadataRequest struct {
	ProfileAppMetadata ProfileAppMetadata `json:"app_metadata"`
}

type ProfileAppMetadata struct {
	ProfileId uuid.UUID `json:"profileId"`
}

type UpdateAppMetadataResponse struct {
	CreatedAt     time.Time `json:"created_at"`
	Email         string    `json:"email"`
	EmailVerified bool      `json:"email_verified"`
}
