package enums

type ShadeEdgeType int16

const (
	SE_FacesShadedEdgesNotHighlighted     = 0
	SE_FacesShadedEdgesHighlightedInBlack = 1
	SE_FacesNotFilledEdgesInEntityColor   = 2
	SE_FacesInEntityColorEdgesInBlack     = 3
)
