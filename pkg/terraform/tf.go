package terraform

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"sync"
)

type StateFile struct {
	Version   string        `json:"version"`
	TfVersion string        `json:"terraform_version"`
	Resources []interface{} `json:"resources"`
	Module    string        `json:",omitempty"`
}

var (
	modules []StateFile
	mu      sync.Mutex
)

func getObjects(sess *session.Session, bucket string) ([]string, error) {
	if len(bucket) == 0 {
		return nil, fmt.Errorf("bucket cannot be empty")
	}

	svc := s3.New(sess)
	input := &s3.ListObjectsV2Input{
		Bucket: aws.String(bucket),
	}

	var keys []string
	var pageNum int
	err := svc.ListObjectsV2Pages(input,
		func(page *s3.ListObjectsV2Output, lastPage bool) bool {
			pageNum++
			for _, v := range page.Contents {
				keys = append(keys, *v.Key)
			}
			return pageNum <= 1000
		})

	if err != nil {
		log.Fatal(err)
	}

	return keys, nil
}

func getStates(sess *session.Session, bucketName string, keys []string) ([]StateFile, error) {
	var wg sync.WaitGroup
	svc := s3.New(sess)

	if len(keys) == 0 {
		return nil, nil
	}

	for _, v := range keys {
		wg.Add(1)
		go func(v string) {
			// dirty hack replace this
			if !strings.HasSuffix(v, "terraform.tfstate") {
				wg.Done()
				return
			}
			defer wg.Done()

			input := &s3.GetObjectInput{
				Bucket: aws.String(bucketName),
				Key:    aws.String(v),
			}
			res, err := svc.GetObject(input)

			if err != nil {
				log.Fatalf("unable to get object: %v", err)
			}
			body, err := ioutil.ReadAll(res.Body)

			defer res.Body.Close()

			mu.Lock()
			mod := StateFile{Module: v}
			json.Unmarshal(body, &mod)
			modules = append(modules, mod)
			mu.Unlock()

			if err != nil {
				log.Fatal(err)
			}
		}(v)
	}

	wg.Wait()

	return modules, nil
}

func GetAWSResources(bucketName string) error {
	sess, err := session.NewSessionWithOptions(session.Options{
		Profile: os.Getenv("AWS_PROFILE"),
	})

	if err != nil {
		return fmt.Errorf("unable to create aws session: %v", err)
	}
	keys, err := getObjects(sess, bucketName)

	if err != nil {
		return fmt.Errorf("unable to get objects from bucket %s: %v", bucketName, err)
	}

	states, err := getStates(sess, bucketName, keys)

	if err != nil {
		return err
	}

	for _, v := range states {
		fmt.Printf("%s has %d resources\n", v.Module, len(v.Resources))
	}

	return nil
}
