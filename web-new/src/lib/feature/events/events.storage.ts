import { EventSourceType } from "./events.slice";

export const eventsStorage = {
  getEventsType: (): EventSourceType => {
    const eventsType = localStorage.getItem("eventsType");
    if (eventsType === "") return EventSourceType.SSE;

    return eventsType as EventSourceType;
  },
  setEventsType: (events: EventSourceType) => {
    localStorage.setItem("eventsType", events);
  },
};
