package middleware

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Idempotency returns middleware that makes a write endpoint safe to retry.
//
// Flow:
//  1. Read the Idempotency-Key header. No key -> pass through untouched.
//  2. Try to CLAIM the key by inserting a 'pending' row. The UNIQUE
//     constraint on idempotency_key makes this atomic across goroutines.
//  3. If the claim succeeds, WE are the first request: run the handler,
//     capture its response, and persist it as 'completed'.
//  4. If the claim fails on unique-violation, someone else already has
//     this key: look up the stored row and replay it (or signal that
//     it's still in flight).
func Idempotency(db *pgxpool.Pool) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			key := r.Header.Get("Idempotency-Key")
			if key == "" {
				// No key supplied -> nothing to dedupe. Pass through.
				next.ServeHTTP(w, r)
				return
			}

			// --- Step 1: attempt to claim the key ---
			claimed, err := claimKey(r.Context(), db, key)
			if err != nil {
				http.Error(w, "internal server error", http.StatusInternalServerError)
				return
			}

			if !claimed {
				// --- Step 2: key already exists -> replay or 409 ---
				replayStored(r.Context(), db, w, key)
				return
			}

			// --- Step 3: we own the key. Run the handler, capturing output. ---
			rec := &responseRecorder{
				ResponseWriter: w,
				status:         http.StatusOK, // default if handler never calls WriteHeader
				body:           &bytes.Buffer{},
			}
			next.ServeHTTP(rec, r)

			// --- Step 4: persist the captured response as 'completed' ---
			_, _ = db.Exec(
				r.Context(),
				`UPDATE idempotency_keys
				    SET status = 'completed',
				        response_code = $2,
				        response_body = $3
				WHERE idempotency_key = $1`,
				key, rec.status, rec.body.Bytes(),
			)
		})
	}
}

// claimKey tries to insert a fresh 'pending' row for this key.
// Returns true if WE claimed it, false if it already existed.
func claimKey(ctx context.Context, db *pgxpool.Pool, key string) (bool, error) {
	// ON CONFLICT DO NOTHING: if the key already exists, the insert affects
	// zero rows instead of erroring. RowsAffected() == 0 means someone beat us.
	tag, err := db.Exec(
		ctx,
		`INSERT INTO idempotency_keys (idempotency_key, status)
		 VALUES ($1, 'pending')
		 ON CONFLICT (idempotency_key) DO NOTHING`,
		key,
	)
	if err != nil {
		return false, err
	}
	return tag.RowsAffected() == 1, nil
}

// replayStored looks up an existing key and writes back its stored response.
// If the row is still 'pending', the original request hasn't finished — we
// return 409 so the client retries later rather than getting a half-answer.
func replayStored(ctx context.Context, db *pgxpool.Pool, w http.ResponseWriter, key string) {
	var status string
	var code *int
	var body []byte

	err := db.QueryRow(
		ctx,
		`SELECT status, response_code, response_body
		   FROM idempotency_keys
		WHERE idempotency_key = $1`,
		key,
	).Scan(&status, &code, &body)

	if errors.Is(err, pgx.ErrNoRows) {
		// Extremely rare race: claimed by another req but not yet visible.
		// Treat as in-flight.
		http.Error(w, "request in progress", http.StatusConflict)
		return
	}
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	if status != "completed" || code == nil {
		// Still pending: the first request is mid-flight. Tell the client to wait.
		http.Error(w, "request in progress", http.StatusConflict)
		return
	}

	// Completed: replay the original response verbatim.
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(*code)
	_, _ = w.Write(body)
}

// responseRecorder wraps http.ResponseWriter so we can capture what the
// handler wrote (status + body) in order to persist it for replay.
type responseRecorder struct {
	http.ResponseWriter
	status int
	body   *bytes.Buffer
}

func (r *responseRecorder) WriteHeader(code int) {
	r.status = code
	r.ResponseWriter.WriteHeader(code)
}

func (r *responseRecorder) Write(b []byte) (int, error) {
	r.body.Write(b)                  // capture a copy for storage
	return r.ResponseWriter.Write(b) // and pass through to the client
}

// (unused import guard — remove once json is used for validation if you add it)
var _ = json.Marshal
