import { useRef, useState, useEffect, useContext } from 'react';
import { useWindowDimensions } from 'react-native';
import { useNavigate } from 'react-router-dom';
import config from 'constants/Config';

export function useSession() {

  const [state, setState] = useState({
    tabbled: null,
    subWidth: '33%',
    baseWidth: '33%',
    cardId: null,
    converstaionId: null,
  });

  const dimensions = useWindowDimensions();
  const navigate = useNavigate();
  const tabbed = useRef(null);

  const updateState = (value) => {
    setState((s) => ({ ...s, ...value }));
  }

  useEffect(() => {
    if (tabbed.current !== true) {
      if (dimensions.width > config.tabbedWidth) {
        const width = Math.floor((dimensions.width * 33) / 100);
        tabbed.current = false;
        if (width > 500) {
          updateStatus({ tabbed: false, baseWidth: 550, subWidth: 500 });
        }
        else {
          updateState({ tabbed: false, baseWidth: width + 50, subWidth: width });
        }
      }
      else {
        tabbed.current = true;
        updateState({ tabbed: true });
      }
    }
  }, [dimensions]);

  const actions = {
    setCardId: (cardId) => {
      updateState({ cardId });
    }
  };

  return { state, actions };
}

