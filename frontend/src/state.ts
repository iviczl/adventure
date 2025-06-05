import { signal } from '@preact/signals-react'
import { GameInfo } from './types/gameInfo'
import { Position } from './types/position'
import { PlayInfo } from './types/play'
import { Attribute } from './types/attribute'

export type AppState = {
  userName: string
  token: string | undefined
  player: string
  playerAttributes: Attribute[]
  games: GameInfo[] | undefined
  plays: PlayInfo[] | undefined
  selectedGameId: string
  adventureId: string
  actualPosition: Position | undefined
}

export const state = signal({
  userName: '',
  token: undefined,
  player: '',
  playerAttributes: [] as Attribute[],
  games: undefined,
  selectedGameId: '',
  adventureId: '',
  actualPosition: undefined,
} as AppState)
