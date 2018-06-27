// Copyright 2018, Google
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package bonfire implements the B2 service.
package bonfire

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"google.golang.org/grpc/metadata"

	"github.com/golang/protobuf/proto"
	"github.com/kurin/blazer/internal/pyre"
)

func getAuth(ctx context.Context) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", errors.New("no metadata")
	}
	data := md.Get("authentication")
	if len(data) == 0 {
		return "", nil
	}
	return data[0], nil
}

type Bonfire struct {
	Root       string
	mu         sync.Mutex
	buckets    map[int][]byte
	nextBucket int
}

func (b *Bonfire) AuthorizeAccount(ctx context.Context, req *pyre.AuthorizeAccountRequest) (*pyre.AuthorizeAccountResponse, error) {
	return &pyre.AuthorizeAccountResponse{
		ApiUrl: b.Root,
	}, nil
}

func (b *Bonfire) ListBuckets(context.Context, *pyre.ListBucketsRequest) (*pyre.ListBucketsResponse, error) {
	resp := &pyre.ListBucketsResponse{}
	b.mu.Lock()
	defer b.mu.Unlock()
	for _, bs := range b.buckets {
		var bucket pyre.Bucket
		if err := proto.Unmarshal(bs, &bucket); err != nil {
			return nil, err
		}
		resp.Buckets = append(resp.Buckets, &bucket)
	}
	return resp, nil
}

func (b *Bonfire) CreateBucket(ctx context.Context, req *pyre.Bucket) (*pyre.Bucket, error) {
	b.mu.Lock()
	defer b.mu.Unlock()

	n := b.nextBucket
	b.nextBucket++
	req.BucketId = fmt.Sprintf("%d", n)
	bs, err := proto.Marshal(req)
	if err != nil {
		return nil, err
	}
	if b.buckets == nil {
		b.buckets = make(map[int][]byte)
	}
	b.buckets[n] = bs
	return req, nil
}

func (b *Bonfire) GetUploadUrl(context.Context, *pyre.GetUploadUrlRequest) (*pyre.GetUploadUrlResponse, error) {
	return &pyre.GetUploadUrlResponse{
		AuthorizationToken: "flooper",
		UploadUrl:          fmt.Sprintf("%s/b2api/v1/b2_upload_file/%s", b.Root, "uploader"),
		BucketId:           "like 4 or whatever",
	}, nil
}

func (b *Bonfire) UploadFile(context.Context, *pyre.UploadFileRequest) (*pyre.UploadFileResponse, error) {
	return &pyre.UploadFileResponse{}, nil
}
