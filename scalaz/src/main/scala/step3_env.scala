import env.Env, Env._
import types._

import scala.annotation.tailrec
import scala.io.StdIn
import scala.util.control.NonFatal

object step3_env {

  def read(str: String): MalType = reader.read_str(str)

  def eval(ast: MalType, env: Env): MalType = ast match {
    case MalNil => MalNil
    case MalAtom(atom) => eval_ast(atom, env)
    case MalList.of(MalSymbol.sp.Def, symbol: MalSymbol, value) =>
      val evaluated = eval(value, env)
      env(symbol) = evaluated
      evaluated
    case MalList.of(MalSymbol.sp.Let, MalColl(bindings), expr) =>
      val inner = env.inner()
      for {
        List(symbol: MalSymbol, value) <- bindings.toList.grouped(2)
      } {
        inner(symbol) = eval(value, inner)
      }
      eval(expr, inner)
    case MalColl(c) => eval_ast(c, env) match {
      case MalColl.cons(MalFunction(pf), cdr) => pf(cdr)
      case other => other
    }
  }

  def eval_ast(ast: MalType, env: Env): MalType = ast match {
    case s: MalSymbol => env(s).get
    case c: MalColl => c.map(eval(_, env))
    case _ => ast
  }

  def print(ast: MalType): String = printer.pr_str(ast)

  def rep(str: String, env: Env): String = print(eval(read(str), env))

  def main(args: Array[String]): Unit = {

    def repl_env: Env = {
      val env = Env()

      env('+) = MalFunction {
        case MalList.of(MalReal(x), MalReal(y)) => MalReal(x + y)
      }
      env('-) = MalFunction {
        case MalList.of(MalReal(x), MalReal(y)) => MalReal(x - y)
      }
      env('*) = MalFunction {
        case MalList.of(MalReal(x), MalReal(y)) => MalReal(x * y)
      }
      env('/) = MalFunction {
        case MalList.of(MalReal(x), MalReal(y)) => MalReal(x / y)
      }

      env
    }

    @tailrec
    def go(env: Env): Unit = Option(StdIn.readLine("user> ")) match {
      case Some(line) =>
        try println(rep(line, env)) catch {
          case NonFatal(e) => e.printStackTrace()
        }
        go(env)
      case None => ()
    }

    go(repl_env)
  }
}
