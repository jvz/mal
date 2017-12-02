import env.Env
import types.MalSymbol.sp._
import types._

import scala.annotation.tailrec
import scala.io.StdIn
import scala.util.control.NonFatal

object step4_if_fn_do {

  def read(str: String): MalType = reader.read_str(str)

  def eval_ast(ast: MalType, env: Env): MalType = ast match {
    case s: MalSymbol => env(s).get
    case c: MalColl => c.map(eval(_, env))
    case _ => ast
  }

  def eval(ast: MalType, env: Env): MalType = ast match {
    case MalList(Nil) => ast

    case MalList(Def :: (bind: MalSymbol) :: expr :: Nil) =>
      val ev = eval(expr, env)
      env(bind) = ev
      ev

    case MalList(Let :: MalColl(binds) :: expr :: Nil) =>
      val letEnv = env.inner()
      for ((s: MalSymbol, value) <- binds.pairs) {
        letEnv(s) = eval(value, letEnv)
      }
      eval(expr, letEnv)

    case MalList(Do :: rest) =>
      eval_ast(MalList(rest), env).asInstanceOf[MalList].value match {
        case Nil => MalNil
        case results => results.last
      }

    case MalList(If :: condition :: ifTrue :: rest) =>
      if (eval(condition, env).?) eval(ifTrue, env)
      else rest match {
        case Nil => MalNil
        case ifFalse :: Nil => eval(ifFalse, env)
        case _ => core.syntax_error
      }

    case MalList(Fn :: MalColl(args) :: expr :: Nil) =>
      MalLambda {
        case params =>
          val binds = args.toSeq.map(_.asInstanceOf[MalSymbol])
          val exprs = params
          eval(expr, env.inner(binds, exprs))
      }

    case MalList(_) =>
      eval_ast(ast, env) match {
        case MalList(MalLambda(f) :: args) => f(args)
        case other => other
//        case _ => core.syntax_error
      }

    case _ => eval_ast(ast, env)
  }

  def print(ast: MalType): String = ast.show()

  def rep(str: String, env: Env): String = print(eval(read(str), env))

  def main(args: Array[String]): Unit = {

    def repl_env(): Env = {
      val env = Env()
      // TODO: it may be useful to make core functions immutable?
      for ((sym, fn) <- core.ns) env(sym) = fn
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

    val not = "(def! not (fn* (a) (if a false true)))"
    val env = repl_env()
    rep(not, env)
    go(env)
  }
}
