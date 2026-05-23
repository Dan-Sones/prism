package experimentResults

type DecisionRecommendation string

const (
	DecisionRecommendationUnspecified  DecisionRecommendation = "DECISION_RECOMMENDATION_UNSPECIFIED"
	DecisionRecommendationRecommend    DecisionRecommendation = "DECISION_RECOMMENDATION_RECOMMEND"
	DecisionRecommendationNotRecommend DecisionRecommendation = "DECISION_RECOMMENDATION_NOT_RECOMMEND"
	DecisionRecommendationInconclusive DecisionRecommendation = "DECISION_RECOMMENDATION_INCONCLUSIVE"
)
