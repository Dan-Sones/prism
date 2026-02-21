export interface TableProps<T> {
  columns: Column<T>[];
  data: T[];
}

export interface Column<T> {
  header: string;
  accessor: string;
  sortable?: boolean;
  // I think key of t will ONLY work if the object is flat.

  // I can just adbie by this or come back and introduce nesting later.
  // for rendering pfps for owners this is going to be a pain, we would need to send in their name and
  // a url but this doesn't allow for that...
  render?: (value: T[keyof T], row: T) => React.ReactNode;
}

export interface TableRowProps<T> {
  row: T;
  columns: Column<T>[];
}
