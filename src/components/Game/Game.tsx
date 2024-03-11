import { Position } from '../../types/position'
import { Signal } from '@preact/signals-react'
import { AppState } from '../../state'
import { takeAction } from '../../services/gameService'
import { MainContainer } from '../../globalStyles'
import { Item } from './gameStyles'
import { useState } from 'react'

const Game = ({ state }: { state: Signal<AppState> }) => {
  const position = state.value.actualPosition as Position
  const [requestProcessing, setRequestProcessing] = useState(false)

  async function callAction(id: string) {
    setRequestProcessing(true)
    await takeAction(id)
    setRequestProcessing(false)
  }

  function home() {
    state.value = {
      ...state.value,
      selectedGameId: '',
    }
  }

  if (!state.value.selectedGameId) {
    return null
  }

  return (
    <MainContainer>
      <p>{position.description}</p>
      <section>
        {position.available_actions.map((a) => (
          <Item
            key={a.code}
            onClick={() => callAction(a.code)}
            disabled={requestProcessing}
          >
            {a.description}
          </Item>
        ))}
        {!position.end_position || (
          <Item onClick={home} disabled={requestProcessing}>
            Home
          </Item>
        )}
      </section>
      {/* <section>
        {position.items.map((i) => (
          <Item
            key={i.id}
            onClick={() => callAction(i.id)}
            disabled={requestProcessing}
          >
            {i.description}
          </Item>
        ))}
      </section> */}
    </MainContainer>
  )
}

export default Game
