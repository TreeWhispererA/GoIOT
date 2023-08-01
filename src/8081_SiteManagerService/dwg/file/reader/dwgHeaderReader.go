package reader

import (
	"io"

	"tracio.com/sitemanagerservice/dwg/document"
	"tracio.com/sitemanagerservice/dwg/document/enums"
	"tracio.com/sitemanagerservice/dwg/entities"
	"tracio.com/sitemanagerservice/dwg/file/header"
	"tracio.com/sitemanagerservice/dwg/file/reader/stream"
	"tracio.com/sitemanagerservice/dwg/file/section"
	"tracio.com/sitemanagerservice/dwg/file/version"
	"tracio.com/sitemanagerservice/dwg/types/units"
)

type DwgHeaderReader struct {
	DwgSectionIO

	maintenanceVersion byte
	sreader            stream.IDwgStreamReader
}

func (this *DwgHeaderReader) Read(cadHeader *document.CadHeader) (*header.DwgHeaderHandlesCollection, error) {
	var bitValue bool
	// var byteValue byte
	var shortValue int16
	var longValue int32
	var doubleValue float64
	var stringValue string

	err := error(nil)
	objectPointers := &header.DwgHeaderHandlesCollection{}

	if err = this.checkSentinel(this.sreader, section.StartSentienls[this.SectionName]); err != nil {
		return nil, err
	}

	//RL : Size of the section.
	longValue, err = this.sreader.ReadRawLong()
	size := longValue

	//R2010/R2013 (only present if the maintenance version is greater than 3!) or R2018+:
	if (this.IsR2010Plus() && this.maintenanceVersion > 3) || this.IsR2018Plus() {
		//Unknown (4 byte long), might be part of a 64-bit size.
		this.sreader.ReadRawLong()
	}

	var initialPos int32
	initialPos, err = this.sreader.GetPositionInBits()

	//+R2007 Only:
	if this.IsR2007Plus() {
		//RL : Size in bits
		var sizeInBits int32
		sizeInBits, err = this.sreader.ReadRawLong()

		lastPositionInBits := initialPos + sizeInBits - 1

		//Setup the text handler for versions 2007 and above
		var textReader stream.IDwgStreamReader
		//Create a copy of the stream
		var newStream io.ReadSeeker
		newStream, err = stream.CloneStream(this.sreader.GetStream())
		textReader, err = stream.NewDwgStreamHandler(this.DwgVersion, newStream)
		//Set the position and use the flag
		textReader.SetPositionByFlag(lastPositionInBits)

		//Setup the handler for the references for versions 2007 and above
		var referenceReader stream.IDwgStreamReader
		//Create a copy of the stream
		newStream, err = stream.CloneStream(this.sreader.GetStream())
		referenceReader, err = stream.NewDwgStreamHandler(this.DwgVersion, newStream)
		//Set the position and jump the flag
		referenceReader.SetPositionInBits(lastPositionInBits + 0b1)

		this.sreader = stream.NewDwgMergedReader(this.sreader, textReader, referenceReader)
	}

	//R2013+:
	if this.IsR2013Plus() {
		//BLL : Variabele REQUIREDVERSIONS, default value 0, read only.
		cadHeader.RequiredVersions, err = this.sreader.ReadBitLongLong()
	}

	//Common:
	//BD : Unknown, default value 412148564080.0
	doubleValue, err = this.sreader.ReadBitDouble()
	//BD: Unknown, default value 1.0
	doubleValue, err = this.sreader.ReadBitDouble()
	//BD: Unknown, default value 1.0
	doubleValue, err = this.sreader.ReadBitDouble()
	//BD: Unknown, default value 1.0
	doubleValue, err = this.sreader.ReadBitDouble()
	//TV: Unknown text string, default "m"
	stringValue, err = this.sreader.ReadVariableText()
	//TV: Unknown text string, default ""
	stringValue, err = this.sreader.ReadVariableText()
	//TV: Unknown text string, default ""
	stringValue, err = this.sreader.ReadVariableText()
	//TV: Unknown text string, default ""
	stringValue, err = this.sreader.ReadVariableText()
	//BL : Unknown long, default value 24L
	longValue, err = this.sreader.ReadBitLong()
	//BL: Unknown long, default value 0L;
	longValue, err = this.sreader.ReadBitLong()

	//R13-R14 Only:
	if this.IsR13_14Only() {
		//BS : Unknown short, default value 0
		shortValue, err = this.sreader.ReadBitShort()
	}

	//Pre-2004 Only:
	if this.IsR2004Pre() {
		// 	//H : Handle of the current viewport entity header (hard pointer)
		// 	this.sreader.HandleReference();
	}

	//Common:
	//B: DIMASO
	cadHeader.AssociatedDimensions, err = this.sreader.ReadBit()
	//B: DIMSHO
	cadHeader.UpdateDimensionsWhileDragging, err = this.sreader.ReadBit()

	//R13-R14 Only:
	if this.IsR13_14Only() {
		//B : DIMSAV Undocumented.
		cadHeader.DIMSAV, err = this.sreader.ReadBit()
	}

	//Common:
	//B: PLINEGEN
	cadHeader.PolylineLineTypeGeneration, err = this.sreader.ReadBit()
	//B : ORTHOMODE
	cadHeader.OrthoMode, err = this.sreader.ReadBit()
	//B: REGENMODE
	cadHeader.RegenerationMode, err = this.sreader.ReadBit()
	//B : FILLMODE
	cadHeader.FillMode, err = this.sreader.ReadBit()
	//B : QTEXTMODE
	cadHeader.QuickTextMode, err = this.sreader.ReadBit()
	//B : PSLTSCALE
	if bitValue, err = this.sreader.ReadBit(); err != nil || bitValue == false {
		cadHeader.PaperSpaceLineTypeScaling = enums.SpaceLineTypeScaling_Viewport
	} else {
		cadHeader.PaperSpaceLineTypeScaling = enums.SpaceLineTypeScaling_Normal
	}
	//B : LIMCHECK
	cadHeader.LimitCheckingOn, err = this.sreader.ReadBit()

	//R13-R14 Only (stored in registry from R15 onwards):
	if this.IsR13_14Only() {
		// 	//B : BLIPMODE
		cadHeader.BlipMode, err = this.sreader.ReadBit()
	}
	//R2004+:
	if this.IsR2004Plus() {
		//B : Undocumented
		this.sreader.ReadBit()
	}

	//Common:
	//B: USRTIMER(User timer on / off).
	cadHeader.UserTimer, err = this.sreader.ReadBit()
	//B : SKPOLY
	cadHeader.SketchPolylines, err = this.sreader.ReadBit()
	//B : ANGDIR
	shortValue, err = this.sreader.ReadBitAsShort()
	cadHeader.AngularDirection = units.AngularDirection(shortValue)
	//B : SPLFRAME
	cadHeader.ShowSplineControlPoints, err = this.sreader.ReadBit()

	//R13-R14 Only (stored in registry from R15 onwards):
	if this.IsR13_14Only() {
		//B : ATTREQ
		this.sreader.ReadBit()
		//B : ATTDIA
		this.sreader.ReadBit()
	}

	//Common:
	//B: MIRRTEXT
	cadHeader.MirrorText, err = this.sreader.ReadBit()
	//B : WORLDVIEW
	cadHeader.WorldView, err = this.sreader.ReadBit()

	//R13 - R14 Only:
	if this.IsR13_14Only() {
		//B: WIREFRAME Undocumented.
		this.sreader.ReadBit()
	}

	//Common:
	//B: TILEMODE
	cadHeader.ShowModelSpace, err = this.sreader.ReadBit()
	//B : PLIMCHECK
	cadHeader.PaperSpaceLimitsChecking, err = this.sreader.ReadBit()
	//B : VISRETAIN
	cadHeader.RetainXRefDependentVisibilitySettings, err = this.sreader.ReadBit()

	//R13 - R14 Only(stored in registry from R15 onwards):
	if this.IsR13_14Only() {
		//B : DELOBJ
		this.sreader.ReadBit()
	}

	//Common:
	//B: DISPSILH
	cadHeader.DisplaySilhouetteCurves, err = this.sreader.ReadBit()
	//B : PELLIPSE(not present in DXF)
	cadHeader.CreateEllipseAsPolyline, err = this.sreader.ReadBit()
	//BS: PROXYGRAPHICS
	cadHeader.ProxyGraphics, err = this.sreader.ReadBitShortAsBool()

	//R13-R14 Only (stored in registry from R15 onwards):
	if this.IsR13_14Only() {
		//BS : DRAGMODE
		this.sreader.ReadBitShort()
	}

	//Common:
	//BS: TREEDEPTH
	cadHeader.SpatialIndexMaxTreeDepth, err = this.sreader.ReadBitShort()
	//BS : LUNITS
	shortValue, err = this.sreader.ReadBitShort()
	cadHeader.LinearUnitFormat = units.LinearUnitFormat(shortValue)
	//BS : LUPREC
	cadHeader.LinearUnitPrecision, err = this.sreader.ReadBitShort()
	//BS : AUNITS
	shortValue, err = this.sreader.ReadBitShort()
	cadHeader.AngularUnit = units.AngularUnitFormat(shortValue)
	//BS : AUPREC
	cadHeader.AngularUnitPrecision, err = this.sreader.ReadBitShort()

	//R13 - R14 Only Only(stored in registry from R15 onwards):
	if this.IsR13_14Only() {
		//BS: OSMODE
		shortValue, err = this.sreader.ReadBitShort()
		cadHeader.ObjectSnapMode = enums.ObjectSnapMode(shortValue)
	}

	//Common:
	//BS: ATTMODE
	shortValue, err = this.sreader.ReadBitShort()
	cadHeader.AttributeVisibility = enums.AttributeVisibilityMode(shortValue)

	//R13 - R14 Only Only(stored in registry from R15 onwards):
	if this.IsR13_14Only() {
		//BS: COORDS
		this.sreader.ReadBitShort()
	}

	//Common:
	//BS: PDMODE
	cadHeader.PointDisplayMode, err = this.sreader.ReadBitShort()

	//R13 - R14 Only Only(stored in registry from R15 onwards):
	if this.IsR13_14Only() {
		//BS: PICKSTYLE
		this.sreader.ReadBitShort()
	}

	//R2004 +:
	if this.IsR2004Plus() {
		//BL: Unknown
		longValue, err = this.sreader.ReadBitLong()
		//BL: Unknown
		longValue, err = this.sreader.ReadBitLong()
		//BL: Unknown
		longValue, err = this.sreader.ReadBitLong()
	}

	//Common:
	//BS : USERI1
	cadHeader.UserShort1, err = this.sreader.ReadBitShort()
	//BS : USERI2
	cadHeader.UserShort2, err = this.sreader.ReadBitShort()
	//BS : USERI3
	cadHeader.UserShort3, err = this.sreader.ReadBitShort()
	//BS : USERI4
	cadHeader.UserShort4, err = this.sreader.ReadBitShort()
	//BS : USERI5
	cadHeader.UserShort5, err = this.sreader.ReadBitShort()

	//BS: SPLINESEGS
	cadHeader.NumberOfSplineSegments, err = this.sreader.ReadBitShort()
	//BS : SURFU
	cadHeader.SurfaceDensityU, err = this.sreader.ReadBitShort()
	//BS : SURFV
	cadHeader.SurfaceDensityV, err = this.sreader.ReadBitShort()
	//BS : SURFTYPE
	cadHeader.SurfaceType, err = this.sreader.ReadBitShort()
	//BS : SURFTAB1
	cadHeader.SurfaceMeshTabulationCount1, err = this.sreader.ReadBitShort()
	//BS : SURFTAB2
	cadHeader.SurfaceMeshTabulationCount2, err = this.sreader.ReadBitShort()
	//BS : SPLINETYPE
	shortValue, err = this.sreader.ReadBitShort()
	cadHeader.SplineType = enums.SplineType(shortValue)
	//BS : SHADEDGE
	shortValue, err = this.sreader.ReadBitShort()
	cadHeader.ShadeEdge = enums.ShadeEdgeType(shortValue)
	//BS : SHADEDIF
	cadHeader.ShadeDiffuseToAmbientPercentage, err = this.sreader.ReadBitShort()
	//BS: UNITMODE
	cadHeader.UnitMode, err = this.sreader.ReadBitShort()
	//BS : MAXACTVP
	cadHeader.MaxViewportCount, err = this.sreader.ReadBitShort()
	//BS : ISOLINES
	cadHeader.SurfaceIsolineCount, err = this.sreader.ReadBitShort()
	//BS : CMLJUST
	shortValue, err = this.sreader.ReadBitShort()
	cadHeader.CurrentMultilineJustification = entities.VerticalAlignmentType(shortValue)
	//BS : TEXTQLTY
	cadHeader.TextQuality, err = this.sreader.ReadBitShort()
	//BD : LTSCALE
	cadHeader.LineTypeScale, err = this.sreader.ReadBitDouble()
	//BD : TEXTSIZE
	cadHeader.TextHeightDefault, err = this.sreader.ReadBitDouble()
	//BD : TRACEWID
	cadHeader.TraceWidthDefault, err = this.sreader.ReadBitDouble()
	//BD : SKETCHINC
	cadHeader.SketchIncrement, err = this.sreader.ReadBitDouble()
	//BD : FILLETRAD
	cadHeader.FilletRadius, err = this.sreader.ReadBitDouble()
	//BD : THICKNESS
	cadHeader.ThicknessDefault, err = this.sreader.ReadBitDouble()
	//BD : ANGBASE
	cadHeader.AngleBase, err = this.sreader.ReadBitDouble()
	//BD : PDSIZE
	cadHeader.PointDisplaySize, err = this.sreader.ReadBitDouble()
	//BD : PLINEWID
	cadHeader.PolylineWidthDefault, err = this.sreader.ReadBitDouble()
	//BD : USERR1
	cadHeader.UserDouble1, err = this.sreader.ReadBitDouble()
	//BD : USERR2
	cadHeader.UserDouble2, err = this.sreader.ReadBitDouble()
	//BD : USERR3
	cadHeader.UserDouble3, err = this.sreader.ReadBitDouble()
	//BD : USERR4
	cadHeader.UserDouble4, err = this.sreader.ReadBitDouble()
	//BD : USERR5
	cadHeader.UserDouble5, err = this.sreader.ReadBitDouble()
	//BD : CHAMFERA
	cadHeader.ChamferDistance1, err = this.sreader.ReadBitDouble()
	//BD : CHAMFERB
	cadHeader.ChamferDistance2, err = this.sreader.ReadBitDouble()
	//BD : CHAMFERC
	cadHeader.ChamferLength, err = this.sreader.ReadBitDouble()
	//BD : CHAMFERD
	cadHeader.ChamferAngle, err = this.sreader.ReadBitDouble()
	//BD : FACETRES
	cadHeader.FacetResolution, err = this.sreader.ReadBitDouble()
	//BD : CMLSCALE
	cadHeader.CurrentMultilineScale, err = this.sreader.ReadBitDouble()
	//BD : CELTSCALE
	cadHeader.CurrentEntityLinetypeScale, err = this.sreader.ReadBitDouble()

	//TV: MENUNAME
	cadHeader.MenuFileName, err = this.sreader.ReadVariableText()

	//Common:
	//BL: TDCREATE(Julian day)
	//BL: TDCREATE(Milliseconds into the day)
	cadHeader.CreateDateTime, err = this.sreader.ReadDateTime()
	//BL: TDUPDATE(Julian day)
	//BL: TDUPDATE(Milliseconds into the day)
	cadHeader.UpdateDateTime, err = this.sreader.ReadDateTime()

	//R2004 +:
	if this.IsR2004Plus() {
		//BL : Unknown
		longValue, err = this.sreader.ReadBitLong()
		//BL : Unknown
		longValue, err = this.sreader.ReadBitLong()
		//BL : Unknown
		longValue, err = this.sreader.ReadBitLong()
	}

	//Common:
	//BL: TDINDWG(Days)
	//BL: TDINDWG(Milliseconds into the day)
	cadHeader.TotalEditingTime, err = this.sreader.ReadTimeSpan()
	//BL: TDUSRTIMER(Days)
	//BL: TDUSRTIMER(Milliseconds into the day)
	cadHeader.UserElapsedTimeSpan, err = this.sreader.ReadTimeSpan()

	// //CMC : CECOLOR
	// cadHeader.CurrentEntityColor = this.sreader.ReadCmColor();

	// //H : HANDSEED The next handle, with an 8-bit length specifier preceding the handle
	// //bytes (standard hex handle form) (code 0). The HANDSEED is not part of the handle
	// //stream, but of the normal data stream (relevant for R21 and later).
	// cadHeader.HandleSeed = mainreader.HandleReference();

	// //H : CLAYER (hard pointer)
	// objectPointers.CLAYER = this.sreader.HandleReference();
	// //H: TEXTSTYLE(hard pointer)
	// objectPointers.TEXTSTYLE = this.sreader.HandleReference();
	// //H: CELTYPE(hard pointer)
	// objectPointers.CELTYPE = this.sreader.HandleReference();

	// //R2007 + Only:
	// if this.IsR2007Plus() {
	// 	//H: CMATERIAL(hard pointer)
	// 	objectPointers.CMATERIAL = this.sreader.HandleReference();
	// }

	// //Common:
	// //H: DIMSTYLE (hard pointer)
	// objectPointers.DIMSTYLE = this.sreader.HandleReference();
	// //H: CMLSTYLE (hard pointer)
	// objectPointers.CMLSTYLE = this.sreader.HandleReference();

	// //R2000+ Only:
	// if this.IsR2000Plus() {
	// 	//BD: PSVPSCALE
	// 	cadHeader.ViewportDefaultViewScaleFactor, err = this.sreader.ReadBitDouble();
	// }

	// //Common:
	// //3BD: INSBASE(PSPACE)
	// cadHeader.PaperSpaceInsertionBase = this.sreader.Read3BitDouble();
	// //3BD: EXTMIN(PSPACE)
	// cadHeader.PaperSpaceExtMin = this.sreader.Read3BitDouble();
	// //3BD: EXTMAX(PSPACE)
	// cadHeader.PaperSpaceExtMax = this.sreader.Read3BitDouble();
	// //2RD: LIMMIN(PSPACE)
	// cadHeader.PaperSpaceLimitsMin = this.sreader.Read2RawDouble();
	// //2RD: LIMMAX(PSPACE)
	// cadHeader.PaperSpaceLimitsMax = this.sreader.Read2RawDouble();
	// //BD: ELEVATION(PSPACE)
	// cadHeader.PaperSpaceElevation = this.sreader.ReadBitDouble();
	// //3BD: UCSORG(PSPACE)
	// cadHeader.PaperSpaceUcsOrigin = this.sreader.Read3BitDouble();
	// //3BD: UCSXDIR(PSPACE)
	// cadHeader.PaperSpaceUcsXAxis = this.sreader.Read3BitDouble();
	// //3BD: UCSYDIR(PSPACE)
	// cadHeader.PaperSpaceUcsYAxis = this.sreader.Read3BitDouble();

	// //H: UCSNAME (PSPACE) (hard pointer)
	// objectPointers.UCSNAME_PSPACE = this.sreader.HandleReference();

	// //R2000+ Only:
	// if this.IsR2000Plus() {
	// 	//H : PUCSORTHOREF (hard pointer)
	// 	objectPointers.PUCSORTHOREF = this.sreader.HandleReference();

	// 	//BS : PUCSORTHOVIEW	??
	// 	this.sreader.ReadBitShort();

	// 	//H: PUCSBASE(hard pointer)
	// 	objectPointers.PUCSBASE = this.sreader.HandleReference();

	// 	//3BD: PUCSORGTOP
	// 	cadHeader.PaperSpaceOrthographicTopDOrigin, err = this.sreader.Read3BitDouble();
	// 	//3BD: PUCSORGBOTTOM
	// 	cadHeader.PaperSpaceOrthographicBottomDOrigin, err = this.sreader.Read3BitDouble();
	// 	//3BD: PUCSORGLEFT
	// 	cadHeader.PaperSpaceOrthographicLeftDOrigin, err = this.sreader.Read3BitDouble();
	// 	//3BD: PUCSORGRIGHT
	// 	cadHeader.PaperSpaceOrthographicRightDOrigin, err = this.sreader.Read3BitDouble();
	// 	//3BD: PUCSORGFRONT
	// 	cadHeader.PaperSpaceOrthographicFrontDOrigin, err = this.sreader.Read3BitDouble();
	// 	//3BD: PUCSORGBACK
	// 	cadHeader.PaperSpaceOrthographicBackDOrigin, err = this.sreader.Read3BitDouble();
	// }

	// //Common:
	// //3BD: INSBASE(MSPACE)
	// cadHeader.ModelSpaceInsertionBase, err = this.sreader.Read3BitDouble();
	// //3BD: EXTMIN(MSPACE)
	// cadHeader.ModelSpaceExtMin, err = this.sreader.Read3BitDouble();
	// //3BD: EXTMAX(MSPACE)
	// cadHeader.ModelSpaceExtMax, err = this.sreader.Read3BitDouble();
	// //2RD: LIMMIN(MSPACE)
	// cadHeader.ModelSpaceLimitsMin , err= this.sreader.Read2RawDouble();
	// //2RD: LIMMAX(MSPACE)
	// cadHeader.ModelSpaceLimitsMax, err = this.sreader.Read2RawDouble();
	// //BD: ELEVATION(MSPACE)
	// cadHeader.Elevation, err = this.sreader.ReadBitDouble();
	// //3BD: UCSORG(MSPACE)
	// cadHeader.ModelSpaceOrigin, err = this.sreader.Read3BitDouble();
	// //3BD: UCSXDIR(MSPACE)
	// cadHeader.ModelSpaceXAxis, err = this.sreader.Read3BitDouble();
	// //3BD: UCSYDIR(MSPACE)
	// cadHeader.ModelSpaceYAxis, err = this.sreader.Read3BitDouble();

	// //H: UCSNAME(MSPACE)(hard pointer)
	// objectPointers.UCSNAME_MSPACE = this.sreader.HandleReference();

	// //R2000 + Only:
	// if this.IsR2000Plus() {
	// 	//H: UCSORTHOREF(hard pointer)
	// 	objectPointers.UCSORTHOREF = this.sreader.HandleReference();

	// 	//BS: UCSORTHOVIEW	??
	// 	this.sreader.ReadBitShort();

	// 	//H : UCSBASE(hard pointer)
	// 	objectPointers.UCSBASE, err = this.sreader.HandleReference();

	// 	//3BD: UCSORGTOP
	// 	cadHeader.ModelSpaceOrthographicTopDOrigin, err = this.sreader.Read3BitDouble();
	// 	//3BD: UCSORGBOTTOM
	// 	cadHeader.ModelSpaceOrthographicBottomDOrigin, err = this.sreader.Read3BitDouble();
	// 	//3BD: UCSORGLEFT
	// 	cadHeader.ModelSpaceOrthographicLeftDOrigin, err = this.sreader.Read3BitDouble();
	// 	//3BD: UCSORGRIGHT
	// 	cadHeader.ModelSpaceOrthographicRightDOrigin, err = this.sreader.Read3BitDouble();
	// 	//3BD: UCSORGFRONT
	// 	cadHeader.ModelSpaceOrthographicFrontDOrigin, err = this.sreader.Read3BitDouble();
	// 	//3BD: UCSORGBACK
	// 	cadHeader.ModelSpaceOrthographicBackDOrigin, err = this.sreader.Read3BitDouble();

	// 	//TV : DIMPOST
	// 	cadHeader.DimensionPostFix, err = this.sreader.ReadVariableText();
	// 	//TV : DIMAPOST
	// 	cadHeader.DimensionAlternateDimensioningSuffix, err = this.sreader.ReadVariableText();
	// }

	// //R13-R14 Only:
	// if this.IsR13_14Only() {
	// 	//B: DIMTOL
	// 	cadHeader.DimensionGenerateTolerances, err = this.sreader.ReadBit();
	// 	//B : DIMLIM
	// 	cadHeader.DimensionLimitsGeneration, err = this.sreader.ReadBit();
	// 	//B : DIMTIH
	// 	cadHeader.DimensionTextInsideHorizontal, err = this.sreader.ReadBit();
	// 	//B : DIMTOH
	// 	cadHeader.DimensionTextOutsideHorizontal , err= this.sreader.ReadBit();
	// 	//B : DIMSE1
	// 	cadHeader.DimensionSuppressFirstExtensionLine, err = this.sreader.ReadBit();
	// 	//B : DIMSE2
	// 	cadHeader.DimensionSuppressSecondExtensionLine, err = this.sreader.ReadBit();
	// 	//B : DIMALT
	// 	cadHeader.DimensionAlternateUnitDimensioning, err = this.sreader.ReadBit();
	// 	//B : DIMTOFL
	// 	cadHeader.DimensionTextOutsideExtensions, err = this.sreader.ReadBit();
	// 	//B : DIMSAH
	// 	cadHeader.DimensionSeparateArrowBlocks, err = this.sreader.ReadBit();
	// 	//B : DIMTIX
	// 	cadHeader.DimensionTextInsideExtensions , err= this.sreader.ReadBit();
	// 	//B : DIMSOXD
	// 	cadHeader.DimensionSuppressOutsideExtensions , err= this.sreader.ReadBit();
	// 	//RC : DIMALTD
	// 	cadHeader.DimensionAlternateUnitDecimalPlaces, err= int16this.sreader.ReadRawChar();
	// 	//RC : DIMZIN
	// 	byteValue, err =this.sreader.ReadRawChar()
	// 	cadHeader.DimensionZeroHandling = tables.ZeroHandling(byteValue);
	// 	//B : DIMSD1
	// 	cadHeader.DimensionSuppressFirstDimensionLine , err= this.sreader.ReadBit();
	// 	//B : DIMSD2
	// 	cadHeader.DimensionSuppressSecondDimensionLine , err= this.sreader.ReadBit();
	// 	//RC : DIMTOLJ
	// 	byteValue, err = this.sreader.ReadRawChar()
	// 	cadHeader.DimensionToleranceAlignment = tables.ToleranceAlignment(byteValue);
	// 	//RC : DIMJUST
	// 	byteValue, err = this.sreader.ReadRawChar()
	// 	cadHeader.DimensionTextHorizontalAlignment = tables.DimensionTextHorizontalAlignment(byteValue);
	// 	//RC : DIMFIT
	// 	cadHeader.DimensionFit = this.sreader.ReadRawChar();
	// 	//B : DIMUPT
	// 	cadHeader.DimensionCursorUpdate = this.sreader.ReadBit();
	// 	//RC : DIMTZIN
	// 	cadHeader.DimensionToleranceZeroHandling = (Tables.ZeroHandling)this.sreader.ReadRawChar();
	// 	//RC: DIMALTZ
	// 	cadHeader.DimensionAlternateUnitZeroHandling = (Tables.ZeroHandling)this.sreader.ReadRawChar();
	// 	//RC : DIMALTTZ
	// 	cadHeader.DimensionAlternateUnitToleranceZeroHandling = (Tables.ZeroHandling)this.sreader.ReadRawChar();
	// 	//RC : DIMTAD
	// 	cadHeader.DimensionTextVerticalAlignment = (Tables.DimensionTextVerticalAlignment)this.sreader.ReadRawChar();
	// 	//BS : DIMUNIT
	// 	cadHeader.DimensionUnit = this.sreader.ReadBitShort();
	// 	//BS : DIMAUNIT
	// 	cadHeader.DimensionAngularDimensionDecimalPlaces = this.sreader.ReadBitShort();
	// 	//BS : DIMDEC
	// 	cadHeader.DimensionDecimalPlaces = this.sreader.ReadBitShort();
	// 	//BS : DIMTDEC
	// 	cadHeader.DimensionToleranceDecimalPlaces = this.sreader.ReadBitShort();
	// 	//BS : DIMALTU
	// 	cadHeader.DimensionAlternateUnitFormat = (LinearUnitFormat)this.sreader.ReadBitShort();
	// 	//BS : DIMALTTD
	// 	cadHeader.DimensionAlternateUnitToleranceDecimalPlaces = this.sreader.ReadBitShort();
	// 	//H : DIMTXSTY(hard pointer)
	// 	objectPointers.DIMTXSTY = this.sreader.HandleReference();
	// }

	// //Common:
	// //BD: DIMSCALE
	// cadHeader.DimensionScaleFactor = this.sreader.ReadBitDouble();
	// //BD : DIMASZ
	// cadHeader.DimensionArrowSize = this.sreader.ReadBitDouble();
	// //BD : DIMEXO
	// cadHeader.DimensionExtensionLineOffset = this.sreader.ReadBitDouble();
	// //BD : DIMDLI
	// cadHeader.DimensionLineIncrement = this.sreader.ReadBitDouble();
	// //BD : DIMEXE
	// cadHeader.DimensionExtensionLineExtension = this.sreader.ReadBitDouble();
	// //BD : DIMRND
	// cadHeader.DimensionRounding = this.sreader.ReadBitDouble();
	// //BD : DIMDLE
	// cadHeader.DimensionLineExtension = this.sreader.ReadBitDouble();
	// //BD : DIMTP
	// cadHeader.DimensionPlusTolerance = this.sreader.ReadBitDouble();
	// //BD : DIMTM
	// cadHeader.DimensionMinusTolerance = this.sreader.ReadBitDouble();

	// //R2007 + Only:
	// if (R2007Plus)
	// {
	// 	//BD: DIMFXL
	// 	cadHeader.DimensionFixedExtensionLineLength = this.sreader.ReadBitDouble();
	// 	//BD : DIMJOGANG
	// 	cadHeader.DimensionJoggedRadiusDimensionTransverseSegmentAngle = this.sreader.ReadBitDouble();
	// 	//BS : DIMTFILL
	// 	cadHeader.DimensionTextBackgroundFillMode = (Tables.DimensionTextBackgroundFillMode)this.sreader.ReadBitShort();
	// 	//CMC : DIMTFILLCLR
	// 	cadHeader.DimensionTextBackgroundColor = this.sreader.ReadCmColor();
	// }

	// //R2000 + Only:
	// if (R2000Plus)
	// {
	// 	//B: DIMTOL
	// 	cadHeader.DimensionGenerateTolerances = this.sreader.ReadBit();
	// 	//B : DIMLIM
	// 	cadHeader.DimensionLimitsGeneration = this.sreader.ReadBit();
	// 	//B : DIMTIH
	// 	cadHeader.DimensionTextInsideHorizontal = this.sreader.ReadBit();
	// 	//B : DIMTOH
	// 	cadHeader.DimensionTextOutsideHorizontal = this.sreader.ReadBit();
	// 	//B : DIMSE1
	// 	cadHeader.DimensionSuppressFirstExtensionLine = this.sreader.ReadBit();
	// 	//B : DIMSE2
	// 	cadHeader.DimensionSuppressSecondExtensionLine = this.sreader.ReadBit();
	// 	//BS : DIMTAD
	// 	cadHeader.DimensionTextVerticalAlignment = (Tables.DimensionTextVerticalAlignment)(char)this.sreader.ReadBitShort();
	// 	//BS : DIMZIN
	// 	cadHeader.DimensionZeroHandling = (Tables.ZeroHandling)(char)this.sreader.ReadBitShort();
	// 	//BS : DIMAZIN
	// 	cadHeader.DimensionAngularZeroHandling = (Tables.ZeroHandling)this.sreader.ReadBitShort();
	// }

	// //R2007 + Only:
	// if this.IsR2007Plus() {
	// 	//BS: DIMARCSYM
	// 	cadHeader.DimensionArcLengthSymbolPosition = (Tables.ArcLengthSymbolPosition)this.sreader.ReadBitShort();
	// }

	// //Common:
	// //BD: DIMTXT
	// cadHeader.DimensionTextHeight = this.sreader.ReadBitDouble();
	// //BD : DIMCEN
	// cadHeader.DimensionCenterMarkSize = this.sreader.ReadBitDouble();
	// //BD: DIMTSZ
	// cadHeader.DimensionTickSize = this.sreader.ReadBitDouble();
	// //BD : DIMALTF
	// cadHeader.DimensionAlternateUnitScaleFactor = this.sreader.ReadBitDouble();
	// //BD : DIMLFAC
	// cadHeader.DimensionLinearScaleFactor = this.sreader.ReadBitDouble();
	// //BD : DIMTVP
	// cadHeader.DimensionTextVerticalPosition = this.sreader.ReadBitDouble();
	// //BD : DIMTFAC
	// cadHeader.DimensionToleranceScaleFactor = this.sreader.ReadBitDouble();
	// //BD : DIMGAP
	// cadHeader.DimensionLineGap = this.sreader.ReadBitDouble();

	// //R13 - R14 Only:
	// if this.IsR13_14Only() {
	// 	//T: DIMPOST
	// 	cadHeader.DimensionPostFix = this.sreader.ReadVariableText();
	// 	//T : DIMAPOST
	// 	cadHeader.DimensionAlternateDimensioningSuffix = this.sreader.ReadVariableText();
	// 	//T : DIMBLK
	// 	cadHeader.DimensionBlockName = this.sreader.ReadVariableText();
	// 	//T : DIMBLK1
	// 	cadHeader.DimensionBlockNameFirst = this.sreader.ReadVariableText();
	// 	//T : DIMBLK2
	// 	cadHeader.DimensionBlockNameSecond = this.sreader.ReadVariableText();
	// }

	// //R2000 + Only:
	// if this.R2000Plus() {
	// 	//BD: DIMALTRND
	// 	cadHeader.DimensionAlternateUnitRounding = this.sreader.ReadBitDouble();
	// 	//B : DIMALT
	// 	cadHeader.DimensionAlternateUnitDimensioning = this.sreader.ReadBit();
	// 	//BS : DIMALTD
	// 	cadHeader.DimensionAlternateUnitDecimalPlaces = int16(char)this.sreader.ReadBitShort();
	// 	//B : DIMTOFL
	// 	cadHeader.DimensionTextOutsideExtensions = this.sreader.ReadBit();
	// 	//B : DIMSAH
	// 	cadHeader.DimensionSeparateArrowBlocks = this.sreader.ReadBit();
	// 	//B : DIMTIX
	// 	cadHeader.DimensionTextInsideExtensions = this.sreader.ReadBit();
	// 	//B : DIMSOXD
	// 	cadHeader.DimensionSuppressOutsideExtensions = this.sreader.ReadBit();
	// }

	// //Common:
	// //CMC: DIMCLRD
	// cadHeader.DimensionLineColor = this.sreader.ReadCmColor();
	// //CMC : DIMCLRE
	// cadHeader.DimensionExtensionLineColor = this.sreader.ReadCmColor();
	// //CMC : DIMCLRT
	// cadHeader.DimensionTextColor = this.sreader.ReadCmColor();

	// //R2000 + Only:
	// if this.IsR2000Plus() {
	// 	//BS: DIMADEC
	// 	cadHeader.DimensionAngularDimensionDecimalPlaces = this.sreader.ReadBitShort();
	// 	//BS : DIMDEC
	// 	cadHeader.DimensionDecimalPlaces = this.sreader.ReadBitShort();
	// 	//BS : DIMTDEC
	// 	cadHeader.DimensionToleranceDecimalPlaces = this.sreader.ReadBitShort();
	// 	//BS : DIMALTU
	// 	cadHeader.DimensionAlternateUnitFormat = (LinearUnitFormat)this.sreader.ReadBitShort();
	// 	//BS : DIMALTTD
	// 	cadHeader.DimensionAlternateUnitToleranceDecimalPlaces = this.sreader.ReadBitShort();
	// 	//BS : DIMAUNIT
	// 	cadHeader.DimensionAngularUnit = (AngularUnitFormat)this.sreader.ReadBitShort();
	// 	//BS : DIMFRAC
	// 	cadHeader.DimensionFractionFormat = (Tables.FractionFormat)this.sreader.ReadBitShort();
	// 	//BS : DIMLUNIT
	// 	cadHeader.DimensionLinearUnitFormat = (LinearUnitFormat)this.sreader.ReadBitShort();
	// 	//BS : DIMDSEP
	// 	cadHeader.DimensionDecimalSeparator = (char)this.sreader.ReadBitShort();
	// 	//BS : DIMTMOVE
	// 	cadHeader.DimensionTextMovement = (Tables.TextMovement)this.sreader.ReadBitShort();
	// 	//BS : DIMJUST
	// 	cadHeader.DimensionTextHorizontalAlignment = (Tables.DimensionTextHorizontalAlignment)(char)this.sreader.ReadBitShort();
	// 	//B : DIMSD1
	// 	cadHeader.DimensionSuppressFirstExtensionLine = this.sreader.ReadBit();
	// 	//B : DIMSD2
	// 	cadHeader.DimensionSuppressSecondExtensionLine = this.sreader.ReadBit();
	// 	//BS : DIMTOLJ
	// 	cadHeader.DimensionToleranceAlignment = (Tables.ToleranceAlignment)(char)this.sreader.ReadBitShort();
	// 	//BS : DIMTZIN
	// 	cadHeader.DimensionToleranceZeroHandling = (Tables.ZeroHandling)(char)this.sreader.ReadBitShort();
	// 	//BS: DIMALTZ
	// 	cadHeader.DimensionAlternateUnitZeroHandling = (Tables.ZeroHandling)(char)this.sreader.ReadBitShort();
	// 	//BS : DIMALTTZ
	// 	cadHeader.DimensionAlternateUnitToleranceZeroHandling = (Tables.ZeroHandling)(char)this.sreader.ReadBitShort();
	// 	//B : DIMUPT
	// 	cadHeader.DimensionCursorUpdate = this.sreader.ReadBit();
	// 	//BS : DIMATFIT
	// 	cadHeader.DimensionDimensionTextArrowFit = this.sreader.ReadBitShort();
	// }

	// //R2007 + Only:
	// if this.IsR2007Plus() {
	// 	//B: DIMFXLON
	// 	cadHeader.DimensionIsExtensionLineLengthFixed = this.sreader.ReadBit();
	// }

	// //R2010 + Only:
	// if this.IsR2010Plus() {
	// 	//B: DIMTXTDIRECTION
	// 	cadHeader.DimensionTextDirection = this.sreader.ReadBit() ? Tables.TextDirection.RightToLeft : Tables.TextDirection.LeftToRight;
	// 	//BD : DIMALTMZF
	// 	cadHeader.DimensionAltMzf = this.sreader.ReadBitDouble();
	// 	//T : DIMALTMZS
	// 	cadHeader.DimensionAltMzs = this.sreader.ReadVariableText();
	// 	//BD : DIMMZF
	// 	cadHeader.DimensionMzf = this.sreader.ReadBitDouble();
	// 	//T : DIMMZS
	// 	cadHeader.DimensionMzs = this.sreader.ReadVariableText();
	// }

	// //R2000 + Only:
	// if this.IsR2000Plus() {
	// 	//H: DIMTXSTY(hard pointer)
	// 	objectPointers.DIMTXSTY = this.sreader.HandleReference();
	// 	//H: DIMLDRBLK(hard pointer)
	// 	objectPointers.DIMLDRBLK = this.sreader.HandleReference();
	// 	//H: DIMBLK(hard pointer)
	// 	objectPointers.DIMBLK = this.sreader.HandleReference();
	// 	//H: DIMBLK1(hard pointer)
	// 	objectPointers.DIMBLK1 = this.sreader.HandleReference();
	// 	//H: DIMBLK2(hard pointer)
	// 	objectPointers.DIMBLK2 = this.sreader.HandleReference();
	// }

	// //R2007+ Only:
	// if this.IsR2007Plus() {
	// 	//H : DIMLTYPE (hard pointer)
	// 	objectPointers.DIMLTYPE = this.sreader.HandleReference();
	// 	//H: DIMLTEX1(hard pointer)
	// 	objectPointers.DIMLTEX1 = this.sreader.HandleReference();
	// 	//H: DIMLTEX2(hard pointer)
	// 	objectPointers.DIMLTEX2 = this.sreader.HandleReference();
	// }

	// //R2000+ Only:
	// if this.IsR2000Plus() {
	// 	//BS: DIMLWD
	// 	cadHeader.DimensionLineWeight = (LineweightType)this.sreader.ReadBitShort();
	// 	//BS : DIMLWE
	// 	cadHeader.ExtensionLineWeight = (LineweightType)this.sreader.ReadBitShort();
	// }

	// //H: BLOCK CONTROL OBJECT(hard owner)
	// objectPointers.BLOCK_CONTROL_OBJECT = this.sreader.HandleReference();
	// //H: LAYER CONTROL OBJECT(hard owner)
	// objectPointers.LAYER_CONTROL_OBJECT = this.sreader.HandleReference();
	// //H: STYLE CONTROL OBJECT(hard owner)
	// objectPointers.STYLE_CONTROL_OBJECT = this.sreader.HandleReference();
	// //H: LINETYPE CONTROL OBJECT(hard owner)
	// objectPointers.LINETYPE_CONTROL_OBJECT = this.sreader.HandleReference();
	// //H: VIEW CONTROL OBJECT(hard owner)
	// objectPointers.VIEW_CONTROL_OBJECT = this.sreader.HandleReference();
	// //H: UCS CONTROL OBJECT(hard owner)
	// objectPointers.UCS_CONTROL_OBJECT = this.sreader.HandleReference();
	// //H: VPORT CONTROL OBJECT(hard owner)
	// objectPointers.VPORT_CONTROL_OBJECT = this.sreader.HandleReference();
	// //H: APPID CONTROL OBJECT(hard owner)
	// objectPointers.APPID_CONTROL_OBJECT = this.sreader.HandleReference();
	// //H: DIMSTYLE CONTROL OBJECT(hard owner)
	// objectPointers.DIMSTYLE_CONTROL_OBJECT = this.sreader.HandleReference();

	// //R13 - R15 Only:
	// if this.IsR13_15Only() {
	// 	//H: VIEWPORT ENTITY HEADER CONTROL OBJECT(hard owner)
	// 	objectPointers.VIEWPORT_ENTITY_HEADER_CONTROL_OBJECT = this.sreader.HandleReference();
	// }

	// //Common:
	// //H: DICTIONARY(ACAD_GROUP)(hard pointer)
	// objectPointers.DICTIONARY_ACAD_GROUP = this.sreader.HandleReference();
	// //H: DICTIONARY(ACAD_MLINESTYLE)(hard pointer)
	// objectPointers.DICTIONARY_ACAD_MLINESTYLE = this.sreader.HandleReference();

	// //H : DICTIONARY (NAMED OBJECTS) (hard owner)
	// objectPointers.DICTIONARY_NAMED_OBJECTS = this.sreader.HandleReference();

	// //R2000+ Only:
	// if this.IsR2000Plus() {
	// 	//BS: TSTACKALIGN, default = 1(not present in DXF)
	// 	cadHeader.StackedTextAlignment = this.sreader.ReadBitShort();
	// 	//BS: TSTACKSIZE, default = 70(not present in DXF)
	// 	cadHeader.StackedTextSizePercentage = this.sreader.ReadBitShort();

	// 	//TV: HYPERLINKBASE
	// 	cadHeader.HyperLinkBase = this.sreader.ReadVariableText();
	// 	//TV : STYLESHEET
	// 	cadHeader.StyleSheetName = this.sreader.ReadVariableText();

	// 	//H : DICTIONARY(LAYOUTS)(hard pointer)
	// 	objectPointers.DICTIONARY_LAYOUTS = this.sreader.HandleReference();
	// 	//H: DICTIONARY(PLOTSETTINGS)(hard pointer)
	// 	objectPointers.DICTIONARY_PLOTSETTINGS = this.sreader.HandleReference();
	// 	//H: DICTIONARY(PLOTSTYLES)(hard pointer)
	// 	objectPointers.DICTIONARY_PLOTSTYLES = this.sreader.HandleReference();
	// }

	// //R2004 +:
	// if this.IsR2004Plus() {
	// 	//H: DICTIONARY (MATERIALS) (hard pointer)
	// 	objectPointers.DICTIONARY_MATERIALS = this.sreader.HandleReference();
	// 	//H: DICTIONARY (COLORS) (hard pointer)
	// 	objectPointers.DICTIONARY_COLORS = this.sreader.HandleReference();
	// }

	// //R2007 +:
	// if this.IsR2007Plus() {
	// 	//H: DICTIONARY(VISUALSTYLE)(hard pointer)
	// 	objectPointers.DICTIONARY_VISUALSTYLE = this.sreader.HandleReference();

	// 	//R2013+:
	// 	if (this.R2013Plus)
	// 		//H : UNKNOWN (hard pointer)
	// 		objectPointers.DICTIONARY_VISUALSTYLE = this.sreader.HandleReference();
	// }

	// //R2000 +:
	// if this.IsR2000Plus() {
	// 	//BL: Flags:
	// 	var flags int32
	// 	flags, err = this.sreader.ReadBitLong()
	// 	//CELWEIGHT Flags & 0x001F
	// 	cadHeader.CurrentEntityLineWeight = entities.LineweightType(flags & 0x1F);
	// 	//ENDCAPS Flags & 0x0060
	// 	cadHeader.EndCaps = int16(flags & 0x60);
	// 	//JOINSTYLE Flags & 0x0180
	// 	cadHeader.JoinStyle = int16(flags & 0x180);
	// 	//LWDISPLAY!(Flags & 0x0200)
	// 	cadHeader.DisplayLineWeight = (flags & 0x200) == 1;
	// 	//XEDIT!(Flags & 0x0400)
	// 	cadHeader.XEdit = int16(flags & 0x400) == 1;
	// 	//EXTNAMES Flags & 0x0800
	// 	cadHeader.ExtendedNames = (flags & 0x800) == 1;
	// 	//PSTYLEMODE Flags & 0x2000
	// 	cadHeader.PlotStyleMode = int16(flags & 0x2000);
	// 	//OLESTARTUP Flags & 0x4000
	// 	cadHeader.LoadOLEObject = (flags & 0x4000) == 1;

	// 	//BS: INSUNITS
	// 	shortValue, err = this.sreader.ReadBitShort()
	// 	cadHeader.InsUnits = units.UnitsType(shortValue);
	// 	//BS : CEPSNTYPE
	// 	shortValue, err = this.sreader.ReadBitShort()
	// 	cadHeader.CurrentEntityPlotStyle = enums.EntityPlotStyleType(shortValue)

	// 	if cadHeader.CurrentEntityPlotStyle == EntityPlotStyleType.ByObjectId {
	// 		//H: CPSNID(present only if CEPSNTYPE == 3) (hard pointer)
	// 		objectPointers.CPSNID = this.sreader.HandleReference();
	// 	}

	// 	//TV: FINGERPRINTGUID
	// 	cadHeader.FingerPrintGuid = this.sreader.ReadVariableText();
	// 	//TV : VERSIONGUID
	// 	cadHeader.VersionGuid = this.sreader.ReadVariableText();
	// }

	// //R2004 +:
	// if this.IsR2004Plus() {
	// 	//RC: SORTENTS
	// 	cadHeader.EntitySortingFlags = (ObjectSortingFlags);
	// 	//RC : INDEXCTL
	// 	cadHeader.IndexCreationFlags = (IndexCreationFlags)this.sreader.ReadByte();
	// 	//RC : HIDETEXT
	// 	cadHeader.HideText = this.sreader.ReadByte();
	// 	//RC : XCLIPFRAME, before R2010 the value can be 0 or 1 only.
	// 	cadHeader.ExternalReferenceClippingBoundaryType = this.sreader.ReadByte();
	// 	//RC : DIMASSOC
	// 	cadHeader.DimensionAssociativity = (DimensionAssociation)this.sreader.ReadByte();
	// 	//RC : HALOGAP
	// 	cadHeader.HaloGapPercentage = this.sreader.ReadByte();
	// 	//BS : OBSCUREDCOLOR
	// 	cadHeader.ObscuredColor = new Color(this.sreader.ReadBitShort());
	// 	//BS : INTERSECTIONCOLOR
	// 	cadHeader.InterfereColor = new Color(this.sreader.ReadBitShort());
	// 	//RC : OBSCUREDLTYPE
	// 	cadHeader.ObscuredType = this.sreader.ReadByte();
	// 	//RC: INTERSECTIONDISPLAY
	// 	cadHeader.IntersectionDisplay = this.sreader.ReadByte();

	// 	//TV : PROJECTNAME
	// 	cadHeader.ProjectName = this.sreader.ReadVariableText();
	// }

	// //Common:
	// //H: BLOCK_RECORD(*PAPER_SPACE)(hard pointer)
	// objectPointers.PAPER_SPACE = this.sreader.HandleReference();
	// //H: BLOCK_RECORD(*MODEL_SPACE)(hard pointer)
	// objectPointers.MODEL_SPACE = this.sreader.HandleReference();
	// //H: LTYPE(BYLAYER)(hard pointer)
	// objectPointers.BYLAYER = this.sreader.HandleReference();
	// //H: LTYPE(BYBLOCK)(hard pointer)
	// objectPointers.BYBLOCK = this.sreader.HandleReference();
	// //H: LTYPE(CONTINUOUS)(hard pointer)
	// objectPointers.CONTINUOUS = this.sreader.HandleReference();

	// //R2007 +:
	// if this.IsR2007Plus() {
	// 	//B: CAMERADISPLAY
	// 	cadHeader.CameraDisplayObjects = this.sreader.ReadBit();

	// 	//BL : unknown
	// 	longValue, err = this.sreader.ReadBitLong();
	// 	//BL : unknown
	// 	longValue, err = this.sreader.ReadBitLong();
	// 	//BD : unknown
	// 	this.sreader.ReadBitDouble();

	// 	//BD : STEPSPERSEC
	// 	cadHeader.StepsPerSecond = this.sreader.ReadBitDouble();
	// 	//BD : STEPSIZE
	// 	cadHeader.StepSize = this.sreader.ReadBitDouble();
	// 	//BD : 3DDWFPREC
	// 	cadHeader.Dw3DPrecision = this.sreader.ReadBitDouble();
	// 	//BD : LENSLENGTH
	// 	cadHeader.LensLength = this.sreader.ReadBitDouble();
	// 	//BD : CAMERAHEIGHT
	// 	cadHeader.CameraHeight = this.sreader.ReadBitDouble();
	// 	//RC : SOLIDHIST
	// 	cadHeader.SolidsRetainHistory = this.sreader.ReadRawChar();
	// 	//RC : SHOWHIST
	// 	cadHeader.ShowSolidsHistory = this.sreader.ReadRawChar();
	// 	//BD : PSOLWIDTH
	// 	cadHeader.SweptSolidWidth = this.sreader.ReadBitDouble();
	// 	//BD : PSOLHEIGHT
	// 	cadHeader.SweptSolidHeight = this.sreader.ReadBitDouble();
	// 	//BD : LOFTANG1
	// 	cadHeader.DraftAngleFirstCrossSection = this.sreader.ReadBitDouble();
	// 	//BD : LOFTANG2
	// 	cadHeader.DraftAngleSecondCrossSection = this.sreader.ReadBitDouble();
	// 	//BD : LOFTMAG1
	// 	cadHeader.DraftMagnitudeFirstCrossSection = this.sreader.ReadBitDouble();
	// 	//BD : LOFTMAG2
	// 	cadHeader.DraftMagnitudeSecondCrossSection = this.sreader.ReadBitDouble();
	// 	//BS : LOFTPARAM
	// 	cadHeader.SolidLoftedShape = this.sreader.ReadBitShort();
	// 	//RC : LOFTNORMALS
	// 	cadHeader.LoftedObjectNormals = this.sreader.ReadRawChar();
	// 	//BD : LATITUDE
	// 	cadHeader.Latitude = this.sreader.ReadBitDouble();
	// 	//BD : LONGITUDE
	// 	cadHeader.Longitude = this.sreader.ReadBitDouble();
	// 	//BD : NORTHDIRECTION
	// 	cadHeader.NorthDirection = this.sreader.ReadBitDouble();
	// 	//BL : TIMEZONE
	// 	cadHeader.TimeZone = this.sreader.ReadBitLong();
	// 	//RC : LIGHTGLYPHDISPLAY
	// 	cadHeader.DisplayLightGlyphs = this.sreader.ReadRawChar();
	// 	//RC : TILEMODELIGHTSYNCH	??
	// 	this.sreader.ReadRawChar();
	// 	//RC : DWFFRAME
	// 	cadHeader.DwgUnderlayFramesVisibility = this.sreader.ReadRawChar();
	// 	//RC : DGNFRAME
	// 	cadHeader.DgnUnderlayFramesVisibility = this.sreader.ReadRawChar();

	// 	//B : unknown
	// 	this.sreader.ReadBit();

	// 	//CMC : INTERFERECOLOR
	// 	cadHeader.InterfereColor = this.sreader.ReadCmColor();

	// 	//H : INTERFEREOBJVS(hard pointer)
	// 	objectPointers.INTERFEREOBJVS = this.sreader.HandleReference();
	// 	//H: INTERFEREVPVS(hard pointer)
	// 	objectPointers.INTERFEREVPVS = this.sreader.HandleReference();
	// 	//H: DRAGVS(hard pointer)
	// 	objectPointers.DRAGVS = this.sreader.HandleReference();

	// 	//RC: CSHADOW
	// 	cadHeader.ShadowMode = (ShadowMode)this.sreader.ReadByte();
	// 	//BD : SHADOWPLANELOCATION
	// 	cadHeader.ShadowPlaneLocation = this.sreader.ReadBitDouble();
	// }

	//Not necessary for the integrity of the data

	//Set the position at the end of the section
	this.sreader.SetPositionInBits(initialPos + size*8)
	this.sreader.ResetShift()

	//Ending sentinel: 0x30,0x84,0xE0,0xDC,0x02,0x21,0xC7,0x56,0xA0,0x83,0x97,0x47,0xB1,0x92,0xCC,0xA0
	if err = this.checkSentinel(this.sreader, section.EndSentinels[this.SectionName]); err != nil {
		return nil, err
	}

	_ = doubleValue
	_ = stringValue

	return objectPointers, nil
}

func NewDwgHeaderReader(dwgVersion version.ACadVersion, maintenanceVersion byte,
	sreader stream.IDwgStreamReader) *DwgHeaderReader {
	return &DwgHeaderReader{
		DwgSectionIO: DwgSectionIO{
			DwgVersion:  dwgVersion,
			SectionName: section.Header,
		},
		maintenanceVersion: maintenanceVersion,
		sreader:            sreader,
	}
}
