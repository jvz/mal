import types._

import scala.annotation.tailrec
import scala.io.StdIn
import scala.util.control.NonFatal

object step2_eval {

  type Env = Map[MalSymbol, MalType]

  def read(str: String): MalType = reader.read_str(str)

  def eval_ast(ast: MalType, env: Env): MalType = ast match {
    case s: MalSymbol => env(s)
    case c: MalColl => c.map(eval(_, env))
    case _ => ast
  }

  def eval(ast: MalType, env: Env): MalType = ast match {
    case MalList(_) =>
      eval_ast(ast, env) match {
        case MalList(MalLambda(f) :: args) => f(args)
        case nil @ MalList(Nil) => nil
        case _ => core.syntax_error
      }
    case _ => eval_ast(ast, env)
  }

  def print(ast: MalType): String = ast.show()

  def rep(str: String, env: Env): String = print(eval(read(str), env))

  def main(args: Array[String]): Unit = {
    @tailrec
    def go(env: Env): Unit = Option(StdIn.readLine("user> ")) match {
      case Some(line) =>
        try println(rep(line, env)) catch {
          case NonFatal(e) => e.printStackTrace()
        }
        go(env)
      case None => ()
    }

    go(core.ns)
  }
}
