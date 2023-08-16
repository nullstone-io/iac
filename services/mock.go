package services

import (
	"gopkg.in/nullstone-io/go-api-client.v0/auth"
	"gopkg.in/nullstone-io/go-api-client.v0/trace"
	"net/http"
	"net/http/httptest"
	"testing"
)

var (
	// Generate this by doing:
	// cd <path-to>/go-api-client
	// go run auth/generate-jwt/main.go
	FakeJwt = `eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJqdGkiOiI0YjcyNGE4OC1kMzFkLTRhN2QtOTFkYy00MmI3M2NlMzJiMmMiLCJhdWQiOiJnby1hcGktY2xpZW50IiwiaXNzIjoibnVsbHN0b25lIiwic3ViIjoiYWNtZSIsImV4cCI6MTY4Nzk5NjIwOSwiaWF0IjoxNjg3OTA5ODA5LCJuYmYiOjE2ODc5MDk4MDksImVtYWlsIjoiYnJhZC5zaWNrbGVzQGFjbWUuY29tIiwicGljdHVyZSI6IiIsImh0dHBzOi8vbnVsbHN0b25lLmlvL3VzZXJuYW1lIjoiYnJhZC5zaWNrbGVzIiwiaHR0cHM6Ly9udWxsc3RvbmUuaW8vcm9sZXMiOnt9fQ.f0tb217YWNFa1ggjplcNwr9wP_sTwetYgUVijnrHap8406dq1S46-L5PlfP_5VXy51mA3Nx9hivORtJRlGixuBqG3QuUGvsxevgokURmhi7JGqEi-2RpUC0JKRXYLKcICh5qzXoDsD8ilDDSFVFDj10UxfWKaC6g6wwZiECRqVTC0VTAdDR9xjdHAq65MMGkdtgDTphOpeDpb7KjoapfEWtpTEpXAuLuWopgjSemT-fnm0Uxa9PxVrbufSsafY4UkeaBxnn1ttzSMJDN3ScUgzblOySthLzkJU8rTqcXWGkvZjNeaSNfTzSyqqpVqRhafkiGCtZUnz_IGYTJEOQ6HBJL8VH0FBmhjHKrV29fb61STpZkPZlIq3aQc_t4lmqBQxiL5NO0GIUt-jarzP3mKiU9PI8jfY06UVWP-LvaOCf-v1qCUHrN1HpyE-Syx4aSGYAq3yf5A4rkMtkXeYmAL27jht67CGXysJROU6sGYik_KWW_I1oe_0EAWmCoch3qGouJMJz50hldVCrMPYcMJiVSAs1pTqFIJf7TeL1hbl-JWEO7wq39RrJ964Apizj_zKO0iKlZgL3y4IReXuNaPQFwQ1-_e_bv4BuuorvI380mUHI3H7kRTXNbT4fCSkFPsApaVsnH4Diu4wET3dUJnssunq8h6UK3bdYJb8Vy9LI`
)

func MockApiHub(t *testing.T, handler http.Handler) ApiHub {
	server := httptest.NewServer(handler)
	cfg := ApiHubConfig{
		MainAddr:       server.URL,
		FurionAddr:     server.URL,
		EnigmaAddr:     server.URL,
		AuthServerAddr: server.URL,
		IsTraceEnabled: trace.IsEnabled(),
	}
	apiHub := NewApiHub(cfg, auth.RawAccessTokenSource{AccessToken: "mock-user"})
	t.Cleanup(server.Close)
	return apiHub
}
