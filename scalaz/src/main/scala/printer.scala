import types._

object printer {
  def pr_str(v: MalType, printReadably: Boolean = true): String = v.show(printReadably)
}
