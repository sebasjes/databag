import styled from 'styled-components';
import Colors from 'constants/Colors';

export const IdentityWrapper = styled.div`
  width: 100%;
  height: 64px;
  display: flex;
  flex-direction: row;
  align-items: center;
  padding-left: 16px;
  padding-right: 16px;
  border-bottom: 2px solid ${Colors.divider};
  background-color: ${Colors.formBackground};
  flex-shrink: 0;

  &:hover {
    cursor: pointer;

    .drop {
      border: 1px solid ${Colors.encircle};
      background-color: ${Colors.formHover};
    }
  }

  .drop {
    padding-left: 4px;
    padding-right: 4px;
    border-radius: 8px;
    border: 1px solid ${Colors.formBackground};
  }

  .label {
    flex-grow: 1;
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;

    .name {
      font-size: 1.2em;
    }

    .handle {
      display: flex;
      flex-direction: row;
      justify-content: center;
      align-items: center;
      font-weight: bold;

      .notice {
        width: 32px;
        display: flex;
        flex-direction: row;
        justify-content: center;
      }
    }
  }
`;

export const ErrorNotice = styled.div`
  color: ${Colors.alert};
`

export const InfoNotice = styled.div`
  color: ${Colors.primary};
`

