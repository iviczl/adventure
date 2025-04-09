import { Action } from './action'
import { Item } from './item'

export type Position = {
  id: string
  description: string
  availableActions: Action[]
  items: Item[]
  endPosition: boolean
}
