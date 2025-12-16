package group

import "tfviz/internal/model"

func isALB(res *model.Resource) bool {
	return res.Type == "aws_lb"
}

func isALBRelated(res *model.Resource) bool {
	switch res.Type {
	case "aws_lb",
		"aws_lb_listener",
		"aws_lb_listener_rule",
		"aws_lb_target_group":
		return true
	default:
		return false
	}
}
