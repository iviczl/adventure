import { Action } from './action'
import { Item } from './item'

export type Position = {
  id: string
  description: string
  available_actions: Action[]
  items: Item[]
  end_position: boolean
}
