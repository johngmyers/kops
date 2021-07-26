package awsmodel

import (
	"github.com/aws/aws-sdk-go/service/ec2"
	"k8s.io/kops/upup/pkg/fi"
	"k8s.io/kops/upup/pkg/fi/cloudup/awstasks"
	"k8s.io/kops/upup/pkg/fi/cloudup/awsup"
)

// SecurityGroupModelBuilder finds orphaned security groups to delete.
type SecurityGroupModelBuilder struct {
	*AWSModelContext
	Lifecycle fi.Lifecycle
}

var _ fi.ModelBuilder = &SecurityGroupModelBuilder{}
var _ fi.HasDeletions = &SecurityGroupModelBuilder{}

func (s SecurityGroupModelBuilder) Build(context *fi.ModelBuilderContext) error {
	return nil
}

func (s SecurityGroupModelBuilder) FindDeletions(context *fi.ModelBuilderContext, cloud fi.Cloud) error {
	awsCloud := cloud.(awsup.AWSCloud)

	request := &ec2.DescribeSecurityGroupsInput{
		Filters: []*ec2.Filter{
			awsup.NewEC2Filter("tag:kubernetes.io/cluster/" + s.Cluster.ObjectMeta.Name, "owned"),
		},
	}

	return awsCloud.EC2().DescribeSecurityGroupsPages(request, func(output *ec2.DescribeSecurityGroupsOutput, lastPage bool) bool {
		for _, securityGroup := range output.SecurityGroups {
			if _, ok := context.Tasks["SecurityGroup/"+fi.StringValue(securityGroup.GroupName)]; !ok {
				for _, rule := range securityGroup.IpPermissionsEgress
				context.AddTask(&awstasks.SecurityGroup{
					ID:        securityGroup.GroupId,
					Name:      securityGroup.GroupName,
					Description: securityGroup.Description,
					VPC: s.LinkToVPC(),
					Lifecycle: s.Lifecycle,
				})
			}
		}
		return true
	})
}
