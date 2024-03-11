import { startGame } from '../../services/gameService'
import { GameInfo } from '../../types/gameInfo'
import { Signal } from '@preact/signals-react'
import { AppState } from '../../state'
import { useEffect, useState } from 'react'
import {
  Button,
  MainHeading,
  Container,
  Input,
  Select,
  Label,
  Paragraph,
  Option,
  SubHeading,
} from '../../globalStyles'
import { Description, Row, StartGrid } from './homeStyles'

export default function Home({ state }: { state: Signal<AppState> }) {
  const [player, setPlayer] = useState('')
  const [selectedGameId, setSelectedGameId] = useState(
    state.value.selectedGameId
  )
  const [requestProcessing, setRequestProcessing] = useState(false)

  useEffect(() => {
    setSelectedGameId(state.value.selectedGameId)
  }, [state.value.selectedGameId])

  async function start() {
    setRequestProcessing(true)
    await startGame(selectedGameId, player)
    setRequestProcessing(false)
  }

  function cancel() {
    setSelectedGameId('')
  }

  function gameItem(game: GameInfo) {
    if (!game) return
    return (
      <Option key={game.id} value={game.id}>
        {game.title}
      </Option>
    )
  }

  function selectedGameTitle() {
    if (selectedGameId) {
      return state.value.games?.find((g) => g.id === selectedGameId)?.title
    }
    return ''
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
        Select a game, state your name and press Start to begin.
      </Paragraph>
      <StartGrid>
        {selectedGameId != '' || (
          <Row>
            <Label $width='10rem'>Games to choose from:</Label>
            <Select
              $width='12rem'
              onChange={(e) => setSelectedGameId(e.target.value)}
            >
              <Option value=''>(Select an option)</Option>
              {state.value.games && state.value.games.map(gameItem)}
            </Select>
          </Row>
        )}
        {!selectedGameId || (
          <div>
            <Row>
              <SubHeading>{selectedGameTitle()}</SubHeading>
            </Row>
            <Row>
              <Description>{selectedGameDescription()}</Description>
            </Row>
            <Row>
              <Button onClick={cancel} disabled={requestProcessing}>
                Cancel
              </Button>
            </Row>
          </div>
        )}
        <Row>
          <Label $width='10rem'>Player name:</Label>
          <Input
            $width='12rem'
            type='text'
            onChange={(e) => setPlayer(e.target.value)}
            value={player}
          />
        </Row>
      </StartGrid>
      <Row>
        <Button
          onClick={start}
          disabled={!selectedGameId || !player || requestProcessing}
        >
          Start
        </Button>
      </Row>
    </Container>
  )
}
