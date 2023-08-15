import React from "react";
import { useAppSelector } from "../../../feature/hooks";
import { selectEventsSourceType } from "../../../feature/events/events.selectors";
import { EventSourceType } from "../../../feature/events/events.slice";
import { SSE } from "./SSE";
import { WS } from "./WS";

export const Events: React.FC = () => {
  const eventsSourceType = useAppSelector((state) =>
    selectEventsSourceType(state),
  );

  return eventsSourceType === EventSourceType.SSE ? <SSE /> : <WS />;
};
