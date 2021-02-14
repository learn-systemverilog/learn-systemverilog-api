package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/learn-systemverilog/learn-systemverilog-api/transpiler"
)

func main() {
	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello, World!",
		})
	})

	logs := make(chan interface{})

	go func() {
		for msg := range logs {
			fmt.Println(msg)
		}
	}()

	err := transpiler.Transpile(`// DESCRIPTION: Verilator: Systemverilog example module
	// with interface to switch buttons, LEDs, LCD and register display
	
	parameter NINSTR_BITS = 32;
	parameter NBITS_TOP = 8, NREGS_TOP = 32, NBITS_LCD = 64;
	module top(input  logic clk_2,
			   input  logic [NBITS_TOP-1:0] SWI,
			   output logic [NBITS_TOP-1:0] LED,
			   output logic [NBITS_TOP-1:0] SEG,
			   output logic [NBITS_LCD-1:0] lcd_a, lcd_b,
			   output logic [NINSTR_BITS-1:0] lcd_instruction,
			   output logic [NBITS_TOP-1:0] lcd_registrador [0:NREGS_TOP-1],
			   output logic [NBITS_TOP-1:0] lcd_pc, lcd_SrcA, lcd_SrcB,
				 lcd_ALUResult, lcd_Result, lcd_WriteData, lcd_ReadData, 
			   output logic lcd_MemWrite, lcd_Branch, lcd_MemtoReg, lcd_RegWrite);
	
	  always_comb begin
		LED <= SWI | clk_2;
		SEG <= SWI;
		lcd_WriteData <= SWI;
		lcd_pc <= 'h12;
		lcd_instruction <= 'h34567890;
		lcd_SrcA <= 'hab;
		lcd_SrcB <= 'hcd;
		lcd_ALUResult <= 'hef;
		lcd_Result <= 'h11;
		lcd_ReadData <= 'h33;
		lcd_MemWrite <= SWI[0];
		lcd_Branch <= SWI[1];
		lcd_MemtoReg <= SWI[2];
		lcd_RegWrite <= SWI[3];
		for(int i=0; i<NREGS_TOP; i++) lcd_registrador[i] <= i+i*16;
		lcd_a <= {56'h1234567890ABCD, SWI};
		lcd_b <= {SWI, 56'hFEDCBA09876543};
	  end
	
	endmodule`, logs)
	if err != nil {
		log.Fatal("Transpile: ", err)
	}

	if err := r.Run(); err != nil {
		log.Fatal("Listen and serve: ", err)
	}
}
