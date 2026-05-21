resource "aws_apigatewayv2_api" "this" {
  name          = "${var.project_name}-gw-${var.env}"
  protocol_type = "HTTP"

  tags = merge(var.tags, {
    Name = "${var.project_name}-gw-${var.env}"
  })
}

resource "aws_apigatewayv2_stage" "this" {
  api_id      = aws_apigatewayv2_api.this.id
  name        = "$default"
  auto_deploy = true
}

# =========================
# authorizer
# =========================

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

# =========================
# lambda integration
# =========================

resource "aws_apigatewayv2_integration" "this" {
  api_id           = aws_apigatewayv2_api.this.id
  integration_type = "AWS_PROXY"

  integration_uri = "arn:aws:apigateway:${var.aws_region}:lambda:path/2015-03-31/functions/${var.lambda_arn}/invocations"

  payload_format_version = "2.0"
  timeout_milliseconds   = 30000
}

# =========================
# public routes
# =========================

resource "aws_apigatewayv2_route" "health" {
  api_id    = aws_apigatewayv2_api.this.id
  route_key = "GET /health"

  target = "integrations/${aws_apigatewayv2_integration.this.id}"
}

resource "aws_apigatewayv2_route" "swagger_docs_root" {
  api_id    = aws_apigatewayv2_api.this.id
  route_key = "GET /api/docs/"

  target = "integrations/${aws_apigatewayv2_integration.this.id}"
}

resource "aws_apigatewayv2_route" "swagger_docs_proxy" {
  api_id    = aws_apigatewayv2_api.this.id
  route_key = "ANY /api/docs/{proxy+}"

  target = "integrations/${aws_apigatewayv2_integration.this.id}"
}

resource "aws_apigatewayv2_route" "swagger_spec" {
  api_id    = aws_apigatewayv2_api.this.id
  route_key = "GET /api/spec/swagger.yml"

  target = "integrations/${aws_apigatewayv2_integration.this.id}"
}

# =========================
# protected routes
# =========================

resource "aws_apigatewayv2_route" "authenticated" {
  api_id    = aws_apigatewayv2_api.this.id
  route_key = "ANY /api/v1/{proxy+}"

  target = "integrations/${aws_apigatewayv2_integration.this.id}"

  authorization_type = "JWT"
  authorizer_id      = aws_apigatewayv2_authorizer.cognito.id
}

# =========================
# lambda permission
# =========================

resource "aws_lambda_permission" "apigw" {
  statement_id = "AllowExecutionFromAPIGateway"

  action        = "lambda:InvokeFunction"
  function_name = var.lambda_function_name
  principal     = "apigateway.amazonaws.com"

  source_arn = "${aws_apigatewayv2_api.this.execution_arn}/*/*"
}
