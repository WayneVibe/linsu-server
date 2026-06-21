package model

type SubscriptionPlan string

const (
	FreePlan       SubscriptionPlan = "free"       // 免费版
	BasicPlan      SubscriptionPlan = "basic"      // 基础版
	ProPlan        SubscriptionPlan = "pro"        // 高级版
	EnterprisePlan SubscriptionPlan = "enterprise" // 企业版
)
