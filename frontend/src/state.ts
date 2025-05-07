import { signal } from '@preact/signals-react'
import { GameInfo } from './types/gameInfo'
import { Position } from './types/position'
import { PlayInfo } from './types/play'

export type AppState = {
  userName: string
  token: string | undefined
  player: string
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
  games: undefined,
  selectedGameId: '',
  adventureId: '',
  actualPosition: undefined,
} as AppState)
