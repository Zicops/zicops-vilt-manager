package bucket

import (
	"context"
	"encoding/base64"
	"fmt"
	"os"
	"time"

	"cloud.google.com/go/storage"
	"github.com/zicops/zicops-cass-pool/redis"
	"github.com/zicops/zicops-vilt-manager/constants"
	"github.com/zicops/zicops-vilt-manager/lib/identity"
	"google.golang.org/api/option"
)

// Client ....
type Client struct {
	projectID string
	client    storage.Client
	bucket    storage.BucketHandle
}

// NewStorageHandler return new database action
func NewStorageHandler() Client {
	return Client{projectID: "", client: storage.Client{}, bucket: storage.BucketHandle{}}
}

// InitializeStorageClient ...........
func (sc *Client) InitializeStorageClient(ctx context.Context, projectID string) error {
	serviceAccountZicops := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")
	if serviceAccountZicops == "" {
		return fmt.Errorf("failed to get right credentials for course creator")
	}
	targetScopes := []string{
		"https://www.googleapis.com/auth/cloud-platform",
		"https://www.googleapis.com/auth/userinfo.email",
	}
	currentCreds, _, err := identity.ReadCredentialsFile(ctx, serviceAccountZicops, targetScopes)
	if err != nil {
		return err
	}
	client, err := storage.NewClient(ctx, option.WithCredentials(currentCreds))
	if err != nil {
		return err
	}
	sc.client = *client
	sc.projectID = projectID
	uBucketClient, _ := sc.CreateBucket(ctx, constants.USERS_BUCKET)
	sc.bucket = *uBucketClient
	return nil
}

// CreateBucket  ...........
func (sc *Client) CreateBucket(ctx context.Context, bucketName string) (*storage.BucketHandle, error) {
	bkt := sc.client.Bucket(bucketName)
	exists, err := bkt.Attrs(ctx)
	if err != nil && exists == nil {
		if err := bkt.Create(ctx, sc.projectID, nil); err != nil {
			return nil, err
		}
	}
	return bkt, nil
}

// UploadToGCS ....
func (sc *Client) UploadToGCS(ctx context.Context, fileName string) (*storage.Writer, error) {
	bucketWriter := sc.bucket.Object(fileName).NewWriter(ctx)
	return bucketWriter, nil
}

func (sc *Client) GetSignedURLForObject(ctx context.Context, object string) string {
	key := "signed_url_" + base64.StdEncoding.EncodeToString([]byte(object))
	res, err := redis.GetRedisValue(ctx, key)
	if err == nil && res != "" {
		return res
	}
	opts := &storage.SignedURLOptions{
		Scheme:  storage.SigningSchemeV4,
		Method:  "GET",
		Expires: time.Now().Add(24 * time.Hour),
	}
	url, err := sc.bucket.SignedURL(object, opts)
	if err != nil {
		return ""
	}
	allBut10Secsto24Hours := 24*time.Hour - 10*time.Second
	redis.SetRedisValue(ctx, key, url)
	redis.SetTTL(ctx, key, int(allBut10Secsto24Hours.Seconds()))
	return url
}

func (sc *Client) DeleteObjectsFromBucket(ctx context.Context, bucketPath string) string {
	o := sc.bucket.Object(bucketPath)

	if err := o.Delete(ctx); err != nil {
		return err.Error()
	}

	return "success"
}
