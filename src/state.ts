import { signal } from '@preact/signals-react'
import { GameInfo } from './types/gameInfo'
import { Position } from './types/position'

export type AppState = {
  player: string
  games: GameInfo[] | undefined
  selectedGameId: string
  actualPosition: Position | undefined
}

export const state = signal({
  player: '',
  games: undefined,
  selectedGameId: '',
  actualPosition: undefined,
} as AppState)
