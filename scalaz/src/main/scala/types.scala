import java.util.concurrent.atomic.AtomicReference

import env.Env

object types {

  sealed trait MalType {
    def eql(that: MalType): Boolean
    def show(pretty: Boolean = true): String
    def ? : Boolean = true
  }
  type MalF = PartialFunction[List[MalType], MalType]
  type MalPair = (MalType, MalType)

  sealed trait MalColl extends MalType with Iterable[MalType] {
    def map(f: MalType => MalType): MalColl
    def flatMap(f: MalType => MalColl): MalColl
    def pairs: Iterator[MalPair] = for {
      List(a, b) <- toList.grouped(2)
    } yield (a, b)
    override def eql(that: MalType): Boolean = if (isEmpty) that match {
      case MalColl(coll) if coll.isEmpty => true
      case _ => false
    } else that match {
      case MalColl(coll) if coll.isEmpty => false
      case MalColl(coll) =>
        val as = this.toSeq
        val bs = coll.toSeq
        as.size == bs.size && as.zip(bs).forall {
          case (a, b) => a eql b
        }
    }
  }
  object MalColl {
    def unapply(arg: MalType): Option[MalColl] = arg match {
      case c: MalColl => Some(c)
      case _ => None
    }
  }

  final case class MalList(value: List[MalType]) extends MalColl {
    override def map(f: MalType => MalType): MalList = copy(value map f)
    override def flatMap(f: MalType => MalColl): MalColl = copy(value.map(f).flatMap(_.toList))
    override def show(pretty: Boolean): String = value.map(_.show(pretty)).mkString("(", " ", ")")
    override def iterator: Iterator[MalType] = value.iterator
  }
  object MalList {
    def apply(args: MalType*): MalList = MalList(args.toList)
  }

  final case class MalVector(value: Vector[MalType]) extends MalColl {
    override def map(f: MalType => MalType): MalVector = copy(value map f)
    override def flatMap(f: MalType => MalColl): MalColl = copy(value.map(f).flatMap(_.toVector))
    override def show(pretty: Boolean): String = value.map(_.show(pretty)).mkString("[", " ", "]")
    override def iterator: Iterator[MalType] = value.iterator
  }
  object MalVector {
    def apply(args: MalType*): MalVector = MalVector(args.toVector)
  }

  final case class MalMap(value: Map[MalType, MalType]) extends MalColl {
    override def map(f: MalType => MalType): MalMap = copy(value mapValues f)
    override def flatMap(f: MalType => MalColl): MalColl = map(f) // TODO
    override def show(pretty: Boolean): String =
      utils.flatten(value).map(_.show(pretty)).mkString("{", " ", "}")
    override def iterator: Iterator[MalType] = utils.flatten(value).iterator
  }
  object MalMap {
    def apply(args: MalPair*): MalMap = MalMap(args.toMap)
  }

  sealed trait MalFn extends MalType {
    def apply(args: List[MalType]): MalType
  }
  object MalFn {
    def unapply(arg: MalType): Option[MalFn] = arg match {
      case f: MalFn => Some(f)
      case _ => None
    }
  }

  final case class MalLambda(lambda: MalF) extends MalFn {
    override def show(pretty: Boolean): String = lambda.toString()
    override def eql(that: MalType): Boolean = that match {
      case MalLambda(other) => lambda == other
      case _ => false
    }
    override def apply(args: List[MalType]): MalType = lambda(args)
  }

  final case class MalFunction(params: Seq[MalSymbol], body: MalType, env: Env, lambda: MalLambda) extends MalFn {
    override def eql(that: MalType): Boolean = this == that
    override def show(pretty: Boolean): String = toString
    def closure(args: Seq[MalType]): Env = env.inner(params, args)
    override def apply(args: List[MalType]): MalType = lambda(args)
  }

  // note that this terminology is lifted from traditional lisp; what mal refers to as an "atom" is called a "ref" here
  sealed trait MalAtom extends MalType {
    override def eql(that: MalType): Boolean = this == that
  }
  object MalAtom {
    def unapply(arg: MalAtom): Option[MalAtom] = Some(arg)
  }

  // what mal calls an atom
  final case class MalRef(value: AtomicReference[MalType]) extends MalType {
    override def eql(that: MalType): Boolean = that match {
      case MalRef(other) => value.get() eql other.get()
      case _ => false
    }
    override def show(pretty: Boolean): String = s"(atom ${value.get().show(pretty)})"
  }

  // this particular implementation is similar to Unit
  final case object MalNil extends MalAtom {
    override def show(pretty: Boolean): String = "nil"
    override def ? : Boolean = false
  }

  sealed abstract class MalBoolean(val value: Boolean) extends MalAtom {
    override def show(pretty: Boolean): String = value.toString
    override def ? : Boolean = value
  }
  final case object MalTrue extends MalBoolean(true)
  final case object MalFalse extends MalBoolean(false)

  final case class MalString(value: String) extends MalAtom {
    override def show(pretty: Boolean): String = if (pretty) utils.escape(value) else value
  }
  final case class MalSymbol(value: Symbol) extends MalAtom {
    override def show(pretty: Boolean): String = value.name
  }
  object MalSymbol {
    def apply(s: String): MalSymbol = MalSymbol(Symbol(s))
    object sp {
      val Def: MalSymbol = MalSymbol("def!")
      val Let: MalSymbol = MalSymbol("let*")
      val Do: MalSymbol = MalSymbol('do)
      val If: MalSymbol = MalSymbol('if)
      val Fn: MalSymbol = MalSymbol("fn*")
      val Variadic: MalSymbol = MalSymbol('&)
      val Args: MalSymbol = MalSymbol("*ARGV*")
    }
  }
  final case class MalKeyword(value: String) extends MalAtom {
    override def show(pretty: Boolean): String = s":$value"
  }

  sealed trait MalNumeric extends MalAtom
  final case class MalInt(value: BigInt) extends MalNumeric {
    override def show(pretty: Boolean): String = value.toString
  }
  object MalInt {
    def apply(s: String): MalInt = MalInt(BigInt(s))
  }

  final case class MalReal(value: BigDecimal) extends MalNumeric {
    override def show(pretty: Boolean): String = value.toString
  }
  object MalReal {
    def apply(s: String): MalReal = MalReal(BigDecimal(s))
  }

  object utils {
    def flatten[A](map: Map[A, A]): Seq[A] =
      map.toSeq.flatMap(tuple => Seq(tuple._1, tuple._2))

    def escape(str: String): String =
      "\"" + str.replace("\\", "\\\\").replace("\"", "\\\"").replace("\n", "\\n") + "\""
  }

  implicit class MalHelper(private val sc: StringContext) extends AnyVal {
    def mal(args: Any*): String = {
      val strings = sc.parts.iterator
      val exprs = args.iterator
      val sb = new StringBuilder(strings.next())
      for ((expr, string) <- exprs.zip(strings)) {
        sb.append {
          expr match {
            case s: String => utils.escape(s)
            case s: Symbol => s.name
            case t: MalType => t.show()
            case _ => throw new IllegalArgumentException(s"Invalid type found: $expr")
          }
        }
          .append(string)
      }
      sb.toString()
    }
  }

}
