package utils

type CodePage int

const (
	CP_Unknown               CodePage = 0x00000000
	CP_Ibm037                CodePage = 37
	CP_Ibm437                CodePage = 437
	CP_Asmo708               CodePage = 708
	CP_Dos720                CodePage = 720
	CP_Ibm737                CodePage = 737
	CP_Ibm775                CodePage = 775
	CP_Ibm850                CodePage = 850
	CP_Ibm852                CodePage = 852
	CP_Ibm855                CodePage = 855
	CP_Ibm857                CodePage = 857
	CP_Ibm860                CodePage = 860
	CP_Ibm861                CodePage = 861
	CP_Dos862                CodePage = 862
	CP_Ibm863                CodePage = 863
	CP_Ibm864                CodePage = 864
	CP_Ibm865                CodePage = 865
	CP_Cp866                 CodePage = 866
	CP_Ibm869                CodePage = 869
	CP_Ibm870                CodePage = 870
	CP_Windows874            CodePage = 874
	CP_Cp875                 CodePage = 875
	CP_Shift_jis             CodePage = 932
	CP_Gb2312                CodePage = 936
	CP_Ksc5601               CodePage = 949
	CP_big5                  CodePage = 950
	CP_Ibm1026               CodePage = 1026
	CP_Ibm01047              CodePage = 1047
	CP_Ibm01140              CodePage = 1140
	CP_Ibm01141              CodePage = 1141
	CP_Ibm01142              CodePage = 1142
	CP_Ibm01143              CodePage = 1143
	CP_Ibm01144              CodePage = 1144
	CP_Ibm01145              CodePage = 1145
	CP_Ibm01146              CodePage = 1146
	CP_Ibm01147              CodePage = 1147
	CP_Ibm01148              CodePage = 1148
	CP_Ibm01149              CodePage = 1149
	CP_Utf16                 CodePage = 1200
	CP_UnicodeFFFE           CodePage = 1201
	CP_Windows1250           CodePage = 1250
	CP_Windows1251           CodePage = 1251
	CP_Windows1252           CodePage = 1252
	CP_Windows1253           CodePage = 1253
	CP_Windows1254           CodePage = 1254
	CP_Windows1255           CodePage = 1255
	CP_Windows1256           CodePage = 1256
	CP_Windows1257           CodePage = 1257
	CP_Windows1258           CodePage = 1258
	CP_Johab                 CodePage = 1361
	CP_Macintosh             CodePage = 10000
	CP_Xmacjapanese          CodePage = 10001
	CP_Xmacchinesetrad       CodePage = 10002
	CP_Xmackorean            CodePage = 10003
	CP_Xmacarabic            CodePage = 10004
	CP_Xmachebrew            CodePage = 10005
	CP_Xmacgreek             CodePage = 10006
	CP_Xmaccyrillic          CodePage = 10007
	CP_Xmacchinesesimp       CodePage = 10008
	CP_Xmacromanian          CodePage = 10010
	CP_Xmacukrainian         CodePage = 10017
	CP_Xmacthai              CodePage = 10021
	CP_Xmacce                CodePage = 10029
	CP_Xmacicelandic         CodePage = 10079
	CP_Xmacturkish           CodePage = 10081
	CP_Xmaccroatian          CodePage = 10082
	CP_Utf32                 CodePage = 12000
	CP_Utf32BE               CodePage = 12001
	CP_XChineseCNS           CodePage = 20000
	CP_Xcp20001              CodePage = 20001
	CP_XChineseEten          CodePage = 20002
	CP_Xcp20003              CodePage = 20003
	CP_Xcp20004              CodePage = 20004
	CP_Xcp20005              CodePage = 20005
	CP_XIA5                  CodePage = 20105
	CP_XIA5German            CodePage = 20106
	CP_XIA5Swedish           CodePage = 20107
	CP_XIA5Norwegian         CodePage = 20108
	CP_Usascii               CodePage = 20127
	CP_Xcp20261              CodePage = 20261
	CP_Xcp20269              CodePage = 20269
	CP_Ibm273                CodePage = 20273
	CP_Ibm277                CodePage = 20277
	CP_Ibm278                CodePage = 20278
	CP_Ibm280                CodePage = 20280
	CP_Ibm284                CodePage = 20284
	CP_Ibm285                CodePage = 20285
	CP_Ibm290                CodePage = 20290
	CP_Ibm297                CodePage = 20297
	CP_Ibm420                CodePage = 20420
	CP_Ibm423                CodePage = 20423
	CP_Ibm424                CodePage = 20424
	CP_XEBCDICKoreanExtended CodePage = 20833
	CP_IbmThai               CodePage = 20838
	CP_Koi8r                 CodePage = 20866
	CP_Ibm871                CodePage = 20871
	CP_Ibm880                CodePage = 20880
	CP_Ibm905                CodePage = 20905
	CP_Ibm00924              CodePage = 20924
	CP_EUCJP                 CodePage = 20932
	CP_Xcp20936              CodePage = 20936
	CP_Xcp20949              CodePage = 20949
	CP_Cp1025                CodePage = 21025
	CP_Koi8u                 CodePage = 21866
	CP_Iso88591              CodePage = 28591
	CP_Iso88592              CodePage = 28592
	CP_Iso88593              CodePage = 28593
	CP_Iso88594              CodePage = 28594
	CP_Iso88595              CodePage = 28595
	CP_Iso88596              CodePage = 28596
	CP_Iso88597              CodePage = 28597
	CP_Iso88598              CodePage = 28598
	CP_Iso88599              CodePage = 28599
	CP_Iso885910             CodePage = 28600
	CP_Iso885913             CodePage = 28603
	CP_Iso885915             CodePage = 28605
	CP_XEuropa               CodePage = 29001
	CP_Iso88598i             CodePage = 38598
	CP_Iso2022jp             CodePage = 50220
	CP_CsISO2022JP           CodePage = 50221
	CP_Iso2022jp_jis         CodePage = 50222
	CP_Iso2022kr             CodePage = 50225
	CP_Xcp50227              CodePage = 50227
	CP_Eucjp                 CodePage = 51932
	CP_EUCCN                 CodePage = 51936
	CP_Euckr                 CodePage = 51949
	CP_Hzgb2312              CodePage = 52936
	CP_Gb18030               CodePage = 54936
	CP_Xisciide              CodePage = 57002
	CP_Xisciibe              CodePage = 57003
	CP_Xisciita              CodePage = 57004
	CP_Xisciite              CodePage = 57005
	CP_Xisciias              CodePage = 57006
	CP_Xisciior              CodePage = 57007
	CP_Xisciika              CodePage = 57008
	CP_Xisciima              CodePage = 57009
	CP_Xisciigu              CodePage = 57010
	CP_Xisciipa              CodePage = 57011
	CP_Utf7                  CodePage = 65000
	CP_Utf8                  CodePage = 65001
)

