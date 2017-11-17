package main

import (
	"encoding/binary"
	"fmt"
	"xcl"
//	"os"
        "testing"
)


func BenchmarkKernel(world xcl.World, 
		krnl *xcl.Kernel,
		B *testing.B,  
		buffShared *xcl.Memory) {

	// Set the pointer to the output buffer
	krnl.SetMemoryArg(0, buffShared)

	// Reset the timer so that we only measure runtime of the kernel
	B.ResetTimer()
	krnl.Run(1, 1, 1)
}

func main() {
	world := xcl.NewWorld()
	defer world.Release()

	krnl := world.Import("kernel_test").GetKernel("reconfigure_io_sdaccel_builder_stub_0_1")
	defer krnl.Release()

        // Allocate a buffer on the FPGA to store the return value of our computation
        // The shared mem locations between host-kernel is a 10-uint32 set for 10 prime nums, 
	// so we need 10 * 4 bytes to store it
        buffShared := world.Malloc(xcl.ReadOnly, 40)
        defer buffShared.Free()


	// Create a function that the benchmarking machinery can call
	f := func(B *testing.B) {
		BenchmarkKernel(world, krnl, B, buffShared)
	}
	// Benchmark it
	result := testing.Benchmark(f)

	// Print the result
	fmt.Printf("%s\n", result.String())

	// Decode that byte slice into the uint32 we're expecting
	var ret [10]uint32
	err := binary.Read(buffShared.Reader(), binary.LittleEndian, &ret)
	if err != nil {
		fmt.Println("binary.Read failed:", err)
	}
        fmt.Printf("Prime numbers %d\n", ret)

}
