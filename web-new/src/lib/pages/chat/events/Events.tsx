import React from "react";
import { useAppSelector } from "../../../feature/hooks";
import { selectEventsSourceType } from "../../../feature/events/events.selectors";
import { EventSourceType } from "../../../feature/events/events.slice";
import { SSE } from "./SSE";
import { WS } from "./WS";

interface EventsProps {
  children: React.ReactNode;
}

export const Events: React.FC<EventsProps> = ({ children }) => {
  const eventsSourceType = useAppSelector((state) =>
    selectEventsSourceType(state),
  );

  return eventsSourceType === EventSourceType.SSE ? (
    <SSE>{children}</SSE>
  ) : (
    <WS>{children}</WS>
  );
};
