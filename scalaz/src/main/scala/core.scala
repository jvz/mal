import java.util.concurrent.atomic.AtomicReference

import types._

import scala.io.Source

// TODO: a better DSL here would be nice (or macros to convert from scala to lisp forms)
object core {
  val ns: Map[MalSymbol, MalLambda] = Map(
    // TODO: make numbers less annoying so we can use ints, too
    MalSymbol('+) -> MalLambda {
      case MalReal(x) :: MalReal(y) :: Nil => MalReal(x + y)
    },
    MalSymbol('-) -> MalLambda {
      case MalReal(x) :: MalReal(y) :: Nil => MalReal(x - y)
    },
    MalSymbol('*) -> MalLambda {
      case MalReal(x) :: MalReal(y) :: Nil => MalReal(x * y)
    },
    MalSymbol('/) -> MalLambda {
      case MalReal(x) :: MalReal(y) :: Nil => MalReal(x / y)
    },
    MalSymbol('list) -> MalLambda {
      case elems => MalList(elems)
    },
    MalSymbol("list?") -> MalLambda {
      case MalList(_) :: Nil => MalTrue
      case _ => MalFalse
    },
    MalSymbol("empty?") -> MalLambda {
      case MalColl(coll) :: Nil if coll.isEmpty => MalTrue
      case _ => MalFalse
    },
    MalSymbol('count) -> MalLambda {
      case MalColl(coll) :: Nil => MalReal(coll.size)
      case MalNil :: Nil => MalReal(0)
    },
    MalSymbol('=) -> MalLambda {
      case a :: b :: Nil if a eql b => MalTrue
      case _ => MalFalse
    },
    MalSymbol('<) -> MalLambda {
      case MalReal(a) :: MalReal(b) :: Nil => if (a < b) MalTrue else MalFalse
    },
    MalSymbol('<=) -> MalLambda {
      case MalReal(a) :: MalReal(b) :: Nil => if (a <= b) MalTrue else MalFalse
    },
    MalSymbol('>) -> MalLambda {
      case MalReal(a) :: MalReal(b) :: Nil => if (a > b) MalTrue else MalFalse
    },
    MalSymbol('>=) -> MalLambda {
      case MalReal(a) :: MalReal(b) :: Nil => if (a >= b) MalTrue else MalFalse
    },
    MalSymbol("pr-str") -> MalLambda {
      case args => MalString(args.map(_.show()).mkString(" "))
    },
    MalSymbol('str) -> MalLambda {
      case args => MalString(args.map(_.show(false)).mkString)
    },
    MalSymbol('prn) -> MalLambda {
      case args =>
        println(args.map(_.show()).mkString(" "))
        MalNil
    },
    MalSymbol('println) -> MalLambda {
      case args =>
        println(args.map(_.show(false)).mkString(" "))
        MalNil
    },
    MalSymbol("read-string") -> MalLambda {
      case MalString(value) :: Nil =>
        reader.read_str(value)
    },
    MalSymbol('slurp) -> MalLambda {
      case MalString(value) :: Nil =>
        MalString(Source.fromFile(value).mkString)
    },
    MalSymbol("type-of") -> MalLambda {
      case arg :: Nil =>
        MalString(arg.getClass.getSimpleName)
    },
    MalSymbol('atom) -> MalLambda {
      case arg :: Nil =>
        MalRef(new AtomicReference(arg))
    },
    MalSymbol("atom?") -> MalLambda {
      case (_: MalRef) :: Nil => MalTrue
      case _ => MalFalse
    },
    MalSymbol('deref) -> MalLambda {
      case MalRef(value) :: Nil =>
        value.get()
    },
    MalSymbol("reset!") -> MalLambda {
      case MalRef(ref) :: value :: Nil =>
        ref.set(value)
        value
    },
    MalSymbol("swap!") -> MalLambda {
      case MalRef(ref) :: MalFn(f) :: args =>
        val result = f(ref.get() :: args)
        ref.set(result)
        result
    }
  )

  def syntax_error: Nothing = sys.error("Invalid syntax")
}
