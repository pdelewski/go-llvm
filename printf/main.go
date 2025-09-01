package main

import (
    "fmt"

    "tinygo.org/x/go-llvm"
)

func main() {
    ctx := llvm.NewContext()
    defer ctx.Dispose()

    mod := ctx.NewModule("demo")
    builder := ctx.NewBuilder()
    defer builder.Dispose()

    i32 := ctx.Int32Type()
    i8 := ctx.Int8Type()
    i8ptr := llvm.PointerType(i8, 0)

    // int printf(char*, ...);
    printfTy := llvm.FunctionType(i32, []llvm.Type{i8ptr}, true)
    printfFn := llvm.AddFunction(mod, "printf", printfTy)

    // main() { printf("%d\n", 42); return 0; }
    mainTy := llvm.FunctionType(i32, nil, false)
    mainFn := llvm.AddFunction(mod, "main", mainTy)
    entry := ctx.AddBasicBlock(mainFn, "entry")
    builder.SetInsertPointAtEnd(entry)

    // IMPORTANT: now that the builder has an insertion point, this is safe:
    fmtPtr := builder.CreateGlobalStringPtr("%d\n", ".fmt")

    val42 := llvm.ConstInt(i32, 42, false)
    builder.CreateCall(printfTy, printfFn, []llvm.Value{fmtPtr, val42}, "")
    builder.CreateRet(llvm.ConstInt(i32, 0, false))

    if err := llvm.VerifyModule(mod, llvm.PrintMessageAction); err != nil {
	panic(err)
    }

    fmt.Println(mod.String())
}
