import { Signal } from '@preact/signals-react'
import { AppState } from '../../state'
import {
  Container,
  Label,
  Paragraph,
  Row,
  Select,
  Option,
  Button,
} from '../../globalStyles'
import { StartGrid } from './homeStyles'
import { useState } from 'react'
import { PlayInfo } from '../../types/play'
import { loadGame } from '../../services/gameService'

export default function Saved({ state }: { state: Signal<AppState> }) {
  const [requestProcessing, setRequestProcessing] = useState(false)
  const [selectedPlayId, setSelectedPlayId] = useState('')

  async function continuePlay() {
    setRequestProcessing(true)
    await loadGame(selectedPlayId)
    setRequestProcessing(false)
  }

  function playItem(play: PlayInfo) {
    return (
      <Option key={play.id} value={play.id}>
        {/* {new Date(play.createdAt).toLocaleDateString(undefined, {
          year: 'numeric',
          month: '2-digit',
          day: '2-digit',
          hour: '2-digit',
          hour12: false,
          minute: '2-digit',
          second: '2-digit',
        })}{' '} */}
        {new Date(play.updatedAt).toLocaleDateString(undefined, {
          year: 'numeric',
          month: '2-digit',
          day: '2-digit',
          hour: '2-digit',
          hour12: false,
          minute: '2-digit',
          second: '2-digit',
        })}{' '}
        {play.adventureTitle}
      </Option>
    )
  }

  if (!state.value.plays) {
    return 'Loading...'
  }

  return (
    <Container>
      <Paragraph>Select a gameplay and press Continue to begin.</Paragraph>
      <StartGrid>
        <Row>
          <Label $width='100%'>Gameplays to choose from:</Label>
          <Select
            $width='12rem'
            onChange={(e) => setSelectedPlayId(e.target.value)}
          >
            <Option value=''>(Select an option)</Option>
            {state.value.plays && state.value.plays.map(playItem)}
          </Select>
        </Row>
        <Row>
          <Button
            onClick={() => continuePlay()}
            disabled={requestProcessing || !selectedPlayId}
          >
            Continue
          </Button>
        </Row>
      </StartGrid>
    </Container>
  )
}
