resource "aws_apigatewayv2_api" "this" {
  name          = "${var.project_name}-gw-${var.env}"
  protocol_type = "HTTP"
}

# CognitoをJWT Authorizerとして登録
resource "aws_apigatewayv2_authorizer" "cognito" {
  api_id           = aws_apigatewayv2_api.this.id
  authorizer_type  = "JWT"
  identity_sources = ["$request.header.Authorization"]
  name             = "cognito-authorizer"

  jwt_configuration {
    audience = [var.client_id]
    issuer   = var.user_pool_endpoint
  }
}

# 認証が必要なルート
resource "aws_apigatewayv2_route" "authenticated" {
  api_id    = aws_apigatewayv2_api.this.id
  route_key = "ANY /{proxy+}"
  target    = "integrations/${aws_apigatewayv2_integration.this.id}"

  # ここで作成したAuthorizerを紐付ける
  authorization_type = "JWT"
  authorizer_id      = aws_apigatewayv2_authorizer.cognito.id
}

resource "aws_apigatewayv2_integration" "this" {
  api_id           = aws_apigatewayv2_api.this.id
  integration_type = "AWS_PROXY"
  integration_uri  = var.lambda_arn
}

# Lambdaへの実行許可（これがないと403になる）
resource "aws_lambda_permission" "apigw" {
  action        = "lambda:InvokeFunction"
  function_name = var.lambda_function_name
  principal     = "apigateway.amazonaws.com"
  source_arn    = "${aws_apigatewayv2_api.this.execution_arn}/*/*"
}