var pageCodes []CodePage = []CodePage{
	CP_Unknown,
	CP_Usascii,
	CP_Iso88591,
	CP_Iso88592,
	CP_Iso88593,
	CP_Iso88594,
	CP_Iso88595,
	CP_Iso88596,
	CP_Iso88597,
	CP_Iso88598,
	CP_Iso88599,
	CP_Ibm437,
	CP_Ibm850,
	CP_Ibm852,
	CP_Ibm855,
	CP_Ibm857,
	CP_Ibm860,
	CP_Ibm861,
	CP_Ibm863,
	CP_Ibm864,
	CP_Ibm865,
	CP_Ibm869,
	CP_Shift_jis,
	CP_Macintosh,
	CP_big5,
	CP_Ksc5601,
	CP_Johab,
	CP_Cp866,
	CP_Windows1250,
	CP_Windows1251,
	CP_Windows1252,
	CP_Gb2312,
	CP_Windows1253,
	CP_Windows1254,
	CP_Windows1255,
	CP_Windows1256,
	CP_Windows1257,
	CP_Windows874,
	CP_Shift_jis,
	CP_Gb2312,
	CP_Ksc5601,
	CP_big5,
	CP_Johab,
	CP_Utf16,
	CP_Windows1258,
}

func GetCodePage(value int) CodePage {
	if value >= len(pageCodes) {
		return CP_Unknown
	}
	return pageCodes[value]
}
