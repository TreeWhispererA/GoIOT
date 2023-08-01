package classes

type ProxyFlags uint16

const (
	PF_None                              ProxyFlags = 0
	PF_EraseAllowed                      ProxyFlags = 1
	PF_TransformAllowed                  ProxyFlags = 2
	PF_ColorChangeAllowed                ProxyFlags = 4
	PF_LayerChangeAllowed                ProxyFlags = 8
	PF_LinetypeChangeAllowed             ProxyFlags = 16
	PF_LinetypeScaleChangeAllowed        ProxyFlags = 32
	PF_VisibilityChangeAllowed           ProxyFlags = 64
	PF_CloningAllowed                    ProxyFlags = 128
	PF_LineweightChangeAllowed           ProxyFlags = 256
	PF_PlotStyleNameChangeAllowed        ProxyFlags = 512
	PF_AllOperationsExceptCloningAllowed ProxyFlags = 895
	PF_AllOperationsAllowed              ProxyFlags = 1023
	PF_DisablesProxyWarningDialog        ProxyFlags = 1024
	PF_R13FormatProxy                    ProxyFlags = 32768
)
