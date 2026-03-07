import type { RouteObject } from "react-router";
import EventsCatalog from "./list/EventsCatalog";
import CreateEvent from "./create-event/CreateEvent";
import Event from "./event/Event";

export const eventsCatalogRoutes: RouteObject[] = [
  { path: "/events-catalog", element: <EventsCatalog /> },
  { path: "/events-catalog/create", element: <CreateEvent /> },
  { path: "/events-catalog/:event_type_key", element: <Event /> },
];
