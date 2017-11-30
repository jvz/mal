import scala.annotation.tailrec

/**
  * While this could be implemented using direct types from Scala, this provides a typesafe DSL style of interpreting
  * MAL code.
  */
object types {

  sealed trait MalType
  type MalF = PartialFunction[MalType, MalType]
  type MalTuple = (MalType, MalType)

  sealed trait MalColl extends MalType {
    // TODO: should this use MalF instead? may require implicit conversions and macros (could be fun)
    def map(f: MalType => MalType): MalColl
    def toList: List[MalType] = {
      @tailrec def go(coll: MalColl, acc: List[MalType]): List[MalType] = coll match {
        case MalColl.cons(car, MalColl(cdr)) => go(cdr, car :: acc)
        case _ => acc.reverse
      }
      go(this, Nil)
    }
  }
  object MalColl {
    object cons {
      def unapply(arg: MalType): Option[MalTuple] = arg match {
        case MalNil => None
        case MalCons(car, cdr) => Some((car, cdr))
        case MalVector(value) if value.nonEmpty => Some((value.head, MalVector(value.tail)))
        case MalMap(value) if value.nonEmpty => Some((MalCons(value.head), MalMap(value.tail)))
          // FIXME: should maps iterate over keyval conses or as a normal list?
        case _ => None
      }
    }
    def unapply(arg: MalType): Option[MalColl] = arg match {
      case c: MalColl => Some(c)
      case _ => None
    }
  }

  sealed trait MalList extends MalColl {
    override def map(f: MalType => MalType): MalList
  }
  case object MalNil extends MalList {
    override def map(f: MalType => MalType): MalList = MalNil
  }
  case class MalCons(car: MalType, cdr: MalType) extends MalList {
    override def map(f: MalType => MalType): MalList = copy(f(car), f(cdr))
    def tupled: MalTuple = (car, cdr)
  }
  object MalCons {
    def apply(tuple: (MalType, MalType)): MalCons = MalCons(tuple._1, tuple._2)
  }
  object MalList {
    private def unfold(list: MalList): Option[List[MalType]] = {
      @tailrec
      def go(l: MalList, acc: List[MalType]): Option[List[MalType]] = l match {
        case MalNil => Some(acc.reverse)
        case MalCons(car, cdr: MalList) => go(cdr, car :: acc)
        case _ => None
      }
      go(list, Nil)
    }

    private def toList(args: Seq[MalType]): MalList = args.foldRight(MalNil: MalList) { (car, cdr) => MalCons(car, cdr) }

    object of {
      def apply(args: MalType*): MalList = toList(args)
      def unapplySeq(list: MalList): Option[Seq[MalType]] = unfold(list)
    }

    def apply(list: Seq[MalType]): MalList = toList(list)
    def unapply(list: MalList): Option[List[MalType]] = unfold(list)
  }

  // general idea here is that the input type will be some sort of MalColl for multiple arguments
  case class MalFunction(pf: PartialFunction[MalType, MalType]) extends MalType

  case class MalVector(value: Vector[MalType]) extends MalColl {
    override def map(f: MalType => MalType): MalVector = copy(value map f)
  }
  // may be simpler to implement as a Seq[MalCons]? or even MalList of MalConses
  case class MalMap(value: Map[MalType, MalType]) extends MalColl {
    override def map(f: MalType => MalType): MalMap = copy(value mapValues f)
  }

  sealed trait MalAtom extends MalType
  object MalAtom {
    def unapply(arg: MalAtom): Option[MalAtom] = Some(arg)
  }

  sealed abstract class MalBoolean(val value: Boolean) extends MalAtom
  case object MalTrue extends MalBoolean(true)
  case object MalFalse extends MalBoolean(false)

  case class MalString(value: String) extends MalAtom
  case class MalSymbol(value: Symbol) extends MalAtom
  object MalSymbol {
    def apply(s: String): MalSymbol = MalSymbol(Symbol(s))
    object sp {
      val Def: MalSymbol = MalSymbol("def!")
      val Let: MalSymbol = MalSymbol("let*")
    }
  }
  case class MalKeyword(value: String) extends MalAtom

  sealed trait MalNumeric extends MalAtom
  case class MalInt(value: BigInt) extends MalNumeric
  object MalInt {
    def apply(s: String): MalInt = MalInt(BigInt(s))
  }

  case class MalReal(value: BigDecimal) extends MalNumeric
  object MalReal {
    def apply(s: String): MalReal = MalReal(BigDecimal(s))
  }

}
