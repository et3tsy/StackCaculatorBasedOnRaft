package calculator

import (
	"fmt"
	"server/models"
)

// To excute the instrucions.
func (c *Calculator) Excution(req models.Request) (result int64, err error) {
	switch len(req.Params) {
	case 0:
		{
			switch req.Instruction {
			case "create", "CREATE":
				{
					return c.Create()
				}
			default:
				{
					err = fmt.Errorf("arguments error")
					return
				}
			}
		}
	case 1:
		{
			switch req.Instruction {
			case "del", "DEL":
				{
					return 0, c.Delete(req.Params[0])
				}

			case "pop", "POP":
				{
					return c.Pop(req.Params[0])
				}

			case "inc", "INC":
				{
					return 0, c.Inc(req.Params[0])
				}
			case "dec", "DEC":
				{
					return 0, c.Dec(req.Params[0])
				}
			case "add", "ADD":
				{
					return 0, c.Add(req.Params[0])
				}
			case "sub", "SUB":
				{
					return 0, c.Sub(req.Params[0])
				}
			case "mul", "MUL":
				{
					return 0, c.Mul(req.Params[0])
				}
			case "div", "DIV":
				{
					return 0, c.Div(req.Params[0])
				}
			default:
				{
					err = fmt.Errorf("arguments error")
					return
				}
			}
		}
	case 2:
		{
			switch req.Instruction {
			case "push", "PUSH":
				{
					return 0, c.Push(req.Params[0], req.Params[1])
				}
			default:
				{
					err = fmt.Errorf("arguments error")
					return
				}
			}
		}
	default:
		{
			err = fmt.Errorf("arguments error")
			return
		}
	}
}
