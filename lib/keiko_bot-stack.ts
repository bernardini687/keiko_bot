import * as apigateway from '@aws-cdk/aws-apigateway'
import * as s3 from '@aws-cdk/aws-s3'
import * as cdk from '@aws-cdk/core'
import * as golang from 'aws-lambda-golang'

export class KeikoBotStack extends cdk.Stack {
  constructor(scope: cdk.Construct, id: string, props?: cdk.StackProps) {
    super(scope, id, props)

    const bucket = new s3.Bucket(this, 'KakebotBucket')

    const handler = new golang.GolangFunction(this, 'handler', {
      environment: {
        BOT_TOKEN: '', // TODO: insert bot token here!
        BUCKET_NAME: bucket.bucketName,
        KEIKO_GOAL: '60'
      },
    })

    bucket.grantReadWrite(handler)

    const api = new apigateway.LambdaRestApi(this, 'kakebot-api', {
      handler,
      proxy: false,
    })

    const bot = api.root.addResource('bot')
    bot.addMethod('POST')
  }
}
