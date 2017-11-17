// build

// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Test basic concurrency: the classic prime sieve.
// Do not run - loops forever.

package main

import (
        // Import the entire framework (including bundled verilog)
        _ "sdaccel"

        // Use the new AXI protocol package
        aximemory "axi/memory"
        axiprotocol "axi/protocol"
)


// Send the sequence 2, 3, 4, ... to channel 'ch'.
/*func Generate(ch chan<- int) {
	for i := 2; i < 10; i++ {
		ch <- i // Send 'i' to channel 'ch'.
	}
}*/


// Copy the values from array 'in' to channel 'ch',
// removing those divisible by 'prime'.
func Filter(in [10]uint32, ch chan<- [10]uint32, prime uint32){

	out := [10]uint32{0}
	for i := 0; i < 10; i++ {		

		val := in[i] // Receive value of new variable 'val' from 'in'.
		if val%prime != 0 {
			out[i] = val // Send 'val' to channel 'out'.
		}
	}
	ch <- out
}

// The prime sieve: Daisy-chain Filter processes together.
func Top(
	addrShared uintptr,

	// The first set of arguments will be the ports for interacting with host 
	// The second set of arguments will be the ports for interacting with memory
	memReadAddr chan<- axiprotocol.Addr,
	memReadData <-chan axiprotocol.ReadData,

	memWriteAddr chan<- axiprotocol.Addr,
	memWriteData chan<- axiprotocol.WriteData,
	memWriteResp <-chan axiprotocol.WriteResp){

	sharedMem := [5][10]uint32{
		 [10]uint32{2,3,4,5,6,7,8,9,10,11},
		 [10]uint32{0},
		 [10]uint32{0},
		 [10]uint32{0},
		 [10]uint32{0}}

//      ch := make(chan int) // Create a new channel.
//      go Generate(ch)      // Start Generate() as a subprocess.
//	sharedMem[0][0] = <-ch
        for i := 0; i < 4; i++{
                //  prime := <-ch
		prime := sharedMem[i][0]
                ch := make(chan [10]uint32)
                go Filter(sharedMem[i], ch, prime)
		//copy out chan to the next vector 
		sharedMem[i+1] = <-ch
	}
	// Write it back to the pointer the host requests
	aximemory.WriteUInt32(
		memWriteAddr, memWriteData, memWriteResp, false, addrShared, uint32(sharedMem[4][0]))

}
