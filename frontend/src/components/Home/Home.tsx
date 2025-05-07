import { Signal } from '@preact/signals-react'
import { AppState } from '../../state'
import { useEffect, useState } from 'react'
import { MainHeading } from '../../globalStyles'
import TabSelect from '../TabSelect/TabSelect'
import Tab from '../TabSelect/Tab'
import Saved from './Saved'
import New from './New'
import { getGames, getPlays } from '../../services/gameService'

export default function Home({ state }: { state: Signal<AppState> }) {
  const [selectedTab, setSelectedTab] = useState('saved')

  useEffect(() => {
    if (!state.value.plays && state.value.token) {
      getPlays()
    }
    if (!state.value.games) {
      getGames()
    }
  }, [state.value.games, state.value.plays, state.value.token])

  function selectTab(id: string) {
    setSelectedTab(id)
  }

  if (
    state.value.userName === '' ||
    state.value.selectedGameId ||
    state.value.actualPosition
  ) {
    return null
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
      {selectedTab === 'saved' ? <Saved state={state} /> : null}
      {selectedTab === 'new' ? <New state={state} /> : null}
    </>
  )
}
