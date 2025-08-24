package helpers

import (
	"encoding/json"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/require"
)

func init() { gofakeit.Seed(0) }

func AssertAllObjectsHaveKeysJSON(t *testing.T, body []byte, requiredKeys ...string) {
	t.Helper()

	var anyJSON any
	require.NoErrorf(t, json.Unmarshal(body, &anyJSON), "invalid JSON: %s", string(body))

	switch v := anyJSON.(type) {
	case map[string]any:
		for _, k := range requiredKeys {
			_, ok := v[k]
			require.Truef(t, ok, "missing key %q in object: %v", k, v)
		}
	case []any:
		for i, elem := range v {
			obj, ok := elem.(map[string]any)
			require.Truef(t, ok, "element %d is not an object: %T", i, elem)
			for _, k := range requiredKeys {
				_, ok := obj[k]
				require.Truef(t, ok, "missing key %q in element %d: %v", k, i, obj)
			}
		}
	default:
		require.Failf(t, "unsupported JSON type", "expected object or array of objects, got %T", v)
	}
}
