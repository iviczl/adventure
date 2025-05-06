import styled from 'styled-components'
import { devices } from '../../globalStyles'

export const Item = styled.button`
  margin: 0.5rem;
  padding: 0.5rem;
  border: none;
  border-radius: 0.8rem;
  max-width: 100%;
  text-wrap: wrap;

  @media ${devices.mobile} {
    background-color: #aa8866;
    color: white;
  }

  @media ${devices.desktop} {
    background-color: #eeccbb;
    color: #331111;
    &:hover {
      background-color: #aa8866;
      color: white;
      cursor: pointer;
    }
  }
`
export const ColumnContainer = styled.section`
  margin: 0;
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
  max-width: 100%;
  max-height: 100%;
`
export const ToolHeadBar = styled.div`
  display: flex;
  flex-direction: row;
  justify-content: right;
  // align-items: right;
  margin: 1rem;
  padding-bottom: 0.5rem;
  width:  calc(100%-2rem);
  max-width: calc(100%-2rem);
  border-bottom: solid 1px #331111;
  background-color: #cc9977;
  color: #331111;
}`
