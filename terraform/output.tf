output "api_endpoint" {
  value       = aws_apigatewayv2_api.weather_api.api_endpoint
  description = "Base URL of the deployed API Gateway"
}