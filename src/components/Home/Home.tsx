import { startGame } from '../../services/gameService'
import { GameInfo } from '../../types/gameInfo'
import { Signal } from '@preact/signals-react'
import { AppState } from '../../state'
import { useState } from 'react'
import {
  Button,
  MainHeading,
  Container,
  Input,
  Select,
  Label,
  Paragraph,
} from '../../globalStyles'
import { Row, StartGrid } from './homeStyles'

export default function Home({ state }: { state: Signal<AppState> }) {
  const [player, setPlayer] = useState('')
  const [selectedGameId, setSelectedGameId] = useState('')

  async function start() {
    await startGame(selectedGameId, player)
  }

  function gameItem(game: GameInfo) {
    if (!game) return
    return (
      <option key={game.id} value={game.id}>
        {game.title}
      </option>
    )
  }

  function selectedGameDescription() {
    if (selectedGameId) {
      return state.value.games?.find((g) => g.id === selectedGameId)
        ?.description
    }

    return ''
  }

  if (state.value.selectedGameId) {
    return null
  }

  if (!state.value.games) {
    return 'Loading...'
  }

  return (
    <Container>
      <MainHeading>Adventure games</MainHeading>
      <Paragraph>
        Select the game, state your name and press Start to begin.
      </Paragraph>
      <StartGrid>
        <Label>Games to choose from:</Label>
        <Select
          $marginLeft='0'
          onChange={(e) => setSelectedGameId(e.target.value)}
        >
          <option value=''>(Select an option)</option>
          {state.value.games && state.value.games.map(gameItem)}
        </Select>
        <div>{selectedGameDescription()}</div>
        <Label>Player name:</Label>
        <Input
          type='text'
          $marginLeft='0'
          onChange={(e) => setPlayer(e.target.value)}
          value={player}
        />
      </StartGrid>
      <Row>
        <Button onClick={start} disabled={!selectedGameId || !player}>
          Start
        </Button>
      </Row>
    </Container>
  )
}
