export interface ErrorData {
  message: string;
  type?: string; // The "?" indicates that the "type" property is optional
}

export interface ErrorMap {
  [key: string]: ErrorData;
}

export interface ErrorResponse {
  errors: ErrorMap;
}

export type ResponseError = {
  data: ErrorResponse;
};
