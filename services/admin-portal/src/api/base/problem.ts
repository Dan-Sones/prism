export type Violation = {
  field: string;
  message: string;
};

export type ProblemDetail = {
  title?: string;
  status?: number;
  detail?: string;
  violations?: Violation[];
};
