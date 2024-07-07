package glfw3

var (
	GlfwVersionMajor              = 3
	GlfwVersionMinor              = 5
	GlfwVersionRevision           = 0
	GlfwTrue                      = 1
	GlfwFalse                     = 0
	GlfwRelease                   = 0
	GlfwPress                     = 1
	GlfwRepeat                    = 2
	GlfwHatCentered               = 0
	GlfwHatUp                     = 1
	GlfwHatRight                  = 2
	GlfwHatDown                   = 4
	GlfwHatLeft                   = 8
	GlfwHatRightUp                = GlfwHatRight | GlfwHatUp
	GlfwHatRightDown              = GlfwHatRight | GlfwHatDown
	GlfwHatLeftUp                 = GlfwHatLeft | GlfwHatUp
	GlfwHatLeftDown               = GlfwHatLeft | GlfwHatDown
	GlfwKeyUnknown                = -1
	GlfwKeySpace                  = 32
	GlfwKeyApostrophe             = 39 /* ' */
	GlfwKeyComma                  = 44 /* , */
	GlfwKeyMinus                  = 45 /* - */
	GlfwKeyPeriod                 = 46 /* . */
	GlfwKeySlash                  = 47 /* / */
	GlfwKey0                      = 48
	GlfwKey1                      = 49
	GlfwKey2                      = 50
	GlfwKey3                      = 51
	GlfwKey4                      = 52
	GlfwKey5                      = 53
	GlfwKey6                      = 54
	GlfwKey7                      = 55
	GlfwKey8                      = 56
	GlfwKey9                      = 57
	GlfwKeySemicolon              = 59 /* ; */
	GlfwKeyEqual                  = 61 /* = */
	GlfwKeyA                      = 65
	GlfwKeyB                      = 66
	GlfwKeyC                      = 67
	GlfwKeyD                      = 68
	GlfwKeyE                      = 69
	GlfwKeyF                      = 70
	GlfwKeyG                      = 71
	GlfwKeyH                      = 72
	GlfwKeyI                      = 73
	GlfwKeyJ                      = 74
	GlfwKeyK                      = 75
	GlfwKeyL                      = 76
	GlfwKeyM                      = 77
	GlfwKeyN                      = 78
	GlfwKeyO                      = 79
	GlfwKeyP                      = 80
	GlfwKeyQ                      = 81
	GlfwKeyR                      = 82
	GlfwKeyS                      = 83
	GlfwKeyT                      = 84
	GlfwKeyU                      = 85
	GlfwKeyV                      = 86
	GlfwKeyW                      = 87
	GlfwKeyX                      = 88
	GlfwKeyY                      = 89
	GlfwKeyZ                      = 90
	GlfwKeyLeftBracket            = 91  /* [ */
	GlfwKeyBackslash              = 92  /*  */
	GlfwKeyRightBracket           = 93  /* ] */
	GlfwKeyGraveAccent            = 96  /* ` */
	GlfwKeyWorld1                 = 161 /* non-US #1 */
	GlfwKeyWorld2                 = 162 /* non-US #2 */
	GlfwKeyEscape                 = 256
	GlfwKeyEnter                  = 257
	GlfwKeyTab                    = 258
	GlfwKeyBackspace              = 259
	GlfwKeyInsert                 = 260
	GlfwKeyDelete                 = 261
	GlfwKeyRight                  = 262
	GlfwKeyLeft                   = 263
	GlfwKeyDown                   = 264
	GlfwKeyUp                     = 265
	GlfwKeyPageUp                 = 266
	GlfwKeyPageDown               = 267
	GlfwKeyHome                   = 268
	GlfwKeyEnd                    = 269
	GlfwKeyCapsLock               = 280
	GlfwKeyScrollLock             = 281
	GlfwKeyNumLock                = 282
	GlfwKeyPrintScreen            = 283
	GlfwKeyPause                  = 284
	GlfwKeyF1                     = 290
	GlfwKeyF2                     = 291
	GlfwKeyF3                     = 292
	GlfwKeyF4                     = 293
	GlfwKeyF5                     = 294
	GlfwKeyF6                     = 295
	GlfwKeyF7                     = 296
	GlfwKeyF8                     = 297
	GlfwKeyF9                     = 298
	GlfwKeyF10                    = 299
	GlfwKeyF11                    = 300
	GlfwKeyF12                    = 301
	GlfwKeyF13                    = 302
	GlfwKeyF14                    = 303
	GlfwKeyF15                    = 304
	GlfwKeyF16                    = 305
	GlfwKeyF17                    = 306
	GlfwKeyF18                    = 307
	GlfwKeyF19                    = 308
	GlfwKeyF20                    = 309
	GlfwKeyF21                    = 310
	GlfwKeyF22                    = 311
	GlfwKeyF23                    = 312
	GlfwKeyF24                    = 313
	GlfwKeyF25                    = 314
	GlfwKeyKp0                    = 320
	GlfwKeyKp1                    = 321
	GlfwKeyKp2                    = 322
	GlfwKeyKp3                    = 323
	GlfwKeyKp4                    = 324
	GlfwKeyKp5                    = 325
	GlfwKeyKp6                    = 326
	GlfwKeyKp7                    = 327
	GlfwKeyKp8                    = 328
	GlfwKeyKp9                    = 329
	GlfwKeyKpDecimal              = 330
	GlfwKeyKpDivide               = 331
	GlfwKeyKpMultiply             = 332
	GlfwKeyKpSubtract             = 333
	GlfwKeyKpAdd                  = 334
	GlfwKeyKpEnter                = 335
	GlfwKeyKpEqual                = 336
	GlfwKeyLeftShift              = 340
	GlfwKeyLeftControl            = 341
	GlfwKeyLeftAlt                = 342
	GlfwKeyLeftSuper              = 343
	GlfwKeyRightShift             = 344
	GlfwKeyRightControl           = 345
	GlfwKeyRightAlt               = 346
	GlfwKeyRightSuper             = 347
	GlfwKeyMenu                   = 348
	GlfwKeyLast                   = GlfwKeyMenu
	GlfwModShift                  = 0x0001
	GlfwModControl                = 0x0002
	GlfwModAlt                    = 0x0004
	GlfwModSuper                  = 0x0008
	GlfwModCapsLock               = 0x0010
	GlfwModNumLock                = 0x0020
	GlfwMouseButton1              = 0
	GlfwMouseButton2              = 1
	GlfwMouseButton3              = 2
	GlfwMouseButton4              = 3
	GlfwMouseButton5              = 4
	GlfwMouseButton6              = 5
	GlfwMouseButton7              = 6
	GlfwMouseButton8              = 7
	GlfwMouseButtonLast           = GlfwMouseButton8
	GlfwMouseButtonLeft           = GlfwMouseButton1
	GlfwMouseButtonRight          = GlfwMouseButton2
	GlfwMouseButtonMiddle         = GlfwMouseButton3
	GlfwJoystick1                 = 0
	GlfwJoystick2                 = 1
	GlfwJoystick3                 = 2
	GlfwJoystick4                 = 3
	GlfwJoystick5                 = 4
	GlfwJoystick6                 = 5
	GlfwJoystick7                 = 6
	GlfwJoystick8                 = 7
	GlfwJoystick9                 = 8
	GlfwJoystick10                = 9
	GlfwJoystick11                = 10
	GlfwJoystick12                = 11
	GlfwJoystick13                = 12
	GlfwJoystick14                = 13
	GlfwJoystick15                = 14
	GlfwJoystick16                = 15
	GlfwJoystickLast              = GlfwJoystick16
	GlfwGamepadButtonA            = 0
	GlfwGamepadButtonB            = 1
	GlfwGamepadButtonX            = 2
	GlfwGamepadButtonY            = 3
	GlfwGamepadButtonLeftBumper   = 4
	GlfwGamepadButtonRightBumper  = 5
	GlfwGamepadButtonBack         = 6
	GlfwGamepadButtonStart        = 7
	GlfwGamepadButtonGuide        = 8
	GlfwGamepadButtonLeftThumb    = 9
	GlfwGamepadButtonRightThumb   = 10
	GlfwGamepadButtonDpadUp       = 11
	GlfwGamepadButtonDpadRight    = 12
	GlfwGamepadButtonDpadDown     = 13
	GlfwGamepadButtonDpadLeft     = 14
	GlfwGamepadButtonLast         = GlfwGamepadButtonDpadLeft
	GlfwGamepadButtonCross        = GlfwGamepadButtonA
	GlfwGamepadButtonCircle       = GlfwGamepadButtonB
	GlfwGamepadButtonSquare       = GlfwGamepadButtonX
	GlfwGamepadButtonTriangle     = GlfwGamepadButtonY
	GlfwGamepadAxisLeftX          = 0
	GlfwGamepadAxisLeftY          = 1
	GlfwGamepadAxisRightX         = 2
	GlfwGamepadAxisRightY         = 3
	GlfwGamepadAxisLeftTrigger    = 4
	GlfwGamepadAxisRightTrigger   = 5
	GlfwGamepadAxisLast           = GlfwGamepadAxisRightTrigger
	GlfwNoError                   = 0
	GlfwNotInitialized            = 0x00010001
	GlfwNoCurrentContext          = 0x00010002
	GlfwInvalidEnum               = 0x00010003
	GlfwInvalidValue              = 0x00010004
	GlfwOutOfMemory               = 0x00010005
	GlfwApiUnavailable            = 0x00010006
	GlfwVersionUnavailable        = 0x00010007
	GlfwPlatformError             = 0x00010008
	GlfwFormatUnavailable         = 0x00010009
	GlfwNoWindowContext           = 0x0001000A
	GlfwCursorUnavailable         = 0x0001000B
	GlfwFeatureUnavailable        = 0x0001000C
	GlfwFeatureUnimplemented      = 0x0001000D
	GlfwPlatformUnavailable       = 0x0001000E
	GlfwFocused                   = 0x00020001
	GlfwIconified                 = 0x00020002
	GlfwResizable                 = 0x00020003
	GlfwVisible                   = 0x00020004
	GlfwDecorated                 = 0x00020005
	GlfwAutoIconify               = 0x00020006
	GlfwFloating                  = 0x00020007
	GlfwMaximized                 = 0x00020008
	GlfwCenterCursor              = 0x00020009
	GlfwTransparentFramebuffer    = 0x0002000A
	GlfwHovered                   = 0x0002000B
	GlfwFocusOnShow               = 0x0002000C
	GlfwMousePassthrough          = 0x0002000D
	GlfwPositionX                 = 0x0002000E
	GlfwPositionY                 = 0x0002000F
	GlfwRedBits                   = 0x00021001
	GlfwGreenBits                 = 0x00021002
	GlfwBlueBits                  = 0x00021003
	GlfwAlphaBits                 = 0x00021004
	GlfwDepthBits                 = 0x00021005
	GlfwStencilBits               = 0x00021006
	GlfwAccumRedBits              = 0x00021007
	GlfwAccumGreenBits            = 0x00021008
	GlfwAccumBlueBits             = 0x00021009
	GlfwAccumAlphaBits            = 0x0002100A
	GlfwAuxBuffers                = 0x0002100B
	GlfwStereo                    = 0x0002100C
	GlfwSamples                   = 0x0002100D
	GlfwSrgbCapable               = 0x0002100E
	GlfwRefreshRate               = 0x0002100F
	GlfwDoublebuffer              = 0x00021010
	GlfwClientApi                 = 0x00022001
	GlfwContextVersionMajor       = 0x00022002
	GlfwContextVersionMinor       = 0x00022003
	GlfwContextRevision           = 0x00022004
	GlfwContextRobustness         = 0x00022005
	GlfwOpenglForwardCompat       = 0x00022006
	GlfwContextDebug              = 0x00022007
	GlfwOpenglDebugContext        = GlfwContextDebug
	GlfwOpenglProfile             = 0x00022008
	GlfwContextReleaseBehavior    = 0x00022009
	GlfwContextNoError            = 0x0002200A
	GlfwContextCreationApi        = 0x0002200B
	GlfwScaleToMonitor            = 0x0002200C
	GlfwScaleFramebuffer          = 0x0002200D
	GlfwCocoaRetinaFramebuffer    = 0x00023001
	GlfwCocoaFrameName            = 0x00023002
	GlfwCocoaGraphicsSwitching    = 0x00023003
	GlfwX11ClassName              = 0x00024001
	GlfwX11InstanceName           = 0x00024002
	GlfwWin32KeyboardMenu         = 0x00025001
	GlfwWin32Showdefault          = 0x00025002
	GlfwWaylandAppId              = 0x00026001
	GlfwNoApi                     = 0
	GlfwOpenglApi                 = 0x00030001
	GlfwOpenglEsApi               = 0x00030002
	GlfwNoRobustness              = 0
	GlfwNoResetNotification       = 0x00031001
	GlfwLoseContextOnReset        = 0x00031002
	GlfwOpenglAnyProfile          = 0
	GlfwOpenglCoreProfile         = 0x00032001
	GlfwOpenglCompatProfile       = 0x00032002
	GlfwCursor                    = 0x00033001
	GlfwStickyKeys                = 0x00033002
	GlfwStickyMouseButtons        = 0x00033003
	GlfwLockKeyMods               = 0x00033004
	GlfwRawMouseMotion            = 0x00033005
	GlfwUnlimitedMouseButtons     = 0x00033006
	GlfwCursorNormal              = 0x00034001
	GlfwCursorHidden              = 0x00034002
	GlfwCursorDisabled            = 0x00034003
	GlfwCursorCaptured            = 0x00034004
	GlfwAnyReleaseBehavior        = 0
	GlfwReleaseBehaviorFlush      = 0x00035001
	GlfwReleaseBehaviorNone       = 0x00035002
	GlfwNativeContextApi          = 0x00036001
	GlfwEglContextApi             = 0x00036002
	GlfwOsmesaContextApi          = 0x00036003
	GlfwAnglePlatformTypeNone     = 0x00037001
	GlfwAnglePlatformTypeOpengl   = 0x00037002
	GlfwAnglePlatformTypeOpengles = 0x00037003
	GlfwAnglePlatformTypeD3D9     = 0x00037004
	GlfwAnglePlatformTypeD3D11    = 0x00037005
	GlfwAnglePlatformTypeVulkan   = 0x00037007
	GlfwAnglePlatformTypeMetal    = 0x00037008
	GlfwWaylandPreferLibdecor     = 0x00038001
	GlfwWaylandDisableLibdecor    = 0x00038002
	GlfwAnyPosition               = 0x80000000
	GlfwArrowCursor               = 0x00036001
	GlfwIbeamCursor               = 0x00036002
	GlfwCrosshairCursor           = 0x00036003
	GlfwPointingHandCursor        = 0x00036004
	GlfwResizeEwCursor            = 0x00036005
	GlfwResizeNsCursor            = 0x00036006
	GlfwResizeNwseCursor          = 0x00036007
	GlfwResizeNeswCursor          = 0x00036008
	GlfwResizeAllCursor           = 0x00036009
	GlfwNotAllowedCursor          = 0x0003600A
	GlfwHresizeCursor             = GlfwResizeEwCursor
	GlfwVresizeCursor             = GlfwResizeNsCursor
	GlfwHandCursor                = GlfwPointingHandCursor
	GlfwConnected                 = 0x00040001
	GlfwDisconnected              = 0x00040002
	GlfwJoystickHatButtons        = 0x00050001
	GlfwAnglePlatformType         = 0x00050002
	GlfwPlatform                  = 0x00050003
	GlfwCocoaChdirResources       = 0x00051001
	GlfwCocoaMenubar              = 0x00051002
	GlfwX11XcbVulkanSurface       = 0x00052001
	GlfwWaylandLibdecor           = 0x00053001
	GlfwAnyPlatform               = 0x00060000
	GlfwPlatformWin32             = 0x00060001
	GlfwPlatformCocoa             = 0x00060002
	GlfwPlatformWayland           = 0x00060003
	GlfwPlatformX11               = 0x00060004
	GlfwPlatformNull              = 0x00060005
	GlfwDontCare                  = -1
)