import { Position } from '../../types/position'
import { Signal } from '@preact/signals-react'
import { AppState } from '../../state'
import { takeAction } from '../../services/gameService'
import { MainContainer } from '../../globalStyles'
import { Item } from './gameStyles'

const Game = ({ state }: { state: Signal<AppState> }) => {
  const position = state.value.actualPosition as Position

  async function callAction(id: string) {
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
          <Item key={a.id} onClick={() => callAction(a.id)}>
            {a.description}
          </Item>
        ))}
      </section>
      <section>
        {position.items.map((i) => (
          <Item key={i.id} onClick={() => callAction(i.id)}>
            {i.description}
          </Item>
        ))}
      </section>
    </MainContainer>
  )
}

export default Game
