package neptune

// This part is a fix for gohook library
// Somekeys in gohook are not decisive, outputing untranslated keys

import "runtime"

var (
	WinRaw2Code = map[uint16]uint16{
		160: 16, // Shift
		161: 16, // Shift Right
		162: 17, // Ctrl
		163: 17, // Ctrl Right
		164: 18, // Alt
		165: 18, // Alt Right
	}
	LinuxRaw2Code = map[uint16]uint16{
		96:    192, // Backquote
		126:   192, // Backquote / grave accent
		41:    48,  // Digit 0
		33:    49,  // Digit 1
		64:    50,  // Digit 2
		35:    51,  // Digit 3
		36:    52,  // Digit 4
		37:    53,  // Digit 5
		94:    54,  // Digit 6
		38:    55,  // Digit 7
		42:    57,  // Digit 8
		40:    57,  // Digit 9
		45:    189, // Dash / Minus
		95:    189, // Dash // Minus
		43:    61,  // Equal
		65289: 9,   // Tab
		65056: 9,   // Tab
		97:    65,  // Key A
		98:    66,  // Key B
		99:    67,  // Key C
		100:   68,  // Key D
		101:   69,  // Key E
		102:   70,  // Key F
		103:   71,  // Key G
		104:   72,  // Key H
		105:   73,  // Key I
		106:   74,  // Key J
		107:   75,  // Key K
		108:   76,  // Key L
		109:   77,  // Key M
		110:   78,  // Key N
		111:   79,  // Key O
		112:   80,  // Key P
		113:   81,  // Key Q
		114:   82,  // Key R
		115:   83,  // Key S
		116:   84,  // Key T
		117:   85,  // Key U
		118:   86,  // Key V
		119:   87,  // Key W
		120:   88,  // Key X
		121:   89,  // Key Y
		122:   90,  // Key Z
		91:    219, // Left bracket
		93:    221, // Right bracket
		123:   219, // Left bracket
		125:   221, // Right bracket
		39:    222, // Quote / Double quote
		34:    222, // Quote / Double quote
		92:    220, // Backslash / Pipe
		124:   220, // Backslash / Pipe
		44:    188, // Comma
		60:    188, // Comma / Arrow left
		46:    190, // Period
		62:    190, // Period / arrow right
		47:    191, // Forward slash
		63:    191, // Forward Slash
		65288: 8,   // backspace
		65509: 20,  // Caps Lock
		65505: 16,  // Shift
		65506: 16,  // Shift Right
		65507: 17,  // Ctrl
		65508: 17,  // Ctrl Left
		65515: 91,  // Meta / Win Key (Left)
		65513: 18,  // Alt left
		65514: 18,  // Alt right
		65293: 13,  // Enter
		65421: 13,  // Enter (Numpad)
		65307: 27,  // Escape
		65470: 112, // F1
		65471: 113, // F2
		65472: 114, // F3
		65473: 115, // F4
		65474: 116, // F5
		65475: 117, // F6
		65476: 118, // F7
		65477: 119, // F8
		65478: 120, // F9
		65479: 121, // F10
		65480: 122, // F11
		65481: 123, // F12
		65409: 112, // Fn + F1
		65297: 113, // Fn + F2
		65299: 114, // Fn + F3
		65298: 115, // Fn + F4
		65301: 116, // Fn + F5
		65302: 117, // Fn + F6
		65300: 118, // Fn + F7
		65303: 119, // Fn + F8
		65305: 120, // Fn + F9
		65304: 121, // Fn + F10
		65309: 123, // Fn + F12
		65377: 44,  // Print screen
		// 65300: 135, // Scroll Lock // Duplicate with Fn + F7
		// 65299: 19, // Pause // Duplicate with Fn + F6
		65379: 45,  // Insert
		65535: 46,  // Delete
		65360: 36,  // Home
		65367: 35,  // End
		65365: 33,  // Page Up
		65366: 34,  // Page down
		65362: 38,  // Arrow up
		65364: 40,  // Arrow Down
		65361: 37,  // Arrow left
		65363: 39,  // Arrow Right
		65407: 144, // Num Lock
		65455: 111, // Devide (Numpad)
		65450: 106, // Multiply (Numpad)
		65453: 109, // Substract (Numpad)
		65451: 107, // Add (Numpad)
		65454: 110, // Decimal point / period (Numpad)
		65456: 96,  // Numpad 0
		65457: 97,  // Numpad 1
		65458: 98,  // Numpad 2
		65459: 99,  // Numpad 3
		65460: 100, // Numpad 4
		65461: 101, // Numpad 5
		65462: 102, // Numpad 6
		65463: 103, // Numpad 7
		65464: 104, // Numpad 8
		65465: 105, // Numpad 9
		65429: 36,  // Shift + Numpad 7 = (Home)
		65430: 37,  // Shift + Numpad 4 = (Arrow left)
		65431: 38,  // Shift + Numpad 8 = (Arrow up)
		65432: 39,  // Shift + Numpad 6 = (Arrow right)
		65433: 40,  // Shift + Numpad 2 = (Arrow down)
		65434: 33,  // Shift + Numpad 9 = (Page Up)
		65435: 34,  // Shift + Numpad 3 = (Page Down)
		65436: 35,  // Shift + Numpad 4 = (End)
		65437: 12,  // Shift + Numpad 5 = (clear)
		65438: 45,  // Shift + Numpad 0 = (Insert)
	}
	LEcode2Char = map[uint16]string{
		65:  "a",
		66:  "b",
		67:  "c",
		68:  "d",
		69:  "e",
		70:  "f",
		71:  "g",
		72:  "h",
		73:  "i",
		74:  "j",
		75:  "k",
		76:  "l",
		77:  "m",
		78:  "n",
		79:  "o",
		80:  "p",
		81:  "q",
		82:  "r",
		83:  "s",
		84:  "t",
		85:  "u",
		86:  "v",
		87:  "w",
		88:  "x",
		89:  "y",
		90:  "z",
		48:  "0",
		49:  "1",
		50:  "2",
		51:  "3",
		52:  "4",
		53:  "5",
		54:  "6",
		55:  "7",
		56:  "8",
		57:  "9",
		189: "dash",
		61:  "equal",
		187: "equal",
		192: "backquote",
		96:  "numpad_0",
		97:  "numpad_1",
		98:  "numpad_2",
		99:  "numpad_3",
		100: "numpad_4",
		101: "numpad_5",
		102: "numpad_6",
		103: "numpad_7",
		104: "numpad_8",
		105: "numpad_9",
		106: "multiply",
		107: "add",
		108: "numpad_period",
		109: "subtract",
		110: "decimal_point",
		111: "divide",
		144: "num_lock",
		8:   "backspace",
		9:   "tab",
		12:  "clear",
		13:  "enter",
		16:  "shift",
		17:  "ctrl",
		18:  "alt",
		91:  "lsuper",
		92:  "rsuper",
		19:  "pause",
		20:  "caps_lock",
		27:  "escape",
		32:  "spacebar",
		33:  "page_up",
		34:  "page_down",
		35:  "end",
		36:  "home",
		37:  "left_arrow",
		38:  "up_arrow",
		39:  "right_arrow",
		40:  "down_arrow",
		44:  "print",
		45:  "insert",
		46:  "delete",
		145: "scroll_lock",
		112: "f1",
		113: "f2",
		114: "f3",
		115: "f4",
		116: "f5",
		117: "f6",
		118: "f7",
		119: "f8",
		120: "f9",
		121: "f10",
		122: "f11",
		123: "f12",
		188: "comma",
		190: "period",
		191: "forward_slash",
		186: "semi-colon",
		59:  "semi-colon",
		222: "quote",
		220: "backslash",
		219: "lbracket",
		221: "rbracket",
	}
)

func CodeToChar(code uint16) string {
	return LEcode2Char[code]
}

func GetKeyCode(code uint16) uint16 {
	switch runtime.GOOS {
	case "linux":
		if keycode := LinuxRaw2Code[code]; keycode != 0 {
			return keycode
		}
	case "windows":
		if keycode := WinRaw2Code[code]; keycode != 0 {
			return keycode
		}
	}
	return code
}
