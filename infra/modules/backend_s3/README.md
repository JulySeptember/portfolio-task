Usage:
- CI はビルド成果物を s3://{bucket}/{key} にアップロードする。
- デフォルト key: lambda/your-project-dev.zip
- ルートの module.lambda に backend_bucket_name = module.backend_s3.backend_bucket_id を渡してください。
