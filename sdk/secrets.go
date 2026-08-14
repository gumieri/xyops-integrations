package sdk

import "context"

// SecretsService
func (s *SecretsService) List(ctx context.Context) ([]Secret, error) {
	body, err := s.client.DoGET(ctx, "get_secrets/v1")
	if err != nil {
		return nil, err
	}

	var resp SecretListResponse
	if err := unmarshalJSON(body, &resp); err != nil {
		return nil, err
	}
	return resp.Rows, nil
}

func (s *SecretsService) Get(ctx context.Context, id string) (*Secret, error) {
	body, err := s.client.DoGETWithQuery(ctx, "get_secret/v1", map[string]string{"id": id})
	if err != nil {
		return nil, err
	}

	var resp SecretResponse
	if err := unmarshalJSON(body, &resp); err != nil {
		return nil, err
	}
	return resp.Secret, nil
}

// Decrypt retrieves the decrypted field values for a secret via
// decrypt_secret/v1 (admin only, access logged). Returns the array of
// {name, value} pairs; the response envelope uses the key "fields".
func (s *SecretsService) Decrypt(ctx context.Context, id string) ([]SecretField, error) {
	body, err := s.client.DoGETWithQuery(ctx, "decrypt_secret/v1", map[string]string{"id": id})
	if err != nil {
		return nil, err
	}

	var resp DecryptSecretResponse
	if err := unmarshalJSON(body, &resp); err != nil {
		return nil, err
	}
	return resp.Fields, nil
}

// Create inserts a secret (admin only). The request accepts a "fields" array
// of {name, value} objects, which are encrypted and stored separately; the
// "names" list is auto-generated from fields. The response is metadata-only
// (no decrypted values).
func (s *SecretsService) Create(ctx context.Context, secret *Secret) (*Secret, error) {
	if secret.ID != "" {
		return s.Update(ctx, secret)
	}

	body, err := s.client.DoPOST(ctx, "create_secret/v1", secret)
	if err != nil {
		return nil, err
	}

	var resp SecretResponse
	if err := unmarshalJSON(body, &resp); err != nil {
		return nil, err
	}
	return resp.Secret, nil
}

// Update performs a shallow merge via update_secret/v1 (admin only), then
// re-fetches the secret. The request may include a new "fields" array to
// re-encrypt the secret; names are regenerated. The update response is only
// {code: 0}.
func (s *SecretsService) Update(ctx context.Context, secret *Secret) (*Secret, error) {
	_, err := s.client.DoPOST(ctx, "update_secret/v1", secret)
	if err != nil {
		return nil, err
	}

	return s.Get(ctx, secret.ID)
}

func (s *SecretsService) Delete(ctx context.Context, id string) error {
	_, err := s.client.DoPOST(ctx, "delete_secret/v1", map[string]string{"id": id})
	return err
}
