resource "aws_iam_role" "iam_for_lambda"{
    name                    = "iam_for_lambda_s3_presigned_url"
    assume_role_policy      = data.aws_iam_policy_document.assume_role.json
}

resource "aws_iam_policy" "lambda_logging" {
    name = "s3_presigned_url_lambda_logging"
    path = "/"
    description = "IAM policy for logging from a lambda"

    policy = <<EOF
    {
        "Version":"2012-10-17",
        "Statement": [
            {
                "Action": [
                    "logs:CreateLogGroup",
                    "logs:CreateLogStream",
                    "logs:PutLogEvents"
                ],
                "Resource": "arn:aws:logs:*:*:*",
                "Effect": "Allow"
            },
            {
                "Action": [
                    "ec2:CreateNetworkInterface",
                    "ec2:DescribeNetworkInterfaces",
                    "ec2:DeleteNetworkInterface"
                ],
                "Resource": "*",
                "Effect": "Allow"
            },
            {
                "Action": [
                    "s3:*"
                ],
                "Resource": "*",
                "Effect": "Allow"
            }
        ]
    }
    EOF
}

resource "aws_iam_role_policy_attachment" "lambda_logs" {
    role = aws_iam_role.iam_for_lambda.name
    policy_arn = aws_iam_policy.lambda_logging.arn
}

resource "aws_lambda_function" "lambda" {
  function_name    = "s3_presigned_url_lambda"
  role             = aws_iam_role.iam_for_lambda.arn
  runtime          = "provided.al2"
  handler          = "bootstrap"
  filename         = "lambda.zip"
  source_code_hash = filebase64sha256("${path.module}/lambda.zip")
  memory_size      = 128
  timeout          = 120
}

resource "aws_cloudwatch_log_group" "example"{
    name = "/aws/lambda/${aws_lambda_function.lambda.function_name}"
    retention_in_days = var.log_retention_days
}

resource "aws_lambda_permission" "apigw_lambda_permission" {
    action = "lambda:InvokeFunction"
    function_name = aws_lambda_function.lambda.function_name
    principal = "apigateway.amazonaws.com"
    statement_id = "AllowExecutionFromAPIGateway"
    source_arn = "arn:aws:execute-api:us-east-1:337328321041:5oq7x8vhsf/*"
}