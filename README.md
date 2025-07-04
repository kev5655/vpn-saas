# vpn-saas



## Deploy container

Get AWS Account ID
```bash
aws sts get-caller-identity --query Account --output text
```


## Delete a stack

```bash
aws cloudformation delete-stack \
  --stack-name CDKToolkit \
  --region us-east-1 \
  --retain-resources LookupRole
```