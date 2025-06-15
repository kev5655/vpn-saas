package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"

	"github.com/aws/aws-cdk-go/awscdk/v2"
	ec2c "github.com/aws/aws-cdk-go/awscdk/v2/awsec2"
	ecs "github.com/aws/aws-cdk-go/awscdk/v2/awsecs"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

func main() {
	// Load .env
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	region := os.Getenv("CDK_REGION")
	if region == "" {
		log.Fatal("CDK_REGION must be set in .env")
	}

	app := awscdk.NewApp(nil)
	NewVpnStack(app, "VpnStack", &VpnStackProps{
		awscdk.StackProps{
			Env: &awscdk.Environment{
				Region: jsii.String(region),
			},
		},
	})
	app.Synth(nil)
}

type VpnStackProps struct {
	awscdk.StackProps
}

func NewVpnStack(scope constructs.Construct, id string, props *VpnStackProps) awscdk.Stack {
	stack := awscdk.NewStack(scope, &id, &props.StackProps)

	vpc := ec2c.NewVpc(stack, jsii.String("VpnVPC"), &ec2c.VpcProps{
		MaxAzs: jsii.Number(2),
	})

	sg := ec2c.NewSecurityGroup(stack, jsii.String("VpnSG"), &ec2c.SecurityGroupProps{
		Vpc:              vpc,
		AllowAllOutbound: jsii.Bool(true),
	})
	sg.AddIngressRule(ec2c.Peer_AnyIpv4(),
		ec2c.Port_Udp(jsii.Number(51820)),
		jsii.String("Allow WireGuard"),
		jsii.Bool(false))

	cluster := ecs.NewCluster(stack, jsii.String("VpnCluster"), &ecs.ClusterProps{
		Vpc: vpc,
	})

	// Spot-backed capacity
	cluster.AddCapacity(jsii.String("SpotCapacity"), &ecs.AddCapacityOptions{
		InstanceType: ec2c.NewInstanceType(jsii.String("t3.micro")),
		SpotPrice:    jsii.String("0.02"),
	})

	taskDef := ecs.NewTaskDefinition(stack, jsii.String("VpnTaskDef"), &ecs.TaskDefinitionProps{
		Compatibility: ecs.Compatibility_EC2,
		Cpu:           jsii.String("256"),
		MemoryMiB:     jsii.String("512"),
	})

	container := taskDef.AddContainer(jsii.String("VpnContainer"), &ecs.ContainerDefinitionOptions{
		Image: ecs.ContainerImage_FromRegistry(
			jsii.String(fmt.Sprintf("%s/%s:%s", os.Getenv("DH_USER"), os.Getenv("DH_REPO"), os.Getenv("DH_TAG"))),
			nil,
		),
		Environment: &map[string]*string{
			"CLIENT_PUBKEY": jsii.String(os.Getenv("CLIENT_PUBKEY")),
		},
		Logging: ecs.LogDrivers_AwsLogs(&ecs.AwsLogDriverProps{
			StreamPrefix: jsii.String("vpn"),
		}),
	})
	container.AddPortMappings(&ecs.PortMapping{
		ContainerPort: jsii.Number(51820),
		Protocol:      ecs.Protocol_UDP,
	})

	ecs.NewEc2Service(stack, jsii.String("VpnService"), &ecs.Ec2ServiceProps{
		Cluster:        cluster,
		TaskDefinition: taskDef,
		SecurityGroups: &[]ec2c.ISecurityGroup{sg},
	})

	return stack
}

func env() *awscdk.Environment {
	return &awscdk.Environment{
		Region: jsii.String(os.Getenv("CDK_REGION")),
	}
}
