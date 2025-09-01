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

    // func add(a i32, b i32) i32 { return a + b }
    addTy := llvm.FunctionType(i32, []llvm.Type{i32, i32}, false)
    addFn := llvm.AddFunction(mod, "add", addTy)
    entry := ctx.AddBasicBlock(addFn, "entry")
    builder.SetInsertPointAtEnd(entry)
    sum := builder.CreateAdd(addFn.Param(0), addFn.Param(1), "sum")
    builder.CreateRet(sum)

    // func main() i32 { return add(2, 3) }
    mainTy := llvm.FunctionType(i32, nil, false)
    mainFn := llvm.AddFunction(mod, "main", mainTy)
    mainEntry := ctx.AddBasicBlock(mainFn, "entry")
    builder.SetInsertPointAtEnd(mainEntry)
    c2 := llvm.ConstInt(i32, 2, false)
    c3 := llvm.ConstInt(i32, 3, false)
    // NOTE: pass the function type as the first arg:
    res := builder.CreateCall(addTy, addFn, []llvm.Value{c2, c3}, "res")
    builder.CreateRet(res)

    // Verify returns error, not bool.
    if err := llvm.VerifyModule(mod, llvm.PrintMessageAction); err != nil {
	panic(err)
    }

    fmt.Println(mod.String())
}
