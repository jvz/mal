import types._

object printer {
  private val escape = raw"""("|\\)""".r
  private val newline = raw"""\R""".r

  def pr_str(v: MalType, printReadably: Boolean = true): String = v match {
    case MalAtom(atom) => atom match {
      case MalTrue => "true"
      case MalFalse => "false"
      case MalString(value) => if (printReadably) escape_str(value) else value
      case MalInt(value) => value.toString
      case MalReal(value) => value.toString
      case MalSymbol(value) => value.name
      case MalKeyword(value) => s":$value"
    }
    case MalList(value) => value.map(pr_str(_)).mkString("(", " ", ")")
    case MalVector(value) => value.map(pr_str(_)).mkString("[", " ", "]")
    case MalMap(value) => (for ((k, v) <- value) yield pr_str(k) + " " + pr_str(v)).mkString("{", " ", "}")
  }

  private def escape_str(str: String): String =
    (for (line <- newline.split(str)) yield escape.replaceAllIn(line, m => "\\" + m.matched)).mkString("\"", "\\n", "\"")
}
