import types._

import scala.collection.mutable

object env {

  class Env private(outer: Option[Env], binds: Iterator[MalSymbol], exprs: Iterator[MalType]) {
    private val data: mutable.Map[MalSymbol, MalType] = mutable.Map.empty
    for (bind <- binds) {
      bind match {
        case MalSymbol.sp.Variadic =>
          val symbol = binds.next()
          val expr = MalList(exprs.toList)
          data(symbol) = expr
        case _ =>
          data(bind) = exprs.next()
      }
    }

    def apply(symbol: MalSymbol): Option[MalType] = for (env <- find(symbol)) yield env.data(symbol)

    def update(symbol: MalSymbol, value: MalType): Unit = data(symbol) = value

    def find(symbol: MalSymbol): Option[Env] = data.get(symbol) match {
      case Some(_) => Some(this)
      case None => outer.flatMap(_.find(symbol))
    }

    def inner(binds: Seq[MalSymbol] = Nil, exprs: Seq[MalType] = Nil): Env =
      new Env(Some(this), binds.iterator, exprs.iterator)
  }

  object Env {
    def apply(binds: Seq[MalSymbol] = Nil, exprs: Seq[MalType] = Nil): Env =
      new Env(None, binds.iterator, exprs.iterator)
    implicit def toSymbol(symbol: Symbol): MalSymbol = MalSymbol(symbol)
  }

}
