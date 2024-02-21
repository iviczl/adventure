import { Position } from '../../types/position'
import { Signal } from '@preact/signals-react'
import { AppState } from '../../state'
import { takeAction } from '../../services/gameService'
import { MainContainer } from '../../globalStyles'
import { Item } from './gameStyles'

const Game = ({ state }: { state: Signal<AppState> }) => {
  const position = state.value.actualPosition as Position

  async function action(id: string) {
    await takeAction(id)
  }

  if (!state.value.selectedGameId) {
    return null
  }

  return (
    <MainContainer>
      <p>{position.description}</p>
      <section>
        {position.available_actions.map((a) => (
          <Item key={a.id}>{a.description}</Item>
        ))}
      </section>
      <section>
        {position.items.map((i) => (
          <p onClick={() => action(i.id)}>{i.description}</p>
        ))}
      </section>
    </MainContainer>
  )
}

export default Game
