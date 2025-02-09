# High Priority Jobs (Priority 3)
curl -X POST http://localhost:8080/jobs \
-H "Content-Type: application/json" \
-d '{"type":"email","priority":3,"payload":{"to":"high1@example.com"}}'

curl -X POST http://localhost:8080/jobs \
-H "Content-Type: application/json" \
-d '{"type":"api","priority":3,"payload":{"url":"high2-api.com"}}'

curl -X POST http://localhost:8080/jobs \
-H "Content-Type: application/json" \
-d '{"type":"data","priority":3,"payload":{"query":"high3-query"}}'

# Medium Priority Jobs (Priority 2)
curl -X POST http://localhost:8080/jobs \
-H "Content-Type: application/json" \
-d '{"type":"email","priority":2,"payload":{"to":"medium1@example.com"}}'

curl -X POST http://localhost:8080/jobs \
-H "Content-Type: application/json" \
-d '{"type":"api","priority":2,"payload":{"url":"medium2-api.com"}}'

curl -X POST http://localhost:8080/jobs \
-H "Content-Type: application/json" \
-d '{"type":"data","priority":2,"payload":{"query":"medium3-query"}}'

curl -X POST http://localhost:8080/jobs \
-H "Content-Type: application/json" \
-d '{"type":"email","priority":2,"payload":{"to":"medium4@example.com"}}'

# Low Priority Jobs (Priority 1)
curl -X POST http://localhost:8080/jobs \
-H "Content-Type: application/json" \
-d '{"type":"email","priority":1,"payload":{"to":"low1@example.com"}}'

curl -X POST http://localhost:8080/jobs \
-H "Content-Type: application/json" \
-d '{"type":"api","priority":1,"payload":{"url":"low2-api.com"}}'

curl -X POST http://localhost:8080/jobs \
-H "Content-Type: application/json" \
-d '{"type":"data","priority":1,"payload":{"query":"low3-query"}}'

curl -X POST http://localhost:8080/jobs \
-H "Content-Type: application/json" \
-d '{"type":"email","priority":1,"payload":{"to":"low4@example.com"}}'

curl -X POST http://localhost:8080/jobs \
-H "Content-Type: application/json" \
-d '{"type":"api","priority":1,"payload":{"url":"low5-api.com"}}'