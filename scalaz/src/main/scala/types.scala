import scala.annotation.tailrec

object types {

  sealed trait MalType

  sealed trait MalColl extends MalType

  sealed trait MalList extends MalColl
  case object MalNil extends MalList
  case class MalCons(car: MalType, cdr: MalType) extends MalList {
    def tupled: (MalType, MalType) = (car, cdr)
  }
  object MalCons {
    def apply(tuple: (MalType, MalType)): MalCons = MalCons(tuple._1, tuple._2)
  }
  object MalList {
    def apply(args: MalType*): MalList = args.foldRight(MalNil: MalList) { (arg, cdr) => MalCons(arg, cdr) }
    def unapply(arg: MalList): Option[List[MalType]] = {
      @tailrec
      def go(l: MalList, acc: List[MalType]): Option[List[MalType]] = l match {
        case MalNil => Some(acc.reverse)
        case MalCons(car, cdr: MalList) => go(cdr, car :: acc)
        case _ => None
      }
      go(arg, Nil)
    }
//    def unapplySeq(arg: MalList): Option[Seq[MalType]] = unapply(arg)
  }

  case class MalVector(value: Vector[MalType]) extends MalColl
  case class MalMap(value: Map[MalType, MalType]) extends MalColl

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
