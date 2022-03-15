(module
  (type $t0 (func))
  (type $t1 (func (param i32) (result i32)))
  (type $t2 (func))
  (import "wasi_unstable" "proc_exit" (func $wasi_unstable.proc_exit (type $t0)))
  (func $calc (type $t1) (param $p0 i32) (result i32)
    i32.const 1
    call $f4
    unreachable)
  (func $f2 (type $t0)
    call $wasi_unstable.proc_exit
    unreachable)
  (func $f3 (type $t2))
  (func $f4 (type $t0)
    call $f3
    call $f3
    call $f2
    unreachable)
  (table $T0 1 1 anyfunc)
  (memory $memory 2)
  (global $g0 (mut i32) (i32.const 66560))
  (export "memory" (memory 0))
  (export "calc" (func $calc)))
