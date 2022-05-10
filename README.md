# keiko bot

a kakebo telegram bot deployed through cdk, using an s3 bucket for persistence and a lambda to handle the bot's requests.  
mainly an educational endeavor but i use it everyday :)

# cdk stuff

the `cdk.json` file tells the cdk toolkit how to execute your app.

## useful commands

 * `npm run build`   compile typescript to js
 * `npm run watch`   watch for changes and compile
 * `npm run test`    perform the jest unit tests
 * `cdk deploy`      deploy this stack to your default aws account/region
 * `cdk diff`        compare deployed stack with current state
 * `cdk synth`       emits the synthesized cloudformation template
