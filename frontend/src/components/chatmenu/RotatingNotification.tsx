import { useState, useEffect, useContext } from 'react';
import styled from 'styled-components';
import { UIContext } from '../../contexts/UIStateContext';

const RotatingNotification = ({
    messages = [],
    interval = 3000,
}: {
    messages?: string[];
    interval?: number;
}) => {
  const state = useContext(UIContext)
  const [currentIndex, setCurrentIndex] = useState(0);

  useEffect(() => {
    if (state.chatHeight != state.chatStartHeight) return;
    if (messages.length <= 1) return;
    
    const timer = setInterval(() => {
      setCurrentIndex((prev) => (prev + 1) % messages.length);
    }, interval);

    return () => clearInterval(timer);
  }, [messages, interval, state]);

  return (
    <NotificationContainer>
      <NotificationText>
        {messages[currentIndex]}
      </NotificationText>
    </NotificationContainer>
  );
};

const NotificationContainer = styled.div`
  overflow: hidden;
  height: 32px;
  display: flex;
  align-items: center;
  margin-right: 10px;
`;

const NotificationText = styled.span`
  color: white;
  font-size: 12px;
  animation: slideIn 0.5s ease-out;
  white-space: nowrap;

  @keyframes slideIn {
    from {
      transform: translateY(20px);
      opacity: 0;
    }
    to {
      transform: translateY(0);
      opacity: 0.8;
    }
  }
`;

export default RotatingNotification;
