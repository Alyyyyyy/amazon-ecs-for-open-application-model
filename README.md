# Amazon ECS for Open Application Model

The oam-ecs CLI is a proof-of-concept that partially implements the [Open Application Model](https://oam.dev/) (OAM) specification, version v1alpha1.

The oam-ecs CLI provisions two of the core OAM workload types as Amazon ECS services running on AWS Fargate using AWS CloudFormation.  A workload of type `core.oam.dev/v1alpha1.Worker` will deploy a CloudFormation stack containing an ECS service running in private VPC subnets with no accessible endpoint.  A workload of type `core.oam.dev/v1alpha1.Server` will deploy a CloudFormation stack containing an ECS service running in private VPC subnets, behind a publicly-accessible network load balancer.

For a full comparison with the OAM specification, see the [Compatibility](COMPATIBILITY.md) page.

See the [full demo here](https://raw.githubusercontent.com/awslabs/amazon-ecs-for-open-application-model/master/demo.gif).

>⚠️ Note that this project is a proof-of-concept and should not be used with production workloads. Use the `--dry-run` option to review all CloudFormation templates generated by this tool before deploying them to your AWS account.

**Table of Contents**

<!-- toc -->

- [Build the CLI](#build--test)
- [Deploy an environment](#deploy-an-oam-ecs-environment)
- [Deploy an application](#deploy-oam-workloads-with-oam-ecs)
- [Upgrade and scale an application](#upgrade-and-scale-oam-workloads-with-oam-ecs)
- [Tear down](#tear-down)
- [Specify AWS credentials and region](#credentials-and-region)
- [Security Disclosures](#security-disclosures)
- [License](#license)

<!-- tocstop -->

## Build & test

```
make
make test
export PATH="$PATH:./bin/local"
oam-ecs --help
```

To customize the CloudFormation templates generated by this tool (for example, to change the VPC configuration or security group rules), edit the following files and re-compile:
* [environment template](templates/environment/cf.yml)
* [workload template](templates/core.oam.dev/cf.yml)

## Deploy an oam-ecs environment

The oam-ecs environment deployment creates a VPC with public and private subnets where OAM workloads can be deployed.

```
oam-ecs env deploy
```

The environment attributes like VPC ID and ECS cluster name can be described.

```
oam-ecs env show
```

The CloudFormation template deployed by this command can be [seen here](templates/environment/cf.yml).

## Deploy OAM workloads with oam-ecs

The dry-run step outputs the CloudFormation template that represents the given OAM workloads.  The CloudFormation templates are written to the `./oam-ecs-dry-run-results` directory.

```
oam-ecs app deploy --dry-run \
  -f examples/example-app.yaml \
  -f examples/worker-component.yaml \
  -f examples/server-component.yaml
```

Then the CloudFormation resources, including load balancers and ECS services running on Fargate, can be deployed:

```
oam-ecs app deploy \
  -f examples/example-app.yaml \
  -f examples/worker-component.yaml \
  -f examples/server-component.yaml
```

The application component instances' attributes like ECS service name and endpoint DNS name can be described.

```
oam-ecs app show -f examples/example-app.yaml
```

## Upgrade and scale OAM workloads with oam-ecs

To change operational settings like the scale of a component instance or to add new component instances to an application, update the application configuration file (e.g. `example-app.yaml`) and re-run the `oam-ecs deploy` command with the same inputs.  The existing CloudFormation stacks for the application will be updated with the new settings.

```
oam-ecs app deploy \
  -f examples/example-app.yaml \
  -f examples/worker-component.yaml \
  -f examples/server-component.yaml
```

To upgrade a component to a new image tag, you can update the image tag in the component schematic file (e.g. `server-component.yaml`), and re-run the `oam-ecs deploy` command with the same inputs as above.  The existing CloudFormation stack for the updated component instance will be updated with the new image tag.

The oam-ecs tool does not require following the [OAM spec guidance](https://github.com/oam-dev/spec/blob/4af9e65769759c408193445baf99eadd93f3426a/6.application_configuration.md#releases) that component schematics be treated as immutable.  To follow the spec guidance when upgrading to a new image tag, create a new component schematic (e.g. `server-component-v2.yaml` with name `server-v2`) and update the component instance in the application configuration (e.g. update the `componentName` to `server-v2` for the instance `example-server` in `example-app.yaml`).  Running `oam-ecs deploy` with the new component schematic will update that component instance's CloudFormation stack with the new image tag.  Do NOT update the `instanceName` in the application configuration, as that will result in creating a new CloudFormation stack and will not delete the previous CloudFormation stack.

```
oam-ecs app deploy \
  -f examples/example-app.yaml \
  -f examples/worker-component.yaml \
  -f examples/server-component-v2.yaml
```

## Tear down

To delete all infrastructure provisioned by oam-ecs, first delete the deployed applications:

```
oam-ecs app delete -f examples/example-app.yaml
```

You can delete the infrastructure for individual component instances by creating an application configuration file containing only that component instance, and running the above `oam-ecs delete` command.  Note that the `oam-ecs deploy` command does NOT comply with the [OAM spec requirement](https://github.com/oam-dev/spec/blob/4af9e65769759c408193445baf99eadd93f3426a/6.application_configuration.md#releases) to automatically delete the infrastructure for component instances that have been removed in an updated application configuration.

To delete the environment infrastructure, once all applications are deleted:

```
oam-ecs env delete
```

## Credentials and Region

oam-ecs will look for credentials in the following order, using the [default provider chain](https://docs.aws.amazon.com/sdk-for-go/v1/developer-guide/configuring-sdk.html#specifying-credentials) in the AWS SDK for Go.

1. Environment variables.
1. Shared credentials file. Profiles can be specified using the `AWS_PROFILE` environment variable.
1. If running on Amazon ECS (with task role) or AWS CodeBuild, IAM role from the container credentials endpoint.
1. If running on an Amazon EC2 instance, IAM role for Amazon EC2.

No credentials are required for dry-runs of the oam-ecs tool.

oam-ecs will determine the region in the following order, using the [default behavior](https://docs.aws.amazon.com/sdk-for-go/v1/developer-guide/configuring-sdk.html#specifying-the-region) in the AWS SDK for Go.

1. From the `AWS_REGION` environment variable.
1. From the `config` file in the `.aws/` folder in your home directory.

## Security Disclosures

If you would like to report a potential security issue in this project, please do not create a GitHub issue.  Instead, please follow the instructions [here](https://aws.amazon.com/security/vulnerability-reporting/) or [email AWS Security directly](mailto:aws-security@amazon.com).

## License

This project is licensed under the Apache-2.0 License.
