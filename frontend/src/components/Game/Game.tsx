import { Position } from '../../types/position'
import { Signal } from '@preact/signals-react'
import { AppState } from '../../state'
import { takeAction } from '../../services/gameService'
import { Container } from '../../globalStyles'
import { Item } from './gameStyles'
import { useState } from 'react'
import { Row } from '../../globalStyles'

const Game = ({ state }: { state: Signal<AppState> }) => {
  const position = state.value.actualPosition as Position
  const [requestProcessing, setRequestProcessing] = useState(false)

  async function callAction(id: string) {
    setRequestProcessing(true)
    await takeAction(id, state.value.player, state.value.adventureId)
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
    <Container>
      <p>{position.description}</p>
      <Row>
        {position.availableActions.map((a) => (
          <Item
            key={a.code}
            onClick={() => callAction(a.code)}
            disabled={requestProcessing}
          >
            {a.description}
          </Item>
        ))}
        {!position.endPosition || (
          <Item onClick={home} disabled={requestProcessing}>
            Home
          </Item>
        )}
      </Row>
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
    </Container>
  )
}

export default Game
