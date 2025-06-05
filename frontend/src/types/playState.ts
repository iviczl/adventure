import { Player } from './player'
import { Position } from './position'

export type PlayState = {
  adventureId: string
  player: Player
  actualPosition: Position
}
