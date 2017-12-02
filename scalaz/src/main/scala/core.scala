import types._

// TODO: a better DSL here would be nice (or macros to convert from scala to lisp forms)
object core {
  val ns: Map[MalSymbol, MalFunction] = Map(
    // TODO: make numbers less annoying so we can use ints, too
    MalSymbol('+) -> MalFunction {
      case MalReal(x) :: MalReal(y) :: Nil => MalReal(x + y)
    },
    MalSymbol('-) -> MalFunction {
      case MalReal(x) :: MalReal(y) :: Nil => MalReal(x - y)
    },
    MalSymbol('*) -> MalFunction {
      case MalReal(x) :: MalReal(y) :: Nil => MalReal(x * y)
    },
    MalSymbol('/) -> MalFunction {
      case MalReal(x) :: MalReal(y) :: Nil => MalReal(x / y)
    },
    MalSymbol('list) -> MalFunction {
      case elems => MalList(elems)
    },
    MalSymbol("list?") -> MalFunction {
      case MalList(_) :: Nil => MalTrue
      case _ => MalFalse
    },
    MalSymbol("empty?") -> MalFunction {
      case MalColl(coll) :: Nil if coll.isEmpty => MalTrue
      case _ => MalFalse
    },
    MalSymbol('count) -> MalFunction {
      case MalColl(coll) :: Nil => MalReal(coll.size)
      case MalNil :: Nil => MalReal(0)
    },
    MalSymbol('=) -> MalFunction {
      case a :: b :: Nil if a eql b => MalTrue
      case _ => MalFalse
    },
    MalSymbol('<) -> MalFunction {
      case MalReal(a) :: MalReal(b) :: Nil => if (a < b) MalTrue else MalFalse
    },
    MalSymbol('<=) -> MalFunction {
      case MalReal(a) :: MalReal(b) :: Nil => if (a <= b) MalTrue else MalFalse
    },
    MalSymbol('>) -> MalFunction {
      case MalReal(a) :: MalReal(b) :: Nil => if (a > b) MalTrue else MalFalse
    },
    MalSymbol('>=) -> MalFunction {
      case MalReal(a) :: MalReal(b) :: Nil => if (a >= b) MalTrue else MalFalse
    },
    MalSymbol("pr-str") -> MalFunction {
      case args => MalString(args.map(_.show()).mkString(" "))
    },
    MalSymbol('str) -> MalFunction {
      case args => MalString(args.map(_.show(false)).mkString)
    },
    MalSymbol('prn) -> MalFunction {
      case args =>
        println(args.map(_.show()).mkString(" "))
        MalNil
    },
    MalSymbol('println) -> MalFunction {
      case args =>
        println(args.map(_.show(false)).mkString(" "))
        MalNil
    }
  )

  def syntax_error: Nothing = sys.error("Invalid syntax")
}
