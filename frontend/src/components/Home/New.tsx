import { Signal } from '@preact/signals-react'
import {
  Container,
  Label,
  Paragraph,
  Row,
  Select,
  Option,
  Button,
  SubHeading,
  Input,
} from '../../globalStyles'
import { AppState } from '../../state'
import { Description, StartGrid } from './homeStyles'
import { useEffect, useState } from 'react'
import { startGame } from '../../services/gameService'
import { GameInfo } from '../../types/gameInfo'

export default function New({ state }: { state: Signal<AppState> }) {
  const [player, setPlayer] = useState(state.value.userName)
  const [selectedGameId, setSelectedGameId] = useState(
    state.value.selectedGameId
  )
  const [requestProcessing, setRequestProcessing] = useState(false)

  useEffect(() => {
    setSelectedGameId(state.value.selectedGameId)
  }, [state.value.selectedGameId])

  async function start() {
    setRequestProcessing(true)
    state.value = {
      ...state.value,
      player,
    }
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

  if (!state.value.games) {
    return 'Loading...'
  }

  return (
    <Container>
      <Paragraph>
        Select a game, state your name and press Start to begin.
      </Paragraph>
      <StartGrid>
        {selectedGameId != '' || (
          <Row>
            <Label $width='100%'>Games to choose from:</Label>
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
          <>
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
          </>
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
