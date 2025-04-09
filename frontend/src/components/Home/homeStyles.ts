import styled from 'styled-components'
import { devices } from '../../globalStyles'

export const Row = styled.div`
  padding: 1rem 0;
  display: flex;
  flex-wrap: wrap;
`
// export const Title = styled.div``

export const Description = styled.div`
  padding: 0;

  @media ${devices.mobile} {
    max-width: 246px;
    max-height: 122px;
    overflow-y: scroll;
  }

  @media ${devices.desktop} {
    max-width: 384px;
  }
  // grid-row-start: 1;
  // grid-row-end: 3;
  // grid-column-start: 2;
  // grid-column-end: 3;
`

export const StartGrid = styled.section`
  margin: 1rem 0 0 0;
  display: grid;
  // grid-template-columns: repeat(auto-fit, clamp(10rem, 15rem, 30rem));
  gap: 1rem;
  // max-width: 100%;
`
