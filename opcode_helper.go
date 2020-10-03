package go_gb

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type operand struct {
	Name      string `json:"name"`
	Immediate bool   `json:"immediate"`
}

func (o operand) String() string {
	if o.Immediate {
		return o.Name
	}
	return fmt.Sprintf("(%s)", o.Name)
}

type opInfo struct {
	Name         string    `json:"mnemonic"`
	Operands     []operand `json:"operands"`
	cachedString string
}

func (o *opInfo) String() string {
	if o.cachedString == "" {
		var sb strings.Builder
		sb.WriteString(o.Name)
		for i, operand := range o.Operands {
			sb.WriteRune(' ')
			sb.WriteString(operand.String())
			if i != len(o.Operands)-1 {
				sb.WriteString(",")
			}
		}
		o.cachedString = sb.String()
	}
	return o.cachedString
}

var Unprefixed = map[byte]*opInfo{}
var Prefixed = map[byte]*opInfo{}

func InitInstructions() {
	ops, err := os.Open("opcodes.json")
	if err != nil {
		panic(err)
	}
	defer ops.Close()

	var document struct {
		Unprefixed map[string]opInfo `json:"unprefixed"`
		Prefixed   map[string]opInfo `json:"cbprefixed"`
	}
	if err = json.NewDecoder(ops).Decode(&document); err != nil {
		panic(err)
	}

	for opcode, op := range document.Unprefixed {
		op := op
		id, err := strconv.ParseUint(opcode[2:], 16, 64)
		if err != nil {
			panic(err)
		}
		Unprefixed[byte(id)] = &op
	}
	for opcode, op := range document.Prefixed {
		op := op
		id, err := strconv.ParseUint(opcode[2:], 16, 64)
		if err != nil {
			panic(err)
		}
		Prefixed[byte(id)] = &op
	}
}
