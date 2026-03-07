import type { RouteObject } from "react-router";
import EventsCatalog from "./list/EventsCatalog";
import CreateEvent from "./create-event/CreateEvent";

export const eventsCatalogRoutes: RouteObject[] = [
  { path: "/events-catalog", element: <EventsCatalog /> },
  { path: "/events-catalog/create", element: <CreateEvent /> },
];
