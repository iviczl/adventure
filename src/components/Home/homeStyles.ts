import styled from 'styled-components'

export const Row = styled.div`
  padding: 1rem 0;
  display: flex;
  flex-wrap: wrap;
`
export const Description = styled.div`
  padding: 1rem 0;
  grid-row-start: 1;
  grid-row-end: 3;
  grid-column-start: 2;
  grid-column-end: 3;
`

export const StartGrid = styled.section`
  margin: 1rem 0 0 0;
  display: grid;
  grid-template-columns: clamp(10rem, 15rem, 30rem) clamp(3rem, 5rem, 10rem);
  gap: 1rem;
  max-width: 100%;
`
