import { getGames } from './services/gameService'
import Home from './components/Home/Home'
import Game from './components/Game/Game'
import { state } from './state'
import { MainContainer } from './globalStyles'
import Login from './components/Login/Login'

function App() {
  async function init() {
    await getGames()
    console.log('app', state.value)
  }

  init()

  return (
    <MainContainer>
      <Login state={state} />
      <Home state={state} />
      <Game state={state} />
    </MainContainer>
  )
}

export default App
