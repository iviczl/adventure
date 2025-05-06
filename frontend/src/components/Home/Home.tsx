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
  Row,
} from '../../globalStyles'
import { Description, StartGrid } from './homeStyles'
import TabSelect from '../TabSelect/TabSelect'
import Tab from '../TabSelect/Tab'

export default function Home({ state }: { state: Signal<AppState> }) {
  const [player, setPlayer] = useState('')
  const [selectedGameId, setSelectedGameId] = useState(
    state.value.selectedGameId
  )
  const [selectedTab, setSelectedTab] = useState('existing')
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
  async function continuePlay() {}

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

  function selectTab(id: string) {
    setSelectedTab(id)
  }

  if (state.value.userName === '' || state.value.selectedGameId) {
    return null
  }

  if (!state.value.games) {
    return 'Loading...'
  }

  return (
    <>
      <MainHeading>Adventure games</MainHeading>
      <TabSelect>
        <Tab
          id='saved'
          title='Saved Gameplays'
          selectedId={selectedTab}
          onClick={selectTab}
        />
        <Tab
          id='new'
          title='New Game'
          selectedId={selectedTab}
          onClick={selectTab}
        />
      </TabSelect>
      {selectedTab === 'saved' ? (
        <Container>
          <Paragraph>Select a gameplay and press Continue to begin.</Paragraph>
          <StartGrid>
            <Row>
              <Label $width='100%'>Gameplays to choose from:</Label>
              <Select
                $width='12rem'
                onChange={(e) => setSelectedGameId(e.target.value)}
              >
                <Option value=''>(Select an option)</Option>
                {state.value.games && state.value.games.map(gameItem)}
              </Select>
            </Row>
            <Row>
              <Button
                onClick={() => continuePlay()}
                disabled={requestProcessing}
              >
                Continue
              </Button>
            </Row>
          </StartGrid>
        </Container>
      ) : null}
      {selectedTab === 'new' ? (
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
      ) : null}
    </>
  )
}
