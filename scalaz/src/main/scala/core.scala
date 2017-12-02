import types._

// TODO: a better DSL here would be nice (or macros to convert from scala to lisp forms)
object core {
  val ns: Map[MalSymbol, MalFn] = Map(
    // TODO: make numbers less annoying so we can use ints, too
    MalSymbol('+) -> MalFn {
      case MalReal(x) :: MalReal(y) :: Nil => MalReal(x + y)
    },
    MalSymbol('-) -> MalFn {
      case MalReal(x) :: MalReal(y) :: Nil => MalReal(x - y)
    },
    MalSymbol('*) -> MalFn {
      case MalReal(x) :: MalReal(y) :: Nil => MalReal(x * y)
    },
    MalSymbol('/) -> MalFn {
      case MalReal(x) :: MalReal(y) :: Nil => MalReal(x / y)
    },
    MalSymbol('list) -> MalFn {
      case elems => MalList(elems)
    },
    MalSymbol("list?") -> MalFn {
      case MalList(_) :: Nil => MalTrue
      case _ => MalFalse
    },
    MalSymbol("empty?") -> MalFn {
      case MalColl(coll) :: Nil if coll.isEmpty => MalTrue
      case _ => MalFalse
    },
    MalSymbol('count) -> MalFn {
      case MalColl(coll) :: Nil => MalReal(coll.size)
      case MalNil :: Nil => MalReal(0)
    },
    MalSymbol('=) -> MalFn {
      case a :: b :: Nil if a eql b => MalTrue
      case _ => MalFalse
    },
    MalSymbol('<) -> MalFn {
      case MalReal(a) :: MalReal(b) :: Nil => if (a < b) MalTrue else MalFalse
    },
    MalSymbol('<=) -> MalFn {
      case MalReal(a) :: MalReal(b) :: Nil => if (a <= b) MalTrue else MalFalse
    },
    MalSymbol('>) -> MalFn {
      case MalReal(a) :: MalReal(b) :: Nil => if (a > b) MalTrue else MalFalse
    },
    MalSymbol('>=) -> MalFn {
      case MalReal(a) :: MalReal(b) :: Nil => if (a >= b) MalTrue else MalFalse
    },
    MalSymbol("pr-str") -> MalFn {
      case args => MalString(args.map(_.show()).mkString(" "))
    },
    MalSymbol('str) -> MalFn {
      case args => MalString(args.map(_.show(false)).mkString)
    },
    MalSymbol('prn) -> MalFn {
      case args =>
        println(args.map(_.show()).mkString(" "))
        MalNil
    },
    MalSymbol('println) -> MalFn {
      case args =>
        println(args.map(_.show(false)).mkString(" "))
        MalNil
    }
  )

  def syntax_error: Nothing = sys.error("Invalid syntax")
}
