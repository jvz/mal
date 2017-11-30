import fastparse.all._
import types._

object reader {

  val ws: P[Unit] = {
    val whitespace = P(CharsWhile(_.isWhitespace))
    val comment = P(";" ~/ CharsWhile(_ != '\n'))
    val comma = P(",")
    P(whitespace | comment | comma).rep
  }

  trait MalParser[T <: MalType] {
    def p: P[T]
  }

  object atom extends MalParser[MalAtom] {
    object numeric extends MalParser[MalNumeric] {
      private val sign = P(CharIn("+-").?)
      private val digits = P(CharsWhile(_.isDigit))
      private val integral = P("0" | CharIn('1' to '9') ~ digits.?)
      val int: P[MalInt] = P(sign ~ integral).!.map(MalInt(_))

      private val exponent = P(CharIn("eE") ~ sign ~ digits)
      private val fractional = P("." ~ digits)
      val real: P[MalReal] = P(sign ~ integral ~ fractional.? ~ exponent.?).!.map(MalReal(_))

      val p: P[MalNumeric] = P(int | real)
    }

    object boolean extends MalParser[MalBoolean] {
      val t: P[MalBoolean] = P("true").map(_ => MalTrue)
      val f: P[MalBoolean] = P("false").map(_ => MalFalse)
      val p: P[MalBoolean] = P(t | f)
    }

    object string extends MalParser[MalAtom] {
      private val escape = P("\\" ~ CharIn("\"\\n"))
      private val strChars = P(CharsWhile(!"\"\\".contains(_: Char)))
      val string: P[MalString] = P("\"" ~/ (strChars | escape).rep.! ~ "\"").map(MalString)

      private val reserved = P(CharIn("[]{}()'`\""))
      private val sym = P(CharsWhile(c => !(c.isWhitespace || ",;[]{}()'`\"".contains(c)))).!
      val symbol: P[MalSymbol] = sym.map(MalSymbol(_))
      val keyword: P[MalKeyword] = P(":" ~/ sym).map(MalKeyword)

      val p: P[MalAtom] = P(string | keyword | symbol)
    }

    val p: P[MalAtom] = P(numeric.p | boolean.p | string.p)
  }

  object coll extends MalParser[MalColl] {
    private def listLike[T <: MalType](start: String, end: String, term: P[T]) =
      P(start ~ ws ~/ term.rep(sep = ws) ~/ ws ~ end)
    val list: P[MalList] = listLike("(", ")", form).map(MalList(_: _*))
    val vector: P[MalVector] = listLike("[", "]", form).map(s => MalVector(s.toVector))

    private val keyVal: P[MalCons] = P(atom.p ~/ ws ~ form).map(MalCons(_))
    private def keyValsToMap(kvs: Seq[MalCons]): MalMap = MalMap(kvs.map(_.tupled).toMap)
    val map: P[MalMap] = listLike("{", "}", keyVal).map(keyValsToMap)

    val p: P[MalColl] = P(list | vector | map)
  }

  object macros extends MalParser[MalList] {
    private def functionLike(prefix: String, name: String) = P(prefix ~/ form).map(MalList(MalSymbol(name), _))
    val splice = functionLike("~@", "splice-unquote")
    val quote = functionLike("'", "quote")
    val quasi = functionLike("`", "quasiquote")
    val unquote = functionLike("~", "unquote")
    val deref = functionLike("@", "deref")

    // TEST: ^{"a" 1} [1 2 3] -> (with-meta [1 2 3] {"a" 1})
//    val meta =

    val p: P[MalList] = P(splice | quote | quasi | unquote | deref)
  }

  val form: P[MalType] = P(macros.p | coll.p | atom.p)

  def read_str(expr: String): MalType = P(ws ~ form ~ ws).parse(expr).get.value

}
