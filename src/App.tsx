import { getGames } from './services/gameService'
import Home from './components/Home/Home'
import Game from './components/Game/Game'
import { state } from './state'

function App() {
  async function init() {
    await getGames()
    console.log('app', state.value)
  }

  init()

  return (
    <>
      <Home state={state} />
      <Game state={state} />
    </>
  )
}

export default App
