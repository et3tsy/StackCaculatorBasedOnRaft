package validate

// to validate the format of the instruction
func Check(cmd string, args []int64) bool {
	switch cmd {
	case "create", "CREATE":
		{
			if len(args) != 0 {
				return false
			}
		}
	case "del", "DEL", "pop", "POP", "inc", "INC", "dec", "DEC",
		"add", "ADD", "sub", "SUB", "mul", "MUL", "div", "DIV":
		{
			if len(args) != 1 {
				return false
			}
		}
	case "push", "PUSH":
		{
			if len(args) != 2 {
				return false
			}
		}
	}
	return true
}
