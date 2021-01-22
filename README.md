# AWS SQS Prometheus Exporter

A Prometheus metrics exporter for AWS SQS queues

> **A few words of Thanks:** Most of the code in this repo is borrowed from [candeemis/sqs-prometheus-exporter](https://github.com/candeemis/sqs-prometheus-exporter) with bundle of thanks and love :pray: :heart:. I didn't submit this as a pull request to the original repository, since I switch code only to meet my need. 

## Metrics

| Metric  | Labels | Description |
| ------  | ------ | ----------- |
| sqs\_messages\_visible | Queue Name | Number of messages available |
| sqs\_messages\_delayed | Queue Name | Number of messages delayed |
| sqs\_messages\_not\_visible | Queue Name | Number of messages in flight |

For more information see the [AWS SQS Documentation](https://docs.aws.amazon.com/AWSSimpleQueueService/latest/SQSDeveloperGuide/sqs-message-attributes.html)

## Configuration

Credentials to AWS are provided in the following order:

- Environment variables (AWS\_ACCESS\_KEY\_ID and AWS\_SECRET\_ACCESS\_KEY)
- Shared credentials file (~/.aws/credentials)
- IAM role for Amazon EC2

For more information see the [AWS SDK Documentation](https://docs.aws.amazon.com/sdk-for-go/v1/developer-guide/configuring-sdk.html)

### AWS IAM permissions

The app needs sqs list and read access to the sqs policies

```json
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Sid": "VisualEditor0",
            "Effect": "Allow",
            "Action": [
                "sqs:ListQueues",
                "sqs:GetQueueUrl",
                "sqs:ListDeadLetterSourceQueues",
                "sqs:ReceiveMessage",
                "sqs:GetQueueAttributes",
                "sqs:ListQueueTags"
            ],
            "Resource": "*"
        }
    ]
}
```

## Environment Variables
| Variable      | Default Value | Description                                                  |
|---------------|:---------|:-------------------------------------------------------------|
| PORT          | 9434     | The port for metrics server                                  |
| ENDPOINT      | metrics  | The metrics endpoint                                         |



## Running

```docker run -d -p 9434:9434 bruceleo1969/sqs-exporter```

You can provide the AWS credentials as environment variables depending upon your security rules configured in AWS;

```docker run -d -p 9434:9434 -e AWS_ACCESS_KEY_ID=<access_key> -e AWS_SECRET_ACCESS_KEY=<secret_key> -e AWS_REGION=<region>  bruceleo1969/sqs-exporter```

