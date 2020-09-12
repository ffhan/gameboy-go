package go_gb

type InterruptBit int

const (
	IE uint16 = 0xFFFF // Interrupt enable (R/W)
	IF uint16 = 0xFF0F // Interrupt flag (R/W)

	BitVBlank InterruptBit = 0 // Bit 0: V-Blank  Interrupt Enable  (INT 40h)  (1=Enable)
	BitLCD    InterruptBit = 1 // Bit 1: LCD STAT Interrupt Enable  (INT 48h)  (1=Enable)
	BitTimer  InterruptBit = 2 // Bit 2: Timer    Interrupt Enable  (INT 50h)  (1=Enable)
	BitSerial InterruptBit = 3 // Bit 3: Serial   Interrupt Enable  (INT 58h)  (1=Enable)
	BitJoypad InterruptBit = 4 // Bit 4: Joypad   Interrupt Enable  (INT 60h)  (1=Enable)
)

type Interrupt struct {
	Bit    InterruptBit
	JpAddr uint16
}

var (
	Interrupts []Interrupt = []Interrupt{
		{BitVBlank, 0x40},
		{BitLCD, 0x48},
		{BitTimer, 0x50},
		{BitSerial, 0x58},
		{BitJoypad, 0x60},
	}
)

// ieR - interrupt enable register, ifR - interrupt flag register
func ShouldServiceInterrupt(ieR, ifR byte, bit InterruptBit) bool {
	mask := byte(1 << bit)
	return ieR&mask != 0 && ifR&mask != 0
}
