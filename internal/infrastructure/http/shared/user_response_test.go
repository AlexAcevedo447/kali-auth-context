package shared

import (
	"encoding/json"
	"strings"
	"testing"

	"kali-auth-context/internal/domain/identity"
)

func TestUserResponse_NeverIncludesPassword(t *testing.T) {
	users := []*identity.User{{Id: "u1", TenantId: "t1", Username: "ana", Email: "ana@x.com", Password: "$argon2id$secreto"}}
	raw, err := json.Marshal(ToUserResponses(users))
	if err != nil {
		t.Fatal(err)
	}
	body := string(raw)
	if strings.Contains(body, "Password") || strings.Contains(body, "argon2") {
		t.Fatalf("user response must not include the password or its hash: %s", body)
	}
	for _, key := range []string{`"Id":"u1"`, `"TenantId":"t1"`, `"Username":"ana"`, `"Email":"ana@x.com"`} {
		if !strings.Contains(body, key) {
			t.Fatalf("missing %s in %s", key, body)
		}
	}
}
