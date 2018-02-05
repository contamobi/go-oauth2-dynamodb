#!/bin/bash

 aws cloudformation \
        deploy \
        --region $DEPLOY_REGION \
        --template-file dynamodb-release.yml \
        --stack-name "oauth2-dynamodb" \
        --capabilities CAPABILITY_NAMED_IAM \
        --parameter-overrides OauthBasicTableName="oauth2_basic" OauthAccessTableName="oauth2_access" \
         OauthRefreshTableName="oauth2_refresh"
