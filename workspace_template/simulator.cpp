#include <iostream>
#include <sstream>

#include <emscripten/bind.h>

#include "Vtop.h"

using namespace std;

Vtop *top;

string top_to_string()
{
    stringstream buffer;

    buffer << "{"
           << "\"led\": " << +top->LED
           << ", \"seg\": " << +top->SEG
           << ", \"lcd_a\": " << +top->lcd_a
           << ", \"lcd_b\": " << +top->lcd_b
           << ", \"lcd_instruction\": " << +top->lcd_instruction;

    buffer << ", \"lcd_registrador\": [" << +top->lcd_registrador[0];
    for (int i = 1; i < 32; i++)
    {
        buffer << ", " << +top->lcd_registrador[i];
    }
    buffer << "]";

    buffer << ", \"lcd_pc\": " << +top->lcd_pc
           << ", \"lcd_src_a\": " << +top->lcd_SrcA
           << ", \"lcd_src_b\": " << +top->lcd_SrcB
           << ", \"lcd_alu_result\": " << +top->lcd_ALUResult
           << ", \"lcd_result\": " << +top->lcd_Result
           << ", \"lcd_write_data\": " << +top->lcd_WriteData
           << ", \"lcd_read_data\": " << +top->lcd_ReadData
           << ", \"lcd_mem_write\": " << +top->lcd_MemWrite
           << ", \"lcd_branch\": " << +top->lcd_Branch
           << ", \"lcd_memto_reg\": " << +top->lcd_MemtoReg
           << ", \"lcd_reg_write\": " << +top->lcd_RegWrite;

    buffer << "}";

    return buffer.str();
}

string simulate(int swi, bool clk)
{
    if (top == nullptr)
        top = new Vtop();

    top->SWI = swi;
    top->clk_2 = clk;

    // Run one simulation step.
    top->eval();

    return top_to_string();
}

void finalize()
{
    top->final();
}

EMSCRIPTEN_BINDINGS(my_module)
{
    emscripten::function("simulate", &simulate);
    emscripten::function("finalize", &finalize);
}
