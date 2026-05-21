# ユーザープール本体
resource "aws_cognito_user_pool" "this" {
  name = "${var.project_name}-${var.env}-user-pool"

  # メールアドレスでログインするように設定
  username_attributes      = ["email"]
  auto_verified_attributes = ["email"]

  password_policy {
    minimum_length    = 8
    require_lowercase = true
    require_numbers   = true
    require_symbols   = true
    require_uppercase = true
  }

  verification_message_template {
    default_email_option = "CONFIRM_WITH_CODE"
    email_message        = "Your verification code is {####}."
    email_subject        = "Verify your email"
  }

  schema {
    attribute_data_type = "String"
    name                = "email"
    required            = true
    mutable             = true
  }

  tags = merge(var.tags, {
    Name = "${var.project_name}-${var.env}-user-pool"
  })
}

# フロントエンドから接続するためのクライアント設定
resource "aws_cognito_user_pool_client" "this" {
  name         = "${var.project_name}-${var.env}-client"
  user_pool_id = aws_cognito_user_pool.this.id

  explicit_auth_flows = [
    "ALLOW_USER_PASSWORD_AUTH",
    "ALLOW_REFRESH_TOKEN_AUTH",
    "ALLOW_USER_SRP_AUTH"
  ]

  # クライアントシークレットはNext.js(フロント)では使わないのでfalse
  generate_secret = false
}

# account id 取得
data "aws_caller_identity" "current" {}

resource "aws_cognito_user_pool_domain" "this" {
  domain = "${var.project_name}-${var.env}-${substr(data.aws_caller_identity.current.account_id, 8, 4)}"

  user_pool_id = aws_cognito_user_pool.this.id
}
