resource "aws_apigatewayv2_api" "this" {
  name          = "${var.project_name}-gw-${var.env}"
  protocol_type = "HTTP"
  tags = merge(var.tags, {
    Name = "${var.project_name}-gw-${var.env}"
  })

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

# Lambda 統合（HTTP API 用の正しい integration_uri 形式）
resource "aws_apigatewayv2_integration" "this" {
  api_id                 = aws_apigatewayv2_api.this.id
  integration_type       = "AWS_PROXY"
  integration_uri        = "arn:aws:apigateway:${var.aws_region}:lambda:path/2015-03-31/functions/${var.lambda_arn}/invocations"
  payload_format_version = "2.0"
  timeout_milliseconds   = 30000
}

# 認証が必要なルート（任意のパスを Lambda にプロキシ）
resource "aws_apigatewayv2_route" "authenticated" {
  api_id    = aws_apigatewayv2_api.this.id
  route_key = "ANY /{proxy+}"
  target    = "integrations/${aws_apigatewayv2_integration.this.id}"

  authorization_type = "JWT"
  authorizer_id      = aws_apigatewayv2_authorizer.cognito.id
}

# Lambdaへの実行許可（API Gateway からの呼び出しを許可）
resource "aws_lambda_permission" "apigw" {
  statement_id  = "AllowExecutionFromAPIGateway"
  action        = "lambda:InvokeFunction"
  function_name = var.lambda_function_name
  principal     = "apigateway.amazonaws.com"
  source_arn    = "${aws_apigatewayv2_api.this.execution_arn}/*/*"
}
