import env.Env
import types.MalSymbol.sp.{Def, Let}
import types._

import scala.annotation.tailrec
import scala.io.StdIn
import scala.util.control.NonFatal

object step3_env {

  def read(str: String): MalType = reader.read_str(str)

  def eval_ast(ast: MalType, env: Env): MalType = ast match {
    case s: MalSymbol => env(s).get
    case c: MalColl => c.map(eval(_, env))
    case _ => ast
  }

  def eval(ast: MalType, env: Env): MalType = ast match {
    case MalList(Def :: (bind: MalSymbol) :: expr :: Nil) =>
      val ev = eval(expr, env)
      env(bind) = ev
      ev

    case MalList(Let :: MalColl(binds) :: expr :: Nil) =>
      val inner = env.inner()
      for ((s: MalSymbol, value) <- binds.pairs) {
        inner(s) = eval(value, inner)
      }
      eval(expr, inner)

    case MalList(_) =>
      eval_ast(ast, env) match {
        case MalList(MalLambda(f) :: args) => f(args)
        case nil @ MalList(Nil) => nil
        case _ => core.syntax_error
      }

    case _ => eval_ast(ast, env)
  }

  def print(ast: MalType): String = printer.pr_str(ast)

  def rep(str: String, env: Env): String = print(eval(read(str), env))

  def main(args: Array[String]): Unit = {

    def repl_env: Env = {
      val env = Env()
      for ((sym, expr) <- core.ns) {
        env(sym) = expr
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
