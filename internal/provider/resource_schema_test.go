package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

// TestWebHookSchema_HeaderValueSensitive locks in the sensitive-data
// guarantee: webhook header values commonly carry auth tokens, so they must
// never be rendered in plain text by plan/apply output.
func TestWebHookSchema_HeaderValueSensitive(t *testing.T) {
	r := NewWebHookResource()

	var resp resource.SchemaResponse
	r.Schema(t.Context(), resource.SchemaRequest{}, &resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("schema diagnostics: %v", resp.Diagnostics)
	}

	headersAttr, ok := resp.Schema.Attributes["headers"]
	if !ok {
		t.Fatal("webhook schema has no headers attribute")
	}
	headers, ok := headersAttr.(schema.ListNestedAttribute)
	if !ok {
		t.Fatalf("headers attribute type = %T, want schema.ListNestedAttribute", headersAttr)
	}

	valueAttr, ok := headers.NestedObject.Attributes["value"]
	if !ok {
		t.Fatal("webhook headers schema has no value attribute")
	}
	value, ok := valueAttr.(schema.StringAttribute)
	if !ok {
		t.Fatalf("headers.value attribute type = %T, want schema.StringAttribute", valueAttr)
	}
	if !value.Sensitive {
		t.Error("webhook headers.value must be Sensitive: header values commonly carry auth tokens")
	}
}

// TestSecretSchema_FieldValueSensitive guards the same guarantee for secret
// field values (encrypted write-only data).
func TestSecretSchema_FieldValueSensitive(t *testing.T) {
	r := NewSecretResource()

	var resp resource.SchemaResponse
	r.Schema(t.Context(), resource.SchemaRequest{}, &resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("schema diagnostics: %v", resp.Diagnostics)
	}

	fieldsAttr, ok := resp.Schema.Attributes["fields"]
	if !ok {
		t.Fatal("secret schema has no fields attribute")
	}
	fields, ok := fieldsAttr.(schema.ListNestedAttribute)
	if !ok {
		t.Fatalf("fields attribute type = %T, want schema.ListNestedAttribute", fieldsAttr)
	}

	valueAttr, ok := fields.NestedObject.Attributes["value"]
	if !ok {
		t.Fatal("secret fields schema has no value attribute")
	}
	value, ok := valueAttr.(schema.StringAttribute)
	if !ok {
		t.Fatalf("fields.value attribute type = %T, want schema.StringAttribute", valueAttr)
	}
	if !value.Sensitive {
		t.Error("secret fields.value must be Sensitive: field values are encrypted write-only data")
	}
}

// TestAPIKeySchema_KeySensitive guards the API key value: it is returned by
// the API only on create and must never be rendered in plain text.
func TestAPIKeySchema_KeySensitive(t *testing.T) {
	r := NewAPIKeyResource()

	var resp resource.SchemaResponse
	r.Schema(t.Context(), resource.SchemaRequest{}, &resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("schema diagnostics: %v", resp.Diagnostics)
	}

	keyAttr, ok := resp.Schema.Attributes["key"]
	if !ok {
		t.Fatal("api_key schema has no key attribute")
	}
	key, ok := keyAttr.(schema.StringAttribute)
	if !ok {
		t.Fatalf("key attribute type = %T, want schema.StringAttribute", keyAttr)
	}
	if !key.Sensitive {
		t.Error("api_key key must be Sensitive: the plain key value is a credential")
	}
	if !key.Computed {
		t.Error("api_key key must be Computed: it is returned by the API on create only")
	}
}
