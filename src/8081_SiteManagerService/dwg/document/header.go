package document

import (
	"time"

	"tracio.com/sitemanagerservice/dwg/document/enums"
	"tracio.com/sitemanagerservice/dwg/entities"
	"tracio.com/sitemanagerservice/dwg/file/version"
	"tracio.com/sitemanagerservice/dwg/types/units"
	"tracio.com/sitemanagerservice/dwg/utils"
)

type CadHeader struct {
	Version                               version.ACadVersion
	MaintenanceVersion                    int16
	CodePage                              utils.CodePage
	RequiredVersions                      int32
	AssociatedDimensions                  bool
	UpdateDimensionsWhileDragging         bool
	DIMSAV                                bool
	MeasurementUnits                      enums.MeasurementUnits
	PolylineLineTypeGeneration            bool
	OrthoMode                             bool
	RegenerationMode                      bool
	FillMode                              bool
	QuickTextMode                         bool
	PaperSpaceLineTypeScaling             enums.SpaceLineTypeScaling
	LimitCheckingOn                       bool
	BlipMode                              bool
	UserTimer                             bool
	SketchPolylines                       bool
	AngularDirection                      units.AngularDirection
	ShowSplineControlPoints               bool
	MirrorText                            bool
	WorldView                             bool
	ShowModelSpace                        bool
	PaperSpaceLimitsChecking              bool
	RetainXRefDependentVisibilitySettings bool
	DisplaySilhouetteCurves               bool
	CreateEllipseAsPolyline               bool
	ProxyGraphics                         bool
	SpatialIndexMaxTreeDepth              int16
	LinearUnitFormat                      units.LinearUnitFormat
	LinearUnitPrecision                   int16
	AngularUnit                           units.AngularUnitFormat
	AngularUnitPrecision                  int16
	ObjectSnapMode                        enums.ObjectSnapMode
	AttributeVisibility                   enums.AttributeVisibilityMode
	PointDisplayMode                      int16
	UserShort1                            int16
	UserShort2                            int16
	UserShort3                            int16
	UserShort4                            int16
	UserShort5                            int16
	NumberOfSplineSegments                int16
	SurfaceDensityU                       int16
	SurfaceDensityV                       int16
	SurfaceType                           int16
	SurfaceMeshTabulationCount1           int16
	SurfaceMeshTabulationCount2           int16
	SplineType                            enums.SplineType
	ShadeEdge                             enums.ShadeEdgeType
	ShadeDiffuseToAmbientPercentage       int16
	UnitMode                              int16
	MaxViewportCount                      int16
	SurfaceIsolineCount                   int16
	CurrentMultilineJustification         entities.VerticalAlignmentType
	TextQuality                           int16
	LineTypeScale                         float64
	TextHeightDefault                     float64
	MultilineStyleName                    string
	TraceWidthDefault                     float64
	SketchIncrement                       float64
	FilletRadius                          float64
	ThicknessDefault                      float64
	AngleBase                             float64
	PointDisplaySize                      float64
	PolylineWidthDefault                  float64
	UserDouble1                           float64
	UserDouble2                           float64
	UserDouble3                           float64
	UserDouble4                           float64
	UserDouble5                           float64
	ChamferDistance1                      float64
	ChamferDistance2                      float64
	ChamferLength                         float64
	ChamferAngle                          float64
	FacetResolution                       float64
	CurrentMultilineScale                 float64
	CurrentEntityLinetypeScale            float64
	MenuFileName                          string
	HandleSeed                            uint32
	CreateDateTime                        time.Time
	UniversalCreateDateTime               time.Time
	UpdateDateTime                        time.Time
	UniversalUpdateDateTime               time.Time
	TotalEditingTime                      time.Duration
	UserElapsedTimeSpan                   time.Duration
}

func NewCadHeader() *CadHeader {
	return &CadHeader{
		Version:                         version.AC1018,
		MaintenanceVersion:              0,
		CodePage:                        utils.CP_Windows1252,
		RequiredVersions:                0,
		AssociatedDimensions:            true,
		UpdateDimensionsWhileDragging:   true,
		MeasurementUnits:                enums.MU_English,
		PolylineLineTypeGeneration:      false,
		FillMode:                        true,
		PaperSpaceLineTypeScaling:       enums.SpaceLineTypeScaling_Normal,
		AngularDirection:                units.AD_CloseWise,
		SpatialIndexMaxTreeDepth:        3020,
		LinearUnitFormat:                units.LU_Decimal,
		LinearUnitPrecision:             4,
		AngularUnit:                     units.AU_DecimalDegrees,
		ObjectSnapMode:                  enums.OSM_None,
		AttributeVisibility:             enums.AVM_None,
		NumberOfSplineSegments:          8,
		SplineType:                      enums.Spline_None,
		ShadeEdge:                       enums.SE_FacesShadedEdgesNotHighlighted,
		ShadeDiffuseToAmbientPercentage: 70,
		CurrentMultilineJustification:   entities.VA_Top,
		LineTypeScale:                   1.0,
		TextHeightDefault:               2.5,
		MultilineStyleName:              "Standard",
		PointDisplaySize:                0.0,
		CurrentMultilineScale:           20,
		CurrentEntityLinetypeScale:      1,
		MenuFileName:                    ".",
		CreateDateTime:                  time.Now(),
		UniversalCreateDateTime:         time.Now().UTC(),
		UpdateDateTime:                  time.Now(),
		UniversalUpdateDateTime:         time.Now().UTC(),
	}
}
