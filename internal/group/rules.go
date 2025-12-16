package group

import "tfviz/internal/model"

// Rule defines an interface for grouping resources.
type Rule interface {
	Match(res *model.Resource) bool
	GroupKey(res *model.Resource) string
	GroupName(res *model.Resource) string
}

// Apply takes a set of rules and a state, and returns resources grouped by those rules.
func Apply(rules []Rule, state *model.State) map[string]*Group {
	groups := make(map[string]*Group)

	for _, res := range state.Resources {
		for _, rule := range rules {
			if rule.Match(res) {
				key := rule.GroupKey(res)
				if _, ok := groups[key]; !ok {
					groups[key] = &Group{
						Key:  key,
						Name: rule.GroupName(res),
					}
				}
				groups[key].Resources = append(groups[key].Resources, res)
				break // First matching rule wins
			}
		}
	}
	return groups
}

// AWSEKSRule groups EKS clusters, ASGs, LTs, and SGs
type AWSEKSRule struct{}

func (r AWSEKSRule) Match(res *model.Resource) bool {
	switch res.Type {
	case "aws_eks_cluster", "aws_autoscaling_group", "aws_launch_template", "aws_launch_configuration", "aws_security_group":
		return true
	default:
		return false
	}
}

func (r AWSEKSRule) GroupKey(res *model.Resource) string {
	return "eks"
}

func (r AWSEKSRule) GroupName(res *model.Resource) string {
	return "EKS Cluster(s)"
}

// AWSS3Rule groups S3 buckets
type AWSS3Rule struct{}

func (r AWSS3Rule) Match(res *model.Resource) bool {
	return res.Type == "aws_s3_bucket"
}

func (r AWSS3Rule) GroupKey(res *model.Resource) string {
	return "s3"
}

func (r AWSS3Rule) GroupName(res *model.Resource) string {
	return "S3 Buckets"
}

// AWSRDSRule groups RDS instances
type AWSRDSRule struct{}

func (r AWSRDSRule) Match(res *model.Resource) bool {
	return res.Type == "aws_db_instance"
}

func (r AWSRDSRule) GroupKey(res *model.Resource) string {
	return "rds"
}

func (r AWSRDSRule) GroupName(res *model.Resource) string {
	return "RDS Instances"
}

// AWSVPCRule groups VPCs and subnets
type AWSVPCRule struct{}

func (r AWSVPCRule) Match(res *model.Resource) bool {
	switch res.Type {
	case "aws_vpc", "aws_subnet":
		return true
	default:
		return false
	}
}

func (r AWSVPCRule) GroupKey(res *model.Resource) string {
	return "vpc"
}

func (r AWSVPCRule) GroupName(res *model.Resource) string {
	return "VPC / Subnets"
}
