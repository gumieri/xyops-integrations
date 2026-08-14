package sdk

import (
	"context"
	"encoding/json"
)

// BucketsService
func (s *BucketsService) List(ctx context.Context) ([]Bucket, error) {
	body, err := s.client.DoGET(ctx, "get_buckets/v1")
	if err != nil {
		return nil, err
	}

	var resp BucketListResponse
	if err := unmarshalJSON(body, &resp); err != nil {
		return nil, err
	}
	return resp.Rows, nil
}

// Get retrieves a bucket's metadata, data, and files via get_bucket/v1.
// Returns the bucket metadata, the data object (as map[string]any), and
// the list of files.
func (s *BucketsService) Get(ctx context.Context, id string) (*Bucket, map[string]any, []File, error) {
	body, err := s.client.DoGETWithQuery(ctx, "get_bucket/v1", map[string]string{"id": id})
	if err != nil {
		return nil, nil, nil, err
	}

	var resp GetBucketResponse
	if err := unmarshalJSON(body, &resp); err != nil {
		return nil, nil, nil, err
	}

	var data map[string]any
	if len(resp.Data) > 0 {
		if err := json.Unmarshal(resp.Data, &data); err != nil {
			return nil, nil, nil, err
		}
	}

	return resp.Bucket, data, resp.Files, nil
}

// Create inserts a bucket. The request body may include "data" to initialize
// the bucket data. The response includes the bucket metadata only (no data or
// files).
func (s *BucketsService) Create(ctx context.Context, bucket *Bucket) (*Bucket, error) {
	if bucket.ID != "" {
		return s.Update(ctx, bucket)
	}

	body, err := s.client.DoPOST(ctx, "create_bucket/v1", bucket)
	if err != nil {
		return nil, err
	}

	var resp BucketResponse
	if err := unmarshalJSON(body, &resp); err != nil {
		return nil, err
	}
	return resp.Bucket, nil
}

// Update performs a shallow merge via update_bucket/v1 (accepts "data"), then
// re-fetches the bucket. The update response is only {code: 0}.
func (s *BucketsService) Update(ctx context.Context, bucket *Bucket) (*Bucket, error) {
	_, err := s.client.DoPOST(ctx, "update_bucket/v1", bucket)
	if err != nil {
		return nil, err
	}

	b, _, _, err := s.Get(ctx, bucket.ID)
	return b, err
}

func (s *BucketsService) Delete(ctx context.Context, id string) error {
	_, err := s.client.DoPOST(ctx, "delete_bucket/v1", map[string]string{"id": id})
	return err
}

// WriteData shallow-merges a JSON object into a bucket's data. If fetch is true,
// the API returns the full merged data object; otherwise it returns only {code: 0}.
func (s *BucketsService) WriteData(ctx context.Context, id string, data map[string]any, fetch bool) (map[string]any, error) {
	payload := map[string]any{"id": id, "data": data}
	if fetch {
		payload["fetch"] = true
	}

	body, err := s.client.DoPOST(ctx, "write_bucket_data/v1", payload)
	if err != nil {
		return nil, err
	}

	if !fetch {
		return nil, nil
	}

	var result map[string]any
	if err := unmarshalJSON(body, &result); err != nil {
		return nil, err
	}
	if dataField, ok := result["data"].(map[string]any); ok {
		return dataField, nil
	}
	return result, nil
}

// DeleteFile deletes files from a bucket by criteria (filename, path, date,
// size, username). All supplied criteria must match the same file.
func (s *BucketsService) DeleteFile(ctx context.Context, id string, criteria map[string]any) error {
	payload := map[string]any{"id": id}
	for k, v := range criteria {
		payload[k] = v
	}
	_, err := s.client.DoPOST(ctx, "delete_bucket_file/v1", payload)
	return err
}

// Empty removes all files, all data, or both from a bucket. At least one of
// deleteFiles or clearData must be true.
func (s *BucketsService) Empty(ctx context.Context, id string, deleteFiles, clearData bool) error {
	if !deleteFiles && !clearData {
		return nil
	}
	payload := map[string]any{"id": id, "files": deleteFiles, "data": clearData}
	_, err := s.client.DoPOST(ctx, "empty_bucket/v1", payload)
	return err
}
