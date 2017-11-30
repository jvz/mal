import types._

import scala.annotation.tailrec
import scala.collection.mutable
object env {

  class Env private (val outer: Option[Env] = None, private val data: mutable.Map[MalSymbol, MalType] = mutable.Map.empty) {
    def apply(symbol: MalSymbol): Option[MalType] = for (env <- find(symbol)) yield env.data(symbol)

    def update(symbol: MalSymbol, value: MalType): Unit = data(symbol) = value

    @tailrec final def find(symbol: MalSymbol): Option[Env] = data.get(symbol) match {
      case Some(_) => Some(this)
      case None => outer match {
        case Some(env) => env.find(symbol)
        case None => None
      }
    }

    def inner(): Env = new Env(Some(this))
  }

  object Env {
    def apply(): Env = new Env()
    implicit def toSymbol(symbol: Symbol): MalSymbol = MalSymbol(symbol)
  }

}
