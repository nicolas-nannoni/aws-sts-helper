# AWS STS (Security Token Service) Helper #
Never set environment variables yourself anymore when using _GetSessionToken_, _AssumeRole_ and temporary credentials from [STS](https://www.google.se/url?sa=t&rct=j&q=&esrc=s&source=web&cd=1&cad=rja&uact=8&ved=0ahUKEwiz1tTdicjOAhWBGywKHQyzCGMQFggdMAA&url=http%3A%2F%2Fdocs.aws.amazon.com%2FSTS%2Flatest%2FAPIReference%2FWelcome.html&usg=AFQjCNHIkvYM6R9tkhrAsp4O9fHqjr0nTw) (Amazon Security Token Service).

## Download ##
Precompiled binaries are available for Linux, macOS and Windows (note: I do not test it on Windows, shell operations support must be limited, but it has been [reported to work](https://github.com/nicolas-nannoni/aws-sts-helper/issues/5#issuecomment-343071645)). 
Check the [latest release](https://github.com/nicolas-nannoni/aws-sts-helper/releases/latest).

## Quick start ##
Set the following environment variables to the values you use the most:

```shell
ROLE_ARN=arn:aws:iam::123456789:role/YourRole # The role ARN you want to assume
MFA_ARN=arn:aws:iam::123456789:mfa/YourUser   # The MFA (Multi-Factor Authentication) ARN attached to your user,
                                                if the role you want to assume requires it
```

Make sure that your main, long-term AWS credentials information is [properly set](http://docs.aws.amazon.com/java-sdk/latest/developer-guide/credentials.html#id6). 
Passing your Access Key and Secret Key as environment variables requires the use of the `--keep-aws-environment` (otherwise, 
these variables get cleared before requesting the token, to avoid reusing previous temporary credentials).

Then, to get a shell with the proper environment variables assuming the role in `ROLE_ARN` and using and MFA device:
```shell
aws-sts-helper get-token --mfa-arn $MFA_ARN \
                         --role-arn $ROLE_ARN \
                         in-new-shell
```

You will get prompted for your MFA token code. Enter it and a new shell with the proper environment variables set will be spawned.

To have a convenient shortcut to that command, you can register an alias in your `.bash_profile` or equivalent:
```
alias sts-shell-somerole="<path>/aws-sts-helper-osx-amd64 get-token --mfa-arn $MFA_ARN --role-arn $ROLE_ARN in-new-shell"
```

If you want to add the variables to your existing shell instead, you have to pass your MFA token code directly in the command invocation:
```shell
eval $(aws-sts-helper get-token --mfa-arn $MFA_ARN \
                                --role-arn $ROLE_ARN \
                                --token-code <TOKEN-CODE-FROM-YOUR-MFA-DEVICE> \
                                and-show-export)
```

### Session tokens ###
AWS supports issuing [session tokens](https://docs.aws.amazon.com/STS/latest/APIReference/API_GetSessionToken.html), which are temporary access
key/secret key/session tokens that can be used in lieu of the permanent IAM or AWS account access keys. They can come in handy when you have to switch role 
multiple times during a given period from the same IAM or AWS user, and those roles need MFA validation. Without session tokens, you would have to retype an 
MFA code every time you assume a role. With a session token, you only type it once and the token you receive is valid (by default) for 12 hours.
 
#### Assuming multiple MFA-protected roles without typing a code every time ###
You can use an MFA-authenticated session token for that and assume the roles using the temporary credentials:
```shell
# Request a session token (giving your MFA code when prompted)
aws-sts-helper get-token --mfa-arn $MFA_ARN \
                         in-new-shell
                         
# Assume a role requiring MFA using the session token in environment variables
# Note that you do not need to specify the MFA ARN that the role needs as it is already authenticated 
# via the session token that will be retained given --keep-aws-environment
aws-sts-helper get-token --role-arn $ROLE_ARN \
                         --keep-aws-environment
                         in-new-shell
``` 

### Credential server ###
It may be useful to serve the temporary credentials retrieved from STS over HTTP as if the credentials came from an 
[EC2 instance metadata](http://docs.aws.amazon.com/AWSEC2/latest/UserGuide/ec2-instance-metadata.html) endpoint. Some 
applications such as [Cyberduck](https://trac.cyberduck.io/wiki/help/en/howto/s3) do not support temporary credentials 
by default but can fetch them from an HTTP endpoint exposing them as an EC2 instance would.

To start a web server serving the the credentials as JSON on port 3000 under `/credentials`:
```shell
aws-sts-helper get-token --mfa-arn $MFA_ARN \
                         --role-arn $ROLE_ARN \
                         and-serve-via-http
```

You can customise the port and the endpoint via `--port` and `--path`.

## Help ##
It is built into the program:

    aws-sts-helper --help


## How to build and run? ##

    go get github.com/nicolas-nannoni/aws-sts-helper
    cd $GOPATH/src/github.com/nicolas-nannoni/aws-sts-helper
    dep ensure
    make

Binaries will be in the `bin` folder

### Cross compile

    make linux
    make osx

On Linux:

    ./bin/aws-sts-helper-linux-amd64

On Mac:

    ./bin/aws-sts-helper-osx-amd64
