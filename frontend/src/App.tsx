import { getGames } from './services/gameService'
import Home from './components/Home/Home'
import Game from './components/Game/Game'
import { state } from './state'
import { MainContainer } from './globalStyles'

function App() {
  async function init() {
    await getGames()
    console.log('app', state.value)
  }

  init()

  return (
    <MainContainer>
      <Home state={state} />
      <Game state={state} />
    </MainContainer>
  )
}

export default App
