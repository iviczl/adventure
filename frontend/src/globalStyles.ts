import styled, { createGlobalStyle } from 'styled-components'

export interface SizeProps {
  $width?: string
  $height?: string
  $marginLeft?: string
  $marginRight?: string
  $marginTop?: string
  $marginBottom?: string
  $paddingLeft?: string
  $paddingRight?: string
  $paddingTop?: string
  $paddingBottom?: string
}

export interface ContainerProps extends SizeProps {
  $visible?: boolean
}

export interface LabelProps extends SizeProps {
  $color?: string
}

export const devices = {
  mobile: `(max-width: 1023px)`,
  desktop: `(min-width: 1024px)`,
}

const GlobalStyle = createGlobalStyle`
*,
*:before,
*:after {
  box-sizing: inherit;
  color: #331111;
  font-family: Inter, system-ui, Avenir, Helvetica, Arial, sans-serif;
  line-height: 1.5;
  font-weight: 400;
  
  font-synthesis: none;
  text-rendering: optimizeLegibility;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  padding: 0;
  margin: 0;
}

#root {
  box-sizing: border-box;
}
`
export default GlobalStyle

export const MainHeading = styled.h1`
  font-size: clamp(1.5rem, 6vw, 5rem);
  margin-bottom: 1rem;
  margin-top: 1rem;
`
export const SubHeading = styled.h5`
  font-size: 1rem;
  font-weight: 500;
  padding: 0;
  margin-bottom: 0;
  margin-top: 0;
`
export const MainContainer = styled.div`
  box-sizing: inherit;
  margin: 0;
  background-color: #cc9977;
  display: flex;
  place-items: center;
  flex-flow: column;
  min-width: 320px;
  max-width: 100%;
  width: 100%;
  height: 100vh;
  max-height: 100%;
  border: solid 1px #331111;
  border-radius: clamp(5px, 20px, 25px);
  // overflow-y: auto;
  // padding: 2rem;
`
export const Container = styled.div<ContainerProps>`
  visibility: ${(props) =>
    props.$visible === undefined
      ? 'visible'
      : props.$visible
      ? 'visible'
      : 'hidden'};
  box-sizing: inherit;
  margin-left: auto;
  margin-right: auto;
  padding-top: ${(props) => props.$paddingTop || '2rem'};
  padding-bottom: ${(props) => props.$paddingBottom || '2rem'};
  padding-left: ${(props) => props.$paddingLeft || '2rem'};
  padding-right: ${(props) => props.$paddingRight || '2rem'};
  display: flex;
  flex-flow: column;
  place-items: center;
  gap: 1rem;
  min-width: 320px;
  max-width: 800px;
  width: ${(props) => props.$width || '100%'};
  height: ${(props) => props.$height || '100%'};
  min-height: 10vh;
  max-height: 100%;
  overflow-y: auto;
  // border: solid 1px #331111;
  border-radius: clamp(5px, 20px, 25px);
  background-color: #cc9977;
`
export const Row = styled.div<SizeProps>`
  padding: 0;
  display: flex;
  flex-wrap: wrap;
  width: ${(props) => props.$width || 'unset'};
  height: ${(props) => props.$height || 'unset'};
  max-width: 100%;
`

export const Button = styled.button`
  border-radius: 8px;
  border: 1px solid transparent;
  padding: 0.6em 1.2em;
  font-size: 1rem;
  font-weight: 500;
  font-family: inherit;
  background-color: #eeccbb;
  color: #331111;
  border: solid 1px #331111;

  cursor: pointer;
  transition: border-color 0.25s;
  &:disabled {
    background-color: #ffddcc;
    color: #bbcccc;
    border: solid 1px #ffddcc;
    cursor: not-allowed;
  }

  &:not(:disabled):hover {
    border-color: white;
  }

  &:focus,
  &:focus-visible {
    outline: 4px auto -webkit-focus-ring-color;
  }
}
`
export const Paragraph = styled.p`
  margin: 1rem 0 0.5rem 0;
`

export const Label = styled.label<LabelProps>`
  font-size: 1rem;
  margin-left: ${(props) => props.$marginLeft || '1rem'};
  font-weight: 500;
  font-family: inherit;
  line-height: 0.8rem; //calc(2.2rem + 2px);
  padding: 0.5rem 0;
  width: ${(props) => props.$width || 'unset'};
  color: ${(props) => props.$color || 'inherited'};
`
export const Option = styled.option`
  padding: 0.6rem 0.6rem 0.6rem 0.6rem;
  font-size: inherit;
  font-weight: 500;
  font-family: inherit;
  background-color: white;
  color: #331111;
  height: 3rem;
`
export const Select = styled.select<SizeProps>`
  border-radius: 8px;
  margin-left: ${(props) => props.$marginLeft || '1rem'};
  padding: 0.6rem .6rem .6rem .6rem;
  font-size: inherit;
  font-weight: 500;
  font-family: inherit;
  background-color: white;
  color: #331111;
  border: solid 1px #331111;
  height: 3rem;
  max-width: 12rem;
  width: ${(props) => props.$width || 'unset'};
  // appearance: none;

  cursor: pointer;
  transition: border-color 0.25s;
  &:disabled {
    color: #bbcccc;
    border: solid 1px #ffddcc;
    cursor: not-allowed;
  }

  &:focus,
  &:focus-visible {
    outline: 4px auto -webkit-focus-ring-color;
  }
}
`
export const Input = styled.input<SizeProps>`
  border-radius: 8px;
  margin-left: ${(props) => props.$marginLeft || '1rem'};

  padding: 0.6rem .6rem;
  font-size: 1rem;
  font-weight: 500;
  font-family: inherit;
  background-color: white;
  color: #331111;
  border: solid 1px #331111;
  max-width: 12rem;
  width: ${(props) => props.$width || 'unset'};
  height: 3rem;

  cursor: pointer;
  transition: border-color 0.25s;
  &:disabled {
    color: #bbcccc;
    border: solid 1px #ffddcc;
    cursor: not-allowed;
  }

  &:focus,
  &:focus-visible {
    outline: 4px auto -webkit-focus-ring-color;
  }
}
`
