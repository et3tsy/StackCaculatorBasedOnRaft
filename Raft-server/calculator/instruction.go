package calculator

import (
	"fmt"
)

// 创建新的栈计算器实例, 返回值为实例编号
func (c *Calculator) Create() (int64, error) {
	c.stackID++
	c.InstanceMap[c.stackID] = &Instance{
		data: []int64{},
	}
	return c.stackID, nil
}

// 删除编号为 instanceID 的实例
func (c *Calculator) Delete(instanceID int64) error {
	_, ok := c.InstanceMap[instanceID]
	if ok {
		delete(c.InstanceMap, instanceID)
		return nil
	}
	return fmt.Errorf("cannot find the instance")
}

// 将整数 x 压入编号为 instanceId 的栈
func (c *Calculator) Push(instanceID int64, x int64) error {
	instance, ok := c.InstanceMap[instanceID]
	if ok {
		instance.data = append(instance.data, x)
		c.InstanceMap[c.stackID] = &Instance{
			data: instance.data,
		}
		return nil
	}
	return fmt.Errorf("cannot find the instance")
}

// 将 instanceID 栈顶整数弹栈, 返回该整数值
func (c *Calculator) Pop(instanceID int64) (int64, error) {
	instance, ok := c.InstanceMap[instanceID]
	if ok {
		if len(instance.data) == 0 {
			return 0, fmt.Errorf("the instance is empty")
		}
		ret := instance.data[len(instance.data)-1]
		instance.data = instance.data[:len(instance.data)-1]
		c.InstanceMap[c.stackID] = instance
		return ret, nil
	}
	return 0, fmt.Errorf("cannot find the instance")
}

// 返回 instanceID 栈顶整数值
func (c *Calculator) Get(instanceID int64) (int64, error) {
	instance, ok := c.InstanceMap[instanceID]
	if ok {
		if len(instance.data) == 0 {
			return 0, fmt.Errorf("the instance is empty")
		}
		return instance.data[len(instance.data)-1], nil
	}
	return 0, fmt.Errorf("cannot find the instance")
}

// 将 instanceID 栈顶整数自增
func (c *Calculator) Inc(instanceID int64) error {
	instance, ok := c.InstanceMap[instanceID]
	if ok {
		if len(instance.data) == 0 {
			return fmt.Errorf("the instance is empty")
		}
		instance.data[len(instance.data)-1]++
		c.InstanceMap[c.stackID] = instance
		return nil
	}
	return fmt.Errorf("cannot find the instance")
}

// 将 instanceID 栈顶整数自减
func (c *Calculator) Dec(instanceID int64) error {
	instance, ok := c.InstanceMap[instanceID]
	if ok {
		if len(instance.data) == 0 {
			return fmt.Errorf("the instance is empty")
		}
		instance.data[len(instance.data)-1]--
		c.InstanceMap[c.stackID] = instance
		return nil
	}
	return fmt.Errorf("cannot find the instance")
}

// 将 instanceID 栈顶两个数分别弹出, 并完成加法, 并压入栈中
func (c *Calculator) Add(instanceID int64) error {
	instance, ok := c.InstanceMap[instanceID]
	if ok {
		if len(instance.data) < 2 {
			return fmt.Errorf("the size of data in the instance is less than 2")
		}
		x := instance.data[len(instance.data)-1]
		y := instance.data[len(instance.data)-2]
		instance.data[len(instance.data)-2] = x + y
		instance.data = instance.data[:len(instance.data)-1]
		c.InstanceMap[c.stackID] = instance
		return nil
	}
	return fmt.Errorf("cannot find the instance")
}

// 将 instanceID 栈顶两个数分别弹出, 并完成减法, 并压入栈中
func (c *Calculator) Sub(instanceID int64) error {
	instance, ok := c.InstanceMap[instanceID]
	if ok {
		if len(instance.data) < 2 {
			return fmt.Errorf("the size of data in the instance is less than 2")
		}
		x := instance.data[len(instance.data)-1]
		y := instance.data[len(instance.data)-2]
		instance.data[len(instance.data)-2] = x - y
		instance.data = instance.data[:len(instance.data)-1]
		c.InstanceMap[c.stackID] = instance
		return nil
	}
	return fmt.Errorf("cannot find the instance")
}

// 将 instanceID 栈顶两个数分别弹出, 并完成乘法, 并压入栈中
func (c *Calculator) Mul(instanceID int64) error {
	instance, ok := c.InstanceMap[instanceID]
	if ok {
		if len(instance.data) < 2 {
			return fmt.Errorf("the size of data in the instance is less than 2")
		}
		x := instance.data[len(instance.data)-1]
		y := instance.data[len(instance.data)-2]
		instance.data[len(instance.data)-2] = x * y
		instance.data = instance.data[:len(instance.data)-1]
		c.InstanceMap[c.stackID] = instance
		return nil
	}
	return fmt.Errorf("cannot find the instance")
}

// 将 instanceID 栈顶两个数分别弹出, 并完成除法, 并压入栈中
func (c *Calculator) Div(instanceID int64) error {
	instance, ok := c.InstanceMap[instanceID]
	if ok {
		if len(instance.data) < 2 {
			return fmt.Errorf("the size of data in the instance is less than 2")
		}
		x := instance.data[len(instance.data)-1]
		y := instance.data[len(instance.data)-2]
		instance.data[len(instance.data)-2] = x / y
		instance.data = instance.data[:len(instance.data)-1]
		c.InstanceMap[c.stackID] = instance
		return nil
	}
	return fmt.Errorf("cannot find the instance")
}
