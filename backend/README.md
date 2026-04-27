# Backend scaffold (no migrations)

## Quickstart local
1. Create .env with DB_DSN, e.g.:
   DB_DSN="user:pass@tcp(localhost:3306)/portfolio?parseTime=true"
2. Run:
   ./run_local.sh
3. Test:
   curl -i http://localhost:8080/api/v1/tasks

## Build Lambda ZIP
./build_lambda.sh
Upload lambda.zip to AWS Lambda.

Notes:
- This scaffold intentionally omits DB migration files.
- In Lambda, set DB_DSN environment variable or use Secrets Manager externally.
