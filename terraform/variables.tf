variable "aws_region" {
  description = "AWS region to deploy to"
  default     = "eu-west-1"
}

variable "dynamo_table_name" {
  description = "DynamoDB table name"
  default     = "WeatherCache"
}
