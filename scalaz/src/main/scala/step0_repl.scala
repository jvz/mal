import scala.annotation.tailrec
import scala.io.StdIn
import scala.util.control.NonFatal

object step0_repl {

  def read(str: String): String = str

  def eval(ast: String, env: String): String = ast

  def print(exp: String): String = exp

  def rep(str: String): String = print(eval(read(str), ""))

  def main(args: Array[String]): Unit = {
    @tailrec
    def go(): Unit = Option(StdIn.readLine("user> ")) match {
      case Some(line) =>
        try println(rep(line)) catch {
          case NonFatal(e) => e.printStackTrace()
        }
        go()
      case None => ()
    }

    go()
  }
}
