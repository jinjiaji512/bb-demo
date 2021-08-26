/**
 * @Author: jinjiaji
 * @Description:
 * @File:  main_test.go
 * @Version: 1.0.0
 * @Date: 2021/8/26 下午4:27
 */

package main

import (
	"bb/engine/world"
	utils "bb/util"
	"fmt"
	"math/rand"
	"testing"
	"time"
)

func TestNewWorld(t *testing.T) {
	world.NewWorld()
}

func TestRand(t *testing.T) {
	for i := 0; i < 100; i++ {
		time.Sleep(time.Millisecond)
		rand.Seed(time.Now().UnixNano())
		fmt.Println(utils.GetFullName())
	}
}

func TestRand2(t *testing.T) {
	for i := 0; i < 100; i++ {
		rand.Seed(int64(i))
		time.Sleep(time.Second)
		fmt.Println(rand.Int())
	}
}
