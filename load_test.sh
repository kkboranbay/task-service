TOKEN=$(curl -s -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin"}' | jq -r '.token')

for i in {1..1000}; do
  curl -s -H "Authorization: Bearer $TOKEN" http://localhost:8080/api/v1/tasks > /dev/null
  sleep 0.1
done