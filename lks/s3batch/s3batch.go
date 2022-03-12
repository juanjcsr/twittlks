package s3batch

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3Client struct {
	Bucket string
	api    S3API
	input  *s3.ListBucketsInput
}

type S3API interface {
	ListBuckets(ctx context.Context,
		params *s3.ListBucketsInput,
		optFns ...func(*s3.Options)) (*s3.ListBucketsOutput, error)
	ListObjectsV2(ctx context.Context,
		params *s3.ListObjectsV2Input,
		optFns ...func(*s3.Options)) (*s3.ListObjectsV2Output, error)
	PutObject(ctx context.Context,
		params *s3.PutObjectInput,
		optFns ...func(*s3.Options)) (*s3.PutObjectOutput, error)
	GetObject(ctx context.Context,
		params *s3.GetObjectInput,
		optFns ...func(*s3.Options)) (*s3.GetObjectOutput, error)
}

func NewAWSClient(bucketName string) (*S3Client, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return nil, err
	}

	client := s3.NewFromConfig(cfg)
	var s3api S3API = client
	return &S3Client{
		api:    s3api,
		input:  &s3.ListBucketsInput{},
		Bucket: bucketName,
	}, nil
}

func (c *S3Client) GetBuckets(ctx context.Context) {
	res, err := c.api.ListBuckets(ctx, c.input)
	if err != nil {
		log.Fatalln(err)
	}

	for _, b := range res.Buckets {
		fmt.Println(*b.Name + " - - " + b.CreationDate.Format("2006-01-02 15:04:05 Monday"))
	}
}

func (c *S3Client) GetFile(ctx context.Context, bucket string, filepath string) (*[]byte, error) {
	input := &s3.GetObjectInput{
		Bucket: &bucket,
		Key:    &filepath,
	}

	obj, err := c.api.GetObject(ctx, input)
	if err != nil {
		return nil, err
	}
	defer obj.Body.Close()

	body, err := ioutil.ReadAll(obj.Body)
	if err != nil {
		return nil, err
	}
	return &body, nil
}

func (c *S3Client) UploadFile(ctx context.Context, bucket string, remotePath string, filename string) error {
	key := remotePath + "/" + filename
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	input := &s3.PutObjectInput{
		Bucket: &bucket,
		Key:    &key,
		Body:   f,
	}

	_, err = c.api.PutObject(ctx, input)
	if err != nil {
		return err
	}
	return nil
}

func (c *S3Client) ListObjects(ctx context.Context, bucketname string) {
	if bucketname == "" {
		bucketname = c.Bucket
	}
	input := &s3.ListObjectsV2Input{
		Bucket: &bucketname,
	}
	res, err := c.api.ListObjectsV2(ctx, input)
	if err != nil {
		log.Println(err)
		return
	}
	for _, item := range res.Contents {
		fmt.Println("Name:          ", *item.Key)
		fmt.Println("Last modified: ", *item.LastModified)
		fmt.Println("Size:          ", item.Size)
		fmt.Println("Storage class: ", item.StorageClass)
		fmt.Println("")
	}
}
