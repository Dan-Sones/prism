export type DataType = "string" | "int" | "float" | "boolean" | "timestamp";

export type EventField = {
  id: string;
  name: string;
  field_key: string;
  data_type: DataType;
};

export type EventType = {
  id: string;
  name: string;
  event_key: string;
  version: number;
  description: string | null;
  created_at: string;
  fields: EventField[];
};

export type EventFieldRequest = {
  name: string;
  field_key: string;
  data_type: DataType;
};

export type CreateEventTypeRequest = {
  name: string;
  event_key: string;
  description?: string;
  fields: EventFieldRequest[];
};

export type FieldKeyAvailableResponse = {
  available: boolean;
};

export type EventKeyAvailableResponse = {
  available: boolean;
};
