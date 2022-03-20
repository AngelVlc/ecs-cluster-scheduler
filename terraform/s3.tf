data "archive_file" "lambda_zip" {
  type = "zip"

  source_dir  = "${path.module}/../.aws-sam/build/EcsClusterSchedulerFunction"
  output_path = "${path.module}/lambda.zip"
}

resource "aws_s3_bucket" "lambda_deploy_bucket" {
  bucket = "lamda-ecs-cluster-scheduler-deploy"

  force_destroy = true
}

resource "aws_s3_bucket_acl" "bucket_acl" {
  bucket = aws_s3_bucket.lambda_deploy_bucket.id
  acl    = "private"
}

resource "aws_s3_object" "lambda_deploy_bucket_object" {
  bucket = aws_s3_bucket.lambda_deploy_bucket.id

  key    = "lambda.zip"
  source = data.archive_file.lambda_zip.output_path

  etag = filemd5(data.archive_file.lambda_zip.output_path)
}
