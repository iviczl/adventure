import { Position } from '../types/position'
import { state } from '../state'
import { GameInfo } from '../types/gameInfo'

let abortController: AbortController

export async function getGames() {
  abortController = new AbortController()
  const result = await doFetch('http://127.0.0.1:5000/')

  if (assertError(result)) {
    return
  }
  const list = (await result.response) as GameInfo[]
  if (list) {
    state.value = {
      ...state.value,
      games: list,
      selectedGameId: '',
    }
  }
}

export async function startGame(gameId: string, player: string) {
  abortController = new AbortController()
  const result = await doFetch('http://127.0.0.1:5000/new', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ gameId, player }),
  })

  if (assertError(result)) {
    return
  }
  // {gameId: string, position: Position}
  const position = (await result.response) as Position
  state.value = {
    ...state.value,
    selectedGameId: gameId,
    player: player,
    actualPosition: position,
  }
}

export async function takeAction(actionId: string) {
  abortController = new AbortController()
  const result = await doFetch('http://127.0.0.1:5000/do', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ gameId: state.value.selectedGameId, actionId }),
  })

  if (assertError(result)) {
    return
  }
  const position = (await result.response) as Position
  state.value = {
    ...state.value,
    actualPosition: position,
  }
}

function assertError(result: { response: Promise<any>; error: unknown }) {
  if (!result.response) {
    console.log(result.error)
    return true
  }
  return false
}

export async function doFetch(url: string, options = {}) {
  let response = null
  let error = null
  try {
    const signal = abortController.signal
    const res = await fetch(url, { ...options, signal })
    response = await res.json()
  } catch (problem) {
    abortController.abort()
    error = problem
  }
  return { response, error }
}
