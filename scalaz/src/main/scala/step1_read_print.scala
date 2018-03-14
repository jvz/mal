import types._

import scala.annotation.tailrec
import scala.io.StdIn
import scala.util.control.NonFatal

object step1_read_print {

  def read(str: String): MalType = reader.read_str(str)

  def eval(ast: MalType, env: String): MalType = ast

  def print(ast: MalType): String = ast.show()

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
