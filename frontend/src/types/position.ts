import { Action } from './action'
import { Item } from './item'

export type Position = {
  id: string
  adventureId: string
  description: string
  availableActions: Action[]
  items: Item[]
  endPosition: boolean
}
