package main

import (
	"encoding/binary"
	"fmt"
	"xcl"
	"os"
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

//	var fpath string
//	input = utils.load_data(fpath)

	//load validations 
//	test := bnn.ReadImage("dataset")
//	fmt.Println(test)

	//reshape image 
//	nw_image:= bnn.ReshapeImage(image)
//	fmt.Println(nw_image)


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
	var ret uint32
	err := binary.Read(buffShared.Reader(), binary.LittleEndian, &ret)
	if err != nil {
		fmt.Println("binary.Read failed:", err)
	}
	// Compute the expected result 
	expected := [4]uint32{2,3,5,7}

	// Exit with an error if the value is not correct
	if expected[0] != ret {
		// Print the value we got from the FPGA
		fmt.Printf("Expected %d, got %d\n", expected[3], ret)
		os.Exit(1)
	}

}
