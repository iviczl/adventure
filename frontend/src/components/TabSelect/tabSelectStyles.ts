import styled from 'styled-components'

export const SelectedTab = styled.div`
  border-radius: 0.5rem 0.5rem 0 0;
  border: 0;
  padding: 0.6em 1.2em;
  font-size: 1rem;
  font-weight: 500;
  font-family: inherit;
  background-color: white;
  color: #331111;
  transition: border-color 0.25s;
`

export const UnselectedTab = styled.div`
  border-radius: 0.5rem 0.5rem 0 0;
  border: 0;
  padding: 0.6em 1.2em;
  font-size: 1rem;
  font-weight: 500;
  font-family: inherit;
  background-color: #eeccbb;
  color: #331111;
  cursor: pointer;
  transition: border-color 0.25s;
`
export const TabContainer = styled.div`
  display: flex;
  flex-direction: row;
  justify-content: center;
  align-items: center;
  gap: 1rem;
  margin: 0;
  padding: 0;
  max-width: 100%;
  max-height: 100%;
  width: 100%;
  border-bottom: solid 0.2rem white;
`
