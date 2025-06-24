package balena

import (
	"fmt"
	"testing"
)

func TestEnvironment(t *testing.T) {
    env := NewEnvironment()
    env.Set("x", &ObjectInteger{Value: 42})
    value, ok := env.Get("x")
    if !ok {
        t.Fatalf("Expected to find variable 'x' in environment")
    }
    if intObj, ok := value.(*ObjectInteger); ok {
        if intObj.Value != 42 {
            t.Errorf("Expected x to be 42, got %d", intObj.Value)
        } else {
            fmt.Println("x is 42 as expected")
        }
    } else {
        t.Errorf("Variable 'x' is not an integer")
    }
    env.Set("y", &ObjectFloat{Value: 3.14})
    value, ok = env.Get("y")
    if !ok {
        t.Fatalf("Expected to find variable 'y' in environment")
    }
    if floatObj, ok := value.(*ObjectFloat); ok {
        if floatObj.Value != 3.14 {
            t.Errorf("Expected y to be 3.14, got %f", floatObj.Value)
        } else {
            fmt.Println("y is 3.14 as expected")
        }
    } else {
        t.Errorf("Variable 'y' is not a float")
    }
    env.Set("z", &ObjectNull{})
    value, ok = env.Get("z")
    if !ok {
        t.Fatalf("Expected to find variable 'z' in environment")
    }
    if _, ok := value.(*ObjectNull); ok {
        fmt.Println("z is null as expected")
    } else {
        t.Errorf("Variable 'z' is not null")
    }
    if _, ok := env.Get("nonexistent"); ok {
        t.Errorf("Expected to not find variable 'nonexistent' in environment")
    } else {
        fmt.Println("nonexistent variable is not found as expected")
    }
}