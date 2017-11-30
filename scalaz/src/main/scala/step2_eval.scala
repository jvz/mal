import types._

import scala.annotation.tailrec
import scala.io.StdIn
import scala.util.control.NonFatal

object step2_eval {

  type Env = Map[MalSymbol, MalType]

  def read(str: String): MalType = reader.read_str(str)

  def eval(ast: MalType, env: Env): MalType = ast match {
    case MalNil => MalNil
    case MalAtom(atom) => eval_ast(atom, env)
    case MalColl(c) => eval_ast(c, env) match {
      case MalColl.cons(MalFunction(pf), cdr) => pf(cdr)
      case other => other
    }
  }

  def eval_ast(ast: MalType, env: Env): MalType = ast match {
    case s: MalSymbol => env(s)
    case c: MalColl => c.map(eval(_, env))
    case _ => ast
  }

  def print(ast: MalType): String = printer.pr_str(ast)

  def rep(str: String, env: Env): String = print(eval(read(str), env))

  def main(args: Array[String]): Unit = {
    val repl_env: Env = Map(
      MalSymbol('+) -> MalFunction {
        case MalList.of(MalReal(x), MalReal(y)) => MalReal(x + y)
      },
      MalSymbol('-) -> MalFunction {
        case MalList.of(MalReal(x), MalReal(y)) => MalReal(x - y)
      },
      MalSymbol('*) -> MalFunction {
        case MalList.of(MalReal(x), MalReal(y)) => MalReal(x * y)
      },
      MalSymbol('/) -> MalFunction {
        case MalList.of(MalReal(x), MalReal(y)) => MalReal(x / y)
      }
    )
    @tailrec
    def go(): Unit = Option(StdIn.readLine("user> ")) match {
      case Some(line) =>
        try println(rep(line, repl_env)) catch {
          case NonFatal(e) => e.printStackTrace()
        }
        go()
      case None => ()
    }

    go()
  }
}
