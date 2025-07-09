provider "aws" {
  region = var.aws_region
}

resource "aws_iam_role" "lambda_exec" {
  name = "weather_lambda_exec"
  assume_role_policy = jsonencode({
    Version = "2012-10-17",
    Statement = [{
      Effect = "Allow",
      Principal = {
        Service = "lambda.amazonaws.com"
      },
      Action = "sts:AssumeRole"
    }]
  })
}

resource "aws_iam_policy_attachment" "lambda_logs" {
  name       = "weather-lambda-log-attachment"
  roles      = [aws_iam_role.lambda_exec.name]
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"
}

resource "aws_lambda_function" "weather_lambda" {
  function_name = "weather_lambda"
  filename      = "${path.module}/../lambda.zip"
  source_code_hash = filebase64sha256("${path.module}/../lambda.zip")
  handler       = "bootstrap" # REQUIRED for custom runtimes
  runtime       = "provided.al2"
  role          = aws_iam_role.lambda_exec.arn

  environment {
    variables = {
      # Add env vars if needed
    }
  }
}

resource "aws_apigatewayv2_api" "weather_api" {
  name          = "weather-api"
  protocol_type = "HTTP"
}

resource "aws_apigatewayv2_integration" "lambda_integration" {
  api_id             = aws_apigatewayv2_api.weather_api.id
  integration_type   = "AWS_PROXY"
  integration_uri    = aws_lambda_function.weather_lambda.invoke_arn
  integration_method = "POST"
  payload_format_version = "2.0"
}

resource "aws_apigatewayv2_route" "weather_route" {
  api_id    = aws_apigatewayv2_api.weather_api.id
  route_key = "GET /weather"
  target    = "integrations/${aws_apigatewayv2_integration.lambda_integration.id}"
}

resource "aws_apigatewayv2_stage" "default_stage" {
  api_id      = aws_apigatewayv2_api.weather_api.id
  name        = "$default"
  auto_deploy = true
}

resource "aws_lambda_permission" "allow_apigw" {
  statement_id  = "AllowAPIGatewayInvoke"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.weather_lambda.arn
  principal     = "apigateway.amazonaws.com"
  source_arn    = "${aws_apigatewayv2_api.weather_api.execution_arn}/*/*"
}
