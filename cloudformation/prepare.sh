#!/bin/bash

 aws cloudformation \
        package \
        --template-file dynamodb-table.yml \
        --output-template-file dynamodb-release.yml \
        --s3-bucket $S3_BUCKET-$DEPLOY_REGION \
        --region $DEPLOY_REGION
