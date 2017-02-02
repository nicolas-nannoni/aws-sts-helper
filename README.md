# AWS STS (Security Token Service) Helper #
Never set environment variables yourself anymore when using _AssumeRole_ and temporary credentials from [STS](https://www.google.se/url?sa=t&rct=j&q=&esrc=s&source=web&cd=1&cad=rja&uact=8&ved=0ahUKEwiz1tTdicjOAhWBGywKHQyzCGMQFggdMAA&url=http%3A%2F%2Fdocs.aws.amazon.com%2FSTS%2Flatest%2FAPIReference%2FWelcome.html&usg=AFQjCNHIkvYM6R9tkhrAsp4O9fHqjr0nTw) (Amazon Security Token Service).

## Download ##
Precompiled binaries are available for Linux and macOS. Check the [latest release](https://github.com/nicolas-nannoni/aws-sts-helper/releases/latest).

## Quick start ##
Set the following environment variables to the values you use the most:

```shell
ROLE_ARN=arn:aws:iam::123456789:role/YourRole # The role ARN you want to assume
MFA_ARN=arn:aws:iam::123456789:mfa/YourUser   # The MFA (Multi-Factor Authentication) ARN attached to your user,
                                                if the role you want to assume requires it
```

Make sure that your main AWS credentials information are [properly set](http://docs.aws.amazon.com/java-sdk/latest/developer-guide/credentials.html#id6). Passing your Access Key and Secret Key as environment variables require the use of the `--keep-aws-environment` (otherwise, these variables get cleared before requesting the token, to avoid reusing previous temporary credentials).

Then, to get a shell with the proper environment variables assuming the role in `ROLE_ARN` and using and MFA device:
```shell
aws-sts-helper get-token --mfa-arn $MFA_ARN \
                         --role-arn $ROLE_ARN \
                         in-new-shell
```

You will get prompted for your MFA token code. Enter it and a new shell with the proper environment variables set will be spawned.

If you want to add the variables to your existing shell instead, you have to pass your MFA token code directly in the command invocation:
```shell
eval $(aws-sts-helper get-token --mfa-arn $MFA_ARN \
                                --role-arn $ROLE_ARN \
                                --token-code <TOKEN-CODE-FROM-YOUR-MFA-DEVICE> \
                                and-show-export)
```


## Help ##
It is built into the program:

    aws-sts-helper help


## How to build and run? ##

    git clone git@github.com:nicolas-nannoni/aws-sts-helper.git
    go get
    make

Binaries will be in the `bin` folder

### Cross compile

    make linux
    make osx
    make windows

On Linux:

    ./bin/aws-sts-helper-linux

On Mac:

    ./bin/aws-sts-helper-osx

On Windows:

    ./bin/aws-sts-helper.exe

Or to just run it without building:

    make run
