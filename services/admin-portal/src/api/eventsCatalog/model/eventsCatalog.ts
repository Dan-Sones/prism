export type DataType = "string" | "int" | "float" | "boolean" | "timestamp";

export type EventField = {
  id: string;
  name: string;
  fieldKey: string;
  dataType: DataType;
};

export type EventType = {
  id: string;
  name: string;
  version: number;
  description: string | null;
  createdAt: string;
  fields: EventField[];
};

export type EventFieldRequest = {
  name: string;
  fieldKey: string;
  dataType: DataType;
};

export type CreateEventTypeRequest = {
  name: string;
  description?: string;
  fields: EventFieldRequest[];
};

export type FieldKeyAvailableResponse = {
  available: boolean;
};
