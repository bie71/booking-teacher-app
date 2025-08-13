package supabase

import (
	"fmt"
	"log"
	"teacher/internal/config"

	"github.com/go-resty/resty/v2"
)

type Client struct {
	Client *resty.Client
	Cfg    *config.Client
}

func InitUploadClient(cfg *config.Client, restyClient *resty.Client) *Client {

	return &Client{Client: restyClient, Cfg: cfg}
}
func (u *Client) UploadToSupabase(objectName, filePath string) error {

	url := fmt.Sprintf("%s/storage/v1/object/%s/%s",
		u.Cfg.Endpoint,
		u.Cfg.BucketName,
		objectName,
	)

	resp, err := u.Client.R().
		SetHeader("Authorization", "Bearer "+u.Cfg.AccessKey).
		SetFile("file", filePath).
		Put(url)

	if err != nil {
		log.Printf("Error uploading to Supabase: %v", err)
		return err
	}

	log.Printf("Response: %s", resp.String())
	log.Println("Successfully uploaded to Supabase")

	if resp.StatusCode() >= 300 {
		return fmt.Errorf("upload failed: %s", resp.String())
	}

	log.Println("Successfully uploaded to Supabase")

	return nil
}
